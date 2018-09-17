package main

import (
	sph "github.com/fpawel/bio3/internal/comport"
	"github.com/fpawel/bio3/internal/device/data"
	prd "github.com/fpawel/bio3/internal/products"
	"github.com/lxn/walk"
	"github.com/tarm/serial"

	"github.com/fpawel/bio3/internal/walkutils"

	"fmt"
	"github.com/lxn/win"
	"log"
	"time"

	"github.com/fpawel/bio3/internal/fetch"
	"github.com/fpawel/bio3/internal/utils"
)

type ReleStateValue struct {
	bool
}

func (x ReleStateValue) String() string {
	return ""
}



type ProductDataReader func(p prd.ProductInfo) error

type WorkProdResults map[string][]prd.ProductInfo

func AddWorkProd(x WorkProdResults, s string, p prd.ProductInfo) WorkProdResults {
	if x == nil {
		x = make(WorkProdResults)
	}
	x[s] = append(x[s], p)
	return x
}


func (x WorkProdResults) WriteLog(tx prd.Tx, app *App, level int) {
	for str, ps := range x {
		var ps_ []int
		for _, p := range ps {
			ps_ = append(ps_, int(p.Addr))
		}

		app.WriteWorkLog(tx, nil, level, fmt.Sprintf("%s: %s", utils.FormatIntRanges(ps_), str) )
	}
}



func (x *App) ProcessEachProduct(productDataReader ProductDataReader) error {
	defer x.mw.tblParty.Synchronize(func() {
		x.tableProductsModel.SetSurveyRow(-1)
	})
	var errs WorkProdResults

	for _, p := range x.db.Products() {
		p :=  p
		x.mw.tblParty.Synchronize(func() {
			x.tableProductsModel.SetSurveyRow(p.Row)
		})
		err := productDataReader(p)

		if err == fetch.ErrorCanceled {
			return fetch.ErrorCanceled
		}
		if err != nil {
			errs = AddWorkProd(errs, err.Error(), p)
		}

		x.mw.tblParty.Synchronize(func() {
			m := walkutils.MessageFromError(err, "успешно")
			mm := x.tableProductsModel
			mm.ProductInfo(p.ProductTime).Connection = &m
		})
	}

	if len(errs) > 0 {
		x.db.Update(func(tx prd.Tx) {
			errs.WriteLog(tx, x, win.NIIF_ERROR)
		})
	}

	return nil
}

func (x *App) Delay(d time.Duration, what string) error {

	satrtTime := time.Now()
	x.workHelper.cancelDelay = 0

	timeKey := prd.TimeToKey(time.Now())

	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, timeKey, 0, fmt.Sprintf("задержка %s, %v", what, d) )
	})

	work := x.tableWorksModel.CurrentWork()
	if work == nil {
		log.Fatal("work == nil")
	}

	x.mw.delayControl.Run(d, what,
		func() {
			check(x.mw.Invalidate())
		}, func() {
			x.db.Update(func(tx prd.Tx) {
				tx.WriteLog(work.PartyLogPath(tx).Path(), timeKey,
					0, fmt.Sprintf("задержка, %s, %v - %v", what, time.Since(satrtTime), d))
				x.mw.Synchronize(func() {
					x.tableLogsModel.PublishRowsReset()
					x.tableWorksModel.PublishRowChanged(x.tableWorksModel.currentWorkIndex)
				})
			})
		} )

	x.db.Update(func(tx prd.Tx) {
		tx.WriteLog(work.PartyLogPath(tx).Path(), timeKey,
			0, fmt.Sprintf("задержка, %s, %v", what, d))
		x.mw.Synchronize(func() {
			x.tableLogsModel.PublishRowsReset()
			x.tableWorksModel.PublishRowChanged(x.tableWorksModel.currentWorkIndex)
		})
	})

	if x.workHelper.port.Canceled() {
		return fetch.ErrorCanceled
	}
	return nil
}


func (x *App) RunMainWorks() {
	x.RunWorks("Выбранные проверки", data.MainWorks)
}

func (x *App) RunWorks(what string, works []*data.Work) {
	portName := x.mw.PortName()
	portConfig := &serial.Config{
		Name:        portName,
		ReadTimeout: time.Millisecond,
		Baud:        9600,
	}

	notifyPortError := func(what string, err error) {
		PrintLog(win.NIIF_ERROR, fmt.Sprintf("%s: COMPORT: \"%s\": %v", what, portName, err))
		x.mw.setNotifyMessage(&NotifyMessage{"Не удалось открыть СОМ порт", err.Error() + ", " + what, win.NIIF_ERROR})
	}
	x.workHelper = &WorkHelper{
		switchKey:  make(map[byte]byte),
	}

	var err error
	var canceled int32
	x.workHelper.port, err = sph.Open(portConfig, x.config.ReadWrite,!(len(works)==1 && works[0] == data.ManualSurvey), &canceled)

	if err != nil {
		notifyPortError("начало настройки", err)
		return
	}
	x.tableProductsModel.Clear()

	tableWorksModel := x.tableWorksModel

	tableWorksModel.SetWorks(works)
	check(x.mw.tblWorks.SetCurrentIndex(0))
	x.mw.setupProductsColumns()
	x.mw.InvalidateWorkRunning(true)

	x.mw.setNotifyMessage(&NotifyMessage{ what, "выполняется", win.NIIF_INFO})

	startTime := time.Now()

	go func() {
		for i, work := range works {
			work  := work

			i := i
			if !x.tableWorksModel.Checked(i) {
				continue
			}
			x.tableWorksModel.SetCurrentWork(work)

			x.mw.tblWorks.Synchronize(func() {
				check(x.mw.tblWorks.SetCurrentIndex(i))
				x.tableWorksModel.PublishRowsReset()
			})

			x.mw.tblParty.Synchronize(func() {
				x.mw.setupProductsColumns()
			})

			if tableWorksModel.WorksIsMainTests() {
				// очистить предыдущий результат проверки
				x.db.Update(func(tx prd.Tx) {
					for _, rele := range data.PartyProductType(tx.Party()).Reles {
						for _, p := range tx.Party().Products() {
							data.SetTestReleResult( p, work, rele, nil)
						}
					}
				})
			}
			var logPath [][]byte
			logKey := prd.TimeToKey(time.Now())
			x.db.Update(func(tx prd.Tx) {
				logPath = work.PartyLogPath(tx).Path()
				tx.ClearLogs(logPath)
				x.WriteWorkLog(tx, logKey, 0, "Выполняется")
			})
			err := work.Work()
			{
				msg := &NotifyMessage{ Title: work.What}
				x.db.Update(func(tx prd.Tx) {
					tx.DeleteLog(logPath, logKey)
					if err != nil {
						x.WriteWorkLog(tx, logKey, win.NIIF_ERROR, err.Error())
						msg.Text = err.Error()
						msg.Level = win.NIIF_ERROR
					} else {
						// взять запись с наибольшим левелем
						if l := tx.MostImportantLogRecord(logPath); l != nil {
							msg.Text = l.Text
							msg.Level = l.Level
						} else {
							x.WriteWorkLog(tx, logKey, win.NIIF_INFO, "выполнено",  )
							msg.Text = "выполнено"
							msg.Level = win.NIIF_INFO
						}
					}
				})
				x.mw.setNotifyMessage(msg)

			}

			if err == fetch.ErrorCanceled {
				break
			}
		}

		if tableWorksModel.WorksIsMainTests() {
			canceled = 0
			if err := x.StopHardware(); err != nil  {
				notifyPortError("остановка оборудования", err)
			}
		}

		if err := x.workHelper.port.Close(); err != nil {
			notifyPortError("не удалось закрыть порт", err)
		}

		x.workHelper = nil
		tableWorksModel.currentWorkIndex = -1

		msg := fmt.Sprintf( "Сценарий %q выполнен без ошибок", what )
		style := walk.MsgBoxOK
		x.db.View(func(tx prd.Tx) {
			for t,r := range tx.Logs(tx.Party().Path(), nil) {
				if t.After(startTime) && r.Level == win.NIIF_ERROR {
					msg = fmt.Sprintf( "Сценарий %q выполнен c ошибками.", what )
					style |= walk.MsgBoxIconError
				}
			}
		})

		x.mw.Synchronize(func() {

			x.tableProductsModel.SetSurveyRow(-1)

			check(x.mw.tblWorks.SetCurrentIndex(-1))
			x.tableWorksModel.PublishRowsReset()

			walk.MsgBox(x.mw, "Отчёт о выполнении сценария", msg , style)
			x.mw.InvalidateWorkRunning(false)
		})
	}()
}

func (x *App) StopHardware() error {
	log.Println("остановка оборудования")


	if x.workHelper.sound {
		if err := x.Sound(false); err != nil {
			return fmt.Errorf("не удалось отключить звук: %v", err)
		}
	}
	if x.workHelper.pneumoPoint > 0 {
		if err := x.SwitchPneumo(0); err != nil {
			return fmt.Errorf("не удалось отключить пневмоблок: %v", err)
		}
	}
	for n,state := range x.workHelper.switchKey{
		if state > 0 {
			if err := x.Switch(n,0); err != nil {
				return fmt.Errorf("не удалось отключить ключ %d: %v", n, err)
			}
		}
	}
	return nil
}
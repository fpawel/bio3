package main

import (
	"github.com/lxn/walk"
	//. "github.com/lxn/walk/declarative"
	"fmt"
	"github.com/fpawel/bio3/device/data"
	"github.com/fpawel/bio3/comport"
	"github.com/fpawel/bio3/walkutils"
	"log"
	"time"
	"github.com/fpawel/bio3/products"
	"github.com/fpawel/bio3/walkutils/delay"
)

type AppMainWindow struct {
	*walk.MainWindow
	app *App
	notifyMessage,
	comportInfo *NotifyMessageTime
	notifyIcon *walk.NotifyIcon
	tblParty,
	tblWorks,
	tblLogs   *walk.TableView
	menuParty  *walk.Menu
	cbProductsType,
	cbPort *walk.ComboBox
	btnRunMenu *walk.SplitButton
	btnTests,
	btnCancel *walk.PushButton
	skipTableWorksCurrentIndexChanged bool
	delayControl *delay.Control
	initialized bool
    userInput	*walk.Action
}

type NotifyMessage struct {
	Title,Text  string
	Level int
}

type NotifyMessageTime struct {
	*NotifyMessage
	Time time.Time
}

func NewAppMainWindow (app *App)(*AppMainWindow) {
	return &AppMainWindow{
		app : app,
		delayControl:delay.New(AssetImage("assets/png16/cancel16.png")),
	}
}

func (x *AppMainWindow) Initialize() {
	x.createNotifyIcon()
	x.refreshSerials()
	check(x.cbProductsType.SetModel(data.ProductTypes))

	x.setupParty()
	cfg := x.app.config
	for i, s := range x.cbPort.Model().([]string) {
		if s == cfg.Port {
			check(x.cbPort.SetCurrentIndex(i))
			break
		}
	}
	x.delayControl.SetVisible(false)

	//check(x.tblWorks.SetCurrentIndex(0))
	x.app.tableWorksModel.SetWorks(data.MainWorks)
	check(x.Invalidate())
	x.initialized = true

}


func (x *AppMainWindow) setupParty() {
	x.app.db.View(func(tx products.Tx) {
		t := time.Time(tx.Party().PartyTime)
		check(x.SetTitle(fmt.Sprintf("%s Партия %s", x.app.Title(), t.Format("02 01 2006, 03:04"))))
		check(x.cbProductsType.SetCurrentIndex(tx.Party().ProductTypeIndex()))
	})
	x.app.tableProductsModel.SetColumns(x.tblParty.Columns())
}

func (x *AppMainWindow) setupProductsColumns() {
	x.app.tableProductsModel.SetColumns(x.tblParty.Columns())
}

func (x *AppMainWindow) setNotifyMessage(m *NotifyMessage) {
	if x.notifyMessage != nil && *x.notifyMessage.NotifyMessage == *m && time.Since(x.notifyMessage.Time) < 5*time.Second {
		return
	}
	x.notifyMessage = &NotifyMessageTime{
		NotifyMessage: m,
		Time:          time.Now(),
	}

	if err := walkutils.ShowNotifyMessage(x.notifyIcon, m.Level)(x.app.ProductName() + ". " + m.Title, m.Text); err != nil {
		log.Panic(err)
	}

}

func (x *AppMainWindow) createNotifyIcon() {
	// We load our icon from a file.
	//icon, err := walk.NewIconFromFile("./img/app.ico")
	//check(err)

	var err error
	x.notifyIcon, err = walk.NewNotifyIcon()
	check(err)

	// Set the icon and a tool tip text.
	check(x.notifyIcon.SetIcon(NewIconFromResourceId(IconAppID)))
	check(x.notifyIcon.SetToolTip(x.app.ProductName()))

	// The notify icon is hidden initially, so we have to make it visible.
	check(x.notifyIcon.SetVisible(true))
}

func (x *AppMainWindow) PortName() string {
	n := x.cbPort.CurrentIndex()
	xs := x.cbPort.Model().([]string)
	if n >= 0 && n < len(xs) {
		return xs[n]
	}
	return ""
}

func (x *AppMainWindow) refreshSerials() {
	xs := comport.AvailablePorts()
	s := x.cbPort.Text()
	if s == "" {
		check(x.cbPort.SetModel(xs))
		return
	}
	for _, v := range xs {
		if v == s {
			check(x.cbPort.SetModel(xs))
			return
		}
	}
	check(x.cbPort.SetModel(append(xs, s)))
}

func (x *AppMainWindow) comboBoxPortMouseDown(_, _ int, _ walk.MouseButton) {
	x.refreshSerials()
}


func (x *AppMainWindow) InvalidateWorkRunning(isRunning bool) {

	x.btnRunMenu.SetEnabled(!isRunning)
	x.btnCancel.SetVisible(isRunning)

	acts := x.btnRunMenu.Menu().Actions()
	for i := 0; i < acts.Len(); i++ {
		check(acts.At(i).SetVisible(!isRunning))
	}
	if isRunning {
		x.btnTests.SetEnabled(false)
	} else {
		x.btnTests.SetEnabled(!x.app.tableWorksModel.WorksIsMainTests())
	}
	x.userInput.SetVisible( prod=="" )
}

func (x *AppMainWindow) SelectedTest() *data.Work {
	if n := x.tblWorks.CurrentIndex(); n >= 0 && n < x.app.tableWorksModel.RowCount() {
		return x.app.tableWorksModel.works[n]
	}
	return nil
}

func (x *AppMainWindow) TableWorksCurrentIndexChanged() {

	if x.skipTableWorksCurrentIndexChanged || x.app.tableWorksModel.RowCount() == 0 {
		return
	}
	x.skipTableWorksCurrentIndexChanged = true
	if x.tblWorks.CurrentIndex() == -1 {
		check(x.tblWorks.SetCurrentIndex(0))
	}
	x.setupProductsColumns()
	x.skipTableWorksCurrentIndexChanged = false

}


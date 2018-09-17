package main

import (
	"fmt"
	"github.com/fpawel/bio3/internal/device/data"
	stnd "github.com/fpawel/bio3/internal/device/stend"
	"github.com/fpawel/bio3/internal/fetch"
	prd "github.com/fpawel/bio3/internal/products"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"time"

	"github.com/fpawel/bio3/internal/walkutils"
)


func (x *App) ManualSurvey() error {
	m := x.tableProductsModel
	defer x.mw.tblParty.Synchronize(func() {
		m.SetSurveyRow(-1)
	})

	for {

		for _, p := range x.db.Products() {
			p :=  p
			x.mw.tblParty.Synchronize(func() {
				m.SetSurveyRow(p.Row)
			})
			b, err := stnd.Rele(x.workHelper.port, p.Addr)

			x.mw.tblParty.Synchronize(func() {
				connection := walkutils.MessageFromError(err, "ok")
				m.ProductInfo(p.ProductTime).Connection = &connection
				for _, rele := range data.Reles {
					if err == nil {
						m.ProductInfo(p.ProductTime).Reles[rele] = data.ParseReleStateValue(rele, b)
					} else {
						delete(m.ProductInfo(p.ProductTime).Reles, rele)
					}
				}
			})

			if err == fetch.ErrorCanceled {
				return fetch.ErrorCanceled
			}
		}
	}
}

func (x *App) Power(value stnd.PowerState) error {

	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil, 0, value.String())
	})

	return x.ProcessEachProduct(func(p prd.ProductInfo) error {
		return stnd.Power(x.workHelper.port, p.Addr, value)
	})
}

func (x *App) Sound(on bool ) error {

	what := "подача звукового сигнала"
	if !on {
		what = "отключение звукового сигнала"
	}
	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil,0, what )
	})

	var state byte
	var ampl,freq float32
	if on {
		state = 1
		ampl = 1
		freq = 400
	}

	for addr:=byte(1); addr<= 20; addr++{
		if err := stnd.Vibration(x.workHelper.port, addr, state, ampl, freq); err != nil {
			return err
		}
	}
	x.workHelper.sound = on
	return nil
}

func (x *App) Reset() error {

	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil, 0, "сброс", )
	})

	if err := x.ProcessEachProduct(func(p prd.ProductInfo) error {
		return  stnd.Switch(x.workHelper.port, p.Addr, 3, 1)
	}); err != nil {
		return err
	}

	if err := x.Delay(time.Second, "сброс"); err != nil {
		return err
	}

	return x.ProcessEachProduct(func(p prd.ProductInfo) error {
		return  stnd.Switch(x.workHelper.port, p.Addr, 3, 0)
	})
}

func (x *App) PowerOn() error {
	return x.Power(stnd.PowerOn)
}
func (x *App) PowerOff() error {
	return x.Power(stnd.PowerOff)
}


func (x *App) LedPlacesOn()  (err error) {

	ns := make( map[byte] struct{} )
	for _, p := range x.db.Products() {
		ns[p.Addr] = struct{}{}
	}

	for i := byte(0); i<20; i++ {
		if err := stnd.LedTurnOff(x.workHelper.port, i+1); err != nil {
			return err
		}
	}

	return x.ProcessEachProduct( func(p prd.ProductInfo) error {
		return stnd.LedGreenTurnOn(x.workHelper.port, p.Addr)
	})
}

func (x *App) SwitchPneumo(point byte) (err error) {

	rate := float32(0.5)
	count := len(x.db.Products())
	what := fmt.Sprintf("подать ПГС%d", point)
	if point == 0 {
		what = 	"отключить пневмоблок"
		count = 0
		rate = 0
	} else {
		rate *= float32(count)
	}
	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil,0, what, )
	})

	if err = stnd.Gas(x.workHelper.port, rate, point , byte(count) ); err == nil || err == fetch.ErrorCanceled {
		x.workHelper.pneumoPoint = point
	}
	return
}

func (x *App) Switch(n byte, state byte) error {
	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil, 0, fmt.Sprintf("ключ %d: %d", n, state), )
	})

	err := x.ProcessEachProduct(func(p prd.ProductInfo) error {
		return stnd.Switch(x.workHelper.port, p.Addr, n, state)
	})
	if err == nil || err == fetch.ErrorCanceled{
		x.workHelper.switchKey[n] = state
	}
	return err
}

func (x *App) FixTest(test *data.Work) error {
	var errs, oks WorkProdResults
	result := x.ProcessEachProduct(func(p prd.ProductInfo) error {

		b, err := stnd.Rele(x.workHelper.port, p.Addr)
		if err != nil {
			return err
		}

		x.db.Update(func(tx prd.Tx) {
			var s string
			for _, rele := range data.PartyProductType(tx.Party()).Reles {
				v := data.ParseReleStateValue(rele, b)
				data.SetTestReleResult(tx.Party().Products()[p.Row], test, rele, &v )
				if releValue :=test.ReleValue(rele); releValue != nil && *releValue != v {
					if s != "" {
						s += ", "
					}
					s += rele.What
				}
			}
			if  s != "" {
				errs = AddWorkProd(errs, fmt.Sprintf("не соответствует требованиям: %s", s), p )
			} else {
				oks = AddWorkProd(oks, "соответствует требованиям", p )
			}
		})
		return err
	})

	if len(oks) > 0 {
		x.db.Update(func(tx prd.Tx) {
			oks.WriteLog(tx, x, win.NIIF_INFO)
		})
	}
	if len(errs) > 0 {
		x.db.Update(func(tx prd.Tx) {
			errs.WriteLog(tx, x, win.NIIF_ERROR)
		})
	}
	return result
}

func (x *App) TestPowerOn()  error {

	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil, 0, "снять питание", )
	})
	if err := x.PowerOff(); err != nil {
		return err
	}

	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil, 0, "подать сброс", )
	})
	if err := x.ProcessEachProduct(func(p prd.ProductInfo) error {
		return  stnd.Switch(x.workHelper.port, p.Addr, 3, 1)
	}); err != nil {
		return err
	}
	if err := x.Delay(time.Second, "подача сброса"); err != nil {
		return err
	}

	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil, 0, "подать питание", )
	})
	if err := x.PowerOn(); err != nil {
		return err
	}
	if err := x.Delay(5 * time.Second, "подача питания"); err != nil {
		return err
	}


	x.db.Update(func(tx prd.Tx) {
		x.WriteWorkLog(tx, nil, 0, "снять сброс", )
	})
	if err := x.ProcessEachProduct(func(p prd.ProductInfo) error {
		return  stnd.Switch(x.workHelper.port, p.Addr, 3, 0)
	}); err != nil {
		return err
	}

	if err := x.Delay(3 * time.Minute, "включение питания"); err != nil {
		return err
	}
	return x.FixTest(data.TestPowerOn)
}

func (x *App) TestIncline() (err error) {
	walk.MsgBox(x.mw, data.TestIncline.What,
		"Произведите наклон проверяемых извещателей и нажмите OK",
		walk.MsgBoxIconInformation)
	if err := x.Delay( 10 * time.Second, "наклон"); err != nil {
		return err
	}

	if err = x.FixTest(data.TestIncline); err != nil {
		return
	}
	walk.MsgBox(x.mw, data.TestIncline.What,
		"Верните излучатели в исходное положение и нажмите OK",
		walk.MsgBoxIconInformation)
	if err := x.Delay( 10 * time.Second, "возврат в исходное положение"); err != nil {
		return err
	}

	return x.Reset()
}

func (x *App) TestVibration() error {

	if err := x.Sound(true); err == fetch.ErrorCanceled {
		return fetch.ErrorCanceled
	}
	if err := x.Delay( 10 * time.Second, "подача звукового сигнала"); err != nil {
		return err
	}

	if err := x.FixTest(data.TestVibration); err == fetch.ErrorCanceled {
		return fetch.ErrorCanceled
	}

	if err := x.Sound(false); err == fetch.ErrorCanceled {
		return fetch.ErrorCanceled
	}

	return x.Reset()
}

func (x *App) TestSmoke()  error {

	walk.MsgBox(x.mw, data.TestSmoke.What, "Установите меры оптические и нажмите OK",
		walk.MsgBoxIconInformation)

	if err := x.Delay( 10 * time.Second, "имитация задымления"); err != nil {
		return err
	}

	if err := x.FixTest(data.TestSmoke); err != nil {
		return fetch.ErrorCanceled
	}

	walk.MsgBox(x.mw, data.TestSmoke.What, "Удалите меры оптические и нажмите OK",
		walk.MsgBoxIconInformation)

	if err := x.Delay( 10 * time.Second, "снятие задымления"); err != nil {
		return err
	}

	return x.Reset()
}

func (x *App) Adjust0() error {

	// подать ПГС1
	if err := x.SwitchPneumo(1); err != nil {
		return err
	}
	if err := x.Delay(30 * time.Second , "продувка ПГС1"); err != nil {
		return err
	}
	// замкнуть ключ 1
	if err := x.Switch(1,1); err != nil {
		return err
	}
	if err := x.Delay(time.Second , "ключ 1"); err != nil {
		return err
	}
	// разомкнуть ключ 1
	if err := x.Switch(1,0); err != nil {
		return err
	}
	// продувка ПГС1
	if err := x.Delay(3 * time.Minute , "продувка ПГС1"); err != nil {
		return err
	}
	// отключить газ
	return x.SwitchPneumo(0)
}

func (x *App) Adjust1() error {

	// замкнуть ключ 2
	if err := x.Switch(2,1); err != nil {
		return err
	}
	// подать газ 3
	if err := x.SwitchPneumo(3); err != nil {
		return err
	}
	if err := x.Delay(time.Minute, "продувка ПГС3"); err != nil {
		return err
	}
	// разомкнуть ключ 2
	return x.Switch(2, 0)

}

func (x *App) TestAir() error {

	// подать газ 1
	if err := x.SwitchPneumo(1); err != nil {
		return err
	}
	if err := x.Delay(time.Minute, "продувка ПГС1"); err != nil {
		return err
	}

	// выключить питание
	if err := x.Power(stnd.PowerOff); err != nil {
		return err
	}
	if err := x.Delay(10 * time.Second, "выключение питания"); err != nil {
		return err
	}

	// включить питание
	if err := x.Power(stnd.PowerOn); err != nil {
		return err
	}
	if err := x.Delay(2 * time.Minute, "включение питания"); err != nil {
		return err
	}
	// фиксировать сотояние реле
	return x.FixTest(data.TestAir)
}

func (x *App) TestGas() error {

	// подать газ 2
	if err := x.SwitchPneumo(2); err != nil {
		return err
	}
	if err := x.Delay(time.Minute, "продувка ПГС2"); err != nil {
		return err
	}
	// фиксировать сотояние реле
	return x.FixTest(data.TestGas)
}


func (x *App) initWorks() {
	data.ManualSurvey.Work = x.ManualSurvey
	data.TestPowerOn.Work = x.TestPowerOn
	data.TestIncline.Work = x.TestIncline
	data.TestVibration.Work = x.TestVibration
	data.TestSmoke.Work = x.TestSmoke
	data.Adjust0.Work = x.Adjust0
	data.Adjust1.Work = x.Adjust1
	data.TestAir.Work = x.TestAir
	data.TestGas.Work = x.TestGas
}


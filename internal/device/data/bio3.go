package data

import (
	"github.com/fpawel/bio3/internal/products"
	"log"
)

type Rele struct {
	What string
}

var ReleInclineVibration = &Rele{
	What: "НАКЛОН/ВИБРАЦИЯ",
}
var ReleSmoke = &Rele{
	What: "ДЫМ",
}
var ReleGas = &Rele{
	What: "ГАЗ",
}

var ReleAlarm = &Rele{
	What: "ТРЕВОГА",
}

var ReleSpec = &Rele{
	What: "СПЕЦ",
}

var ReleAlarm1 = &Rele{
	What: "ТРЕВОГА 1",
}

type ProductType struct {
	What  string
	Reles []*Rele
}

type Work struct {
	What  string
	Reles RelesValues
	Work func() error
}


type RelesValues map[*Rele]bool

var ProductTypes = []*ProductType{
	{
		What:  "БИО-3",
		Reles: [] *Rele {
			ReleAlarm,
			ReleInclineVibration,
			ReleSmoke,
			ReleGas,
		},
	},
	{
		What:  "БИО-3-01",
		Reles: []*Rele {
			ReleAlarm,
			ReleAlarm1,
		},
	},
}

var ManualSurvey = newWork("Опрос")

var TestPowerOn = &Work{
	What: "Проверка после подачи питания",
	Reles: RelesValues{
		ReleSmoke:            true,
		ReleGas:              true,
		ReleInclineVibration: true,
		ReleAlarm:            true,
		ReleAlarm1:           true,
	},
}

var TestIncline = &Work{
	What: "Проверка датчика наклона",
	Reles: RelesValues{
		ReleSmoke:            true,
		ReleGas:              true,
		ReleInclineVibration: false,
		ReleAlarm:            false,
		ReleAlarm1:           false,
	},
}
var TestVibration = &Work{
	What:  "Проверка датчика вибрации",
	Reles: TestIncline.Reles,
}
var TestSmoke = &Work{
	What: "Проверка датчика задымления",
	Reles: RelesValues{
		ReleSmoke:            false,
		ReleGas:              true,
		ReleInclineVibration: true,
		ReleAlarm:            false,
		ReleAlarm1:           false,
	},
}

var Adjust0 = &Work{
	What:  "Регулировка нулевых показаний датчика загазованности",
	Reles: RelesValues{},
}

var Adjust1 = &Work{
	What:  "Регулировка чувствительности датчика загазованности",
	Reles: RelesValues{},
}

var TestAir = &Work{
	What:  "Проверка отсутствия сигнализации ТРЕВОГА-ГАЗ",
	Reles: RelesValues{
		ReleSmoke:            true,
		ReleGas:              true,
		ReleInclineVibration: true,
		ReleAlarm:            true,
		ReleAlarm1:           true,
	},
}

var TestGas = &Work{
	What: "Проверка срабатывания сигнализации ТРЕВОГА-ГАЗ",
	Reles: RelesValues{
		ReleSmoke:            true,
		ReleGas:              false,
		ReleInclineVibration: true,
		ReleAlarm:            false,
		ReleAlarm1:           false,
	},
}

var MainWorks = []*Work{
	TestPowerOn, TestIncline, TestVibration, TestSmoke, Adjust0, Adjust1, TestAir, TestGas,
}

var MainTestsRele = []*Work{
	TestPowerOn, TestIncline, TestVibration, TestSmoke, TestAir, TestGas,
}



var Reles = []*Rele{
	ReleInclineVibration, ReleSmoke, ReleGas, ReleAlarm, ReleAlarm1,
}

func newWork(s string) *Work {
	return &Work{
		What:  s,
		Reles: make(RelesValues),
	}
}

func (x *Rele) Index() []byte {
	for i,y := range Reles {
		if x == y {
			return []byte{byte(i)}
		}
	}
	return nil
}

func (x *Work) RelesList(t *ProductType) (r []*Rele) {
	for _,a := range t.Reles {
		if _,has := x.Reles[a]; has {
			r = append(r,a)
		}

	}
	return
}

func (x *Work) Index() int {
	for n,y := range MainWorks {
		if y == x {
			return n
		}
	}
	return -1
}

func (x *Work) ReleValue(rele *Rele) *bool {
	if value,f := x.Reles[rele]; f {
		return &value
	}
	return nil
}



func ParseReleStateValue(rele *Rele, b []byte) bool {

	if rele  == nil {
		log.Fatal("rele must be not nil")
	}

	if rele == ReleAlarm1 {
		rele = ReleGas
	}

	for n,r := range releStateIndex {
		if r == rele {
			return b[n] != 0
		}
	}

	log.Fatalln("wrong rele", *rele)

	return false
}

var releStateIndex = []*Rele{
	ReleInclineVibration,
	ReleSmoke,
	ReleGas,
	ReleSpec,
	ReleAlarm,
}

func PartyProductType(party products.Party) *ProductType {
	n := party.ProductTypeIndex()
	if n < len(ProductTypes) {
		return ProductTypes[n]
	}
	return nil
}


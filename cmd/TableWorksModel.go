package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"

	"github.com/fpawel/bio3/internal/device/data"
	"github.com/fpawel/bio3/internal/products"
	"log"
)


type TableWorksModel struct {
	walk.ReflectTableModelBase
	unchecked map[int] struct{}
	works []*data.Work
	currentWorkIndex int
	db products.DB
}

func NewTableWorksViewModel(db products.DB) *TableWorksModel {
	return &TableWorksModel{
		db:db,
		unchecked:make(map[int] struct{}),
		currentWorkIndex: -1,
	}
}

func (x *TableWorksModel) CurrentWork() *data.Work {
	if x.currentWorkIndex >=0 && x.currentWorkIndex < len(x.works) {
		return x.works[x.currentWorkIndex]
	}
	return nil
}

func (x *TableWorksModel) SetCurrentWork(w *data.Work)  {
	if x.CurrentWork() == w {
		return
	}
	if w == nil {
		x.currentWorkIndex = -1
	} else {
		for i,work := range x.works{
			if w == work {
				x.currentWorkIndex = i
				break
			}
		}
	}
	x.PublishRowsReset()
}

func (x *TableWorksModel) RowCount() int {
	return len(x.works)
}

func (x *TableWorksModel) Value(row, col int) (result interface{}) {

	x.db.View(func(tx products.Tx) {
		work := x.works[row]

		switch col {
		case 0:
			result = work.What
			return
		case 1:
			r := tx.MostImportantLogRecord(work.PartyLogPath(tx).Path())
			if r == nil {
				result = ""
				return
			}
			result = r.Text
		default:
			log.Panicln("col out of range", col, row)
			result = ""
		}
	})
	return
}

func (x *TableWorksModel) Checked(row int) bool {
	_,f := x.unchecked[row]
	return !f
}

func (x *TableWorksModel) SetChecked(row int, checked bool) error {
	if checked {
		delete(x.unchecked, row)
	} else {
		x.unchecked[row] = struct{}{}
	}
	return nil
}



func (x *TableWorksModel) StyleCell(c *walk.CellStyle) {

	x.db.View(func(tx products.Tx) {
		switch c.Col() {
		case 1:
			if x.currentWorkIndex == c.Row() {
				c.Image =  mustImg("assets/png16/forward.png")
				return
			}
			if r := tx.MostImportantLogRecord(x.works[c.Row()].PartyLogPath(tx).Path()); r != nil {
				switch r.Level {
				case win.NIIF_ERROR:
					c.TextColor = walk.RGB(255, 0, 0)
					c.Image = ImgErrorPng16
				case win.NIIF_INFO:
					c.TextColor = walk.RGB(0, 32, 128)
					c.Image = ImgCheckmarkPng16
				}
			}
		}
	})
}

func (x *TableWorksModel) WorksIs(works []*data.Work) bool {
	return len(x.works) == len(works) && x.works[0] == works[0]
}

func (x *TableWorksModel) WorksIsMainTests() bool {
	return x.WorksIs(data.MainWorks)
}

func (x *TableWorksModel) SetWorks(works []*data.Work) {
	if x.WorksIs(works) {
		return
	}
	x.works = works
	x.unchecked = make(map[int]struct{})
	x.currentWorkIndex = -1
	x.PublishRowsReset()
}
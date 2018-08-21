package productsview

import (
	"github.com/fpawel/bio3/products"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"time"
)

type TableLogsViewModel struct {
	walk.ReflectTableModelBase
	db products.DB
	path func (tx products.Tx) [][]byte
}

func NewTableLogsViewModel(db products.DB, path func (tx products.Tx) [][]byte) *TableLogsViewModel {
	return &TableLogsViewModel{
		db:db,
		path : path,
	}
}


func (x *TableLogsViewModel) RowCount() (result int) {
	x.db.View(func(tx products.Tx) {
		result = len(tx.Logs(x.path(tx), nil))
	})
	return
}

func (x *TableLogsViewModel) Data(row, col int) (text string, image walk.Image, textColor *walk.Color, backgroundColor *walk.Color) {



	x.db.View(func(tx products.Tx) {
		var logRec *products.LogRecord
		var logTm time.Time
		logs := tx.Logs(x.path(tx), nil)
		for i,t := range  logs.Times() {
			if i == row {
				logRec = logs[t]
				logTm = t
				break
			}
		}
		if logRec == nil {
			//log.Println("TableLogViewModel Value row out of range", len(x.logs()), row)
			return
		}

		switch col {
		case 0:
			text = logTm.Format("02.01 15:04:05")
		case 1:
			for i := 0; i < len(logRec.Path) -1; i++{
				if string(logRec.Path[i]) == "tests"{
					text = string(logRec.Path[i+1])
					return
				}
			}
		case 2:
			text = logRec.Text
			switch logRec.Level {
			case win.NIIF_INFO:
				image = ImgCheckmarkPng16
			case win.NIIF_ERROR:
				image = ImgErrorPng16
				textColor = new(walk.Color)
				*textColor = walk.RGB(255, 0, 0)
			}
		}

	})
	return
}

func (x *TableLogsViewModel) Value(row, col int) (r interface{}) {
	str,_,_,_ := x.Data(row,col)
	return str
}

func (x *TableLogsViewModel) StyleCell(c *walk.CellStyle) {
	_, image , textColor , backgroundColor := x.Data(c.Row(),c.Col())
	if image != nil {
		c.Image = image
	}
	if textColor != nil {
		c.TextColor = *textColor
	}
	if backgroundColor != nil {
		c.BackgroundColor = *backgroundColor
	}
}



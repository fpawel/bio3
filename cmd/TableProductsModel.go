package main

import (
	"fmt"
	"github.com/fpawel/bio3/internal/device/data"
	"github.com/fpawel/bio3/internal/products"
	"github.com/fpawel/bio3/internal/walkutils"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type TableProductsModel struct {
	walk.ReflectTableModelBase
	db          products.DB
	surveyRow      int
	items map[products.ProductTime] *ProductInfo
	test func ()  *data.Work
}



type ProductInfo struct {
	Reles    map[*data.Rele]bool
	Connection *walkutils.Message
}

func NewTableProductsModel(db  products.DB, test func () *data.Work) (x *TableProductsModel) {
	return &TableProductsModel{
		db:          db,
		surveyRow:   -1,
		test:test,
		items:make(map[products.ProductTime] *ProductInfo ),
	}
}

func (x *TableProductsModel) Clear(){
	x.items = make(map[products.ProductTime] *ProductInfo )
	x.PublishRowsReset()
}

func (x *TableProductsModel) SurveyRow() int {
	return x.surveyRow
}

func (x *TableProductsModel) SetSurveyRow(row int) {
	if row != x.surveyRow {
		x.surveyRow = row
		x.PublishRowsReset()
	}
}

func (x *TableProductsModel) ProductInfo(p products.ProductTime) *ProductInfo{
	c,ok := x.items[p]
	if !ok {
		c = &ProductInfo{
			Reles: make(map[*data.Rele]bool),
		}
		x.items[p] = c
	}
	return c
}


func (x *TableProductsModel) Data(row, col int) (text string, image walk.Image, textColor *walk.Color, backgroundColor *walk.Color ) {

	x.db.View(func(tx products.Tx) {
		if row >= len(tx.Party().Products()){
			return
		}
		p := tx.Party().Products()[row]
		pi := x.ProductInfo(p.ProductTime)
		if col == 0{
			text = fmt.Sprintf("%02d: %d", p.Addr(), p.Serial())
			//println("Data", p.ProductTime.Time().Unix(), ":", p.Serial())
			return
		}

		if col==1{
			if pi.Connection == nil {
				return
			}
			if row == x.surveyRow {
				//c.BackgroundColor = walk.RGB(128, 255, 128)
				image = ImgForwardPng16
				backgroundColor = new(walk.Color)
				*backgroundColor = walk.RGB(204, 255, 255)
				return
			}
			switch pi.Connection.Level {
			case win.NIIF_ERROR:
				text = pi.Connection.Text
				textColor = new(walk.Color)
				*textColor = walk.RGB(255, 0, 0)
				image = ImgErrorPng16
			case win.NIIF_INFO:
				textColor = new(walk.Color)
				*textColor = walk.RGB(0, 32, 128)
				image = ImgCheckmarkPng16
			}
			return
		}
		test := x.test()
		if test == nil {
			return
		}
		productType := data.ProductTypes[tx.Party().ProductTypeIndex()]

		if test == data.ManualSurvey {
			rele := productType.Reles[col-2]
			if v,ok := pi.Reles[rele]; ok {
				if v {
					image = ImgReleOnPng16
				} else {
					image = ImgReleOffPng16
				}
			}
			return
		}

		reles := test.RelesList(productType)
		if test.Index() == -1 || len (reles) == 0 {
			return
		}
		if col - 2 == len(reles){
			for _,rele := range reles {
				v := data.TestReleResult(p, test, rele)
				if v == nil {
					image = ImgQuestionPng16
					return
				}
				if releValue := test.ReleValue(rele); releValue != nil {
					if *v != *releValue {
						image = ImgErrorPng16
						return
					}
				}
			}
			image = ImgCheckmarkPng16
			return
		}

		if col - 2 < 0 && col - 2 >= len(reles) {
			return
		}

		if col == -1 {
			return
		}

		rele := reles[col-2]

		productReleValue := data.TestReleResult(p, test, rele)
		if productReleValue == nil {
			return
		}
		okReleValue := test.ReleValue(rele)
		if okReleValue == nil {
			return
		}

		if *productReleValue == *okReleValue {
			if *productReleValue {
				image = ImgReleOnPng16
			} else {
				image = ImgReleOffPng16
			}
		} else {
			if *productReleValue {
				image = ImgReleOnError
			} else {
				image = ImgReleOffError
			}
		}
	})
	return
}

func (x *TableProductsModel) SetColumns(l *walk.TableViewColumnList) {

	x.db.View(func(tx products.Tx) {
		for l.Len()>2{
			check(l.RemoveAt(l.Len()-1))
		}
		cols := []string{}
		ws := []int{}
		productsType := data.ProductTypes[tx.Party().ProductTypeIndex()]
		test := x.test()
		if test == nil {
			return
		}

		if test == data.ManualSurvey {
			for _, r := range productsType.Reles {
				cols = append(cols, r.What)
				ws = append(ws, 110)
			}
		} else {
			reles := test.RelesList(productsType)
			for _,r := range reles {
				cols = append(cols, r.What)
				ws = append(ws, 110)
			}
			if test.Index() > -1 {
				cols = append(cols, x.test().What)
				ws = append(ws, 350)
			}
		}
		for i, s := range cols {
			c := walk.NewTableViewColumn()
			check(c.SetTitle(s))
			check(c.SetWidth(ws[i]))
			check(l.Add(c))
		}
	})
}

func (x *TableProductsModel) StyleCell(c *walk.CellStyle) {
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

func (x *TableProductsModel) Value (row, col int) interface{}{
	str,_,_,_ := x.Data(row,col)
	return str
}

func (x *TableProductsModel) RowCount() (result int) {
	x.db.View(func(tx products.Tx) {
		result =  len(tx.Party().Products())
	})
	return


}

/*
var imgError = newImage("error.png")
var imgQuestion = newImage("question.png")
var imgSurvey = newImage("forward.png")
var imgConnect = newImage("checkmark.png")
var imgReleOn = newImage("rele-on.ico")
var imgReleOff = newImage("rele-off.ico")

var imgReleOnError = newImage("on-error.ico")
var imgReleOffError = newImage("off-error.ico")
*/


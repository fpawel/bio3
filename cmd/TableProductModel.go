package main

import (
	"github.com/fpawel/bio3/internal/device/data"
	"github.com/fpawel/bio3/internal/products"
	"github.com/lxn/walk"
)

type TableProductModel struct {
	walk.ReflectTableModelBase
	product products.ProductInfo
	db products.DB
}


func (x *TableProductModel) Data(row, col int) (text string ,
	image walk.Image, textColor *walk.Color, backgroundColor *walk.Color ) {

	if col < 0 {
		return
	}

	test := data.MainTestsRele[row]

	if col == 0 {
		text = test.What
		return
	}

	productType := data.ProductTypes[x.product.Party.ProductTypeIndex]
	reles := productType.Reles
	rele := reles[col-1]

	var productReleValue *bool

	x.db.View(func(tx products.Tx) {
		p,ok := tx.GetProductByProductTime(x.product.Party.PartyTime, x.product.ProductTime)
		if ok {
			productReleValue = data.TestReleResult(p, test, rele)
		}
	})



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
	return
}

func SetTableProductColumns(productsType *data.ProductType, l *walk.TableViewColumnList) {
	check(l.Clear())

	c := walk.NewTableViewColumn()
	check(c.SetTitle("Проверка"))
	check(c.SetWidth(300))
	check(l.Add(c))
	for _, r := range productsType.Reles {
		c := walk.NewTableViewColumn()
		check(c.SetTitle(r.What))
		check(c.SetWidth(110))
		check(l.Add(c))
	}
	return
}

func (x *TableProductModel) StyleCell(c *walk.CellStyle) {
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

func (x *TableProductModel) Value (row, col int) interface{}{
	str,_,_,_ := x.Data(row,col)
	return str
}

func (x *TableProductModel) RowCount() int {
	return len(data.MainTestsRele)
}
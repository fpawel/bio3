package main

import (
	"github.com/fpawel/bio3/internal/products"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"strconv"
)

func (x *AppMainWindow) runEditPartyDialog() {

	var dlg *walk.Dialog
	var acceptPB *walk.PushButton
	itemsCount := len(x.app.db.Products())

	font := Font{
		Family:    "Arial",
		PointSize: 10,
	}

	children := []Widget{
		Composite{
			Layout: Grid{},
			Children: []Widget{
				Label{
					Text: "Адрес",
					Font: font,
				},
			},
		},
		Composite{
			Layout: Grid{},
			Children: []Widget{
				Label{
					Text: "Серийный №",
					Font: font,
				},
			},
		},
	}

	leSerial := make([]*walk.LineEdit, itemsCount)

	for i, p := range x.app.db.Products() {
		i := i
		p := p
		children = append(children,
			Label{
				Text: strconv.Itoa(int(p.Addr)),
				Font: font,
			},

			LineEdit{
				Text:     strconv.Itoa(int(p.Serial)),
				AssignTo: &leSerial[i],
				Font:     font,
				OnTextChanged: func() {
					v, err := strconv.Atoi(leSerial[i].Text())

					if err == nil && v >= 0 {

						value := uint64(v)

						x.app.db.Update(func(tx products.Tx) {
							product,ok := tx.GetProductByProductTime(p.Party.PartyTime, p.ProductTime)
							if ok && product.Serial() != value {
								product.SetSerial(value)
								x.app.tableProductsModel.PublishRowChanged(p.Row)
							}
						})
					}
				},
			},
		)

	}

	_, err := Dialog{
		Icon:          mustImg("rc/settings.ico"),
		AssignTo:      &dlg,
		Title:         "Серийные номера приборов",
		DefaultButton: &acceptPB,
		MinSize:       Size{300, 300},
		Layout:        HBox{},

		Children: []Widget{
			ScrollView{
				HorizontalFixed: true,
				Layout: Grid{
					Columns:     2,
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: children,
			},
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{},
				Children: []Widget{
					PushButton{
						Text: "Закрыть",
						OnClicked: func() {
							dlg.Accept()
						},
						AssignTo: &acceptPB,
					},
				},
			},
		},
	}.Run(x)
	if err != nil {
		log.Panic(err)
	}


}

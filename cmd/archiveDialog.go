package main

import (
	"github.com/fpawel/bio3/productsview"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/fpawel/bio3/device/data"
	"github.com/fpawel/bio3/products"
	"github.com/fpawel/bio3/walkutils"
	"github.com/lxn/win"
	"log"
)


type ArchiveDialog struct {
	*walk.MainWindow
	tblLogs, tblProduct  *walk.TableView
	treeItem walk.TreeItem
	treeView *walk.TreeView
	btnDelete *walk.PushButton
	app *App
}

func (x *App) ExecuteArchiveDialog() {
	dlg := &ArchiveDialog{ app:x }
	x.mw.Hide()
	check(dlg.Markup().Create())
	dlg.Invalidate()
	win.ShowWindow(dlg.Handle(), win.SW_MAXIMIZE)
	if windowResult := dlg.Run(); windowResult != 0 {
		log.Println("PartiesDialog  windowResult:", windowResult)
	}
	x.mw.Show()
}

func Bio3Info() productsview.DeviceInfoProvider{
	return productsview.DeviceInfoProvider{
		data.GoodProduct,
		data.BadProduct,
		data.FormatProductType,
	}
}

func (x *ArchiveDialog) Invalidate() {
	m := productsview.NewPartiesTreeViewModel(x.app.db, Bio3Info())
	check(x.treeView.SetModel(m))

	var currentPartytime products.PartyTime
	x.app.db.View(func(tx products.Tx) {
		currentPartytime = tx.Party().PartyTime
	})
	itemContainsCurrentParty := func (it walk.TreeItem) (result bool) {
		switch it := it.(type) {
		case productsview.Node:
			result = it.ContainsParty(currentPartytime)
		}
		return
	}
	walkutils.SetExpandedTreeviewItems(x.treeView, true, itemContainsCurrentParty)
}

func (x *ArchiveDialog) deleteSelectedNode() {
	if x.treeItem == nil {
		return
	}
	x.app.db.Update(func(tx products.Tx) {
		productsview.DeleteNode( x.treeItem.(productsview.Node), tx )
	})
	x.Invalidate()
}

func (x *ArchiveDialog) Markup() MainWindow {

	tblLogsModel := productsview.NewTableLogsViewModel(x.app.db, func(tx products.Tx) [][]byte {
		switch node := x.treeItem.(type) {
		case *productsview.NodeParty:
			return node.Party().PartyTime.Path()
		case *productsview.NodeProduct:
			return node.Product().ProductTime.Path(node.Party().PartyTime)

		default:
			return nil
		}

	})

	const fontPointSize = 10

	return MainWindow{
		Icon:     NewIconFromResourceId(IconDBID),
		AssignTo: &x.MainWindow,
		Title: "Обзор партий",
		Layout:   HBox{},

		Children: []Widget{
			TreeView{
				Model:    productsview.NewPartiesTreeViewModel(x.app.db, Bio3Info()),
				AssignTo: &x.treeView,
				MaxSize:  Size{200, 0},

				Font: Font{
					Family:    "Arial",
					PointSize: 12,
				},
				OnCurrentItemChanged: x.treeItemChanged,
			},

			TableView{
				AssignTo:              &x.tblLogs,
				AlternatingRowBGColor: walk.RGB(239, 239, 239),
				CheckBoxes:            false,
				ColumnsOrderable:      false,
				MultiSelection:        true,
				Font:                  Font{PointSize: fontPointSize},
				Model:                 tblLogsModel,
				StyleCell:             tblLogsModel.StyleCell,
				Columns: []TableViewColumn{
					{Title: "Время", Width: 100},
					{Title: "Проверка", Width: 400},
					{Title: "Сообщение", Width: 500},
				},
			},

			TableView{
				AssignTo:              &x.tblProduct,
				AlternatingRowBGColor: walk.RGB(239, 239, 239),
				CheckBoxes:            false,
				ColumnsOrderable:      false,
				MultiSelection:        true,
				Font:                  Font{PointSize: fontPointSize},
			},

			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{},
				Children: []Widget{
					PushButton{
						AssignTo:  &x.btnDelete,
						Text:      "Удалить",
						OnClicked: x.deleteSelectedNode,
					},

				},
			},
		},
	}
}

func (x *ArchiveDialog) treeItemChanged() {
	x.treeItem = x.treeView.CurrentItem()

	switch node := x.treeItem.(type) {
	case *productsview.NodeParty:
		x.tblLogs.Model().(* productsview.TableLogsViewModel).PublishRowsReset()
		x.tblLogs.SetVisible(true)
		x.tblProduct.SetVisible(false)
		x.btnDelete.SetVisible(true)
	case *productsview.NodeProduct:
		t := data.ProductTypes[node.Product().Party.ProductTypeIndex]
		SetTableProductColumns(t, x.tblProduct.Columns() )

		//m := &TableProductModel{product: node.Product(), db: x.app.db}
		m := &TableProductModel{product: node.Product(), db: x.app.db}
		check(x.tblProduct.SetModel(m))
		x.tblProduct.SetCellStyler(m)
		x.tblLogs.SetVisible(false)
		x.tblProduct.SetVisible(true)
		x.btnDelete.SetVisible(false)
	default:
		x.tblLogs.SetVisible(false)
		x.btnDelete.SetVisible(true)
		check(x.tblProduct.Columns().Clear())
		x.tblProduct.SetVisible(true)
	}
}

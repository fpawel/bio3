package main

import (
	. "github.com/lxn/walk/declarative"

	"fmt"
	"github.com/lxn/walk"

	"github.com/fpawel/bio3/device/data"
	"github.com/fpawel/bio3/device/stend"

	"github.com/fpawel/bio3/products"
	"log"
)

const FontPointSize = 10


func NewMainwindow(x *AppMainWindow) MainWindow {
	var mMewPartyItems [20]MenuItem
	for i := 0; i < 20; i++ {
		n := i + 1
		mMewPartyItems[i] =
			Action{

				Text: fmt.Sprintf("%d", n),
				OnTriggered: func() {
					x.app.db.Update(func(tx products.Tx) {
						for j, p := range tx.NewParty(n).Products() {
							p.SetAddr(byte(stend.Places[n-1][j]))
						}
					})

					x.app.tableProductsModel.Clear()
					x.setupParty()
					x.app.tableLogsModel.PublishRowsReset()
					x.runEditPartyDialog()
					x.app.tableProductsModel.PublishRowsReset()
				},
			}
	}

	singleWork := func (visible bool, work *data.Work ) Action{
		return Action{
			Text: work.What,
			Visible:visible,
			OnTriggered: func() {
				x.app.RunWorks(work.What, []*data.Work{work})
			},
		}
	}

	singleWorkAction := func (s string, visible bool, a func() error ) Action{
		return singleWork(visible, &data.Work{
			What:  s,
			Reles: data.RelesValues{},
			Work:  a,
		})
	}



	tableLog := TableView{
		AssignTo:&x.tblLogs,
		AlternatingRowBGColor: walk.RGB(239, 239, 239),
		CheckBoxes:            false,
		ColumnsOrderable:      false,
		MultiSelection:        true,
		Font: Font{ PointSize: FontPointSize },
		Model:     x.app.tableLogsModel,
		MaxSize:Size{0,180},
		MinSize:Size{0,180},

		Columns:[]TableViewColumn{
			{Title : "Время", Width:100 },
			{Title : "Проверка" , Width:400},
			{Title : "Сообщение", Width:500},
		},
	}

	tableParty := TableView{
		AssignTo:              &x.tblParty,
		AlternatingRowBGColor: walk.RGB(239, 239, 239),
		CheckBoxes:            false,
		ColumnsOrderable:      false,
		MultiSelection:        true,
		Font: Font{ PointSize: FontPointSize },
		Model:     x.app.tableProductsModel,
		//StyleCell: tblProductsModel.StyleCell,

		Columns:[]TableViewColumn{
			{Title : "Номер"},
			{Title : "Связь"},
		},
	}

	tableWorks := TableView{
		Model:x.app.tableWorksModel,
		//StyleCell:x.tableWorksModel.StyleCell,
		OnCurrentIndexChanged: x.TableWorksCurrentIndexChanged,
		MaxSize:          Size{0, 170},
		MinSize:          Size{0, 170},
		AssignTo:         &x.tblWorks,
		CheckBoxes:       true,
		ColumnsOrderable: false,
		MultiSelection:   false,
		Font: Font{ PointSize: FontPointSize },
		Columns: []TableViewColumn{
			{Title: "Проверка", Width: 400},
			{Title: "Статус", Width: 50},
		},
	}


	return MainWindow{
		Title:    x.app.Title(),
		Icon:  NewIconFromResourceId(IconAppID),
		Name:     "AppMainWindow",
		Size:     Size{800, 600},
		Layout:   HBox{MarginsZero: true, SpacingZero: true},
		AssignTo: &x.MainWindow,
		Font: Font{ PointSize: FontPointSize },

		MenuItems: []MenuItem{

			Menu{
				Text: "Помощь",

				Items: []MenuItem{
					Action{
						Text: "О программе",
						OnTriggered: func() {

							walk.MsgBox(x, "О программе",
								fmt.Sprintf("Версия ПО: %d.%d.%d.%s", majorVersion, minorVersion, bugfixVersion, buildtime),
								walk.MsgBoxIconInformation)
						},
					},
				},
			},
		},

		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					tableWorks,
					tableParty,
					tableLog,
					x.delayControl.Markup(),
				},

			},

			ScrollView{
				Name:            "RightPan",
				Layout:          VBox{},
				HorizontalFixed: true,
				Children: []Widget{

					SplitButton{
						ImageAboveText: true,
						Image:    AssetImage("assets/png16/right-arrow.png"),
						Text:     "Выполнить",
						AssignTo: &x.btnRunMenu,
						MenuItems: []MenuItem{
							Action{
								Text: "Проверки",
								OnTriggered: func() {
									x.app.tableWorksModel.SetWorks(data.MainWorks)
								},
							},
							Separator{},
							singleWork(true, data.ManualSurvey),
							singleWorkAction ("Показать места", true,  x.app.LedPlacesOn),
							singleWorkAction ("Включить питание", true, x.app.PowerOn),
							singleWorkAction ("Выключить питание", true, x.app.PowerOff),

							Action{
								Text:    "Ввод",
								Visible: prod=="",
								AssignTo:&x.userInput,
								OnTriggered: func() {
									x.app.RunWorks("Ввод", []*data.Work{
										{
											What:  "Ввод",
											Reles: data.RelesValues{},
											Work:  x.app.consoleInput,
										},

									})
								},
							},
						},
						OnClicked: x.app.RunMainWorks,
					},
					PushButton{
						ImageAboveText: true,
						Image:          AssetImage("assets/png16/cancel.png"),
						Text:           "Прервать",
						OnClicked: func() {
							if  x.app.workHelper == nil {
								log.Fatal("cancel work click: x.app.workHelper == nil")
							}
							x.delayControl.Cancel()
							x.app.workHelper.port.Cancel()
						},
						Visible:  false,
						AssignTo: &x.btnCancel,
					},

					leftAlignedTitleLabel("СОМ порт"),
					ComboBox{
						Name:        "ComboBoxSerialPort",
						AssignTo:    &x.cbPort,
						OnMouseDown: x.comboBoxPortMouseDown,
						OnCurrentIndexChanged: func() {
							x.app.config.Port = x.cbPort.Text()
							x.app.SaveConfig()
						},
					},


					leftAlignedTitleLabel("Исполнение"),

					ComboBox{
						AssignTo:      &x.cbProductsType,
						Model:         data.ProductTypes,
						DisplayMember: "What",
						OnCurrentIndexChanged: func() {
							x.app.db.Update(func(tx products.Tx) {
								if tx.Party().ProductTypeIndex() != x.cbProductsType.CurrentIndex() {
									tx.Party().SetProductTypeIndex(x.cbProductsType.CurrentIndex())
									x.setupProductsColumns()
								}
							})
						},
					},

					PushButton{
						AssignTo:&x.btnTests,
						ImageAboveText: true,
						Text:           "Проверки",
						Image:          AssetImage("assets/png16/if_script.png"),
						OnClicked: func() {
							x.app.tableWorksModel.SetWorks(data.MainWorks)
							x.setupProductsColumns()
							x.btnTests.SetEnabled(false)
						},
						Enabled:false,
					},

					PushButton{
						ImageAboveText: true,
						Text:           "Настройки",
						Image:          AssetImage("assets/png16/preference.png"),
						OnClicked: x.RunSetsDialog,
					},

					SplitButton{
						ImageAboveText: true,
						Text:"Новая партия",
						Image:AssetImage("assets/png16/add.png"),
						MenuItems: mMewPartyItems[:],
					},

					PushButton{
						ImageAboveText: true,
						Text:           "Архив",
						Image:          AssetImage("assets/png16/db.png"),
						OnClicked: x.app.ExecuteArchiveDialog,
					},
				},
			},

		},
	}
}


func leftAlignedTitleLabel(text string) ScrollView {
	return ScrollView{
		Layout:        HBox{MarginsZero: true, SpacingZero: true},
		VerticalFixed: true,
		Children: []Widget{
			Label{
				Text: text,
				Font: Font{Bold: true, PointSize: FontPointSize},
			},
		},
	}
}



package main

import (
	"github.com/fpawel/bio3/internal/fetch"
	"github.com/lxn/walk"
	"time"
)

func (x *AppMainWindow) RunSetsDialog() {
	var acceptPB, cancelPB *walk.PushButton
	var edReadTimeout,
		edReadByteTimeout,
		edAttempts *walk.NumberEdit
	var dlg *walk.Dialog
	cfg := &x.app.config

	check(Dialog{
		Icon: NewIconFromResourceId(IconSettingsID),
		AssignTo:      &dlg,
		Title:         "Настройки",
		DefaultButton: &acceptPB,
		MinSize:       Size{300, 300},
		Layout: HBox{},
		Font:Font{PointSize:10},

		Children: []Widget{
			ScrollView{
				HorizontalFixed:true,
				Layout: Grid{
					Columns:     2,
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: []Widget{
					Label{Text: "Таймаут ответа, мс"},
					NumberEdit{
						AssignTo: &edReadTimeout,
						MinValue: 10,
						MaxValue: 10000,
						MinSize:Size{150,0},
						MaxSize:Size{150,0},
					},
					Label{Text: "Таймаут байта, мс"},
					NumberEdit{
						AssignTo: &edReadByteTimeout,
						MinValue: 1,
						MaxValue: 100,
					},
					Label{Text: "Макс.повторов"},
					NumberEdit{
						Name:     "NumberEditAttempts",
						AssignTo: &edAttempts,
						MinValue: 1,
						MaxValue: 5,
					},
				},
			},
			ScrollView{
				HorizontalFixed:true,
				Layout:VBox{},
				Children:[]Widget{
					PushButton{
						Text:"Применить",
						AssignTo:&acceptPB,
						OnClicked: func() {
							cfg.ReadWrite = fetch.Config{
								ReadTimeout:     time.Millisecond * time.Duration(edReadTimeout.Value()),
								ReadByteTimeout: time.Millisecond * time.Duration(edReadByteTimeout.Value()),
								MaxAttemptsRead: int(edAttempts.Value()),
							}
							x.app.SaveConfig()
							dlg.Accept()
						},
					},
					PushButton{
						Text:"Отмена",
						AssignTo:&cancelPB,
						OnClicked:func(){
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Create(x))



	check(edReadByteTimeout.SetValue(cfg.ReadWrite.ReadByteTimeout.Seconds() * 1000))
	check(edReadTimeout.SetValue(cfg.ReadWrite.ReadTimeout.Seconds() * 1000))
	check(edAttempts.SetValue(float64(cfg.ReadWrite.MaxAttemptsRead)))
	dlg.Run()

}

package delay

func (x *Control) Markup() Composite {
	return Composite{
		AssignTo: &x.Composite,
		Layout:VBox{ SpacingZero:true, MarginsZero:true},
		Visible:false,
		Children:[]Widget{
			ScrollView{
				VerticalFixed:true,
				Layout:HBox{SpacingZero:true, MarginsZero:true}	,
				Children:[]Widget{
					Label{ AssignTo: &x.label},
				},
			},
			Composite{
				Layout:   HBox{MarginsZero:true},
				Children: []Widget{
					ProgressBar{
						AssignTo: &x.progressBar,
					},
					PushButton{
						Image: x.imageSkip,
						Text : "Продолжить без задержки",
						OnClicked: x.Cancel,
					},
				},
			},
		},
	}
}

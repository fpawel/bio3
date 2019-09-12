package walkutils

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type ProgressDialog struct {
	*walk.Dialog
	Label       *walk.Label
	ProgressBar *walk.ProgressBar
	closed      bool
}

type ProgressDialogConfig struct {
	Title string
	ProgrssMin, ProgrssMax, ProgressValue int
	Icon walk.Property
}

func (x *ProgressDialog) Closed() bool {
	return x.closed
}

func (x *ProgressDialog) Run() int {
	r := x.Dialog.Run()
	x.closed = true
	return r
}

func (x *ProgressDialog) Cancel() {
	x.closed = true
	x.Dialog.Cancel()
}

func (x *ProgressDialog) SetText( m  Message ) {
	if err :=  x.Label.SetText(m.Text); err != nil {
		return
	}
	x.Label.SetTextColor(m.Color())
}

func (x *ProgressDialog) Markup(c ProgressDialogConfig) Dialog {
	return Dialog{
		Title:c.Title,
		Icon:c.Icon,
		FixedSize:true,
		AssignTo:&x.Dialog,
		Layout:VBox{},
		MinSize:Size{400, 100},
		Children:[]Widget{
			ScrollView{
				VerticalFixed:true,
				Layout:HBox{SpacingZero:true, MarginsZero:true}	,
				Children:[]Widget{
					Label{ AssignTo: &x.Label, Font:Font{PointSize:12}},
				},
			},
			Composite{
				Layout:   HBox{MarginsZero:true},
				Children: []Widget{
					ProgressBar{
						AssignTo: &x.ProgressBar,
						Value:c.ProgressValue,
						MinValue:c.ProgrssMin,
						MaxValue:c.ProgrssMax,
					},
				},
			},
		},
	}
}

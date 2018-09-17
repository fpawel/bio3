package walkutils

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

func ErrorLevel(err error, okMessage string) (string, int) {
	if err == nil {
		return okMessage, win.NIIF_INFO
	}
	return err.Error(), win.NIIF_ERROR
}

func LeftAlignedTitleLabel(text string) ScrollView {
	return ScrollView{
		Layout:        HBox{MarginsZero: true, SpacingZero: true},
		VerticalFixed: true,
		Children: []Widget{
			Label{
				Text: text,
				Font: Font{Bold: true},
			},
		},
	}
}

func ScrollDownTableView( tableView *walk.TableView ){
	tableView.Synchronize(func() {
		tableView.SendMessage(win.WM_VSCROLL, win.SB_PAGEDOWN, 0)
	})
}

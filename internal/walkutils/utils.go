package walkutils

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"log"
	"path/filepath"
	"runtime"
)

func ShowNotifyMessage(n *walk.NotifyIcon, level int ) func(title, info string) error {
	switch level {
	case win.NIIF_INFO:
		return n.ShowInfo
	case win.NIIF_WARNING:
		return n.ShowWarning
	case win.NIIF_ERROR:
		return n.ShowError
	case win.NIIF_USER:
		return n.ShowInfo
	default:
		return n.ShowMessage
	}
}



func SetExpandedTreeviewItems(tv *walk.TreeView, expanded bool, pred func(item walk.TreeItem) bool) {
	for i := 0; i < tv.Model().RootCount(); i++ {
		it := tv.Model().RootAt(i)
		if pred(it) {
			setExpandedTreeviewItem(it, tv, expanded, pred)
		}

	}
}

func setExpandedTreeviewItem(it walk.TreeItem, tv *walk.TreeView, expanded bool, pred func(item walk.TreeItem) bool) {
	if it.ChildCount() == 0 {
		return
	}
	if err :=tv.SetExpanded(it, expanded); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < it.ChildCount(); i++ {
		itt := it.ChildAt(i)
		if pred(itt) {
			setExpandedTreeviewItem(itt, tv, expanded, pred)
		}

	}
}

func check(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Panicf("%s:%d %v\n", filepath.Base(file), line, err)
	}
}

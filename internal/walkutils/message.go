package walkutils

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"time"
)

type Message struct{
	Level int
	Text string
}

type MessageTime struct{
	Message
	time.Time
}



func (x Message) String() string{
	return x.Text
}


func (x Message) Color() walk.Color{
	switch x.Level{
	case win.NIIF_ERROR:
		return walk.RGB(255,0,0)
	case win.NIIF_INFO:
		return walk.RGB(0,0,153)
	default:
		return walk.RGB(0,0,0)
	}
}

func MessageFromError(err error, okText string) (x Message){
	if err == nil {
		x.Level = win.NIIF_INFO
		x.Text = okText
	} else {
		x.Text = err.Error()
		x.Level = win.NIIF_ERROR
	}
	return
}


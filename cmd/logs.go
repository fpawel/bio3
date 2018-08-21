package main

import (
	"github.com/fpawel/bio3/products"
	"os"
	"github.com/lxn/win"
	"github.com/daviddengcn/go-colortext"
	"fmt"
	"github.com/fpawel/bio3/device/data"
)

type SerialPortLogger struct {
	app *App
}

func (x SerialPortLogger) Write(p []byte) (n int, err error){
	if x.app.tableWorksModel.CurrentWork() != data.ManualSurvey {
		return os.Stdout.Write(p)
	}
	return len(p), nil
}

func (x *App) WriteWorkLog(tx products.Tx, timeKey []byte, level int, text string)  {
	if work := x.tableWorksModel.CurrentWork(); work!=nil {
		work.WriteLog(tx, timeKey, level, text)
		PrintLog(level,fmt.Sprintf("%s: %s", work.What, text))
		x.mw.Synchronize(func() {
			x.tableWorksModel.PublishRowsReset()
			x.tableLogsModel.PublishRowsReset()
			x.mw.tblLogs.SendMessage(win.WM_VSCROLL, win.SB_PAGEDOWN, 0)
		})
	} else {
		PrintLog(level,text)
	}
}

func PrintLog(level int, text string)   {
	switch level {
	case win.NIIF_ERROR:
		ct.Foreground(ct.Red,true)
		fmt.Println( "ERROR:", text)
		ct.ResetColor()
	case win.NIIF_INFO:
		ct.Foreground(ct.Blue,true)
		fmt.Println( text)
		ct.ResetColor()
	default:
		fmt.Println( text)
	}
	return
}




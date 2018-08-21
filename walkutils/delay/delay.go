package delay

import (
	"github.com/lxn/walk"
	"fmt"
	"time"
	"sync/atomic"
	"runtime"
	"log"
	"path/filepath"
)

type Control struct {
	*walk.Composite
	progressBar *walk.ProgressBar
	label *walk.Label
	imageSkip walk.Image
	canceled int32
}


func New(imageSkip walk.Image) *Control{
	return &Control{ imageSkip : imageSkip }
}

func (x *Control) Cancel() {
	atomic.AddInt32(&x.canceled, 1)
}

func (x *Control) Run(d time.Duration, what string, start, update func()) error {

	satrtTime := time.Now()
	x.canceled = 0

	x.Synchronize(func() {
		x.SetVisible(true)
		x.progressBar.SetRange(0, int(d.Nanoseconds()/1000000))
		x.progressBar.SetValue(0)
		check(x.label.SetText(fmt.Sprintf("задержка %s, %v", what, d)))
		start()
	})

	for  time.Since(satrtTime) < d && atomic.LoadInt32(&x.canceled) == 0 {
		time.Sleep(time.Millisecond * 500)
		x.Synchronize(func() {
			x.progressBar.SetValue(int(time.Since(satrtTime).Nanoseconds() / 1000000))
			str := fmt.Sprintf("задержка, %s, %v - %v", what, time.Since(satrtTime), d)
			check(x.label.SetText(str))
		})
		update()
	}

	x.Synchronize(func() {
		x.SetVisible(false)
	})

	return nil
}

func check(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Panicf("%s:%d %v\n", filepath.Base(file), line, err)
	}
}



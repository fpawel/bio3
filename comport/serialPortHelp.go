package comport

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"github.com/fpawel/bio3/fetch"
	"sync/atomic"

)

type Port struct {
	*serial.Port
	config      *serial.Config
	configFetch fetch.Config
	canceled    *int32
	logging     bool
}

func (x *Port) Canceled() bool {
	return atomic.LoadInt32(x.canceled) != 0
}

func (x *Port) Cancel() {
	atomic.AddInt32(x.canceled, 1)
}


func (x *Port) Write(buf []byte) (int, error){
	if x.Canceled() {
		return 0, fetch.ErrorCanceled
	}
	if err := x.Port.Flush(); err != nil {
		return 0, err
	}
	return x.Port.Write(buf)
}

func (x *Port) Fetch(request []byte) ([]byte, error){

	str1 := fmt.Sprintf("%s:% X",  x.config.Name, request)

	if x.Canceled() {
		if x.logging{
			fmt.Println(str1, "canceled")
		}
		return nil, fetch.ErrorCanceled
	}

	b,t,at,err := fetch.Fetch(request, x.configFetch, x, x.canceled)

	if x.logging{
		if err == nil {
			fmt.Printf( "%s --> % X %v\n", str1, b, t)
		} else {
			fmt.Printf( "%s --> % X %v at %d: %v\n", str1, b, t, at, err )
		}
	}

	return b,err
}

func Open(config *serial.Config, configFetch fetch.Config, logging bool,
	cancelation *int32) (*Port, error) {

	tarmSerialPort,err := serial.OpenPort(config)
	if err != nil {
		err = fmt.Errorf("open port: %s", err.Error())
		return nil,err
	}
	x := &Port{
		Port:tarmSerialPort,
		config : config,
		canceled : cancelation ,
		configFetch : configFetch,
		logging : logging,
	}

	if x.canceled == nil {
		x.canceled = new(int32)
	}

	return x, nil
}


func AvailablePorts () (serials []string) {
	for i := 1; i < 100; i++ {
		s := fmt.Sprintf("COM%d", i)
		port, err := serial.OpenPort(&serial.Config{Name: s, Baud: 9600})
		if err == nil {
			err = port.Close()
			if err != nil {
				log.Panic(err)
			}
			serials = append(serials, s)
		}
	}
	return
}
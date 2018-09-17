package fetch

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)



type Config struct {
	ReadTimeout     time.Duration `json:"read_timeout"`
	ReadByteTimeout time.Duration `json:"read_byte_timeout"`
	MaxAttemptsRead int           `json:"max_attempts_read"`
}

var ErrorNoAnswer = fmt.Errorf("не отвечает")
var ErrorCanceled = fmt.Errorf("прервано")


func isCanceled(v *int32) bool{
	return atomic.LoadInt32(v) != 0
}

func Fetch(request []byte, config Config, rw io.ReadWriter, canceled *int32) (response []byte, delay time.Duration, attempt int, err error) {
	if canceled == nil {
		canceled = new(int32)
	}
	if config.MaxAttemptsRead > 10 {
		config.MaxAttemptsRead = 10
	}
	if config.MaxAttemptsRead < 1 {
		config.MaxAttemptsRead = 1
	}
	for attempt = 1; ; attempt++ {
		t := time.Now()
		response, err = tryFetch(request, config, rw, canceled)
		delay = time.Since(t)
		if err == nil || err == ErrorCanceled || attempt == config.MaxAttemptsRead {
			return
		}
	}
	return
}

func write(w io.Writer, b []byte) (err error) {

	var writenCount int
	writenCount, err = w.Write(b)
	if err != nil {
		return
	}
	if writenCount != len(b) {
		err = fmt.Errorf("%d bytes writen of %d", writenCount, len(b))
	}
	return
}

func tryFetch(request []byte, config Config, rw io.ReadWriter, canceled *int32) (response []byte, err error) {

	if isCanceled(canceled){
		return nil, ErrorCanceled
	}
	// записать и проверить ошибку записи
	if err = write(rw, request); err != nil {
		err = fmt.Errorf("write: %s", err.Error())
		return
	}

	// от этой временной метки отсчитывается таймаут считывания
	t := time.Now()

	for {
		if isCanceled(canceled){
			return nil, ErrorCanceled
		}

		// проверка таймаута считывания
		if time.Since(t) > config.ReadTimeout {
			err = ErrorNoAnswer
			return
		}

		// пытаться считать первый байт ответа
		var readedCount int
		tmp := []byte{0}
		if readedCount, err = rw.Read(tmp); err != nil {
			err = fmt.Errorf("read first byte: %s", err.Error())
			return
		}
		if readedCount == 0 {
			continue // не удалось считать первый байт ответа
		}

		// считан первый байт ответа

		// добавить считанный байт в результат
		response = append(response, tmp[0])

		t = time.Now() // временная метка таймаута байтов в ответе
		for {

			if isCanceled(canceled){
				return nil, ErrorCanceled
			}

			// пытаться считать один байт
			if readedCount, err = rw.Read(tmp); err != nil {
				err = fmt.Errorf("read: %s", err.Error())
				return
			}
			if readedCount == 0 {
				if time.Since(t) > config.ReadByteTimeout { // проверка таймаута байтов в ответе
					return // больше байтов ответа не будет, ответ считан полностью
				}
				continue
			}
			// добавить считанный байт в результат
			response = append(response, tmp[0])
			t = time.Now()
			// продолжить до достижения таймаута байтов в ответе
		}
	}
}



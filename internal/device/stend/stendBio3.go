package stend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/fpawel/bio3/internal/comport"
	"log"
)

const (
	CMD_CONNECT byte = 1 + iota
	CMD_POWER
	CMD_ADDR
	CMD_LED
	CMD_SWITCH
	CMD_RELE
	CMD_VIBRATION
	CMD_ADDR_PINS
	CMD_ADDR_LED
)

const (
	CMD_GAS_CONNECT byte = 161
	CMD_GAS byte = 162

)


const (
	LED_ALL byte = iota
	LED_RED
	LED_YELLOW
	LED_GREEN
)

const (
	LED_OFF byte = iota
	LED_ON
	LED_BLINK
)

const (
	SWITCH_ALL byte = 0
)

const (
	GAS_OFF byte = iota
	GAS1
	GAS2
	GAS3
)


var Places = [20][]int{
	{1},																				//1
	{2, 3},																				//2
	{1, 2, 3},																			//3	
	{4, 5, 6, 7},																		//4		
	{1, 4, 5, 6, 7},																	//5
	{2, 3, 4, 5, 6, 7},																	//6			
	{1, 2, 3, 4, 5, 6, 7},																//7
	{8, 9, 10, 11, 12, 13, 14, 15},														//8
	{1, 8, 9, 10, 11, 12, 13, 14, 15},													//9			
	{2, 3, 8, 9, 10, 11, 12, 13, 14, 15},												//10	
	{1, 2, 3, 8, 9, 10, 11, 12, 13, 14, 15},											//11
	{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},											//12
	{1, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},										//13
	{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},									//14
	{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},								//15
	{1, 2, 3, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},						//16
	{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},						//17
	{1, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},					//18
	{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},				//19
	{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},			//20
}

type PowerState byte

const (
	PowerOn PowerState = 1
	PowerOff PowerState = 0
)



func (x PowerState) String() string{
	switch x {
	case PowerOn:
		return "включение питания"
	case PowerOff:
		return "выключение питания"
	default:
		return fmt.Sprintf("питание: %d", x)
	}
}



func Led(diode, mode byte, on, off uint16) []byte {
	return []byte{
		diode, mode,
		byte(on >> 8), byte(on),
		byte(off >> 8), byte(off),
	}
}

func LedTurnOff(port *comport.Port, addr byte) (err error) {
	b := Request(addr, CMD_LED, Led(LED_ALL, LED_OFF, 0, 0))
	_, err = fetch(port, b)
	return
}

func LedGreenTurnOn(port *comport.Port, addr byte) (err error) {
	b := Request(addr, CMD_LED, Led(LED_GREEN, LED_ON, 1000, 0))
	_, err = fetch(port, b)
	return
}

func Rele(port *comport.Port, addr byte) ( r []byte, err error) {
	b := Request(addr, CMD_RELE, nil)
	b, err = port.Fetch(b)
	if err != nil {
		return
	}
	if len(b) != 10{
		err = fmt.Errorf("ожидалось 10 байт")
	}
	r = make([]byte,5)
	copy(r[:], b[4:9])
	return
}


func Power(port *comport.Port, addr byte, powerState PowerState) (err error) {
	_, err = fetch(port, Request(addr, CMD_POWER, []byte{ byte(powerState) }))
	return
}

func Switch(port *comport.Port, addr byte, n byte, state byte) (err error) {
	_, err = fetch(port, Request(addr, CMD_SWITCH, []byte{ n, state }))
	return
}

func Vibration(port *comport.Port, addr byte, state byte, amplitude float32, freq float32) (err error) {
	b := make([]byte, 9)
	b[0] = state

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, amplitude); err!= nil {
		log.Panicln(err)
	}
	copy(b[1:], buf.Bytes())

	buf = new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, freq); err!= nil {
		log.Panicln(err)
	}
	copy(b[5:], buf.Bytes())


	_, err = fetch(port, Request(addr, CMD_VIBRATION, b))
	return
}

func Gas(port *comport.Port, rate float32, point byte, count byte) (err error) {
	b := make([]byte, 6)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, rate); err!= nil {
		log.Panicln(err)
	}


	b[4] = point
	b[5] = count
	copy(b, buf.Bytes())
	_, err = fetch(port, Request(100, CMD_GAS, b))
	return
}




func Request(addr, cmd byte, b []byte) (r []byte) {
	n := 4 + len(b)
	r = make([]byte, 4+len(b))
	r[0] = byte(n)
	r[1] = addr
	r[2] = cmd
	copy(r[3:], b)
	for _, x := range r[:n-1] {
		r[n-1] += x
	}
	return
}

func fetch(port *comport.Port, b []byte)( r []byte, err error ){
	r,err = port.Fetch(b)
	if err != nil{
		return
	}
	err = checkResponse(b[1], b[2], r)
	return
}

func checkResponse(addr byte, cmd byte, b []byte) error {
	n := len(b)

	if n < 5 {
		return fmt.Errorf("длина ответа %d, должна быть более 4", n)
	}

	if b[1] != addr {
		return fmt.Errorf("второй байт не равен адресу контроллера запроса")
	}
	if b[2] != cmd {
		return fmt.Errorf("третий байт не равен коду команды запроса")
	}
	crc := byte(0)
	for _, x := range b[:n-1] {
		crc += x
	}
	if b[n-1] != crc {
		return fmt.Errorf("несовпадение контрольной суммы")
	}
	if b[3] != 0 {
		return fmt.Errorf("код ошибки %d", b[3])
	}

	return nil
}

package main

import (
	"log"
	//_ "runtime/cgo"

)

var buildtime = "(undefined)"
var majorVersion = 0
var minorVersion = 0
var bugfixVersion = 0
var debug = ""
var prod = ""


func main() {

	println("debug:", debug, "prod:", prod)

	log.SetFlags(log.Lshortfile | log.Ltime  )

	app := NewApp()
	app.mw.Run()
	app.Close()
}

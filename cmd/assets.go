package main

import (
	"github.com/lxn/walk"
	"image"
	"log"
	"bytes"

	//_ "github.com/dkua/go-ico"
	_ "image/png"
	"github.com/fpawel/bio3/productsview"
)

var ImgReleOnPng16 = AssetImage("assets/png16/rele_on.png")
var ImgReleOffPng16 = AssetImage("assets/png16/rele_off.png")

var ImgReleOnError = AssetImage("assets/png16/on_error.png")
var ImgReleOffError = AssetImage("assets/png16/off_error.png")



var ImgCheckmarkPng16 = AssetImage("assets/png16/checkmark.png")
var ImgCloudPng16 = AssetImage("assets/png16/cloud.png")
var ImgErrorPng16 = AssetImage("assets/png16/error.png")
var ImgQuestionPng16 = AssetImage("assets/png16/question.png")
var ImgForwardPng16 = AssetImage("assets/png16/forward.png")


func AssetImage(path string)  walk.Image {

	b, err := Asset(path)
	if err != nil {
		log.Fatalln(err, path)
	}

	x, s, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatalln(err, s, path)
	}
	r, err := walk.NewBitmapFromImage(x)
	if err != nil {
		log.Fatalln(err, s, path)
	}
	return r

}

func init(){

	productsview.ImgCalendarYearPng16 = AssetImage("assets/png16/calendar-year.png")
	productsview.ImgCalendarMonthPng16 = AssetImage("assets/png16/calendar-month.png")
	productsview.ImgCalendarDayPng16 = AssetImage("assets/png16/calendar-day.png")
	productsview.ImgErrorPng16 = ImgErrorPng16
	productsview.ImgPartyNodePng16 = AssetImage("assets/png16/folder2.png")
	productsview.ImgCheckmarkPng16 = ImgCheckmarkPng16
	productsview.ImgProductNodePng16 = AssetImage("assets/png16/folder1.png")
	productsview.ImgWindowIcon = NewIconFromResourceId(IconDBID)
}

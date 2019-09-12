package main

import (
	"github.com/fpawel/bio3/internal/productsview"
	"github.com/lxn/walk"
	//_ "github.com/dkua/go-ico"
	_ "image/png"
)

var ImgReleOnPng16 = mustImg("assets/png16/rele_on.png")
var ImgReleOffPng16 = mustImg("assets/png16/rele_off.png")

var ImgReleOnError = mustImg("assets/png16/on_error.png")
var ImgReleOffError = mustImg("assets/png16/off_error.png")



var ImgCheckmarkPng16 = mustImg("assets/png16/checkmark.png")
var ImgCloudPng16 = mustImg("assets/png16/cloud.png")
var ImgErrorPng16 = mustImg("assets/png16/error.png")
var ImgQuestionPng16 = mustImg("assets/png16/question.png")
var ImgForwardPng16 = mustImg("assets/png16/forward.png")


func mustImg(path string)  walk.Image {
	img,err := walk.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return img
}

func init(){

	productsview.ImgCalendarYearPng16 = mustImg("assets/png16/calendar-year.png")
	productsview.ImgCalendarMonthPng16 = mustImg("assets/png16/calendar-month.png")
	productsview.ImgCalendarDayPng16 = mustImg("assets/png16/calendar-day.png")
	productsview.ImgErrorPng16 = ImgErrorPng16
	productsview.ImgPartyNodePng16 = mustImg("assets/png16/folder2.png")
	productsview.ImgCheckmarkPng16 = ImgCheckmarkPng16
	productsview.ImgProductNodePng16 = mustImg("assets/png16/folder1.png")
	productsview.ImgWindowIcon = mustImg("assets/rc/app.ico")
}

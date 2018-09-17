package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/fpawel/bio3/internal/device/stend"
	"strings"
)

//consoleInputCmd определяет команду консольного ввода с параметрами
type consoleInputCommand struct {
	s, d string // s -имя команды в консоли, d - описание

	n int                // количество аргументов
	a func([]byte) error // функция, которая будут выполнена
}

func (x *App) ConsoleInputCommands() []consoleInputCommand {
	xs := []consoleInputCommand{
		{
			"power",
			"[0 - вкл. | 1 - выкл.] - управление питанием",
			1,
			func(b []byte) error {
				return x.Power(stend.PowerState(b[0]))
			},
		},
		{
			"gas",
			"[0 - выкл. | вкл. клапан №] - управление пневмоблоком",
			1,
			func(b []byte) error {
				return x.SwitchPneumo(b[0])
			},
		},
		{
			"key",
			"[№ ключа] [0 - вкл. | 1 - выкл.] - комутация ключа",
			2,
			func(b []byte) error {
				return x.Switch(b[0], b[1])
			},
		},
		{
			"sound",
			"[0 - вкл. | 1 - выкл.] - звук",
			1,
			func(b []byte) error {
				return x.Sound(b[0]!=0)
			},
		},
	}
	xs = append(xs, consoleInputCommand{"help", "описание команд", 0, func(b []byte) error {
		for _, y := range xs {
			fmt.Printf("%-5s - %s\n", y.s, y.d)
		}
		println("[q | quit]  - выход")
		return nil
	}})
	return xs

}

func (x *App) consoleInput() error {

	x.mw.Synchronize(x.mw.Hide)
	defer x.mw.Synchronize(x.mw.Show)

	var input string
	printlnError := func(err error) {
		if err == nil {
			return
		}
		println(err.Error())
	}
	ct.Foreground(ct.Cyan,true)
	fmt.Println("Ввод: начало")
	ct.ResetColor()
Loop1:
	for {
		print("> ")
		for _, err := fmt.Scanf("%s", &input); err != nil; _, err = fmt.Scanf("%s", &input) {

		}
		input = strings.ToLower(input)
		if input == "q" || input == "quit" {
			ct.Foreground(ct.Cyan,true)
			fmt.Println("Ввод: окончание")
			ct.ResetColor()
			return nil
		}
		for _,cmd := range x.consoleInputCommands{
			b := make([]byte,cmd.n)
			if cmd.s == input{
				for i:=0; i<cmd.n; i++{
					if _, err := fmt.Scanf("%d", &b[i]); err != nil {
						printlnError(err)
					    continue Loop1
					}
				}
				printlnError(cmd.a(b))
				continue Loop1
			}
		}
		println("unexpected:", input)
	}
}

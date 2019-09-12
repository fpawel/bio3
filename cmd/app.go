package main

import (
	"fmt"

	"github.com/boltdb/bolt"
	sph "github.com/fpawel/bio3/internal/comport"
	"github.com/fpawel/bio3/internal/products"
	"github.com/lxn/walk"

	"log"
	"os"
	"path/filepath"

	"github.com/fpawel/bio3/internal/fetch"
	"runtime"

	"time"

	"encoding/json"
	"github.com/fpawel/bio3/internal/productsview"
	"io/ioutil"
)

type App struct {
	*walk.Application
	db                   products.DB
	mw                   *AppMainWindow
	consoleInputCommands []consoleInputCommand
	config               AppConfig
	workHelper           *WorkHelper
	tableWorksModel      *TableWorksModel
	tableProductsModel   *TableProductsModel
	tableLogsModel       *productsview.TableLogsViewModel
}


type WorkHelper struct {
	port             *sph.Port
	cancelDelay      int32
	switchKey   map[byte]byte
	pneumoPoint byte
	sound       bool
}

type AppConfig struct {
	Port      string
	ReadWrite fetch.Config
}



func check(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Panicf("%s:%d %v\n", filepath.Base(file), line, err)
	}
}

func NewApp() *App {

	x := &App{
		Application:      walk.App(),
		config: AppConfig{
			Port: "COM1",
			ReadWrite: fetch.Config{
				MaxAttemptsRead: 1,
				ReadTimeout:     1000 * time.Millisecond,
				ReadByteTimeout: 50 * time.Millisecond,
			},
		},
	}
	x.mw = NewAppMainWindow(x)

	// инициализация тестов
	x.initWorks()
	// команд консольного ввода
	x.consoleInputCommands = x.ConsoleInputCommands()

	x.SetOrganizationName("Аналитприбор")
	x.SetProductName("БИО-3")

	{
		sets := walk.NewIniFileSettings("settings.ini")
		if err := sets.Load(); err != nil {
			println("load settings.ini error:", err)
		}
		x.SetSettings(sets)
	}


	fmt.Println(x.FolderPath())

	// считать настройки приложения из сохранённого файла json
	{
		b, err := ioutil.ReadFile(x.ConfigPath())
		if err != nil {
			fmt.Print("config.json error:", err)
		} else {
			if err := json.Unmarshal(b, &x.config); err != nil {
				fmt.Print("config.json content error:", err)
			}
		}
	}

	// создать каталог с данными и настройками программы если его нет
	if _, err := os.Stat(x.FolderPath()); os.IsNotExist(err) {
		x.SaveIni()
	}

	var err error
	x.db.DB, err = bolt.Open(x.FolderPath()+"/products.db", 0600, nil)
	check(err)

	x.db.Update(func(tx products.Tx) {
		if len(tx.Party().Products())==0{
			tx.Party().AddProducts(20)
		}
	})

	x.tableProductsModel = NewTableProductsModel(x.db, x.mw.SelectedTest)
	x.tableWorksModel = NewTableWorksViewModel(x.db)
	x.tableLogsModel = productsview.NewTableLogsViewModel(x.db, func(tx products.Tx) (r [][]byte) {
		x.db.View(func(tx products.Tx) {
			r = tx.Party().Path()
		})
		return
	})

	check( NewMainwindow(x.mw).Create() )
	x.mw.Initialize()

	return x
}

func (x *App) ProductValue(p products.Product, row int, col int) interface{} {
	switch col {
	case 0:
		return fmt.Sprintf("%02d: %d", p.Addr(), p.Serial())
	default:
		return nil
	}
}

func (x *App) Title() string {
	return x.OrganizationName() + ". " + x.ProductName()
}


func (x *App) FolderPath() string {

	appDataPath, err := walk.AppDataPath()
	check(err)
	return filepath.Join(
		appDataPath,
		x.OrganizationName(),
		x.ProductName() )
}

func (x *App) ConfigPath() string  {
	return filepath.Join( x.FolderPath(), "config.json" )
}

func (x *App) SaveIni() {
	check(x.Settings().(*walk.IniFileSettings).Save())
}

func (x *App) SaveConfig() {
	b,err := json.Marshal(x.config)
	check(err)
	check( ioutil.WriteFile(x.ConfigPath(), b, 0644) )
}

func (x *App) Close() {
	x.SaveIni()
	x.SaveConfig()
	check( x.db.DB.Close() )
	//check(x.mw.notifyIcon.Dispose())
}





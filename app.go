package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

type Result struct {
	Data []Item
}

// App struct
type App struct {
	ctx        context.Context
	filePath   string
	results    []Item
	cache      []Item
	originData []Item
	num        int
	total      int
}

func (a *App) GetTotal() int {
	return a.total
}

func (a *App) GetNum() int {
	return a.num
}

func (a *App) SetNum(v int) {
	a.num = v
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		num: 10,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectFile() (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "excel Files (*.XLSX, *.xlsx)",
				Pattern:     "*.XLSX;*.xlsx",
			},
		},
	})
	if err != nil || filePath == "" {
		return "", nil
	}
	a.filePath = filePath
	return filePath, nil
}

func (a *App) CleanCache() {
	a.results = []Item{}
	a.cache = []Item{}
	a.originData = []Item{}
}

func (a *App) CleanFile() {
	a.filePath = ""
	a.total = 10
	a.CleanCache()
}

func (a *App) openFile() ([][]string, error) {
	f, err := excelize.OpenFile(a.filePath)
	defer func() {
		if f == nil {
			return
		}
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		return [][]string{}, err
	}
	rows, err := f.GetRows("Sheet1")
	return rows, nil
}

type Item struct {
	Name string
	Url  string
}

func (a *App) getResult() ([]Item, error) {
	if len(a.originData) < a.num {
		return []Item{}, errors.New("可抽奖人数不足")
	}
	// 设置随机种子
	rand.NewSource(time.Now().UnixNano())
	// 创建一个用于存储抽奖结果的切片
	winners := make([]Item, 0)
	// 执行抽奖
	for len(winners) < a.num {
		// 生成一个随机索引
		index := rand.Intn(len(a.originData))
		// 从原始数据中取出对应索引的元素
		winner := a.originData[index]

		// 将中奖者添加到结果中
		winners = append(winners, winner)
		// 从原始数据中移除已经中奖的元素
		a.originData = append(a.originData[:index], a.originData[index+1:]...)
	}
	return winners, nil
}

func (a *App) Analyse() error {
	rows, err := a.openFile()
	if err != nil {
		return err
	}
	if len(rows) <= 1 {
		return errors.New("表格无内容或格式错误")
	}
	list := make([]Item, 0)
	for _, item := range rows[1:] {
		list = append(list, Item{
			Name: item[0],
			Url:  item[1],
		})
	}
	a.originData = list
	a.total = len(list)
	return nil
}
func (a *App) Lottery() (*Result, error) {
	resList, err := a.getResult()
	a.cache = append(a.cache, resList...)
	if err != nil {
		return &Result{
			Data: resList,
		}, err
	}
	a.results = resList
	return &Result{
		Data: resList,
	}, nil
}

func (a *App) SaveAs() error {
	str, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "中奖名单.xlsx",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "excel Files (*.XLSX, *.xlsx)",
				Pattern:     "*.XLSX;*.xlsx",
			},
		},
	})
	if str == "" || err != nil {
		return errors.New("保存文件出错")
	}
	err = a.export(str)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) export(filename string) error {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	var names [][]interface{}
	names = append(names, []interface{}{"推特名称", "地址"})
	for _, item := range a.cache {
		names = append(names, []interface{}{item.Name, item.Url})
	}
	for index, item := range names {
		cell, _ := excelize.CoordinatesToCellName(1, index+1)
		if err := streamWriter.SetRow(cell, item); err != nil {
			continue
		}
	}
	if err := streamWriter.Flush(); err != nil {
		return err
	}
	if err := file.SaveAs(filename); err != nil {
		return err
	}
	return nil
}

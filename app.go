package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
	"math/rand"
	"time"
)

type Result struct {
	Data []Item
}

// App struct
type App struct {
	ctx      context.Context
	filePath string
	results  []Item
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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
		return "", errors.New("选择文件出错")
	}
	a.filePath = filePath
	return filePath, nil
}

func (a *App) CleanCache() {
	a.results = []Item{}
}

func (a *App) CleanFile() {
	a.filePath = ""
	a.results = []Item{}
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

func (a *App) getResult(data [][]string) []Item {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 创建一个用于存储抽奖结果的切片
	winners := make([]Item, 0)
	// 执行抽奖
	for len(winners) < 10 {
		// 生成一个随机索引
		index := rand.Intn(len(data))
		// 从原始数据中取出对应索引的元素
		winner := data[index]
		item := Item{
			Name: winner[0],
			Url:  winner[1],
		}
		// 将中奖者添加到结果中
		winners = append(winners, item)
		// 从原始数据中移除已经中奖的元素
		data = append(data[:index], data[index+1:]...)
	}
	return winners
}

func (a *App) Analyse() (*Result, error) {
	rows, err := a.openFile()
	if err != nil {
		return nil, err
	}
	if len(rows) <= 1 {
		return nil, errors.New("表格无内容或格式错误")
	}
	resList := a.getResult(rows[1:])
	a.results = resList
	return &Result{
		Data: resList,
	}, nil
}

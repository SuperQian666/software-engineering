package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
	edit  *walk.TextEdit
	inTE  *walk.TextEdit
	outTE *walk.TextEdit
}

// 声明两个文本域控件
func main() {
	//var //声明两个文本域控件
	mw := &MyMainWindow{}
	//配置主窗口，并运行起来
	err := MainWindow{
		AssignTo: &mw.MainWindow, //窗口重定向至mw，重定向后可由重定向变量控制控件
		Icon:     "test.ico",     //窗体图标
		Title:    "文件备份（01）",     //标题
		MinSize:  Size{Width: 400, Height: 300},
		Size:     Size{600, 400},
		MenuItems: []MenuItem{
			Menu{
				Text: "文件",
				Items: []MenuItem{
					Action{
						Text: "打开文件",
						Shortcut: Shortcut{ //定义快捷键后会有响应提示显示
							Modifiers: walk.ModControl,
							Key:       walk.KeyO,
						},
						OnTriggered: mw.selectFile, //openFileActionTriggered, //点击动作触发响应函数
					},
					Action{
						Text: "另存为",
						Shortcut: Shortcut{
							Modifiers: walk.ModControl | walk.ModShift,
							Key:       walk.KeyS,
						},
						OnTriggered: mw.saveFile,
					},
					Action{
						Text: "退出",
						OnTriggered: func() {
							mw.Close()
						},
					},
				},
			},
			Menu{
				Text: "帮助",
				Items: []MenuItem{
					Action{
						Text: "关于",
						OnTriggered: func() {
							walk.MsgBox(mw, "关于", "这是一个菜单和工具栏的实例",
								walk.MsgBoxIconInformation|walk.MsgBoxDefButton1)
						},
					},
				},
			},
		},

		Layout: VBox{}, //样式，纵向
		OnDropFiles: func(files []string) {
			mw.inTE.SetText(strings.Join(files, "\r\n"))
		},
		Children: []Widget{ //控件组
			TextEdit{
				AssignTo: &mw.edit,
			},
			TextEdit{AssignTo: &mw.inTE,
				ReadOnly: true,
				Text:     "Drop files here, from windows explorer...",
			},
			HSplitter{
				Children: []Widget{
					Label{
						Text: "待输入路径",
					},
					TextEdit{AssignTo: &mw.inTE},
					PushButton{
						Text:      "打开",
						OnClicked: mw.selectFile, //点击事件响应函数
					},
				},
			},
			HSplitter{
				Children: []Widget{

					Label{
						Text: "待输出路径",
					},
					TextEdit{AssignTo: &mw.outTE},
					PushButton{
						Text:      "另存为",
						OnClicked: mw.saveFile,
					},
				},
			},
		},
	}.Create() //创建

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mw.Run() //运行
}

func (mw *MyMainWindow) selectFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "文本文件 (*.txt)|*.txt|所有文件 (*.*)|*.*"

	mw.inTE.SetText("") //通过重定向变量设置TextEdit的Text
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("错误 : 打开文件时\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("用户取消\r\n")
		return
	}
	s := fmt.Sprintf("选择了: %s\r\n", dlg.FilePath)
	mw.edit.AppendText(s)
	s1 := fmt.Sprintf("%s\r\n", dlg.FilePath)
	mw.inTE.AppendText(s1)
}

func (mw *MyMainWindow) saveFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "另存为"
	dlg.Filter = "文本文件 (*.txt)|*.txt|所有文件 (*.*)|*.*"

	if ok, err := dlg.ShowSave(mw); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	} else if !ok {
		fmt.Fprintln(os.Stderr, "取消")
		return
	}

	data := mw.inTE.Text()
	filename := dlg.FilePath
	mw.outTE.AppendText(filename)
	f, err := os.Open(filename)
	if err != nil {
		f, _ = os.Create(filename)
	} else {
		f.Close()
		//打开文件，参数：文件路径及名称，打开方式，控制权限
		f, err = os.OpenFile(filename, os.O_WRONLY, 0x666)
	}
	if len(data) == 0 {
		f.Close()
		return
	}
	//copy
	io.Copy(f, strings.NewReader(data))
	f.Close()
}

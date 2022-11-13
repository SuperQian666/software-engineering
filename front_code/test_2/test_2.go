package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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
		Title:    "文件备份02",       //标题
		MinSize:  Size{Width: 450, Height: 600},
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
							walk.MsgBox(mw, "关于", "文件备份软件",
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
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text:      "加密",
						OnClicked: mw.file_Encrypt_lrx, //点击事件响应函数
					},
					PushButton{
						Text:      "解密",
						OnClicked: mw.file_Decrypt_lrx, //点击事件响应函数
					},
					PushButton{
						Text:      "备份",
						OnClicked: mw.file_Copy_lrx, //点击事件响应函数
					},
					PushButton{
						Text:      "还原",
						OnClicked: mw.file_Restore_lrx, //点击事件响应函数
					},
				},
			},

			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text:      "压缩",
						OnClicked: mw.file_Compress_lrx, //点击事件响应函数
					},
					PushButton{
						Text:      "解压",
						OnClicked: mw.file_Decompress_lrx, //点击事件响应函数
					},
					PushButton{
						Text:      "云备份",
						OnClicked: mw.file_Copy_Cloud_lrx, //点击事件响应函数
					},
					PushButton{
						Text:      "云还原",
						OnClicked: mw.file_Restore_Cloud_lrx, //点击事件响应函数
					},
				},
			},

			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					TextEdit{
						AssignTo: &mw.inTE,
						ReadOnly: true,
						Text:     "Drop files here, from windows explorer...",
					},
				},
			},
			GroupBox{
				Layout: HBox{},
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
			GroupBox{
				Layout: HBox{},
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
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					TableView{
						Name:             "tableView", // Name is needed for settings persistence
						AlternatingRowBG: true,
						ColumnsOrderable: true,
						Columns: []TableViewColumn{
							// Name is needed for settings persistence
							{Name: "#", DataMember: "Id", Width: 40}, // Use DataMember, if names differ
							{Name: "FileName"},
							{Name: "SourcePath", Width: 180},
							{Name: "TargetPath", Width: 180},
							{Name: "SaveTime", Format: "2006-01-02 15:04:05", Width: 150},
						},
						Model: NewFooModel(),
					},
				},
			},
			TextEdit{
				AssignTo: &mw.edit,
			},
			PushButton{
				Text: "退出",
				OnClicked: func() {
					mw.Close()
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

// 选择需要保存的文件
func (mw *MyMainWindow) selectFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "所有文件 (*.*)|*.*|图片(*.gif;*.jpg;*.jpeg;*.bmp;*.png)|*.gif;*.jpg;*.jpeg;*.bmp;*.png;|word文件(*.doc)|*.doc|excel文件(*.xls)|*.xls|文本文件 (*.txt)|*.txt"

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

// 选择保存路径并保存待输入路径
func (mw *MyMainWindow) saveFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "另存为"
	dlg.Filter = "所有文件 (*.*)|*.*|图片(*.gif;*.jpg;*.jpeg;*.bmp;*.png)|*.gif;*.jpg;*.jpeg;*.bmp;*.png;|word文件(*.doc)|*.doc|excel文件(*.xls)|*.xls|文本文件 (*.txt)|*.txt"

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

// 加密
func (mw *MyMainWindow) file_Encrypt_lrx() {
	mw.showNoneMessage("加密成功")

}

// 解密
func (mw *MyMainWindow) file_Decrypt_lrx() {
	mw.showNoneMessage("解密成功")
}

// 备份
func (mw *MyMainWindow) file_Copy_lrx() {
	mw.showNoneMessage("备份成功")
}

// 还原
func (mw *MyMainWindow) file_Restore_lrx() {
	mw.showNoneMessage("还原成功")
}

// 压缩
func (mw *MyMainWindow) file_Compress_lrx() {
	mw.showNoneMessage("压缩成功")
}

// 解压
func (mw *MyMainWindow) file_Decompress_lrx() {
	mw.showNoneMessage("解压成功")
}

// 云备份
func (mw *MyMainWindow) file_Copy_Cloud_lrx() {
	mw.showNoneMessage("云备份成功")
}

// 云解压
func (mw *MyMainWindow) file_Restore_Cloud_lrx() {
	mw.showNoneMessage("云解压成功")
}

// 提示框
func (mw *MyMainWindow) showNoneMessage(message string) {
	walk.MsgBox(mw, "提示", message, walk.MsgBoxIconInformation)
}

// 显示已备份的文件
func NewFooModel() *FooModel {
	//now := time.Now()

	//rand.Seed(now.UnixNano())

	m := &FooModel{items: make([]*Foo, 20)}

	for i := range m.items {
		m.items[i] = &Foo{
			Id:         i,
			FileName:   "文件名",
			SourcePath: "源路径",
			TargetPath: "目标路径",
			SaveTime:   time.Now(), //time.Unix(rand.Int63n(now.Unix()), 0),
		}
	}

	return m
}

type FooModel struct {
	walk.SortedReflectTableModelBase
	items []*Foo
}

func (m *FooModel) Items() interface{} {
	return m.items
}

type Foo struct {
	Id         int
	FileName   string
	SourcePath string
	TargetPath string
	SaveTime   time.Time
}

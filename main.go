package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"soft/compress"
	"soft/copy"
	"soft/encrypt"
	"soft/network"
	"strings"
)

type MyMainWindow struct {
	*walk.MainWindow
	edit      *walk.TextEdit
	inTE      *walk.TextEdit
	outTE     *walk.TextEdit
	secretKey *walk.LineEdit
}

// 声明两个文本域控件
func main() {
	//var //声明两个文本域控件
	mw := &MyMainWindow{}
	//配置主窗口，并运行起来
	err := MainWindow{
		AssignTo: &mw.MainWindow, //窗口重定向至mw，重定向后可由重定向变量控制控件
		//Icon:     "./frontcode/test.ico", //窗体图标
		Title:   "文件备份", //标题
		MinSize: Size{Width: 450, Height: 600},
		//Size:     Size{600, 400},
		MenuItems: []MenuItem{
			Menu{
				Text: "文件",
				Items: []MenuItem{
					Action{
						Text: "输入路径",
						Shortcut: Shortcut{ //定义快捷键后会有响应提示显示
							Modifiers: walk.ModControl,
							Key:       walk.KeyO,
						},
						OnTriggered: mw.inFile, //openFileActionTriggered, //点击动作触发响应函数
					},
					Action{
						Text: "输出路径",
						Shortcut: Shortcut{
							Modifiers: walk.ModControl | walk.ModShift,
							Key:       walk.KeyS,
						},
						OnTriggered: mw.outFile,
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
						OnClicked: mw.fileEncrypt, //点击事件响应函数
					},
					PushButton{
						Text:      "解密",
						OnClicked: mw.fileDecrypt, //点击事件响应函数
					},
					PushButton{
						Text:      "备份",
						OnClicked: mw.fileCopy, //点击事件响应函数
					},
					PushButton{
						Text:      "还原",
						OnClicked: mw.fileRestore, //点击事件响应函数
					},
				},
			},

			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text:      "压缩",
						OnClicked: mw.fileCompress, //点击事件响应函数
					},
					PushButton{
						Text:      "解压",
						OnClicked: mw.fileDecompress, //点击事件响应函数
					},
					PushButton{
						Text:      "云备份",
						OnClicked: mw.fileCopyToCloud, //点击事件响应函数
					},
					PushButton{
						Text:      "云取回",
						OnClicked: mw.fileRestoreFromCloud, //点击事件响应函数
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
						Text:      "选择文件",
						OnClicked: mw.inFile,
					},
					PushButton{
						Text:      "选择文件夹",
						OnClicked: mw.inFolder,
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
						Text:      "选择文件",
						OnClicked: mw.outFile,
					},
					PushButton{
						Text:      "选择文件夹",
						OnClicked: mw.outFolder,
					},
				},
			},
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "加密/解密密钥",
					},
					LineEdit{
						AssignTo:     &mw.secretKey,
						PasswordMode: true,
					},
				},
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
func (mw *MyMainWindow) inFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "所有文件 (*.*)|*.*|图片(*.gif;*.jpg;*.jpeg;*.bmp;*.png)|*.gif;*.jpg;*.jpeg;*.bmp;*.png;|word文件(*.doc)|*.doc|excel文件(*.xls)|*.xls|文本文件 (*.txt)|*.txt"
	mw.inTE.SetText("") //通过重定向变量设置TextEdit的Text
	if ok, err := dlg.ShowOpenMultiple(mw); err != nil {
		mw.edit.AppendText("错误 : 打开文件时\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("用户取消\r\n")
		return
	}
	var s string
	for _, f := range dlg.FilePaths {
		s += f + ";"
	}
	mw.edit.AppendText(fmt.Sprintf("选择了: %s\r\n", s))
	s1 := fmt.Sprintf("%s", s)
	mw.inTE.AppendText(s1)
}

func (mw *MyMainWindow) inFolder() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件夹"

	mw.inTE.SetText("") //通过重定向变量设置TextEdit的Text
	if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
		mw.edit.AppendText("错误 : 打开文件时\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("用户取消\r\n")
		return
	}
	s := fmt.Sprintf("选择了: %s\r\n", dlg.FilePath)
	mw.edit.AppendText(s)
	s1 := fmt.Sprintf("%s", dlg.FilePath)
	mw.inTE.AppendText(s1)
}

// 选择保存路径并保存待输入路径
func (mw *MyMainWindow) outFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "文件夹(*)|所有文件 (*.*)|*.*|图片(*.gif;*.jpg;*.jpeg;*.bmp;*.png)|*.gif;*.jpg;*.jpeg;*.bmp;*.png;|word文件(*.doc)|*.doc|excel文件(*.xls)|*.xls|文本文件 (*.txt)|*.txt"

	mw.outTE.SetText("")
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("错误 : 打开文件时\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("用户取消\r\n")
		return
	}
	s := fmt.Sprintf("选择了: %s\r\n", dlg.FilePath)
	mw.edit.AppendText(s)
	s1 := fmt.Sprintf("%s", dlg.FilePath)
	mw.outTE.AppendText(s1)
}

func (mw *MyMainWindow) outFolder() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"

	mw.outTE.SetText("")
	if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
		mw.edit.AppendText("错误 : 打开文件时\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("用户取消\r\n")
		return
	}
	s := fmt.Sprintf("选择了: %s\r\n", dlg.FilePath)
	mw.edit.AppendText(s)
	s1 := fmt.Sprintf("%s", dlg.FilePath)
	mw.outTE.AppendText(s1)
}

// 加密
func (mw *MyMainWindow) fileEncrypt() {
	if mw.secretKey.Text() == "" || len(mw.secretKey.Text()) < 5 {
		mw.showNoneMessage("密码至少为5位")
		return
	}
	if err := encrypt.EncryptFile(mw.inTE.Text(), mw.outTE.Text(), mw.secretKey.Text()); err != nil {
		mw.showNoneMessage(fmt.Sprintf("%s:加密失败：%s\r\n", mw.inTE.Text(), err))
	} else {
		mw.showNoneMessage(mw.inTE.Text() + ":" + "加密成功")
	}
}

// 解密
func (mw *MyMainWindow) fileDecrypt() {
	if mw.secretKey.Text() == "" || len(mw.secretKey.Text()) < 5 {
		mw.showNoneMessage("密码至少为5位")
		return
	}
	if err := encrypt.DecryptFile(mw.inTE.Text(), mw.outTE.Text(), mw.secretKey.Text()); err != nil {
		mw.showNoneMessage(fmt.Sprintf("%s:解密失败：%s\r\n", mw.inTE.Text(), err))
	} else {
		mw.showNoneMessage(mw.inTE.Text() + ":" + "解密成功")
	}

}

// 备份
func (mw *MyMainWindow) fileCopy() {
	addArr := strings.Split(mw.inTE.Text(), ";")
	msg := ""
	for _, i := range addArr {
		if err := copy.Copy(i, mw.outTE.Text()); err != nil {
			msg += fmt.Sprintf("%s:备份失败：%s\r\n", i, err)

		} else {
			msg += fmt.Sprintf("%s:备份成功\r\n", i)
		}
	}
	mw.showNoneMessage(msg)
}

// 还原
func (mw *MyMainWindow) fileRestore() {
	addArr := strings.Split(mw.inTE.Text(), ";")
	msg := ""
	for _, i := range addArr {
		if err := copy.Copy(i, mw.outTE.Text()); err != nil {
			msg += fmt.Sprintf("%s:还原失败：%s\r\n", i, err)

		} else {
			msg += fmt.Sprintf("%s:还原成功\r\n", i)
		}
	}
	mw.showNoneMessage(msg)
}

// 压缩
func (mw *MyMainWindow) fileCompress() {

	var msg = ""
	if err := compress.Zip(mw.inTE.Text(), mw.outTE.Text()); err != nil {
		msg += fmt.Sprintf("%s:压缩失败：%s\r\n", mw.inTE.Text(), err)

	} else {
		msg += fmt.Sprintf("%s:压缩成功\r\n", mw.inTE.Text())
	}
	mw.showNoneMessage(msg)
}

// 解压
func (mw *MyMainWindow) fileDecompress() {
	var msg = ""
	if err := compress.UnZip(mw.inTE.Text(), mw.outTE.Text()); err != nil {
		msg += fmt.Sprintf("%s:解压失败：%s\r\n", mw.inTE.Text(), err)

	} else {
		msg += fmt.Sprintf("%s:解压成功\r\n", mw.inTE.Text())
	}
	mw.showNoneMessage(msg)
}

// 云备份
func (mw *MyMainWindow) fileCopyToCloud() {
	var msg = ""
	if err := network.Upload(mw.inTE.Text(), mw.outTE.Text()); err != nil {
		msg += fmt.Sprintf("%s:云备份失败：%s\r\n", mw.inTE.Text(), err)
	} else {
		msg += fmt.Sprintf("%s:云备份成功\r\n", mw.inTE.Text())
	}
	mw.showNoneMessage(msg)
}

// 云取回
func (mw *MyMainWindow) fileRestoreFromCloud() {
	var msg = ""
	if err := network.Download(mw.inTE.Text(), mw.outTE.Text()); err != nil {
		msg += fmt.Sprintf("%s:云取回失败：%s\r\n", mw.inTE.Text(), err)
	} else {
		msg += fmt.Sprintf("%s:云取回成功\r\n", mw.inTE.Text())
	}
	mw.showNoneMessage(msg)
}

// 提示框
func (mw *MyMainWindow) showNoneMessage(message string) {
	walk.MsgBox(mw, "提示", message, walk.MsgBoxIconInformation)
}

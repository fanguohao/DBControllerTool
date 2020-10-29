package main

import (
	"fmt"
	"github.com/lxn/win"
	//"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	DBcontroller "walkTest/DBController"
)

type MyMainWindow struct {
	*walk.MainWindow
	edit         *walk.TextEdit
	dbAdd        *walk.LineEdit
	inPutaddress *walk.LineEdit
	db           *walk.DataBinder
	combox       *walk.ComboBox
}

//绑定信息
type Cfg struct {
	DbAddress        string   // 数据库地址；
	DatabaseName     string   //数据库名;
	UserName         string   //用户名；
	Password         string   // 密码；
	FileAddress      string   //文件地址;
	ToVersion        float32  //要升级的版本号；
	VersionToDisplay []string // 用于列举需要升级的版本；
}

var (
	cfg = &Cfg{
		FileAddress:      "请输入地址! eg: c://fileAddress",
		ToVersion:        0.0,
		DbAddress:        "127.0.0.1:3306",
		DatabaseName:     "elmcms",
		UserName:         "root",
		Password:         "Fan003174",
		VersionToDisplay: []string{"v1", "v2", "v3", "v4"},
	}
)

var (
	mw         = &MyMainWindow{}
	windowMain *walk.MainWindow
)

var (
	ConnectInfoBool bool
	DBVersion       float32
)

const (
	SIZE_W = 600
	SIZE_H = 400
)

func main() {
	MainWindow{
		Title:    "数据库版本管理工具 v1.0",
		Size:     Size{Width: 700, Height: 300},
		Layout:   VBox{},
		AssignTo: &windowMain,
		Children: widget,

		DataBinder: DataBinder{
			AssignTo:   &mw.db,
			DataSource: cfg,
		},
	}.Create()
	win.SetWindowLong(
		windowMain.Handle(), win.GWL_STYLE,
		win.GetWindowLong(windowMain.Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME,
	)
	//xScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
	//yScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
	//win.SetWindowPos(
	//	mw.Handle(),
	//	0,
	//	(xScreen-SIZE_W)/2,
	//	(yScreen-SIZE_H)/2,
	//	SIZE_W,
	//	SIZE_H,
	//	win.SWP_FRAMECHANGED,
	//)
	windowMain.Run()
}
func ShowMsgBox(title, msg string) int {
	return walk.MsgBox(windowMain, title, msg, walk.MsgBoxOK)
}

// 控件
var widget = []Widget{
	Composite{
		Layout: HBox{},
		Children: []Widget{
			Label{Text: "数据库地址:"},
			LineEditDBaddress,
			Label{Text: "数据库名:"},
			LineEditDBName,
			Label{Text: "用户名:"},
			LineEditUserName,
			Label{Text: "密码:"},
			LineEditPassword,
			PushButtonVerify,
			//Label{Text: "当前版本："},
			//Label{Text: v},

			//PushButton{
			//	StretchFactor: 8,
			//	Text:          "摸我有惊喜",
			//	OnClicked: func() {
			//		ShowMsgBox("天降惊喜", "测试")
			//	},
			//},
		},
	},
	Composite{
		Layout: HBox{},
		Children: []Widget{
			Label{Text: "选择升级文件:"},
			MyLineEdit,
			PushButtonSelect,
			//Label{Text: "升级到:"},
			//Combox,
			PushButtonOK,
		},
	},
	MyTextEdit,
}

var PushButtonVerify = PushButton{
	Text:      "连接验证",
	OnClicked: mw.VerifyDB,
}

var LineEditDBaddress = LineEdit{
	Text:     Bind("DbAddress"),
	AssignTo: &mw.dbAdd,
}

var LineEditDBName = LineEdit{
	Text:     Bind("DatabaseName"),
	AssignTo: &mw.dbAdd,
}
var LineEditUserName = LineEdit{
	Text:     Bind("UserName"),
	AssignTo: &mw.dbAdd,
}

var LineEditPassword = LineEdit{

	Text:     Bind("Password"),
	AssignTo: &mw.dbAdd,
}

var PushButtonOK = PushButton{
	Text:      "升级",
	OnClicked: mw.update,
}
var MyLineEdit = LineEdit{
	AssignTo: &mw.inPutaddress,
}

var MyTextEdit = TextEdit{
	AssignTo: &mw.edit,
}

var PushButtonSelect = PushButton{
	Text:      "选择",
	OnClicked: mw.selectFile,
}

func (mw *MyMainWindow) lb_CurrentIndexChanged() {
	i := mw.combox.CurrentIndex()
	//item := &mw.model.items[i]

	//mw.te.SetText(item.value)

	fmt.Println("CurrentIndex: ", i)
	//fmt.Println("CurrentEnvVarName: ", item.name)
}

var Combox = ComboBox{
	AssignTo: &mw.combox,
	//MinSize:  Size{300, 0},
	Value: Bind("ToVersion"),
	// todo:动态加载版本；
	Model: cfg.VersionToDisplay,
	//OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
	//OnKeyPress: df
}

func dialogUpdataOK(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	//var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	msg := fmt.Sprintf("当前数据库将升级到：v%.2f \r\n", cfg.ToVersion)
	return Dialog{
		AssignTo:      &dlg,
		Title:         "升级提示",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{100, 150},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					Label{Text: msg},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "确定",
						OnClicked: func() {
							dlg.Accept()
							s := DBcontroller.Update(cfg.FileAddress, cfg.ToVersion)
							mw.edit.AppendText(s + "\r\n")
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}

func dialogUpdataFalse(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	//var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	s := fmt.Sprintf("err : %s\r\n", "请先连接数据库！")
	/*mw.edit.AppendText(s)*/
	return Dialog{
		AssignTo:      &dlg,
		Title:         "升级提示",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{100, 150},
		Layout:        VBox{},
		Children: []Widget{
			Label{Text: s},
			PushButton{
				Text: "确定",
				OnClicked: func() {
					dlg.Accept()
				},
			},
		},
	}.Run(owner)
}

func dialogVerifyTrue(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	//var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	s := fmt.Sprintf("     数据库连接成功！\n 当前数据库版本为：v%.2f", DBVersion)
	return Dialog{
		AssignTo:      &dlg,
		Title:         "验证提示",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{100, 150},
		Layout:        VBox{},
		Children: []Widget{
			Label{Text: s},
			PushButton{
				Text: "确定",
				OnClicked: func() {
					dlg.Accept()
				},
			},
		},
	}.Run(owner)
}

func dialogVerifyFalse(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	//var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		Title:         "提示",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{100, 150},
		Layout:        VBox{},
		Children: []Widget{
			Label{Text: "  数据库连接失败，\n    请核查信息！"},
			PushButton{
				Text: "确定",
				OnClicked: func() {
					dlg.Accept()
				},
			},
		},
	}.Run(owner)
}

func (mw *MyMainWindow) VerifyDB() {
	mw.db.Submit()
	ConnectInfoBool = DBcontroller.GetConnectionIfo(cfg.UserName, cfg.Password, cfg.DatabaseName)
	DBVersion = DBcontroller.GetDBVersion()
	if ConnectInfoBool == true {
		//fmt.Println("连接成功")
		//s := fmt.Sprintf("数据库连接成功！当前数据库版本为：v%.2f \r\n", DBVersion)
		//mw.edit.AppendText(s)
		dialogVerifyTrue(windowMain)
	} else {
		dialogVerifyFalse(windowMain)
	}
}

func (mw *MyMainWindow) update() {
	mw.db.Submit()
	// 连接状态判断；
	if ConnectInfoBool != true {

		dialogUpdataFalse(windowMain)

	} else {
		//fmt.Println(cfg.ToVersion)
		//v, _ := strconv.Atoi(cfg.ToVersion)
		//v, _ := strconv.ParseFloat(cfg.ToVersion, 32)
		//s := DBcontroller.Update(cfg.FileAddress, cfg.ToVersion)
		//mw.edit.AppendText(s + "\r\n")
		dialogUpdataOK(windowMain)
	}
}

func (mw *MyMainWindow) selectFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "可执行文件 (*.sql)|*.sql|所有文件 (*.*)|*.*"

	//mw.edit.SetText("") //通过重定向变量设置TextEdit的Text
	if ok, err := dlg.ShowOpen(windowMain); err != nil {
		mw.edit.AppendText("Error : File Open\r\n")
		return
	} else if !ok {
		//mw.edit.AppendText("Cancel\r\n")
		return
	}
	s := fmt.Sprintf("%s\r\n", dlg.FilePath)

	//cfg.VersionToDisplay = DBcontroller.GetVersionDisplay()
	cfg.ToVersion = DBcontroller.GetToversion(dlg.FilePath)
	//s2 := fmt.Sprintf("版本将升级到：v%.2f \r\n", cfg.ToVersion)
	//mw.edit.AppendText(s2)
	mw.inPutaddress.SetText(s)
	cfg.FileAddress = dlg.FilePath
}

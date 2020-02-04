package Manager

import (
	"fmt"
)

type GUI struct {
	FileList []string
	FileLen []int
	Len int
	LocalIP []string
	func_list map[string]func()
}

var port string = ":10240"

var doc = map[string]string{
	"ls_file" : "显示当前目录文件",
	"ls_dir" : "显示当前目录文件夹",
	"accept" : "进入文件接收模式，在这个模式下等待对方通过你的 ID 向你传输文件",
	"send" : "发送文件",
	"list" : "显示指令",
	"help" : "显示指令说明",
}

func checkErr(err error, msg string) bool {
	if err != nil{
		fmt.Println("[red]" + msg + err.Error())
		return false
	}
	return true
}

// 功能起始入口
func (gui *GUI) Render() {
	gui.func_list = map[string]func(){
		"ls_file" : LsFile,
		"ls_dir" : LsDir,
		"list" : listFunc,
		"accept" : gui.Accept,
		"send" : gui.fileTransport,
		"help" : gui.helpInfo,
	}
	for {
		var cmd string
		_, err := fmt.Scanln(&cmd)
		if cmd == "quit" || err != nil{
			return
		}
		if fun, ok := gui.func_list[cmd]; ok {
			fun()
		}else {
			fmt.Println("无此命令")
		}
	}
}

// 显示帮助信息
func (gui *GUI) helpInfo() {
	fmt.Println("想要发送文件需要先让接收方进入接收模式，然后根据对方的有效ID进行传输")
	for cmd, info := range doc {
		fmt.Println(cmd + ":\t" + info)
	}
}

// 接收文件入口函数
func (gui *GUI) Accept() {
	Accept()
}

// 文件传输入口函数
func (gui *GUI) fileTransport()  {
	LsFile()
	fmt.Println(len(gui.FileList), ". 取消")
	var n int
	_, err := fmt.Scanf("%d\n", &n)
	if !checkErr(err, "输入错误: ") {
		return
	}
	if n < 0 || n > len(gui.FileList){
		fmt.Println("[red]文件选择错误！")
		return
	}else if n == len(gui.FileList){
		return
	}
	fmt.Println("请输入对方ID")
	var ip string
	_, _ = fmt.Scanln(&ip)
	Send(gui.FileList[n], ip)
}

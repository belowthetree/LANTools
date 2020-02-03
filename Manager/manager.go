package Manager

import (
	"../NetworkManager"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type GUI struct {
	FileList []string
	FileLen []int
	Len int
	LocalIP []string
	func_list map[string]func()
}

var doc = map[string]string{
	"ls_file" : "显示当前目录文件",
	"ls_dir" : "显示当前目录文件夹",
	"accept" : "进入文件接收模式，在这个模式下等待对方通过你的 ID 向你传输文件",
	"send" : "发送文件",
	"list" : "显示指令",
	"help" : "显示指令说明",
}

var port string = ":10240"

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
		"ls_file" : gui.showFileList,
		"ls_dir" : gui.showDir,
		"list" : gui.listFunc,
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
// 显示路径
func (gui *GUI) showDir() {

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
	gui.showFileList()
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
	transFile(gui.FileList[n])
}
// 发送文件的入口函数
func transFile(filename string)  {
	fmt.Println("请输入对方ID")
	var ip string
	_, _ = fmt.Scanln(&ip)
	Send(filename, ip)
}
// 启动 TCP 连接并发送
func sendFile(filename string, conn *net.TCPConn) bool {
	file, err := os.Open(filename)
	if err == nil{
		bytes := make(chan []byte, 1024)
		// 开线程获取文件字节
		go getFileBytes(file, bytes)
		// 发送文件名
		if !NetworkManager.SendBytes([]byte(filename), conn) {
			return false
		}
		info, _ := file.Stat()
		time.Sleep(time.Second)
		size := info.Size()
		if !NetworkManager.SendBytes([]byte(strconv.FormatInt(size, 10)), conn) {
			return false
		}
		//var result = true
		// 验证是否接收错误
		//go func() {
		//	tmp := make([]byte, 100)
		//	n, _ := conn.Read(tmp)
		//	if string(tmp[:n]) == "no"{
		//		result = false
		//	}
		//}()

		var waiter sync.WaitGroup
		//var num = 0
		totle := int64(0)
		// 每个字节序列都独立用线程发送
		for ctx := range bytes{
			//result = NetworkManager.SendBytes(ctx, conn)
			_, err := conn.Write(ctx)
			//if num < 0{
			//	waiter.Add(1)
			//	go func() {
			//		result = NetworkManager.SendBytes(ctx, conn)
			//		waiter.Done()
			//	}()
			//	num++
			//}else {
			//	result = NetworkManager.SendBytes(ctx, conn)
			//	num = 0
			//}
			totle += int64(len(ctx))

			if err != nil {
				return false
			}
			go func() {
				fmt.Printf("进度：%f%%\r", float64(totle) / float64(size) * 100)
			}()
		}
		fmt.Println("")
		time.Sleep(time.Microsecond*100)
		fmt.Println("发送完毕，等待结束")
		waiter.Wait()
		return true
	}else{
		fmt.Println("文件读取错误")
		return false
	}
}
// 显示当前文件列表
func (gui *GUI) showFileList() {
	gui.getFileList(false)
	fmt.Println("文件列表")
	for num, file := range gui.FileList{
		fmt.Printf("%d. %s\n", num, file)
	}
	fmt.Println()
}

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
	res := acceptFile()
	if res {
		fmt.Println("保存成功")
	}else{
		fmt.Println("保存失败")
	}
}

func acceptFile() bool {
	var listener *net.TCPListener
	var err error
	host, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0" + port)
	listener, err = net.ListenTCP("tcp", host)
	fmt.Println(host)
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer listener.Close()
	client, err := listener.AcceptTCP()
	if !checkErr(err, "接收客户端失败") {
		return false
	}
	fmt.Println("开始接收")
	defer client.Close()
	var tmp [] byte
	var result = make(chan []byte, 1024)
	var res = true
	var filename = ""
	var file *os.File
	var size, cnt int64
	size = 0
	cnt = 0
	var totle = int64(0)
	for tmp, res = NetworkManager.ReceiveBytes(client);res && tmp != nil;tmp, res = NetworkManager.ReceiveBytes(client) {
		if filename == "" {
			filename = string(tmp)
			file, err = os.Create(filename)
			if err != nil {
				return false
			}
			fmt.Println("接收文件：" + filename)
			go func() {
				res = writeFile(file, result)
			}()
			continue
		}
		if size == 0{
			size, _ = strconv.ParseInt(string(tmp), 10, 64)
			fmt.Println(size)
			fmt.Println("文件大小：", float64(size) / 1024 / 1024, "MB")
			continue
		}
		totle += int64(len(tmp))
		cnt++
		// 显示进度
		go func() {
			fmt.Printf("进度：%f%%\r", float64(cnt * 64) / float64(size) * 100)
		}()
		if !res {
			return false
		}
		result <- tmp
	}
	fmt.Println()
	fmt.Println("结束，总接收", totle,"字节")
	close(result)
	return true
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
	ip = ip + port
	host, _ := net.ResolveTCPAddr("tcp4", ip)
	conn, err := net.DialTCP("tcp", nil, host)
	if !checkErr(err, "连接失败"){
		return
	}
	defer conn.Close()
	fmt.Println("连接成功，开始发送")
	if !sendFile(filename, conn) {
		fmt.Println("发送失败")
	}
	fmt.Println("发送结束")
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
		fmt.Println(size)
		if !NetworkManager.SendBytes([]byte(strconv.FormatInt(size, 10)), conn) {
			return false
		}
		var result = true
		// 验证是否接收错误
		go func() {
			tmp := make([]byte, 100)
			n, _ := conn.Read(tmp)
			if string(tmp[:n]) == "no"{
				result = false
			}
		}()

		var waiter sync.WaitGroup
		cnt := int64(0)
		var num = 0
		totle := int64(0)
		// 每个字节序列都独立用线程发送
		for ctx := range bytes{
			if num < 0{
				waiter.Add(1)
				go func() {
					result = NetworkManager.SendBytes(ctx, conn)
					waiter.Done()
				}()
				num++
			}else {
				result = NetworkManager.SendBytes(ctx, conn)
				num = 0
			}
			totle += int64(len(ctx))

			if !result {
				return false
			}
			cnt++
			go func() {
				fmt.Printf("进度：%f%%\r", float64(cnt * 64) / float64(size) * 100)
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

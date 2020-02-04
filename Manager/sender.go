package Manager

import (
	"github.com/color"
	"net"
	"os"
	"strconv"
)

func Send(filename string, ip string) bool {
	host, _ := net.ResolveTCPAddr("tcp4", ip + filePort)
	client, err := net.DialTCP("tcp", nil, host)
	if err != nil {
		color.Red("连接对方主机失败", err)
		return false
	}
	defer client.Close()
	
	// 成功分割线----------------------
	
	file, err := os.Open(filename)
	if err != nil {
		color.Red("发送前文件打开失败", err)
		return false
	}
	info, _ := file.Stat()
	size := info.Size()
	if !sendName(filename, client){
		return false
	}
	if !sendSize(size, client){
		return false
	}
	file.Close()

	readerResult := make(chan bool)
	senderResult := make(chan bool)
	counter := make(chan int64)
	data := make(chan []byte, 1024)
	go func() {
		readerResult <- FileReader(filename, data)
	}()
	go func() {
		senderResult <- Sender(client, data, true, counter)
	}()

	go DisplayCounter(size, counter)

	if <-readerResult && <-senderResult {
		color.Yellow("发送成功")
	}else{
		color.Red("发送失败")
		return false
	}

	return true
}

func sendName(filename string, client *net.TCPConn) bool {
	tmp := []byte(filename)
	_, err := client.Write(tmp)
	if err != nil {
		color.Red("发送文件名失败", err)
		return false
	}
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success"{
		color.Red("对方接收文件名失败")
		return false
	}
	return true
}

func sendSize(size int64, client *net.TCPConn) bool {
	tmp := []byte(strconv.FormatInt(size, 10))
	_, err := client.Write(tmp)
	if err != nil {
		color.Red("发送文件大小失败", err)
		return false
	}
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success"{
		color.Red("对方接收文件大小失败")
		return false
	}
	return true
}

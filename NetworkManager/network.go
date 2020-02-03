package NetworkManager

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"net"
)
// 工具函数
func checkCRC32(data []byte) bool {
	n := len(data)
	if crc(data[:n-1]) != data[n-1]{
		return false
	}
	return true
}

func crc(data []byte) byte {
	res := crc32.ChecksumIEEE(data)
	res = (res >> 16) & (res & 0xffff)
	res = (res >> 8) & (res & 0xff)
	return byte(res)
}

func checkReceiveErr(err error) bool {
	if err != io.EOF && err != nil {
		fmt.Println("接收错误：" + err.Error())
		return false
	}
	return true
}

func checkSendErr(err error) bool {
	if err != io.EOF && err != nil {
		fmt.Println("发送错误：" + err.Error())
		return false
	}
	return true
}

// 功能函数
// 接收文件字节序列
func ReceiveBytes(conn *net.TCPConn) ([]byte, bool) {
	//if n > 64 {
	//	tmp := make([]byte, n - 64)
	//	_, _ = conn.Read(tmp)
	//	for _, i := range tmp{
	//		data = append(data, i)
	//	}
	//}
	data := make([]byte, 64)
	n, err := conn.Read(data)
	if err != nil{
		if err == io.EOF{
			return nil, true
		}
		checkReceiveErr(err)
		return nil, false
	}
	//n := int(data[0])
	return data[:n], true
	//if checkReceiveErr(err){
	//	//return data[1:n+1], true
	//	//采用 CRC32 8位累加和进行验证
	//		//var tmp []byte
	//		//tmp = []byte("no")
	//		//fmt.Println("数据验证错误")
	//		//_, _ =conn.Write(tmp)
	//	var tmp []byte
	//	tmp = []byte("yes")
	//	_, _ =conn.Write(tmp)
	//	return data[1:n+1], true
	//}else{
	//	var tmp []byte
	//	tmp = []byte("no")
	//	_, _ =conn.Write(tmp)
	//	//return nil, false
	//}
	//return nil, false
}
// 发送字节序列
func SendBytes(data []byte, conn *net.TCPConn) bool {
	var buffer bytes.Buffer
	buffer.WriteByte(byte(len(data)))
	if len(data) > 64{
		fmt.Println("文件大小:", len(data))
		fmt.Println(len(data))
	}
	buffer.Write(data)
	sender := make([]byte, 65)
	_,_ = buffer.Read(sender)
	for i := len(data) + 1;i < 65;i++{
		sender[i] = 1
	}

	//res := crc(sender)
	//sender = append(sender, res)

	_, err := conn.Write(sender)
	if !checkSendErr(err){
		return false
	}
	return true
}

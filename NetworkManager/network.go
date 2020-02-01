package NetworkManager

import (
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
	data := make([]byte, 65)
	n, err := conn.Read(data)
	if err == io.EOF{
		return nil, true
	}
	if checkReceiveErr(err){
		// 采用 CRC32 8位累加和进行验证
		if !checkCRC32(data[:n]){
			var tmp []byte
			tmp = []byte("no")
			fmt.Println("数据验证错误")
			_, _ =conn.Write(tmp)
			return nil, false
		}
	}else{
		return nil, false
	}
	return data[:n-1], true
}
// 发送字节序列
func SendBytes(data []byte, conn *net.TCPConn) bool {
	res := crc(data)
	data = append(data, res)
	_, err := conn.Write(data)
	if !checkSendErr(err){
		return false
	}
	return true
}

package Manager

import (
	"github.com/color"
	"io"
	"net"
)

func Sender(conn *net.TCPConn, data chan []byte, isDisplay bool, counter chan int64) bool {
	defer close(counter)

	for tmp := range data{
		_, err := conn.Write(tmp)
		if err != nil {
			color.Red("发送失败", err)
			return false
		}
		if isDisplay {
			counter <- int64(len(tmp))
		}
	}

	return true
}

func Receiver(conn *net.TCPConn, data chan []byte, isDisplay bool, counter chan int64) bool {
	defer close(data)
	defer close(counter)

	for {
		tmp := make([]byte, blockSize)
		n, err := conn.Read(tmp)
		if err != nil && err != io.EOF{
			color.Red("接收失败", err)
			return false
		}else if err == io.EOF{
			return true
		}
		data <- tmp[:n]
		if isDisplay {
			counter <- int64(n)
		}
	}
}

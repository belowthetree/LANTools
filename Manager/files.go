package Manager

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var blockSize = 4096

// 判断文件是否存在
func IsExit(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func FileReader(filename string, data chan []byte) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("打开文件出错", err)
		return false
	}
	defer file.Close()
	defer close(data)

	reader := bufio.NewReader(file)
	for ;;{
		tmp := make([]byte, blockSize)
		n, err := reader.Read(tmp)
		if err != nil && err != io.EOF{
			fmt.Println("文件读取错误", err)
			return false
		}
		if n == 0{
			return true
		}
		data <- tmp
	}
}

func FileWriter(filename string, data chan []byte) bool {
	for ;IsExit(filename);{
		filename += "-副本"
	}
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("文件创建失败", err)
		return false
	}
	defer file.Close()

	for bytes := range data{
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Println("文件写入错误", err)
			return false
		}
	}

	return true
}
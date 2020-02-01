package Manager

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 显示功能
func (gui * GUI) listFunc()  {
	fmt.Println("以下是本程序指令")
	for name := range gui.func_list{
		fmt.Println("	" + name)
	}
}
// 获取当前目录文件列表
func (gui *GUI) getFileList(isDir bool) {
	files, _ := ioutil.ReadDir("./")
	gui.Len = 0
	gui.FileList = make([]string, 0)
	gui.FileLen = make([]int, 0)
	for _, file := range files{
		if file.IsDir() == isDir {
			if gui.Len < len(file.Name()){
				gui.Len = len(file.Name())
			}
			gui.FileList = append(gui.FileList, file.Name())
			gui.FileLen = append(gui.FileLen, len(file.Name()))
		}
	}
	if gui.Len < 10{
		gui.Len = 10
	}
}
// 读取文件并通过管道发送给 TCP 传输
func getFileBytes(file *os.File, bytes chan []byte)  {
	ctx := make([]byte, 64)
	var err error = nil
	st,_ := file.Stat()
	size := st.Size()
	for i := int64(0); i < size;i += 64{
		_, _ = file.Seek(i, 0)
		if size - i < 64{
			ctx = make([]byte, size - i)
		}else{
			ctx = make([]byte, 64)
		}
		_, err = file.Read(ctx)
		if err != nil{
			return
		}
		bytes <- ctx
	}
	close(bytes)
}
// 数据通过管道从这个函数并发写入文件中
func writeFile(file *os.File, data chan []byte) bool {
	for bytes := range data {
		_, err := file.Write(bytes)
		if err != nil {
			fmt.Println("保存文件出错" + err.Error())
			return false
		}
	}
	fmt.Println("文件保存完毕")
	file.Close()
	return true
}
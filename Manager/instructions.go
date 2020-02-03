package Manager

import (
	"fmt"
	"github.com/color"
	"io/ioutil"
	"os"
)

func Cd(path string)  {
	err := os.Chdir(path)
	if err != nil {
		color.Red("无法进入目录", err)
	}
}

func Ls(isDir bool)  {
	names, _ := ioutil.ReadDir("./")
	var files []string
	for _, file := range names{
		if file.IsDir() == isDir {
			files = append(files, file.Name())
		}
	}
	fmt.Println("文件列表：")
	for num, file := range files{
		fmt.Printf("%d. %s\n", num, file)
	}
	fmt.Println()
}

func LsFile() {
	Ls(false)
}

func LsDir() {
	Ls(true)
}

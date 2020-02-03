package main

import (
	"./Manager"
	"fmt"
)

var ch = make(chan []byte, 10)

func main() {
	fmt.Println("开始启动！")
	LocalIps := GetIntranetIp()
	fmt.Print("你的ID们是 ")
	for _, LocalIp := range LocalIps{
		fmt.Println(LocalIp)
	}
	fmt.Println("输入 help 以获取更多帮助")

	gui := Manager.GUI{}
	gui.LocalIP = LocalIps
	gui.Render()
}
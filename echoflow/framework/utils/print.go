package utils

import (
	"fmt"
)

// 在控制台打印出不同颜色的信息
// 0 - 黑色
// 1 - 红色
// 2 - 绿色
// 3 - 黄色
// 4 - 蓝色
// 5 - 紫红色
// 6 - 青蓝色;
func colorPrint(color int, message string) {
	fmt.Printf("\033[0;%dm%s\033[0m", color+30, message)
}

func PrintInfo(message ...interface{}) {
	msgstr := fmt.Sprintln(message...)
	colorPrint(2, msgstr)
}

func PrintError(message ...interface{}) {
	msgstr := fmt.Sprintln(message...)
	colorPrint(1, msgstr)
}
func PrintWarn(message ...interface{}) {
	msgstr := fmt.Sprintln(message...)
	colorPrint(3, msgstr)
}

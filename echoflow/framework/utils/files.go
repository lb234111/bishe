package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// 读文件
// @param fileName 文件名
func ReadFileWithBufio(fileName string) []string {
	ls_result := []string{}

	fileObj, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("open file failed, err:%v", err)
		return ls_result
	}
	// 关闭文件
	defer fileObj.Close()

	// 创建一个用来从文件中读内容的对象
	reader := bufio.NewReader(fileObj)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(line) != 0 {
				// 去掉末尾的\n
				if line[len(line)-1] == '\n' {
					line = line[0 : len(line)-1]
				}
				ls_result = append(ls_result, line)
			}
			return ls_result
		}
		if err != nil {
			fmt.Printf("read ffraom file failed, err:%v\n", err)
			return ls_result
		}
		// 去掉末尾的\n
		if line[len(line)-1] == '\n' {
			line = line[0 : len(line)-1]
		}
		ls_result = append(ls_result, line)
	}
}

// 覆盖写入
// @param fileName 写入的文件名
// @param ls_content 写入的内容
func RewriteFile(fileName string, ls_content []string) {
	fileObj, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("open file failed, err:%v\n", err)
		return
	}
	defer fileObj.Close()

	// 写入
	for _, line := range ls_content {
		fileObj.WriteString(line + "\n")
	}
}

// 判断文件是否存在
// @param filename 文件名
// @return true 文件存在
// @return false 文件不存在
func FileExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在
			return false
		}
		return false
	}
	// 是目录也不行
	if info.IsDir() {
		return false
	}
	return true
}

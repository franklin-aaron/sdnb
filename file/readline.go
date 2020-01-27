package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/**
 * @Description: 逐行读取
 * @File: readline.go
 * @Time: 2020/1/23 10:57
 */
type readline struct {
	FileInfo os.FileInfo   //文件的信息
	File     *os.File      //文件
	Reader   *bufio.Reader //文件缓冲器
}

var instance *readline

/**
 * @Description: 文件操作单例
 * @filePath: 文件路径
 * @File: readline.go
 * @Time: 2020/1/23 10:57
 */
func GetInstance(filePath string) (*readline, error) {
	if instance == nil {

		//读取文件信息
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return instance, err
		}

		//读取文件
		file, err := os.Open(filePath)
		if err != nil {
			return instance, err
		}

		//创建缓冲器
		reader := bufio.NewReader(file)
		instance = &readline{
			FileInfo: fileInfo,
			File:     file,
			Reader:   reader,
		}
	}
	return instance, nil
}

/**
 * @Description: 文件逐行读取
 * @arr: 读取结果存放切片， 否则为nil
 * @data: 文件结果存放管道， 否则为nil
 * @Time: 2020/1/23 10:57
 */
func (r *readline) ReadLine(arr [] string, data chan string) error {
	if arr == nil && data == nil {
		return fmt.Errorf("parameter error : You can't have two arguments that are nil")
	}

	i := 0
	for {
		//逐行读取
		lineContent, err := r.Reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		//存放到结果
		if arr != nil {
			arr[i] = strings.Replace(strings.Replace(string(lineContent), "\n", "", -1), "\r", "", -1)
			i++
		} else {
			data <- strings.Replace(strings.Replace(string(lineContent), "\n", "", -1), "\r", "", -1)
		}

		//如果读取完毕
		if err == io.EOF {
			//如果是管道方式读取，读取完毕关闭管道
			if data != nil {
				close(data)
			}
			break
		}
	}

	r.File.Close()

	return nil
}

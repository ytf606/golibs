package stringx

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func Read(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Write(content, file string) error {
	f, err := os.Create(file) //创建文件
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f) //创建新的 Writer 对象
	if _, err := writer.WriteString(content); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func ReadSlice(file string) ([]string, error) {
	// 创建句柄
	fi, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	//func NewReader(rd io.Reader) *Reader {}，返回的是bufio.Reader结构体
	r := bufio.NewReader(fi) // 创建 Reader
	lines := make([]string, 0)
	for {
		//func (b *Reader) ReadBytes(delim byte) ([]byte, error) {}
		lineBytes, err := r.ReadBytes('\n')
		//去掉字符串首尾空白字符，返回字符串
		line := strings.TrimSpace(string(lineBytes))
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		lines = append(lines, line)
	}
	return lines, nil
}

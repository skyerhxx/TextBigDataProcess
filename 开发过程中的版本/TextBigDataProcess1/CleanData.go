package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	//基于磁盘和缓存的读取
	file,_ := os.Open("d:/golang/src/go_code/TextBigDataProcess/text/kaifangX.txt")
	defer file.Close()

	//准备一个优质文件
	goodFile,_ := os.OpenFile("d:/golang/src/go_code/TextBigDataProcess/text/kaifang_good.txt", os.O_WRONLY | os.O_CREATE|os.O_APPEND, 0644)
	defer goodFile.Close()

	//准备一个劣质文件
	badFile,_ := os.OpenFile("d:/golang/src/go_code/TextBigDataProcess/text/kaifang_bad.txt", os.O_WRONLY | os.O_CREATE|os.O_APPEND, 0644)
	defer badFile.Close()

	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		lineStr := string(lineBytes)

		fields := strings.Split(lineStr, "，")
		if len(fields) >1 && len(fields[1])==18{  //防止有空白行
		                                         //身份证号18位，这里偷懒了，仔细的话应该用正则判断
			//摘取到优质文件中
			goodFile.WriteString(lineStr+"\n")
			fmt.Println("Good: ",lineStr)
		}else{
			badFile.WriteString(lineStr+"\n")
			fmt.Println("Bad: ",lineStr)
		}
	}

}
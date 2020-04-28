package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//一次性将全部数据读入内存
func main1() {

	contentBytes,err := ioutil.ReadFile("d:/golang/src/go_code/TextBigDataProcess/kaifangX.txt")
	if err != nil {  //如果有错
		fmt.Println("读入失败",err)
	}
	contentStr := string(contentBytes)

	//逐行打印
	lineStrs := strings.Split(contentStr,"\n\r")
	for _,lineStr := range lineStrs {
		fmt.Println(lineStr)
	}

}

//基于磁盘和缓存的读取
func main() {

	file,_ := os.Open("d:/golang/src/go_code/TextBigDataProcess/text/kaifangX.txt")
	//file,_ := os.Open("d:/golang/src/go_code/TextBigDataProcess/rumors.txt")

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
	fmt.Println(string(lineBytes))
	}
}
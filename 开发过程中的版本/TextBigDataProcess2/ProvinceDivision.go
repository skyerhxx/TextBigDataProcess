package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Province struct {
	Id string	//身份证号前两位
	Name string
	File *os.File  //黑龙江.txt
	               //指针
	chanData chan string //本省文件的数据管道
}

var wg sync.WaitGroup
func writeFile(province *Province){

	//死循环读取管道，管道关闭时循环结束
	for lineStr := range province.chanData{
		province.File.WriteString(lineStr)
		fmt.Print(province.Name,"写入",lineStr)
	}
	//标记协程结束
	wg.Done()
}

func main(){
	//创建34个省份的实例
	pMap := make(map[string]*Province)
	ps := []string{"北京市11","天津市12","河北省13","山西省14","内蒙古自治区15","辽宁省21","吉林省22","黑龙江省23","上海市31","江苏省32","浙江省33","安徽省34","福建省35","江西省36","山东省37","河南省41","湖北省42","湖南省43","广东省44","广西壮族自治区45","海南省46","重庆市50","四川省51","贵州省52","云南省53","西藏自治区54","陕西省61","甘肃省62","青海省63","宁夏回族自治区64","新疆维吾尔自治区65","台湾省71","香港特别行政区81","澳门特别行政区91"}
	for _,p := range ps{
		name := p[:len(p)-2]
		id := p[len(p)-2:]
		province := Province{Id: id, Name: name}
		pMap[id] = &province

		//为每个省份打开一个文件
		file, _ := os.OpenFile("d:/golang/src/go_code/TextBigDataProcess/text/"+province.Name+".txt", os.O_WRONLY|os.O_CREATE | os.O_TRUNC,0644)
		province.File = file
		defer file.Close()

		//创建每个省的数据管道
		province.chanData = make(chan string)
		fmt.Println("管道已经创建")
	}

	//每个省份各起一条协程
	for _,province := range pMap{
		wg.Add(1)
		go writeFile(province)
	}

	//读入优质数据
	file, _ := os.Open("d:/golang/src/go_code/TextBigDataProcess/text/kaifang_good.txt")
	defer file.Close()
	reader := bufio.NewReader(file)

	//逐行判断身份证号的前两位
	for{
		lineBytes, _, err := reader.ReadLine()
		//读取完毕时，关闭所有数据管道
		if err == io.EOF {

			for _,province := range pMap{
				close(province.chanData)
				fmt.Println("管道已关闭")
			}
			break
		}
		//拿出省份ID
		lineStr := string(lineBytes)
		fieldsSlice := strings.Split(lineStr, "，")
		id := fieldsSlice[1][0:2]

		//对号入座，写入响应的管道
		if province,ok := pMap[id];ok{
			province.chanData <- (lineStr + "\n")
		}else{
			fmt.Println("莫名其妙的省",id)
		}
	}

	//阻塞等待协程结束
	wg.Wait()
}


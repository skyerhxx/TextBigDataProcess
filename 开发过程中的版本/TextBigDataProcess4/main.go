package main

import (
	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"strings"
	"time"
)

const CACHE_LEN = 2

var(
	kfMap map[string]QueryResult
	//chanSema chan int  //信号量管道
	chanData chan *KfPerson
	db *sqlx.DB
)


//错误处理函数
func HandleError(err error,why string){
	if err != nil {
		fmt.Println("ERROR OCCURED!!!",err,why)
	}
}


//将文本大数据入库
//入库成功后，做一个文件标记，下一次见到标记就不再执行入库操作
func init() {
	//如果数据库已经初始化过了，就直接退出
	_, err1 := os.Stat("d:/golang/src/go_code/TextBigDataProcess/text/kaifanggood_dbok.mark")
	if err1 == nil {
		fmt.Println("数据库业已初始化")
		return
	}


	//打开数据库
	var err error
	db,err = sqlx.Open("mysql","root:root@tcp(127.0.0.1:3306)/kaifang")
	HandleError(err,"sqlx.Open")
	defer db.Close()
	fmt.Println("数据库已打开")

	//必要时建表
	_,err = db.Exec("create table if not exists kfperson(id int primary key auto_increment,name varchar(20),idcard char(18),sex char(1));")
	HandleError(err,"db.Exec create table")
	fmt.Println("数据表已创建")


	//初始化信号量管道(控制并发数)
	//chanSema = make(chan int,100)
	chanData = make(chan *KfPerson, 10000000)  //相当于我现在管道的缓存能力是无限的
	fmt.Println("管道已初始化")

	//开辟协程，源源不断地从数据管道获取信息，插入数据库
	for i:= 0;i<100;i++{
		go insertKfPerson()
		fmt.Printf("协程%d已开辟\n",i)
	}
	//打开大数据文件
	file,e := os.Open("d:/golang/src/go_code/TextBigDataProcess/text/kaifang_good.txt")
	HandleError(e,"os.Open")
	defer file.Close()
	reader := bufio.NewReader(file)
	fmt.Println("大数据文本已打开")

	//分批次读入大数据文本
	//还是要基于缓存的读取
	for{
		lineBytes, _, err := reader.ReadLine()
		//如果读到了文件尾
		if err == io.EOF {
			//关闭数据管道
			close(chanData)
			break
		}

		HandleError(err,"reader.ReadLine")

		//逐条入库（并发）
		lineStr := string(lineBytes)
		fields := strings.Split(lineStr, "，")
		name,idcard := fields[0],fields[1]
		//抛弃过长的名字
		name = strings.TrimSpace(name)
		if len(strings.Split(name,"")) > 20{
			fmt.Println("%s 名字过长，已经抛弃了")
			continue
		}


		//方案一：开2000万协程，行不通，耗尽了资源，程序崩溃
		//go insertKfPerson(db,&kfPerson)

		//方案二：开有限条协程，从管道中读取数据
		kfPerson := KfPerson{Name:name, Idcard: idcard}
		chanData <- &kfPerson //主协程管着往里面送数据
		//fmt.Println(&kfPerson)
		//fmt.Println(<-chanData)

	}

	fmt.Println("数据初始化成功！")

	//创建一个标记文件，标记数据库已经初始化成功
	_, err = os.Create("d:/golang/src/go_code/TextBigDataProcess/text/kaifanggood_dbok.mark")
	if err == nil {
		fmt.Println("初始化标记已经创建!")
	}
}

//源源不断地从管道里拿数据，然后源源不断地往数据库里送
func insertKfPerson(){
	for kfPerson := range chanData{
		result,err := db.Exec("insert into kfperson(name,idcard) values(?,?);",kfPerson.Name,kfPerson.Idcard)
		HandleError(err,"db.Exec insert")
		//fmt.Println(result)
		if n,e := result.RowsAffected(); e == nil && n > 0{
			fmt.Printf("插入 %s 成功!\n",kfPerson.Name)
		}
	}
}


func main() {

	//打开数据库
	db,err := sqlx.Open("mysql","root:root@tcp(127.0.0.1:3306)/kaifang")
	HandleError(err,"sqlx.Open")
	defer db.Close()

	//初始化缓存
	kfMap = make(map[string]QueryResult,0)

	var name string
	//循环接收用户想要查询的姓名
	//循环查询
	for{
		fmt.Print("请输入要查询的开房者姓名: ")
		fmt.Scanf("%s",&name)

		//如果用户想退出
		if name == "exit"{
			break
		}

		//查看所有缓存
		if name == "cache"{
			fmt.Println("共缓存了%d条结果: \n",len(kfMap))
			for key := range kfMap{
				fmt.Println(key)
			}
		}

		//先查看内存中是否有结果
		if qr,ok := kfMap[name];ok{
			fmt.Println("查询到%d条结果", len(qr.value))
			fmt.Println(qr.value)
			continue
		}

		//内存中没有，查数据库
		kfpeople := make([]KfPerson,0)
		e := db.Select(&kfpeople,"select id,name,idcard from kfperson where name like ?;", name)//在这里name是string，带入select语句的时候会自动加上引号，否则select语句 like '西门庆' 是应该加引号的
		HandleError(e,"db.Select")
		fmt.Println("查询到%d条结果", len(kfpeople))
		fmt.Println(kfpeople)

		//查到的结果丢入内存
		queryResult := QueryResult{value:kfpeople}
		queryResult.cacheTime = time.Now().UnixNano()
		queryResult.count = 1
		kfMap[name] = queryResult

		//有必要时淘汰一些缓存

		if len(kfMap) > CACHE_LEN{
			delKey := UpdateCache(&kfMap)
			fmt.Printf("%s已经被淘汰出缓存!\n",delKey)
		}
	}

	fmt.Println("ALL OVER!")

}


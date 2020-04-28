# TextBigDataProcess
文本大数据处理项目(Go语言)

***

### 开发环境

* 版本：Go SDK 1.13.5

* IDE：GoLand/Vscode



### 文件说明&运行方式

>##### ==*<u>main.go, cache.go, model.go需要一起运行</u>*==
>
>*  **main.go是主要程序**
>
>* model.go用于定义一些数据结构
>
>* cache.go是将用于缓存处理相关的函数和接口单独拿了出来



> ReadData.go、CleanData.go、ProvinceDivision.go、AgeDivision.go都是前期处理阶段的程序，分别用于
>
> * ReadData.go——文本大数据读取
>
> * CleanData.go——数据清洗
>
> * ProvinceDivision.go——省份划分
>
> * AgeDivision.go——年龄划分

详见博客



### 实现的功能
* 使用缓冲区读取大数据文本
* 多协程并发将数据写入文件
* 单协程将数据写入mysql数据库+多协程并发写入mysql数据库
* 二级缓存查询数据（精确查询和模糊查询）
* 内存缓存+动态清理



### 项目具体过程：

- https://blog.csdn.net/hxxjxw/article/details/105618767



### 关于kfX.txt文件的说明

​      数据是我自己通过OCR从视频中识别了一些自制了一个70K左右的文件，数据格式什么都是一样的，用于学习这个项目。里面的一些信息也做过修改，不是真实的，不是真实的，不是真实的，仅用于学习所用。如果引发什么问题请联系我删掉。

​     text里的kf_good.txt是因为我复制了好几遍，所以比较大
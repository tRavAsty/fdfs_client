# Fdfs_client go
The Golang interface to the Fastdfs Ver 5.08.
## Notice:Only realized the upload,download, delete functions
作者1年前就没有commit了，尝试修复一下bug并且新增一点功能

第一次commit：


	1、增加common.go和common_test.go文件，用于松耦合读取配置文件的功能


	2、connection_test.go client_test.go改成读取配置文件的ip地址，而不是hard code到测试文件中


	3、设置连接池最大，最小连接数量可变

第二次commit


	1、修复客户端泄露的可能错误

	2、增加日志记录，使得更加方便debug

第三次commit


	1、修复for循环没有break的bug
	
	2、优化日志
	
第四次commit


	1、完善append by file name操作的truncate， append和modify功能
	
	2、测试完整append by file name功能并通过
	 go test github.com/tRavAsty/fdfs_client -test.run "TestUploadAppenderByFilename" -v

## Installation

	$ go get github.com/tRavAsty/fdfs_client
	
## Getting Started

	go test github.com/tRavAsty/fdfs_client -v

 运行某一测试函数
 
 	go test github.com/tRavAsty/fdfs_client -test.run "Test*" -v

 或者看看 client_test.go

# Author
 我是要做毕设的大四狗
 
 有什么问题请联系我 lchuyi@mail.ustc.edu.cn

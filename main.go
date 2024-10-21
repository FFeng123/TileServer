package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	SourceDir string `short:"s" long:"source" description:"用于构建瓦片地图服务数据的源数据目录" required:"false"`
	DbFile    string `short:"d" long:"db" description:"瓦片地图服务数据存储的数据库文件" required:"true"`
	Port      int    `short:"p" long:"port" description:"服务端口号" required:"false"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	if opts.SourceDir != "" {
		fmt.Println("开始构建瓦片地图服务数据...")
		err := makeDb(opts.SourceDir, opts.DbFile)
		if err != nil {
			fmt.Println("构建瓦片地图服务数据失败：", err)
			return
		}
		fmt.Println("构建瓦片地图服务数据成功！")
	} else {
		fmt.Println("启动瓦片地图服务...")
		err := serve(opts.Port, opts.DbFile)
		if err != nil {
			fmt.Println("启动瓦片地图服务失败：", err)
			return
		}
	}
}

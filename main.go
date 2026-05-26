package main

import (
	"log"

	"rename4scanner/util"
	"rename4scanner/core"

	"github.com/spf13/cobra"
)

func main() {
	// root：待重命名的文件夹路径
	// startAt：起始序号偏移量
	var root string
	var startAt int
	util.SetLog("r4s.log")
	// 使用 cobra 构建命令行工具
	var rootCmd = &cobra.Command{
		Use:   "r4s",
		Short: "重命名扫描仪文件，使其从指定序号开始",
		Long:  `rename4scanner 是一个用于处理扫描仪分批扫描时文件序号连续性的工具`,
		Run: func(cmd *cobra.Command, args []string) {
			// 调用核心逻辑，并在发生错误时退出程序
			err := core.Rename4Scanner(startAt, root)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	/*
		cli 使用说明：
		go build -o r4s main.go
		./r4s -d <文件夹路径> -s <起始序号>
		或
		./r4s --dir <文件夹路径> --start-at <起始序号>
	*/
	rootCmd.Flags().StringVarP(&root, "dir", "d", "", "待重命名的文件夹路径")
	rootCmd.Flags().IntVarP(&startAt, "start-at", "s", 0, "起始序号偏移量（例如：如果第一批最后一个是Scan_0012.jpg，则设置为12）")

	// 标记参数为必填
	if err := rootCmd.MarkFlagRequired("dir"); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.MarkFlagRequired("start-at"); err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

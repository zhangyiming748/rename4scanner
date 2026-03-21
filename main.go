package main

import (
	"log"

	"rename4scanner/core"

	"github.com/spf13/cobra"
)

func main() {
	// failedPath：待重命名的失败文件夹路径
	// successFile：上次成功的文件路径，作为下一次序号的起点
	var failedPath, successFile string

	// 使用 cobra 构建命令行工具
	var rootCmd = &cobra.Command{
		Use:   "rename4scanner",
		Short: "重命名扫描失败后的文件，使其从上次成功的文件序号继续",
		Long:  `rename4scanner 是一个用于处理扫描仪中断后文件序号连续性的工具`,
		Run: func(cmd *cobra.Command, args []string) {
			// 调用核心逻辑，并在发生错误时退出程序
			err := core.Rename4Scanner(failedPath, successFile)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	/*
		cli 使用说明：
		go build -o rename4scanner main.go
		./rename4scanner -f <失败文件夹路径> -s <上一个成功文件路径>
		或
		./rename4scanner --failed <失败文件夹路径> --success <上一个成功文件路径>
	*/
	rootCmd.Flags().StringVarP(&failedPath, "failed", "f", "", "失败文件夹路径")
	rootCmd.Flags().StringVarP(&successFile, "success", "s", "", "上一个成功文件路径")

	// 标记参数为必填
	if err := rootCmd.MarkFlagRequired("failed"); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.MarkFlagRequired("success"); err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

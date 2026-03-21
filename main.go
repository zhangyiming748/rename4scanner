package main

import (
	"log"

	"github.com/spf13/cobra"
	"rename4scanner/core"
)

func main() {
	var failedPath, successFile string

	var rootCmd = &cobra.Command{
		Use:   "rename4scanner",
		Short: "重命名扫描失败后的文件，使其从上次成功的文件序号继续",
		Long:  `rename4scanner 是一个用于处理扫描仪中断后文件序号连续性的工具`,
		Run: func(cmd *cobra.Command, args []string) {
			err := core.Rename4Scanner(failedPath, successFile)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	/*
		这里使用 cobra 方式实现通过cli使用这个程序
		用户通过go builg -o rename4scanner main.go
		然后在命令行中运行：
		./rename4scanner -f <失败文件夹路径> -s <上一个成功文件路径>
		或
		./rename4scanner --failed <失败文件夹路径> --success <上一个成功文件路径>
		*/
	rootCmd.Flags().StringVarP(&failedPath, "failed", "f", "", "失败文件夹路径")
	rootCmd.Flags().StringVarP(&successFile, "success", "s", "", "上一个成功文件路径")

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
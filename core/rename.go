package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// Rename4Scanner 重命名扫描失败后的文件，使其从上次成功的文件序号继续
func Rename4Scanner(failedPath, successFile string) error {
	// 检查失败文件夹是否存在
	if _, err := os.Stat(failedPath); os.IsNotExist(err) {
		return fmt.Errorf("错误：文件夹不存在 - %s", failedPath)
	}

	// 检查最后一个成功文件是否存在
	if _, err := os.Stat(successFile); os.IsNotExist(err) {
		return fmt.Errorf("错误：文件不存在 - %s", successFile)
	}

	// 从成功文件名中提取基准数字
	successFileName := filepath.Base(successFile)

	// 使用正则表达式提取数字部分（假设格式为 Scan_0014.jpg）
	re := regexp.MustCompile(`_([0-9]+)\.([a-zA-Z]+)$`)
	matches := re.FindStringSubmatch(successFileName)

	if len(matches) < 2 {
		return fmt.Errorf("错误：无法从文件名中提取数字 - %s", successFileName)
	}

	baseNum, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("错误：无法解析数字 - %s", matches[1])
	}

	fmt.Printf("基准数字：%d\n", baseNum)
	fmt.Printf("开始处理文件夹：%s\n", failedPath)

	// 获取失败文件夹中的所有文件
	files, err := os.ReadDir(failedPath)
	if err != nil {
		return fmt.Errorf("错误：无法读取文件夹 - %s", failedPath)
	}

	// 过滤出与成功文件具有相同扩展名的文件
	ext := filepath.Ext(successFileName)
	var targetFiles []os.FileInfo

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ext {
			// 排除基准文件本身，只处理失败批次的文件
			if file.Name() != successFileName {
				info, err := file.Info()
				if err != nil {
					return fmt.Errorf("错误：无法获取文件信息 - %s: %v", file.Name(), err)
				}
				targetFiles = append(targetFiles, info)
			}
		}
	}

	if len(targetFiles) == 0 {
		return fmt.Errorf("错误：文件夹中没有扩展名为 %s 的文件（排除基准文件）", ext)
	}

	// 按名称排序以确保重命名的一致性
	sort.Slice(targetFiles, func(i, j int) bool {
		return targetFiles[i].Name() < targetFiles[j].Name()
	})

	// 计数器从 base + 1 开始
	counter := baseNum + 1

	// 遍历文件并重命名
	renamedCount := 0
	for _, file := range targetFiles {
		oldPath := filepath.Join(failedPath, file.Name())

		// 生成新的文件名
		newFilename := fmt.Sprintf("Scan_%04d%s", counter, ext)
		newPath := filepath.Join(failedPath, newFilename)

		fmt.Printf("重命名：%s -> %s\n", file.Name(), newFilename)

		err := os.Rename(oldPath, newPath)
		if err != nil {
			return fmt.Errorf("错误：重命名失败 %s 到 %s: %v", oldPath, newPath, err)
		}

		counter++
		renamedCount++
	}

	fmt.Printf("处理完成！总共重命名了 %d 个文件\n", renamedCount)
	fmt.Printf("起始编号：%d\n", baseNum+1)
	fmt.Printf("结束编号：%d\n", counter-1)

	return nil
}

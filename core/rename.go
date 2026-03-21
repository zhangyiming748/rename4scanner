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
//
// failedPath: 失败批次所在目录，该目录下可能有多个需要重新编号的文件。
// successFile: 上一个成功扫描文件的全路径（用于提取当前已经使用的最大序号）。
//
// 该函数的行为：
// 1. 验证输入路径是否存在
// 2. 从 successFile 提取基准序号
// 3. 遍历 failedPath 中与成功文件同扩展名的文件，排除基准文件自身
// 4. 统一排序后按 Scan_XXXX.ext 规则进行连续命名（从基准序号+1开始）
// 5. 逐个重命名，遇到错误立即返回。
func Rename4Scanner(failedPath, successFile string) error {
	// 检查失败文件夹是否存在
	if _, err := os.Stat(failedPath); os.IsNotExist(err) {
		return fmt.Errorf("错误：文件夹不存在 - %s", failedPath)
	}

	// 检查最后一个成功文件是否存在
	if _, err := os.Stat(successFile); os.IsNotExist(err) {
		return fmt.Errorf("错误：文件不存在 - %s", successFile)
	}

	// 提取 successFile 的基本文件名（去除路径）
	successFileName := filepath.Base(successFile)

	// 使用正则提取文件名中的数字部分，期望匹配类似"Scan_0014.jpg"的格式
	// matches[1] 为数字部分，matches[2] 为扩展名后缀（不包含点）
	re := regexp.MustCompile(`_([0-9]+)\.([a-zA-Z]+)$`)
	matches := re.FindStringSubmatch(successFileName)
	if len(matches) < 2 {
		return fmt.Errorf("错误：无法从文件名中提取数字 - %s", successFileName)
	}

	// 解析数字基准序号
	baseNum, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("错误：无法解析数字 - %s", matches[1])
	}

	fmt.Printf("基准数字：%d\n", baseNum)
	fmt.Printf("开始处理文件夹：%s\n", failedPath)

	// 读取失败目录下的所有项
	files, err := os.ReadDir(failedPath)
	if err != nil {
		return fmt.Errorf("错误：无法读取文件夹 - %s", failedPath)
	}

	// 只处理与基准文件同扩展名的文件
	ext := filepath.Ext(successFileName)
	var targetFiles []os.FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue // 跳过子目录
		}

		if filepath.Ext(file.Name()) != ext {
			continue // 相异扩展名跳过
		}

		// 排除成功文件本身，防止重复处理
		if file.Name() == successFileName {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return fmt.Errorf("错误：无法获取文件信息 - %s: %v", file.Name(), err)
		}
		targetFiles = append(targetFiles, info)
	}

	if len(targetFiles) == 0 {
		return fmt.Errorf("错误：文件夹中没有扩展名为 %s 的文件（排除基准文件）", ext)
	}

	// 按文件名排序以保证重命名顺序可预测、稳定
	sort.Slice(targetFiles, func(i, j int) bool {
		return targetFiles[i].Name() < targetFiles[j].Name()
	})

	// 从基准号 + 1 开始重命名
	counter := baseNum + 1
	renamedCount := 0

	for _, file := range targetFiles {
		oldPath := filepath.Join(failedPath, file.Name())
		newFilename := fmt.Sprintf("Scan_%04d%s", counter, ext)
		newPath := filepath.Join(failedPath, newFilename)

		fmt.Printf("重命名：%s -> %s\n", file.Name(), newFilename)

		if err := os.Rename(oldPath, newPath); err != nil {
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

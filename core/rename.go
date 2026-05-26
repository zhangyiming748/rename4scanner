package core

import (
	"fmt"
	"github.com/zhangyiming748/finder"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

/*
现在的逻辑是
我的扫描仪扫描出来的文件会自动重命名图片为Scan_0001.jpg到Scan_9999.jpg(实际上最多300张 扫描仪不能支持短时大批量扫描)
现在第二批扫描的图片想合并到第一批扫描图片的文件夹里就会出现同名文件
所以我想用startAt这个变量表示第二批文件夹应该被重命名的起始点
重命名之后就应该是Scan_0001+{startAT}
比如第一批文件夹里
最后一个文件是Scan_0012.jpg
startAt=12
那么第二批文件夹里的文件应该重命名为Scan_0001+{12}.jpg 即 Scan_0013.jpg

简化一下以上逻辑
现在不谈第一个文件夹
就第二个文件夹
假设Scan_0001.jpg到Scan_0999.jpg
这些jpg图片
从头遍历这些文件名
*/
func Rename4Scanner(startAT int, root string) error {
	// 读取目录中的所有文件

	// 过滤出符合 Scan_XXXX.ext 格式的图片文件

	images := finder.FindAllImages(root)
	// 按文件名排序，确保重命名顺序一致
	sort.Strings(images)

	log.Printf("找到 %d 个图片文件", len(images))
	for i, image := range images {
		fmt.Println(i, image)
		//在这里将每一个文件的绝对路径处理一下
		//将文件名中的数字部分加startAT
		//然后重新组成为绝对路径形式的新文件名newNmae
		//log打印旧文件名和新文件名

		// 获取文件名（不含路径）
		fileName := filepath.Base(image)

		// 从文件名中提取数字部分
		re := regexp.MustCompile(`Scan_(\d+)`)
		matches := re.FindStringSubmatch(fileName)
		if len(matches) < 2 {
			log.Printf("跳过不符合格式的文件: %s", fileName)
			continue
		}

		// 解析原有数字
		oldNum, err := strconv.Atoi(matches[1])
		if err != nil {
			return fmt.Errorf("解析文件序号失败 %s: %w", fileName, err)
		}

		// 计算新序号
		newNum := oldNum + startAT

		// 生成新文件名，保持原有格式（4位数字）
		newFileName := fmt.Sprintf("Scan_%04d%s", newNum, filepath.Ext(fileName))

		// 组装新文件的绝对路径
		newName := filepath.Join(filepath.Dir(image), newFileName)

		// 打印旧文件名和新文件名
		log.Printf("重命名: %s -> %s", image, newName)

		// 执行重命名
		if err := os.Rename(image, newName); err != nil {
			return fmt.Errorf("重命名失败 %s -> %s: %w", image, newName, err)
		}
	}

	return nil
}

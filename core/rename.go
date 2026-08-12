package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/zhangyiming748/finder"
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

	re := regexp.MustCompile(`Scan_(\d+)`)

	// 第一遍遍历：只计算每个文件的目标名，不执行任何重命名
	type renamePlan struct {
		src string
		dst string
	}
	var plans []renamePlan
	for _, image := range images {
		// 获取文件名（不含路径）
		fileName := filepath.Base(image)

		// 从文件名中提取数字部分
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

		plans = append(plans, renamePlan{src: image, dst: newName})
	}

	// 冲突预检：发现任何文件会被覆盖就发出警告并终止，不执行重命名
	sourceSet := make(map[string]bool, len(plans))
	for _, p := range plans {
		sourceSet[p.src] = true
	}
	conflicts := 0
	firstSrc := make(map[string]string, len(plans))
	for _, p := range plans {
		// 多个文件被重命名到同一个目标名，后执行的会覆盖先执行的
		if src, ok := firstSrc[p.dst]; ok {
			log.Printf("警告: %s 和 %s 都会被重命名为 %s，其中一个文件会被覆盖", src, p.src, p.dst)
			conflicts++
			continue
		}
		firstSrc[p.dst] = p.src
		// 目标名已被一个不参与重命名的文件占用，执行会将其覆盖
		if !sourceSet[p.dst] {
			if _, err := os.Stat(p.dst); err == nil {
				log.Printf("警告: 目标文件 %s 已存在且不在本次重命名范围内，%s 会覆盖它", p.dst, p.src)
				conflicts++
			}
		}
	}
	if conflicts > 0 {
		return fmt.Errorf("预检发现 %d 处重命名冲突，已终止操作以防止文件被覆盖，请检查 --start-at 参数是否正确", conflicts)
	}

	// 从最大序号开始倒序重命名，避免文件被重复重命名
	for i := len(plans) - 1; i >= 0; i-- {
		fmt.Println(i, plans[i].src)

		// 打印旧文件名和新文件名
		log.Printf("重命名: %s -> %s", plans[i].src, plans[i].dst)

		// 执行重命名
		if err := os.Rename(plans[i].src, plans[i].dst); err != nil {
			return fmt.Errorf("重命名失败 %s -> %s: %w", plans[i].src, plans[i].dst, err)
		}
	}

	return nil
}

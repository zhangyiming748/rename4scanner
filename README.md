# rename4scanner

解决扫描仪在扫描过程中因卡纸等错误中断后，重新开始扫描时会从 `Scan_0001.jpg` 重新计数的问题，导致文件名重复，无法按顺序连续管理扫描文件。

## 功能

- 根据上一次最后一个成功文件的序号（如 `Scan_0014.jpg`），提取基准数字（0014）
- 对新一批次的扫描文件（如 `Scan_0001.jpg`, `Scan_0002.jpg`...）进行重命名，起始序号为 base + 1（即 0015 开始）
- 批量处理目标文件夹中的所有扫描图像文件，确保命名连续

## 使用方法

### Shell脚本方式

```bash
./rename.sh /path/to/failed/folder /path/to/last_success/Scan_0014.jpg
```

### Go库方式

```go
import "rename4scanner/core"

err := core.Rename4Scanner("/path/to/failed/folder", "/path/to/last_success/Scan_0014.jpg")
if err != nil {
    log.Fatal(err)
}
```

## 安装

```bash
go mod tidy
```

## 编译和运行

```bash
go run main.go /path/to/failed/folder /path/to/last_success/Scan_0014.jpg
```

## 测试

```bash
go test ./core/
```

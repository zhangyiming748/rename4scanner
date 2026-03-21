#!/bin/bash

# 检查参数数量
if [ $# -ne 2 ]; then
    echo "用法：$0 <失败文件夹路径> <上一个成功文件路径>"
    echo "示例：$0 /path/to/failed/folder /path/to/last_success/Scan_0014.jpg"
    exit 1
fi

FAILED_FOLDER="$1"
LAST_SUCCESS_FILE="$2"

# 检查文件夹是否存在
if [ ! -d "$FAILED_FOLDER" ]; then
    echo "错误：文件夹不存在 - $FAILED_FOLDER"
    exit 1
fi

# 检查最后一个成功文件是否存在
if [ ! -f "$LAST_SUCCESS_FILE" ]; then
    echo "错误：文件不存在 - $LAST_SUCCESS_FILE"
    exit 1
fi

# 从文件名中提取基准数字
BASENAME=$(basename "$LAST_SUCCESS_FILE")
# 使用正则表达式提取数字部分（假设格式为 Scan_0014.jpg）
if [[ $BASENAME =~ _([0-9]+)\.[a-zA-Z]+$ ]]; then
    BASE_NUM=${BASH_REMATCH[1]}
    # 去除前导零，转换为十进制数
    BASE_NUM=$((10#$BASE_NUM))
else
    echo "错误：无法从文件名中提取数字 - $BASENAME"
    echo "期望格式：Scan_0014.jpg"
    exit 1
fi

echo "基准数字：$BASE_NUM"
echo "开始处理文件夹：$FAILED_FOLDER"

# 获取文件扩展名（假设所有文件扩展名相同）
EXTENSION=""
for file in "$FAILED_FOLDER"/*; do
    if [ -f "$file" ]; then
        EXTENSION=$(echo "$file" | grep -oE '\.[a-zA-Z]+$')
        break
    fi
done

if [ -z "$EXTENSION" ]; then
    echo "错误：文件夹中没有文件"
    exit 1
fi

echo "文件扩展名：$EXTENSION"

# 计数器从 base + 1 开始
COUNTER=$((BASE_NUM + 1))

# 遍历文件夹中的所有文件并重命名
for file in $(ls -v "$FAILED_FOLDER"/*$EXTENSION 2>/dev/null); do
    if [ -f "$file" ]; then
        # 生成新的文件名
        DIR=$(dirname "$file")
        NEW_FILENAME=$(printf "Scan_%04d%s" $COUNTER "$EXTENSION")
        NEW_PATH="$DIR/$NEW_FILENAME"
        
        echo "重命名：$(basename "$file") -> $NEW_FILENAME"
        mv "$file" "$NEW_PATH"
        
        COUNTER=$((COUNTER + 1))
    fi
done

echo "处理完成！总共重命名了 $((COUNTER - BASE_NUM - 1)) 个文件"
echo "起始编号：$((BASE_NUM + 1))"
echo "结束编号：$((COUNTER - 1))"

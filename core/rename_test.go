package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRename4Scanner(t *testing.T) {
	// 创建临时目录用于测试环境，避免对真实数据产生副作用
	tempDir, err := os.MkdirTemp("", "rename4scanner_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 在临时目录中创建一些测试文件（模拟第二批扫描的文件）
	testFiles := []string{
		"Scan_0001.jpg",
		"Scan_0002.jpg",
		"Scan_0003.jpg",
	}

	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		err := os.WriteFile(filePath, []byte("dummy content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// 从成功文件中提取起始序号 (14)
	// 假设第一批最后一个文件是 Scan_0014.jpg
	startAT := 14

	// 调用函数并确认不返回错误
	err = Rename4Scanner(startAT, tempDir)
	if err != nil {
		t.Error(err)
	}

	// 验证重命名结果：旧文件按顺序重命名为基准号后续编号
	expectedFiles := []string{
		"Scan_0015.jpg", // 原来的 Scan_0001.jpg (1 + 14)
		"Scan_0016.jpg", // 原来的 Scan_0002.jpg (2 + 14)
		"Scan_0017.jpg", // 原来的 Scan_0003.jpg (3 + 14)
	}

	for _, expectedFile := range expectedFiles {
		expectedPath := filepath.Join(tempDir, expectedFile)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected file does not exist: %s", expectedPath)
		} else {
			t.Logf("File exists as expected: %s", expectedPath)
		}
	}

	// 验证旧文件名不再存在，确保文件已成功重命名
	oldFiles := []string{
		"Scan_0001.jpg",
		"Scan_0002.jpg",
		"Scan_0003.jpg",
	}

	for _, oldFile := range oldFiles {
		oldPath := filepath.Join(tempDir, oldFile)
		if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
			t.Errorf("Old file still exists: %s", oldPath)
		} else {
			t.Logf("Old file correctly removed: %s", oldPath)
		}
	}
}

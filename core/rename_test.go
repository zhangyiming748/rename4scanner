package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestRename4Scanner(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := ioutil.TempDir("", "rename4scanner_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 在临时目录中创建一些测试文件（模拟失败批次的文件）
	testFiles := []string{
		"Scan_0001.jpg",
		"Scan_0002.jpg",
		"Scan_0003.jpg",
	}

	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		err := ioutil.WriteFile(filePath, []byte("dummy content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// 创建一个"成功"文件作为基准（在另一个目录中，模拟实际场景）
	successFilePath := filepath.Join(tempDir, "Scan_0014.jpg")
	err = ioutil.WriteFile(successFilePath, []byte("dummy content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 调用Rename4Scanner函数
	err = Rename4Scanner(tempDir, successFilePath)
	if err != nil {
		t.Error(err)
	}

	// 验证重命名结果 - 应该只有失败批次的文件被重命名，不包括基准文件
	expectedFiles := []string{
		"Scan_0015.jpg", // 原来的 Scan_0001.jpg
		"Scan_0016.jpg", // 原来的 Scan_0002.jpg
		"Scan_0017.jpg", // 原来的 Scan_0003.jpg
		"Scan_0014.jpg", // 基准文件应该保持不变
	}

	for _, expectedFile := range expectedFiles {
		expectedPath := filepath.Join(tempDir, expectedFile)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected file does not exist: %s", expectedPath)
		} else {
			t.Logf("File exists as expected: %s", expectedPath)
		}
	}

	// 验证旧的文件名不再存在
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

package writer

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestShouldRotateByTime(t *testing.T) {
	tests := []struct {
		name         string
		lastModified time.Time
		currentTime  time.Time
		expected     bool
	}{
		{
			name:         "Same day",
			lastModified: time.Date(2024, time.December, 17, 10, 0, 0, 0, time.UTC),
			currentTime:  time.Date(2024, time.December, 17, 15, 0, 0, 0, time.UTC),
			expected:     false, // Same day, no rotation
		},
		{
			name:         "Different day, same year",
			lastModified: time.Date(2024, time.December, 17, 10, 0, 0, 0, time.UTC),
			currentTime:  time.Date(2024, time.December, 18, 10, 0, 0, 0, time.UTC),
			expected:     true, // Different day, needs rotation
		},
		{
			name:         "Different month, same year",
			lastModified: time.Date(2024, time.December, 17, 10, 0, 0, 0, time.UTC),
			currentTime:  time.Date(2025, time.January, 17, 10, 0, 0, 0, time.UTC),
			expected:     true, // Different month, needs rotation
		},
		{
			name:         "Different year",
			lastModified: time.Date(2024, time.December, 17, 10, 0, 0, 0, time.UTC),
			currentTime:  time.Date(2025, time.December, 17, 10, 0, 0, 0, time.UTC),
			expected:     true, // Different year, needs rotation
		},
	}

	// Loop through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fw := &FileWriter{lastModified: tt.lastModified}

			// Call the method with the currentTime and check the result
			result := fw.shouldRotateByTime(tt.currentTime)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
func TestShouldRotateBySize(t *testing.T) {
	// 创建一个临时文件
	tempFile, err := os.CreateTemp("test", "test_file_")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // 确保测试结束后删除临时文件

	tests := []struct {
		name        string
		maxSize     int64
		currentFile *os.File
		expected    bool
		setup       func() // 用于在测试中设置文件大小
	}{
		{
			name:        "maxSize <= 0",
			maxSize:     0,
			currentFile: tempFile,
			expected:    false, // maxSize <= 0 不需要切割
		},
		{
			name:        "no currentFile",
			maxSize:     100,
			currentFile: nil,
			expected:    false, // 没有文件时不需要切割
		},
		{
			name:        "file size exceeds maxSize",
			maxSize:     100,
			currentFile: tempFile,
			expected:    true, // 文件大小超出最大限制
			setup: func() {
				// 设置文件的大小超过 100 字节
				tempFile.WriteString("This is a test content:sadhajdnaksjndkandckasnjkasncasnkdcnaskdnaksndkassndkjasndkjasndjkankjdnaskndkjasndkandknaskjdnaskndkasndkja")
			},
		},
		{
			name:        "file size does not exceed maxSize",
			maxSize:     100,
			currentFile: tempFile,
			expected:    false, // 文件大小不超过最大限制
			setup: func() {
				// 设置文件的大小小于 100 字节
				tempFile.WriteString("Small file.")
			},
		},
	}

	// Loop through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置文件内容
			if tt.setup != nil {
				tt.setup()
			}

			// 创建 FileWriter 实例
			fw := &FileWriter{
				currentFile: tt.currentFile,
				maxSize:     tt.maxSize,
			}

			// 调用 shouldRotateBySize 函数并验证结果
			result := fw.shouldRotateBySize()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
func createTestFiles() error {
	// 创建 test 目录，如果它不存在
	dir := "test"
	err := os.MkdirAll(dir, 0755) // 0755 是常用的权限
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// 创建 test/test1.log 文件
	file1Path := fmt.Sprintf("%s/test1.log", dir)
	file1, err := os.Create(file1Path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", file1Path, err)
	}
	defer file1.Close() // 确保文件关闭

	// 在 test/test1.log 中写入一些内容
	_, err = file1.WriteString("This is test1.log content.")
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v", file1Path, err)
	}

	// 创建 test/test2-(1).log 文件
	file2Path := fmt.Sprintf("%s/test2-(1).log", dir)
	file2, err := os.Create(file2Path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", file2Path, err)
	}
	defer file2.Close() // 确保文件关闭

	// 在 test/test2-(1).log 中写入一些内容
	_, err = file2.WriteString("This is test2-(1).log content.")
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v", file2Path, err)
	}

	return nil
}

func TestFileWriter_getRotateNameBySize(t *testing.T) {
	err1 := createTestFiles()
	if err1 != nil {
		t.Fatalf("Failed to create test files: %v", err1)
	}
	type args struct {
		oldname string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{oldname: "test/test1"}, "test/test1-(1)"},
		{"test2", args{oldname: "test/test2"}, "test/test2-(2)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileWriter{}
			if got := f.getRotateNameBySize(tt.args.oldname); got != tt.want {
				t.Errorf("getRotateNameBySize() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestConcat(t *testing.T) {
	tests := []struct {
		name     string
		oldname  string
		newpart  string
		expected string
	}{
		{
			name:     "Normal concatenation",
			oldname:  "test",
			newpart:  "v1",
			expected: "test-v1",
		},
		{
			name:     "Empty oldname",
			oldname:  "",
			newpart:  "v1",
			expected: "-v1", // 空字符串与 newpart 连接后应该是 "-v1"
		},
		{
			name:     "Empty newpart",
			oldname:  "test",
			newpart:  "",
			expected: "test-", // oldname 和 空字符串连接后应该是 "test-"
		},
		{
			name:     "Both are empty",
			oldname:  "",
			newpart:  "",
			expected: "-", // 两个空字符串连接后是 "-"
		},
		{
			name:     "Special characters",
			oldname:  "test$&",
			newpart:  "v1@#",
			expected: "test$&-v1@#",
		},
		{
			name:     "With spaces",
			oldname:  "hello world",
			newpart:  "v2",
			expected: "hello world-v2", // 有空格的字符串也会正确拼接
		},
	}

	// Loop through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 FileWriter 实例
			fw := &FileWriter{}

			// 调用 concat 方法
			result := fw.concat(tt.oldname, tt.newpart)

			// 验证结果
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
func TestCreateFile(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		expectedErr    bool   // 是否预期发生错误
		expectedFile   string // 期望文件名
		expectedOpened bool   // 是否期望文件被成功打开
	}{
		{
			name:           "Successful file creation",
			fileName:       "test/testfile.log",
			expectedErr:    false,
			expectedFile:   "test/testfile.log",
			expectedOpened: true,
		},
	}

	// Loop through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 FileWriter 实例
			fw := &FileWriter{}

			// 尝试创建文件
			err := fw.createFile(tt.fileName)

			// 检查是否预期发生错误
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			// 检查是否正确赋值 lastFileName
			if fw.lastFileName != tt.expectedFile {
				t.Errorf("expected lastFileName to be %v, got %v", tt.expectedFile, fw.lastFileName)
			}

			// 检查文件是否成功打开
			if tt.expectedOpened && fw.currentFile == nil {
				t.Errorf("expected file to be opened, but it was not")
			}
			if !tt.expectedOpened && fw.currentFile != nil {
				t.Errorf("expected file to not be opened, but it was")
			}

			// 清理文件，确保测试后的环境干净
			if fw.currentFile != nil {
				fw.currentFile.Close()
				os.Remove(tt.fileName) // 删除文件
			}
		})
	}
}

func TestFileWriter_Write(t *testing.T) {
	type fields struct {
		filePath     string
		filename     string
		maxSize      int64
		rotateByTime bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"test1", fields{
			filePath:     "test",
			filename:     "app",
			maxSize:      0,
			rotateByTime: false,
		}, args{[]byte("this is a test")}, false},
		{"test2", fields{
			filePath:     "test",
			filename:     "app",
			maxSize:      0,
			rotateByTime: true,
		}, args{[]byte("this is a test")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileWriter{
				filePath:     tt.fields.filePath,
				filename:     tt.fields.filename,
				maxSize:      tt.fields.maxSize,
				rotateByTime: tt.fields.rotateByTime,
			}
			if err := f.Write(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

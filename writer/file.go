package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileWriter 结构体，包含文件路径、大小限制和时间格式
type FileWriter struct {
	//文件后缀默认为.log
	filePath     string    // 文件路径
	filename     string    //文件名称
	maxSize      int64     // 最大文件大小（字节），小于等于0为永不切割
	rotateByTime bool      //是否根据时间来进行切割
	lastModified time.Time //上次修改的时间
	lastFileName string    //上一次编辑的文件
	currentFile  *os.File  // 当前打开的文件

	mutex sync.Mutex
}

// Write 将数据写入文件，支持时间和大小切割
func (f *FileWriter) Write(data []byte) error {
	//写操作,上锁
	f.mutex.Lock()
	defer f.mutex.Unlock()

	currentTime := time.Now()
	newFileName := filepath.Join(f.filePath, f.filename)
	if f.rotateByTime && f.shouldRotateByTime(currentTime) {
		newFileName = f.concat(newFileName, currentTime.Format("2006_01_02"))
	}
	if f.shouldRotateBySize() {
		newFileName = f.getRotateNameBySize(newFileName)
	}
	//如果currentFile是空指针，要创建文件
	//如果这次创建的文件和上次创建的文件名不一样,也需要创建文件
	if fn := newFileName + ".log"; f.currentFile == nil || fn != f.lastFileName {
		err := f.createFile(fn)
		if err != nil {
			return err
		}
	}
	// 写入数据
	_, err := f.currentFile.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}
	return nil
}

// shouldRotateByTime 判断是否需要根据时间切割文件
func (f *FileWriter) shouldRotateByTime(currentTime time.Time) bool {
	// 如果当前时间与上次写入时间不同，则需要切割
	return f.lastModified.Year() != currentTime.Year() ||
		f.lastModified.Month() != currentTime.Month() ||
		f.lastModified.Day() != currentTime.Day()
}

// shouldRotateBySize 判断文件大小是否超出限制
func (f *FileWriter) shouldRotateBySize() bool {
	if f.maxSize <= 0 {
		return false
	}
	//如果开始没有文件
	if f.currentFile == nil {
		return false
	}

	// 获取当前文件的大小
	fileInfo, err := f.currentFile.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Size() > f.maxSize
}

// 为因为大小分割的文件命名
func (f *FileWriter) getRotateNameBySize(oldname string) string {
	// 检查文件是否已存在
	fileName := oldname
	count := 1
	for {
		//oldname + count
		//注意不是filename + count
		fileName = f.concat(oldname, fmt.Sprintf("(%d)", count))
		// 尝试创建文件
		_, err := os.Stat(fileName + ".log")
		if os.IsNotExist(err) {
			//如果不存在说明合法,直接返回
			return fileName
		}
		count++
	}
}

// 连接
func (f *FileWriter) concat(oldname string, newpart string) string {
	return fmt.Sprintf("%s-%s", oldname, newpart)
}

// 创建文件,并给f.lastFileName赋值
func (f *FileWriter) createFile(fn string) error {
	var err error
	f.currentFile, err = os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	f.lastFileName = fn
	return nil
}

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	src := "C:\\Users\\xxx\\Downloads\\win7.iso"
	// C:\Users\xxx\Downloads\win7.iso
	dst := "win7.zip"

	sourceFileAbsPath, _ := filepath.Abs(src)
	destFileAbsPath, _ := filepath.Abs(dst)

	fmt.Println("sourceFileAbsPath: ", sourceFileAbsPath)
	fmt.Println("destFileAbsPath: ", destFileAbsPath)

	fmt.Println(filepath.Base(sourceFileAbsPath))
	fmt.Println(filepath.Base(destFileAbsPath))
	start := time.Now()
	createZipFile(sourceFileAbsPath, destFileAbsPath)
	fmt.Println("耗时：", time.Since(start))
}

func createZipFile(filePath, destZipPath string) error {

	// 如果目标文件已存在，则删除
	if _, err := os.Stat(destZipPath); err == nil {
		if err := os.Remove(destZipPath); err != nil {
			return err
		}
	}

	// 创建 zip 文件
	zipFile, err := os.Create(destZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建 zip.Writer 对象
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 打开要压缩的文件
	fileToCompress, err := os.Open(filePath)
	if err != nil {
		return err
	}
	fmt.Println("fileToCompress.Name(): ", fileToCompress.Name())
	defer fileToCompress.Close()

	// 获取文件信息
	fileInfo, err := fileToCompress.Stat()
	if err != nil {
		return err
	}

	// 创建 zip 文件中的文件头信息
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	// 指定文件名为压缩文件的基本名
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	// 向 zip 文件中添加文件
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 将文件内容复制到 zip 文件中
	_, err = io.Copy(writer, fileToCompress)
	if err != nil {
		return err
	}

	return nil
}

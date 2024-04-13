package _zip

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 全部压缩
func ZipAll(src, dst string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err := fw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// 获取相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(relPath, string(filepath.Separator))

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		if err != nil {
			return
		}
		defer func() {
			if err := fr.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		// 将打开的文件 Copy 到 w
		_, err1 := io.Copy(w, fr)
		if err1 != nil {
			return
		}

		return nil
	})
}

// 排除指定的目录
func shouldExcludeDirectory(dir string, excludeDir []string) bool {
	for _, excluded := range excludeDir {
		if dir == excluded {
			return true
		}
	}
	return false
}

// 排除指定的文件
func shouldExcludeFile(filename string, excludeExt []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, excluded := range excludeExt {
		if ext == "."+excluded {
			return true
		}
	}
	return false
}

// 排除指定的文件、目录压缩
func ZipExclude(src, dst string, excludeExt, excludeDir []string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err := fw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		// 检查是否要排除该目录
		if fi.IsDir() && shouldExcludeDirectory(filepath.Base(path), excludeDir) {
			return filepath.SkipDir
		}

		// 检查是否要排除该文件
		if !fi.IsDir() && shouldExcludeFile(fi.Name(), excludeExt) {
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(relPath, string(filepath.Separator))

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		// 将打开的文件 Copy 到 w
		_, err = io.Copy(w, fr)
		return err
	})
}

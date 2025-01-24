package pull

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Package 打包为可load镜像
func Package(path string, imageName string, tag string, arch string, operateSystem string, outputFile *string) (*string, error) {
	// 确定打包文件名
	if outputFile == nil {
		defaultOutputFile := strings.ReplaceAll(imageName, "/", "_") + "_" + tag + "_" + arch + "_" + operateSystem + ".tar.gz"
		outputFile = &defaultOutputFile
	}
	// 创建镜像包
	imagePath := filepath.Join(filepath.Join(filepath.Dir(path), *outputFile))
	file, err := os.Create(imagePath)
	if err != nil {
		fmt.Println("Error while creating tar file", err)
		return nil, err
	}
	defer file.Close()
	return outputFile, tarImage(path, imagePath)
}

// tarImage 将目录打包为tar文件
func tarImage(path string, dest string) error {
	var writer io.WriteCloser
	// 根据后缀名判断创建不同的写入器
	switch filepath.Ext(dest) {
	case ".gz":
		f, err := os.Create(dest)
		if err != nil {
			return err
		}
		writer = f
		gzw := gzip.NewWriter(f)
		writer = gzw
		defer gzw.Close()
	case ".tar":
		f, err := os.Create(dest)
		if err != nil {
			return err
		}
		writer = f
	default:
		return fmt.Errorf("unknown file extension: %s", filepath.Ext(dest))
	}
	defer writer.Close()

	tw := tar.NewWriter(writer)
	defer tw.Close()

	// 计算总文件数
	totalFiles := 0
	err := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 创建进度条
	bar := progressbar.NewOptions(totalFiles, progressbar.OptionSetDescription("Packaging..."))

	return filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(path, file)
		if err != nil {
			return err
		}
		header.Name = rel

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()
			_, err = io.Copy(tw, data)
			if err != nil {
				return err
			}
			bar.Add(1)
		}
		return nil
	})
}

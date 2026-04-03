package pull

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// Package 打包为可load镜像
func Package(path string, imageName string, tag string, arch string, operateSystem string, outputPath string, outputFile *string) (*string, error) {
	if outputFile == nil {
		defaultOutputFile := strings.ReplaceAll(imageName, "/", "_") + "_" + tag + "_" + arch + "_" + operateSystem + ".tar.gz"
		outputFile = &defaultOutputFile
	}
	if outputPath == "" {
		outputPath = "."
	}
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		fmt.Println("Error while creating output directory", err)
		return nil, err
	}

	imagePath := filepath.Join(outputPath, *outputFile)
	if err := tarImage(path, imagePath); err != nil {
		return nil, err
	}
	return &imagePath, nil
}

// tarImage 将目录打包为tar文件
func tarImage(path string, dest string) (err error) {
	var writer io.Writer
	var closers []io.Closer

	switch filepath.Ext(dest) {
	case ".gz":
		f, createErr := os.Create(dest)
		if createErr != nil {
			return createErr
		}
		gzw := gzip.NewWriter(f)
		writer = gzw
		closers = append(closers, gzw, f)
	case ".tar":
		f, createErr := os.Create(dest)
		if createErr != nil {
			return createErr
		}
		writer = f
		closers = append(closers, f)
	default:
		return fmt.Errorf("unknown file extension: %s", filepath.Ext(dest))
	}
	defer func() {
		for _, closer := range closers {
			if closeErr := closer.Close(); err == nil && closeErr != nil {
				err = closeErr
			}
		}
	}()

	tw := tar.NewWriter(writer)
	defer func() {
		if closeErr := tw.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	totalFiles := 0
	err = filepath.Walk(path, func(file string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if fi.Mode().IsRegular() {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions(totalFiles, progressbar.OptionSetDescription("Packaging..."))

	return filepath.Walk(path, func(file string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		header, headerErr := tar.FileInfoHeader(fi, "")
		if headerErr != nil {
			return headerErr
		}
		rel, relErr := filepath.Rel(path, file)
		if relErr != nil {
			return relErr
		}
		header.Name = rel

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			data, openErr := os.Open(file)
			if openErr != nil {
				return openErr
			}
			defer data.Close()
			if _, err = io.Copy(tw, data); err != nil {
				return err
			}
			_ = bar.Add(1)
		}
		return nil
	})
}

package storage

import (
	"fmt"
	"funtabs-server/config"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type localStorage struct {
	basePath string
}

func newLocalStorage() Storage {
	base := config.Cfg.Storage.LocalPath
	if base == "" {
		base = "./uploads"
	}
	if err := os.MkdirAll(base, 0755); err != nil {
		panic(fmt.Sprintf("创建本地存储目录失败: %v", err))
	}
	return &localStorage{basePath: base}
}

func (s *localStorage) Save(file multipart.File, header *multipart.FileHeader, dir string) (string, error) {
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".bin"
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	subDir := filepath.Join(s.basePath, dir)
	if err := os.MkdirAll(subDir, 0755); err != nil {
		return "", fmt.Errorf("创建子目录失败: %w", err)
	}

	dst, err := os.Create(filepath.Join(subDir, filename))
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err = file.Seek(0, 0); err != nil {
		return "", err
	}
	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 返回可访问的 URL 路径（Gin 会挂载 /uploads 静态目录）
	urlPath := "/uploads/" + dir + "/" + filename
	if cdn := config.Cfg.Storage.CDN; cdn != "" {
		return strings.TrimRight(cdn, "/") + urlPath, nil
	}
	return urlPath, nil
}

func (s *localStorage) Delete(urlPath string) error {
	if urlPath == "" {
		return nil
	}
	// 把 URL 路径还原为本地路径：/uploads/avatars/xxx.jpg → basePath/avatars/xxx.jpg
	rel := strings.TrimPrefix(urlPath, "/uploads/")
	localPath := filepath.Join(s.basePath, rel)
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

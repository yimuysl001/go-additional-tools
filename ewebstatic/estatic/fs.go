package estatic

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

func TarGzipEmbedFS(fsys fs.FS, root string) ([]byte, error) {

	// 打开源文件
	//file, err := fsys.Open(root)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to open file: %w", err)
	//}
	//defer file.Close()
	//
	//// 创建内存缓冲区用于存储压缩后的数据
	//var buf bytes.Buffer
	//gzipWriter := gzip.NewWriter(&buf)
	//
	//// 将源文件内容复制到 gzip writer
	//_, err = io.Copy(gzipWriter, file)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to compress file: %w", err)
	//}
	//
	//// 关闭 gzip writer 以完成压缩
	//err = gzipWriter.Close()
	//if err != nil {
	//	return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	//}
	//
	//return buf.Bytes(), nil
	//
	//
	//
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	//tw := tar.NewWriter(gw)
	defer g.Try(context.TODO(), func(ctx context.Context) {
		//tw.Close()
		gw.Close()
	})

	// 遍历文件系统
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath := strings.TrimPrefix(path, root+"/")
		if relPath == path { // 处理根目录情况
			relPath = filepath.Base(path)
		}

		// 获取文件信息
		info, err := d.Info()
		if err != nil {
			return err
		}

		// 创建tar头
		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}
		header.Name = relPath

		//// 写入tar头
		//if err := tw.WriteHeader(header); err != nil {
		//	return err
		//}

		// 如果是普通文件，写入内容
		if !d.IsDir() {
			file, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(gw, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk embedded FS: %w", err)
	}

	return buf.Bytes(), nil
}

func InitPublic(ctx context.Context, staticFiles fs.FS, dir ...string) bool {
	serverRoot := g.Cfg().MustGet(ctx, "server.serverRoot").String()
	if serverRoot == "" { // 未配置路径，不需要处理swagger
		g.Log().Error(ctx, "未开启serverRoot")
		return false
	}

	//embedFS, err := TarGzipEmbedFS(staticFiles, fspath)
	//if err != nil {
	//	g.Log().Error(ctx, err)
	//	return false
	//}
	var (
		err error
	)
	if len(dir) > 0 && dir[0] != "" {
		staticFiles, err = fs.Sub(staticFiles, dir[0])
		if err != nil {
			g.Log().Error(ctx, err)
			return false
		}

	}

	embedFS, err := GzipFSAndHex(staticFiles)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}

	if !strings.Contains(serverRoot, ":") { // 全路径不处理
		var pwd = gfile.Pwd()
		pwd = strings.ReplaceAll(pwd, "\\", "/")
		serverRoot = path.Join(pwd, serverRoot)
	}
	if !strings.HasSuffix(serverRoot, "/") {
		serverRoot = serverRoot + "/"
	}

	err = gres.Add(embedFS, serverRoot)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	g.Log().Debug(ctx, "静态资源添加完成 添加完成")
	return true

}

// GzipFSAndHex 将 fs.FS 中的所有文件压缩成 gzip 并返回 hex 字符串
func GzipFSAndHex(fsys fs.FS) (string, error) {
	// 创建一个缓冲区来存储 gzip 数据
	var buf bytes.Buffer

	// 创建 gzip writer
	gw := gzip.NewWriter(&buf)
	defer gw.Close()

	// 创建 zip writer
	zw := zip.NewWriter(gw)
	defer zw.Close()

	// 遍历 fs.FS 中的所有文件
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 获取文件信息
		info, err := d.Info()
		if err != nil {
			return err
		}

		// 跳过根目录
		if path == "." {
			return nil
		}

		// 创建 zip 文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置正确的文件名（包含路径）
		header.Name = filepath.ToSlash(path)

		// 如果是目录，确保以 / 结尾
		if info.IsDir() {
			if !strings.HasSuffix(header.Name, "/") {
				header.Name += "/"
			}
		}

		// 创建文件写入器
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果不是目录，写入文件内容
		if !info.IsDir() {
			file, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(writer, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("遍历文件系统失败: %w", err)
	}

	// 确保所有数据都被写入
	if err := zw.Close(); err != nil {
		return "", fmt.Errorf("关闭 zip writer 失败: %w", err)
	}
	if err := gw.Close(); err != nil {
		return "", fmt.Errorf("关闭 gzip writer 失败: %w", err)
	}

	// 转换为十六进制字符串
	hexStr := hex.EncodeToString(buf.Bytes())
	return hexStr, nil
}

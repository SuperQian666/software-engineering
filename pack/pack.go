package pack

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//Tar 打包
func Tar(src, dst string) (err error) {
	// 创建文件
	fw, err := os.Create(dst)
	if err != nil {
		return
	}
	defer fw.Close()

	// 创建 Tar.Writer 结构
	tw := tar.NewWriter(fw)
	// 如果需要启用 gzip 将上面代码注释，换成下面的

	defer tw.Close()

	// 递归处理目录及目录下的所有文件和目录
	return filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {
		// 这个闭包会返回个 error ，所以先要处理一下
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}

		hdr.Name = fileName[strings.LastIndex(fileName, "\\")+1:]

		// 写入文件信息
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// 判断下文件是否是标准文件，如果不是就不处理了，
		// 如： 目录，这里就只记录了文件信息，不会执行下面的 copy
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		fr, err := os.Open(fileName)
		defer fr.Close()
		if err != nil {
			return err
		}

		// copy 文件数据到 tw
		n, err := io.Copy(tw, fr)
		if err != nil {
			return err
		}

		log.Printf("成功打包 %s ，共写入了 %d 字节的数据\n", fileName, n)

		return nil
	})

}

func UnTar(src, dst string) (err error) {
	// 打开准备解压的 tar 包
	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer fr.Close()

	tr := tar.NewReader(fr)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		dstFileDir := filepath.Join(dst, hdr.Name)

		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeDir:
			// 判断下目录是否存在，不存在就创建
			if b := ExistDir(dstFileDir); !b {
				// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
				// 可以递归创建每一级目录
				if err := os.MkdirAll(dstFileDir, 0775); err != nil {
					return err
				}
			}
		case tar.TypeReg:

			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			n, err := io.Copy(file, tr)
			if err != nil {
				return err
			}
			// 将解压结果输出显示
			fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}

	return nil
}

// 判断目录是否存在
func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

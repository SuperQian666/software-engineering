package encrypt

import (
	"encoding/binary"
	"io"
	"os"
)

func EncryptFile(src, dest, pwd string) (err error) {
	realpwd := generateKey([]byte(pwd))

	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	pwd = string(realpwd)

	var FrameSize = 2048
	//加密处理
	frameBuf := make([]byte, FrameSize) //一次读取多少个字节
	sizeBuf := make([]byte, 2)
	for {
		n, err := srcfile.Read(frameBuf)
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}
		if n <= 0 {
			break
		}

		encryptedFrame, err := encryptByAes(frameBuf[:n], pwd)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(sizeBuf[:], uint16(len(encryptedFrame)))
		_, err = destfile.Write(sizeBuf)
		if err != nil {
			return err
		}
		_, err = destfile.Write([]byte(encryptedFrame))
		if err != nil {
			return err
		}
	}

	return nil
}

func DecryptFile(src, dest, pwd string) error {
	correctPwd := generateKey([]byte(pwd))
	pwd = string(correctPwd)

	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	sizeBuf := make([]byte, 2)
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()
	for {
		n, err := srcfile.Read(sizeBuf)
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}
		if n <= 0 {
			break
		}
		encryptedFrameSize := binary.BigEndian.Uint16(sizeBuf)
		buf := make([]byte, encryptedFrameSize)
		n, err = srcfile.Read(buf)
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}

		if n <= 0 {
			break
		}

		d, err := decryptByAes(string(buf[:n]), pwd)
		if err != nil {
			return err
		}
		destfile.Write(d)
	}

	return nil
}

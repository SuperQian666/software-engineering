package encrypt

import (
	"errors"
	"os"
)

func EncryptFile(filepath, fName string) (err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	fInfo, _ := f.Stat()
	//the size of the file
	fSize := fInfo.Size()
	maxSize := 1024 * 10

	//encrypted file
	ff, err := os.OpenFile("encryptFile_"+fName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("文件写入出错")
	}
	defer ff.Close()

	for i := 0; i < int(fSize/int64(maxSize+1)); i++ {

	}
	return nil
}

package encrypt

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var data = "wcdicwicnewi"
var key = "acwcqe"

func TestAes(t *testing.T) {
	//key := "acwcqe"
	//加密
	str, _ := EncryptByAes([]byte(data), key)
	//解密
	str1, _ := DecryptByAes(str, key)
	//打印
	fmt.Printf(" 加密：%v\n 解密：%s\n ",
		str, str1,
	)

}

//测试速度
func TestAesTime(t *testing.T) {
	startTime := time.Now()
	count := 1000000
	for i := 0; i < count; i++ {
		str, _ := EncryptByAes([]byte(data), key)
		DecryptByAes(str, key)
	}
	fmt.Printf("%v次 - %v", count, time.Since(startTime))
}

func TestFunction(t *testing.T) {
	filePath := "../src.txt"
	f, _ := os.Open(filePath)
	finfo, _ := f.Stat()
	name := finfo.Name()
	size := finfo.Size()
	fmt.Printf("%v\n", name)
	fmt.Printf("%v\n", size)
}

func TestEncryptFile(t *testing.T) {
	EncryptFile("../src.txt", "../test", "12345")

}

func TestDecryptFile(t *testing.T) {
	DecryptFile("../test", "../testDecrypt", "12345")
}

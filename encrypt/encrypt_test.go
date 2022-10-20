package encrypt

import (
	"fmt"
	"testing"
	"time"
)

var data = "wcdicwicnewi"

func TestAes(t *testing.T) {
	//加密

	str, _ := EncryptByAes([]byte(data))
	//解密
	str1, _ := DecryptByAes(str)
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
		str, _ := EncryptByAes([]byte(data))
		DecryptByAes(str)
	}
	fmt.Printf("%v次 - %v", count, time.Since(startTime))
}

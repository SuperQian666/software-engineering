package copy

import (
	"errors"
	"fmt"
	"io/ioutil"
)

func FileCopy(src, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//判断目标文件是否已存在，若存在不可备份
	if _, err := ioutil.ReadFile(dest); err == nil {
		return errors.New("目标文件已存在")
	}
	ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		fmt.Println("Error creating", dest)
		fmt.Println(err)
		return err
	}
	return nil
}

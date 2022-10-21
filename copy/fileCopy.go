package copy

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func FileCopy(src, dest, file string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return err
	}

	toPath := strings.Split(dest, file)
	_, err = os.Lstat(toPath[0])
	if err != nil {
		err = os.MkdirAll(toPath[0], 0777)
		if err != nil {
			return errors.New("创建文件夹失败")
		}
	}

	//判断目标文件是否已存在，若存在不可备份
	if _, err = ioutil.ReadFile(dest); err == nil {
		return errors.New("目标文件已存在")
	}

	err = ioutil.WriteFile(dest, input, 0777)

	if err != nil {
		fmt.Println("Error creating", dest)
		fmt.Println(err)
		return err
	}
	return nil
}

func Copy(src, dest string) error {
	files, _ := ioutil.ReadDir(src)
	pathSeparator := "/"
	for _, file := range files {
		if file.IsDir() {
			if err := Copy(src+pathSeparator+file.Name(), dest+pathSeparator+file.Name()); err != nil {
				return err
			}
		} else {
			oldPath := fmt.Sprintf("%s%s%s", src, pathSeparator, file.Name())
			newPath := fmt.Sprintf("%s%s%s", dest, pathSeparator, file.Name())
			if err := FileCopy(oldPath, newPath, file.Name()); err != nil {
				return err
			}
		}
	}
	_, err := os.Lstat(dest)
	if err != nil {
		os.MkdirAll(dest, 0777)
	}

	return nil
}

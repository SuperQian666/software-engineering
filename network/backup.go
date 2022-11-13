package network

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var BaseUploadPath = "/var/file"
var user = "soft"
var passwd = "software_user"
var hostAddr = "121.4.113.117:22"

//todo 增加文件夹备份功能

//利用ssh连接云服务器
func connect(user, passwd, hostAddr string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(passwd))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         15 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if sshClient, err = ssh.Dial("tcp", hostAddr, clientConfig); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func uploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) {
	//打开本地文件夹流
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("路径错误 ", err)
	}
	//先创建最外层文件夹
	sftpClient.Mkdir(remotePath)
	//遍历文件夹内容
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		//判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if backupDir.IsDir() {
			if err := sftpClient.Mkdir(remoteFilePath); err != nil {
				log.Fatal("makeDir failed", err)
			}
			uploadDirectory(sftpClient, localFilePath, remoteFilePath)
		} else {
			uploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
		}
	}

	fmt.Println(localPath + "  copy directory to remote server finished!")
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	//打开本地文件流
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)

	}
	//关闭文件流
	defer srcFile.Close()
	//上传到远端服务器的文件名,与本地路径末尾相同
	var remoteFileName = path.Base(localFilePath)
	//打开远程文件,如果不存在就创建一个
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)

	}
	//关闭远程文件
	defer dstFile.Close()
	//读取本地文件,写入到远程文件中(这里没有分快穿,自己写的话可以改一下,防止内存溢出)
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)

	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + "  copy file to remote server finished!")
}

func Upload(localPath string, remotePath string) error {
	sftpClient, err := connect(user, passwd, hostAddr)
	if err != nil {
		return errors.New("ssh连接远程服务器失败")
	}
	//获取路径的属性
	s, err := os.Stat(localPath)
	if err != nil {
		return errors.New("文件路径不存在")
	}
	//判断是否是文件夹
	if s.IsDir() {
		uploadDirectory(sftpClient, localPath, remotePath)
	} else {
		uploadFile(sftpClient, localPath, remotePath)
	}

	return nil
}

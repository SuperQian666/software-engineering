package network

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var BaseUploadPath = "/var/file"
var user = "soft"
var passwd = "software_user"
var hostAddr = "121.4.113.117:22"

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

func uploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) error {
	//打开本地文件夹流
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return errors.New(fmt.Sprintf("路径错误: %s", err))
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
				return errors.New(fmt.Sprintf("makeDir failed: %s", err))
			}
			if err := uploadDirectory(sftpClient, localFilePath, remoteFilePath); err != nil {
				return err
			}
		} else {
			return uploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
		}
	}

	return nil
}
func downloadDirectory(sftpClient *sftp.Client, localPath, remotePath string) error {
	remoteFiles, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		return errors.New(fmt.Sprintf("远端路径错误: %s", err))
	}
	//先创建最外层文件夹
	os.Mkdir(localPath, 0755)
	//遍历文件夹内容
	for _, getBackDir := range remoteFiles {
		localFilePath := path.Join(localPath, getBackDir.Name())
		remoteFilePath := path.Join(remotePath, getBackDir.Name())
		//判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if getBackDir.IsDir() {
			if err := os.Mkdir(localFilePath, 0755); err != nil {
				return errors.New(fmt.Sprintf("makeDir failed: %s", err))
			}
			return downloadDirectory(sftpClient, localFilePath, remoteFilePath)
		} else {
			return downloadFile(sftpClient, path.Join(localPath, getBackDir.Name()), remotePath)
		}
	}
	return nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) error {
	//打开本地文件流
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("os.Open error: %s ", localFilePath))

	}
	//关闭文件流
	defer srcFile.Close()
	//上传到远端服务器的文件名,与本地路径末尾相同
	var remoteFileName = localFilePath[strings.LastIndex(localFilePath, "\\")+1:]
	//打开远程文件,如果不存在就创建一个
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		return errors.New(fmt.Sprintf("sftpClient.Create error : %s", path.Join(remotePath, remoteFileName)))

	}
	//关闭远程文件
	defer dstFile.Close()
	//读取本地文件,写入到远程文件中(这里没有分快穿,自己写的话可以改一下,防止内存溢出)
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return errors.New(fmt.Sprintf("ReadAll error: %s", localFilePath))

	}
	dstFile.Write(ff)
	return nil
}

func downloadFile(sftpClient *sftp.Client, localPath, remotePath string) error {
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return errors.New(fmt.Sprintf("remote open error: %s ", remotePath))
	}
	defer remoteFile.Close()

	var localFileName = remotePath[strings.LastIndex(remotePath, "/")+1:]
	localFile, err := os.Create(path.Join(localPath, localFileName))
	if err != nil {
		return errors.New(fmt.Sprintf("local Create file error : %s", path.Join(localPath, localFileName)))
	}
	defer localFile.Close()
	if _, err := remoteFile.WriteTo(localFile); err != nil {
		return err
	}
	return nil
}

func Upload(localPath, remotePath string) error {
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
		return uploadDirectory(sftpClient, localPath, remotePath)
	} else {
		return uploadFile(sftpClient, localPath, remotePath)
	}

	return nil
}

func Download(localPath, remotePath string) error {
	sftpClient, err := connect(user, passwd, hostAddr)
	if err != nil {
		return errors.New("ssh连接远程服务器失败")
	}

	s, err := sftpClient.Stat(remotePath)
	if err != nil {
		return errors.New("远程文件不存在")
	}

	if s.IsDir() {
		return downloadDirectory(sftpClient, localPath, remotePath)
	} else {
		return downloadFile(sftpClient, localPath, remotePath)
	}
	return nil
}

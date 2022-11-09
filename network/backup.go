package network

import (
	"errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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

func Upload(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return errors.New("文件不存在")
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(src)
	sftpClient, err := connect(user, passwd, hostAddr)
	defer sftpClient.Close()

	if err != nil {
		return err
	}

	destFile, err := sftpClient.Create(path.Join(dest, remoteFileName))
	if err != nil {
		return err
	}
	defer destFile.Close()

	buf := make([]byte, 1024)

	for {
		count, _ := srcFile.Read(buf)
		if count == 0 {
			break
		}
		if _, err = destFile.Write(buf); err != nil {
			return errors.New("目标文件写入失败")
		}
	}

	return nil
}

func Download(local, remote string) error {
	sftpClient, err := connect(user, passwd, hostAddr)
	defer sftpClient.Close()

	if err != nil {
		return errors.New("文件不存在")
	}

	remoteFile, err := sftpClient.Open(remote)
	if err != nil {
		return errors.New("云端文件获取失败")
	}
	defer remoteFile.Close()

	_, err = os.Open(local)
	if err != nil {
		return errors.New("本地存储路径不存在")
	}

	localFilePath := path.Join(local, path.Base(remote))
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	if _, err = remoteFile.WriteTo(localFile); err != nil {
		return err
	}

	return nil

}

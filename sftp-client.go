package main

// link: http://blog.ralch.com/tutorial/golang-ssh-connection/
import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	CERT_PASSWORD        = 1
	CERT_PUBLIC_KEY_FILE = 2
	DEFAULT_TIMEOUT      = 3 // second
)

type SSH struct {
	Ip     string
	User   string
	Cert   string //password or key file path
	Port   int
	client *ssh.Client
	sftpc  *sftp.Client
}

func (ssh_client *SSH) readPublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func (ssh_client *SSH) Connect(mode int) {

	var ssh_config *ssh.ClientConfig
	var auth []ssh.AuthMethod
	if mode == CERT_PASSWORD {
		auth = []ssh.AuthMethod{ssh.Password(ssh_client.Cert)}
	} else if mode == CERT_PUBLIC_KEY_FILE {
		auth = []ssh.AuthMethod{ssh_client.readPublicKeyFile(ssh_client.Cert)}
	} else {
		log.Println("does not support mode: ", mode)
		return
	}

	ssh_config = &ssh.ClientConfig{
		User: ssh_client.User,
		Auth: auth,
		//需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * DEFAULT_TIMEOUT,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ssh_client.Ip, ssh_client.Port), ssh_config)
	if err != nil {
		fmt.Println(err)
		return
	}

	sftpc, err := sftp.NewClient(client)
	if err != nil {
		fmt.Println(err)
		client.Close()
	}

	ssh_client.client = client
	ssh_client.sftpc = sftpc
}

func (ssh_client *SSH) Close() {
	ssh_client.sftpc.Close()
	ssh_client.client.Close()
}

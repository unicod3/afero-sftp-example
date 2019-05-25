package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/afero/sftpfs"
)

func main() {
	scli := &SSH{
		Ip:   "host",
		User: "user",
		Port: 22,
		Cert: os.Getenv("HOME") + "/.ssh/id_rsa",
	}
	scli.Connect(CERT_PUBLIC_KEY_FILE)
	defer scli.Close()
	fmt.Println("Connected to the server ")

	appFs := sftpfs.New(scli.sftpc)

	f1, err := appFs.Open("/home/cnone/test.xxx")
	if err != nil {
		log.Fatalln("open: %v", err)
	}
	defer f1.Close()

	b := make([]byte, 100)

	_, err = f1.Read(b)
	fmt.Println(string(b))

	path := "/home/cnone/xxx.test"
	f, err := appFs.Create(path)
	if err != nil {
		log.Fatalln(appFs.Name(), "Create failed:", err)
		f.Close()
	}
	io.WriteString(f, "initial")
	f.Close()

	fmt.Println("File created at: ", path)

}

package lib

import (
	"bufio"
	"io"
	"os"
	"strings"
)

/*用于读取配置文件每一行的结构体*/
type DialAccount struct {
	Account  string
	Password string
	IsActive bool
	Ppp string
	Ip string
	Dst string
	MacVLan string
	Mac string
}

/*读取配置文件*/
func ReadConfFile(confFile string) ([]DialAccount, error) {
	f, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var dialArray []DialAccount
	var account, password string
	for {
		b, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		lineString := strings.TrimSpace(string(b))
		startPos := strings.Index(lineString, " ")
		if startPos == -1 {
			break
		}
		account = strings.TrimSpace(lineString[0:startPos])
		password = ""
		if len(lineString) > startPos {
			password = strings.TrimSpace(lineString[startPos:])
		}
		dialArray = append(dialArray, DialAccount{
			Account:  account,
			Password: password,
		})
	}
	return dialArray, nil
}

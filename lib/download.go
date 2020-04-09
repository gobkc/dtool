package lib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func Download(downloadUrl string, outIp string, proxyAddress string, fb func(length, downLen int64) bool) error {
	var (
		fSize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)

	httpTransport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			lAddr, err := net.ResolveTCPAddr(network, outIp+":0")
			if err != nil {
				return nil, err
			}
			//被请求的地址
			rAddr, err := net.ResolveTCPAddr(network, addr)
			if err != nil {
				return nil, err
			}
			conn, err = net.DialTCP(network, lAddr, rAddr)
			if err != nil {
				return nil, err
			}
			deadline := time.Now().Add(35 * time.Second)
			conn.SetDeadline(deadline)
			return conn, nil
		},
	}

	//创建一个http client
	client := new(http.Client)
	client.Transport = httpTransport

	if proxyAddress != "" {
		client.Transport = httpTransport
		proxy, err := url.Parse(proxyAddress)
		if err != nil {
			return nil
		}
		httpTransport.Proxy = http.ProxyURL(proxy)
	}

	//get方法获取资源
	resp, err := client.Get(downloadUrl)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fSize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//创建文件
	file, err := os.Create("/dev/null")
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()

	over := false
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		if over = fb(fSize, written); over == true {
			return nil
		}
	}

	/*如果还没到超时时间，就已经结束了。这里删除一次防止意外*/
	if over == false {
	}
	return err
}

func IsFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(info)
		return false
	}
	if filesize == info.Size() {
		fmt.Println("安装包已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}

func Download2(url string, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	//tmpFilePath := localPath + ".download"
	//fmt.Println(tmpFilePath)
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	//if IsFileExist(localPath, fsize) {
	//	return err
	//}
	fmt.Println("fsize", fsize)
	//创建文件
	file, err := os.Create("/dev/null")
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		fb(fsize, written)
	}
	if err == nil {
		file.Close()
		//err = os.Rename(tmpFilePath, localPath)
		//fmt.Println(err)
	}else{
		fmt.Println(err)
	}
	return err
}

func B2M(bytes int64) float64 {
	value := float64(bytes) / 1024 / 1024
	valFloat, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return valFloat
}

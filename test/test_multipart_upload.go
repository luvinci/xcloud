package main

import (
	"bufio"
	"bytes"
	"fmt"
	jsonit "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	username := "luvinci"
	token := "6d19ca982ddd2b990bec0f161d719674"
	filehash := "3d8dcdfed355e259084749e55656b1513d13b3c6"
	filesize := "92313318"
	filepath := "D:\\Pd\\Go从0到1实战微服务版抢红包系统\\01 - 第一章\\1-1 抢红包系统项目演示&导学.mp4"
	filename := "1-1 抢红包系统项目演示&导学.mp4"

	// 请求初始化分块接口
	resp, err := http.PostForm(
		"http://localhost:28080/file/mpupload/upinit" + "?username=" + username + "&token=" + token,
		url.Values{
			"filehash": {filehash},
			"filesize": {filesize},
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 得到uploadId以及服务端指定的分块大小chunkSize
	uploadId := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d\n", uploadId, chunkSize)

	// 请求分块上传接口
	upChunkUrl := "http://localhost:28080/file/mpupload/upchunk" +
		"?username=" + username + "&token=" + token + "&uploadid=" + uploadId
	multipartUpload(upChunkUrl, filepath, chunkSize)

	// 请求分块上传完成接口
	resp, err = http.PostForm(
		"http://localhost:28080/file/mpupload/upcomplete" + "?username=" + username + "&token=" + token,
		url.Values{
			"filehash": {filehash},
			"filesize": {filesize},
			"filename": {filename},
			"uploadid": {uploadId},
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	msg := jsonit.Get(body, "msg").ToString()
	fmt.Println(msg)
}

func multipartUpload(upChunkUrl string, filepath string, chunkSize int) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) // 每次读取chunkSize大小的文件内容
	for {
		n, err := reader.Read(buf)
		if n <= 0 {
			break
		}
		index++

		bufCopied := make([]byte, chunkSize)
		copy(bufCopied, buf)

		go func(b []byte, currIdx int) {
			fmt.Printf("上传编号:%d  上传大小:%d\n", currIdx, len(b))

			resp, err := http.Post(
				upChunkUrl+"&chunkidx="+strconv.Itoa(currIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			msg := jsonit.Get(body, "msg").ToString()
			fmt.Printf("%s -- 编号:%d\n", msg, currIdx)
			resp.Body.Close()

			ch <- currIdx
		}(bufCopied[:n], index)

		// 遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
			}
		}
	}
	for idx := 1; idx <= index; idx++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		}
	}
}

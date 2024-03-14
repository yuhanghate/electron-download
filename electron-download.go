package main

import (
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	baseURL := "https://npmmirror.com/mirrors/electron/"
	var version, _ = readPackageFile()

	// 生成需要下载的electron二进制包
	downloadURLs := buildDownloadURL(baseURL, version)

	downloadElectron(downloadURLs)

}

// 获取缓存目录
func getCacheDir() string {
	usr, _ := user.Current()
	platform := runtime.GOOS
	var cacheDir string

	if platform == "linux" {
		cacheDir = filepath.Join(usr.HomeDir, ".cache", "electron")
	} else if platform == "darwin" {
		cacheDir = filepath.Join(usr.HomeDir, "Library", "Caches", "electron")
	} else if platform == "windows" {
		cacheDir = filepath.Join(usr.HomeDir, "AppData", "Local", "electron", "Cache")
	}

	return cacheDir
}

// 生成下载地址
func buildDownloadURL(baseURL string, version string) []string {

	urls := []string{}

	urls = append(urls, fmt.Sprintf("%s%s/electron-v%s-%s-%s.zip", baseURL, version, version, "win32", "ia32"))
	urls = append(urls, fmt.Sprintf("%s%s/electron-v%s-%s-%s.zip", baseURL, version, version, "win32", "x64"))
	urls = append(urls, fmt.Sprintf("%s%s/electron-v%s-%s-%s.zip", baseURL, version, version, "win32", "arm64"))

	urls = append(urls, fmt.Sprintf("%s%s/electron-v%s-%s-%s.zip", baseURL, version, version, "linux", "x64"))
	urls = append(urls, fmt.Sprintf("%s%s/electron-v%s-%s-%s.zip", baseURL, version, version, "linux", "arm64"))

	urls = append(urls, fmt.Sprintf("%s%s/electron-v%s-%s-%s.zip", baseURL, version, version, "darwin", "x64"))
	urls = append(urls, fmt.Sprintf("%s%s/electron-v%s-%s-%s.zip", baseURL, version, version, "darwin", "arm64"))

	return urls
}

// 下载Electron
func downloadElectron(downloadUrls []string) {
	cacheDir := getCacheDir()

	for _, url := range downloadUrls {

		// 使用 path 包获取文件名
		filename := path.Base(url)

		// 如果 URL 中包含查询参数，需要额外处理
		if strings.ContainsRune(filename, '?') {
			filename = strings.Split(filename, "?")[0]
		}
		savePath := filepath.Join(cacheDir, filename)

		file := downloadFile(url, savePath)
		if file {
			fmt.Printf("下载完成，保存至： %s\n", savePath)
		}

	}

}

// 下载文件
func downloadFile(url string, savePath string) bool {
	// 检查本地文件是否已存在
	if _, err := os.Stat(savePath); err == nil {
		fmt.Printf("\n文件已存在，无需下载:%s\n本地文件路径：%s \n\n", url, savePath)
		return false
	}

	client := http.DefaultClient
	client.Timeout = 60 * 10 * time.Second
	reps, err := client.Get(url)
	if err != nil {
		log.Panic(err.Error())
	}
	if reps.StatusCode == http.StatusOK {
		//保存文件
		file, err := os.Create(savePath)
		if err != nil {
			log.Panic(err.Error())
		}
		defer file.Close() //关闭文件
		//获取下载文件的大小
		length := reps.Header.Get("Content-Length")
		size, _ := strconv.ParseInt(length, 10, 64)
		body := reps.Body //获取文件内容
		bar := pb.Full.Start64(size)
		bar.SetWidth(120)                         //设置进度条宽度
		bar.SetRefreshRate(10 * time.Millisecond) //设置刷新速率
		defer bar.Finish()
		// create proxy reader
		barReader := bar.NewProxyReader(body)
		//写入文件
		writer := io.Writer(file)
		io.Copy(writer, barReader)
		//defer fmt.Printf("\n下载完成，保存至： %s\n", savePath)
	}
	return true
}

// 读取当前目录下的文件
func readPackageFile() (string, error) {

	// 读取 package.json 文件
	filePath := "package.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var parsedData map[string]interface{}
	err = json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		return "", err
	}

	electronVersion := parsedData["devDependencies"].(map[string]interface{})["electron"].(string)
	re := regexp.MustCompile(`\d+\.\d+\.\d+`)
	versionNumbers := re.FindString(electronVersion)
	return versionNumbers, nil
}

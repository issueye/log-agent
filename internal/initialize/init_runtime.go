package initialize

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/pkg/utils"
)

func InitRuntime() {
	// 检查本地是否存在runtime文件夹
	// 获取当前程序的路径
	path := utils.GetWorkDir()
	rtPath := isExistsCreatePath(path, "runtime")
	isExistsCreatePath(rtPath, "data")
	isExistsCreatePath(rtPath, "logs")
	staticPath := isExistsCreatePath(rtPath, "static")
	isExistsCreatePath(staticPath, "admin")
}

func isExistsCreatePath(path, name string) string {
	p := filepath.Join(path, name)
	exists, err := utils.PathExists(p)
	if err != nil {
		panic(err.Error())
	}

	if !exists {
		panic(fmt.Errorf("创建【%s】文件夹失败", name))
	}

	return p
}

func isExistsCreateFile(path, assetName, name string) {
	_, _ = utils.PathExists(path)
	p := fmt.Sprintf("%s/%s", path, name)
	exists := checkFile(p)
	if !exists {
		file, err := global.Asset.ReadFile(assetName)
		if err != nil {
			return
		}

		// 写入文件
		bufferWrite(p, file)
	}
}

func checkFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil

}

func bufferWrite(path string, data []byte) {
	fileHandle, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriter(fileHandle)
	// 字节写入
	_, err = buf.Write(data)
	if err != nil {
		log.Printf("写入文件失败，失败原因：%s", err.Error())
		return
	}
	// 将缓冲中的数据写入
	err = buf.Flush()
	if err != nil {
		log.Println("flush error :", err)
	}
}

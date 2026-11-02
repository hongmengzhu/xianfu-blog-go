package dataFilePg

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/datetimePg"
)

// MakeContent 账号内容
func MakeContent(account, pwd, module string) string {
	return fmt.Sprintf("模块：%s\n账号：%s\n密码：%s", module, account, pwd)
}

func NewAccountFileRecord(pg configPg.Pg, content string) *AccountFileRecord {
	return &AccountFileRecord{
		pg:      pg,
		content: content,
	}
}

// AccountFileRecord 账号文件记录
type AccountFileRecord struct {
	pg      configPg.Pg `value:"${pg}"`
	content string
}

func (s *AccountFileRecord) Write() {
	dataPath, err := filepath.Abs(path.Join(s.pg.Directory.Data))
	if err != nil {
		fmt.Printf("获取 [%s] 的绝对路径失败：%v\n", path.Join(s.pg.Directory.Data), err)
		dataPath = "."
	}
	//fmt.Println("绝对路径:" + dataPath)
	//根据绝对路径 md5 作为文件名
	salt := cryptPg.Md5(dataPath)
	//文件名
	filePath := path.Join(s.pg.Directory.Data, "account-"+salt+".md")
	//
	logStr := fmt.Sprintf("\n生成时间：%s\n%s\n", datetimePg.Now(), s.content)
	//fmt.Println("账号文件记录:" + logStr)
	//
	content := []byte(logStr)

	//获取目录
	dir := path.Dir(filePath)
	//如果目录不存在则创建
	if !s.ExistsObject(dir) {
		// 创建目录的权限 os.ModePerm
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("文件上传失败，创建目录错误：%v\n", err)
			return
		}
	}
	// 以追加模式打开文件（不存在则创建）
	// 打开模式说明：
	// - os.O_APPEND：追加模式
	// - os.O_WRONLY：只写模式
	// - os.O_CREATE：文件不存在则创建
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("打开文件失败：%v\n", err)
		return
	}
	// 延迟关闭文件（必须！否则会导致文件句柄泄漏）
	defer file.Close()

	// 写入内容
	_, err = file.Write(content)
	if err != nil {
		fmt.Printf("追加写入失败：%v\n", err)
		return
	}

	fmt.Println("内容追加成功！")
}

func (s *AccountFileRecord) ExistsObject(name string) bool {

	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

package drive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/h2non/filetype"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/modAttachment"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/types"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/types"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
)

var _ types.FileProvider = (*Local)(nil)

// Local 本地文件上传
// @Description:
type Local struct {
	pg     configPg.Pg     `value:"${pg}"`
	server configPg.Server `value:"${server}"`
	Bus    event.Bus       `autowire:"?"`
}

func (s *Local) PutObject(r io.Reader, put modAttachment.PutFileDto, ext modAttachment.Ext) (modAttachment.Attachment, error) {
	//获取文件名带后缀
	filenameWithSuffix := path.Base(put.Name)
	//获取文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	//去除问号
	if strings.Contains(fileSuffix, "?") {
		fileSuffix = strPg.Replace(fileSuffix, "?", "")
	}
	attachmentCfg := s.pg.Attachment
	domain := attachmentCfg.Domain
	if strPg.IsBlank(domain) {
		domain = s.server.Domain
	}
	//修改为正确后缀
	if strings.Contains(fileSuffix, "awebp") {
		fileSuffix = strPg.Replace(fileSuffix, "awebp", "webp")
	}
	//获取文件名
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	//md5
	fileNewName := datetimePg.NowNotFormat() + "-" + cryptPg.Md5(filenameOnly) + fileSuffix
	attachment := modAttachment.Attachment{
		OriginalName: put.Name,
		Name:         fileNewName,
		Size:         put.Size,
		Ext:          fileSuffix,
		Domain:       domain,
	}
	out := path.Join(attachmentCfg.Dir, datetimePg.YearMonth(), fileNewName)
	out_root := out
	//是否存在跟目录
	if strPg.IsNotBlank(attachmentCfg.DirRoot) {
		out_root = path.Join(attachmentCfg.DirRoot, out)
		if s.ExistsObject(out_root) {
			return attachment, errors.New("文件已存在，请勿重复上传")
		}
	} else {
		if s.ExistsObject(out) {
			return attachment, errors.New("文件已存在，请勿重复上传")
		}
	}
	dir := path.Dir(out_root)
	if !s.ExistsObject(dir) {
		// 创建目录的权限 os.ModePerm
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return attachment, errors.New("文件上传失败，创建目录错误")
		}
	}
	dst, err := os.OpenFile(out_root, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return attachment, errors.New("文件保存失败，请稍后重试")
	}
	defer func() {
		_ = dst.Close()
	}()
	//去除最后末尾 /
	if "/" != strPg.Substr(out, 0, 1) {
		out = "/" + out
	}
	attachment.File = out
	attachment.Url = domain + out
	_, err = io.Copy(dst, r)
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "err=%+v\n", err)
		return attachment, errors.New("文件保存失败")
	}
	buf, _ := os.ReadFile(out_root)
	if filetype.IsImage(buf) {
		fmt.Println("File is an image")
	} else {
		fmt.Println("Not an image")
	}
	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		fmt.Println("Unknown file type")
	} else {
		fmt.Printf("File type matched: %s\n", kind.Extension)
	}

	//保存到数据库
	err = s.Bus.Publish(context.Background(), entityBasic.BasicAttachmentEntity{
		Name:         attachment.Name,
		OriginalName: attachment.OriginalName,
		Description:  attachment.Description,
		Sort:         attachment.Sort,
		Size:         attachment.Size,
		Module:       attachment.Module,
		//Value:         attachment.Value,
		Tag:           attachment.Tag,
		Label:         attachment.Label,
		File:          attachment.File,
		Domain:        attachment.Domain,
		Url:           attachment.Url,
		No:            attachment.No,
		Method:        attachment.Method,
		Ext:           attachment.Ext,
		Category:      attachment.Category,
		Client:        attachment.Client,
		ProtocolSpace: attachment.ProtocolSpace,
	})
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "err=%+v\n", err)
	}
	return attachment, nil
}

func (s *Local) ExistsObject(name string) bool {

	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

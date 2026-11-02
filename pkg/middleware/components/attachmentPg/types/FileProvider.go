package types

import (
	"io"

	modAttachment2 "github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/attachmentPg/modAttachment"
)

type FileProvider interface {
	//
	// PutObject 上传保存
	//  @Description:
	//  @param name 原始文件名称
	//  @param r 文件流
	//  @param size 文件大小
	//  @return Attachment 对象
	//  @return error 报错
	//
	//
	PutObject(r io.Reader, put modAttachment2.PutFileDto, ext modAttachment2.Ext) (modAttachment2.Attachment, error)

	ExistsObject(name string) bool
}

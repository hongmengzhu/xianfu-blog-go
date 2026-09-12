package types

import (
	"io"

	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/modAttachment"
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
	PutObject(r io.Reader, put modAttachment.PutFileDto, ext modAttachment.Ext) (modAttachment.Attachment, error)

	ExistsObject(name string) bool
}

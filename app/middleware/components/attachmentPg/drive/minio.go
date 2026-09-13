package drive

import (
	"context"
	"errors"
	"io"
	"path"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/modAttachment"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/types"
	"github.com/minio/minio-go/v7"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"go-spring.org/log"
)

type Minio struct {
	// 自动注入 minio-client
	Client *minio.Client `autowire:""`
	// 存储桶
	Bucket string `value:"${minio.bucket}"`
	// 存储路径
	Dir string `value:"${file.dir}"`
}

func (s *Minio) PutObject(r io.Reader, put modAttachment.PutFileDto, ext modAttachment.Ext) (modAttachment.Attachment, error) {
	//获取文件名带后缀
	filenameWithSuffix := path.Base(put.Name)
	//获取文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	//获取文件名
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	//md5
	fileNewName := datetimePg.NowNotFormat() + "-" + cryptPg.Md5(filenameOnly)
	attachment := modAttachment.Attachment{
		OriginalName: put.Name,
		Name:         fileNewName,
		Size:         put.Size,
	}
	out := path.Join(s.Dir, put.Name)
	if s.ExistsObject(out) {
		return attachment, errors.New("文件已存在")
	}

	_, err := s.Client.PutObject(context.Background(), s.Bucket, out, r, put.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "minio upload error: %v", err)
		return attachment, errors.New("文件上传失败")
	}

	return attachment, nil
}

func (s *Minio) ExistsObject(name string) bool {
	_, err := s.Client.StatObject(context.Background(), s.Bucket, name, minio.StatObjectOptions{})
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "err=%+v", err)
		if err.Error() == "The specified key does not exist." {
			return false
		}
	}

	return true
}

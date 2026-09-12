package service

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/modAttachment"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/types"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicAttachment"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"

	"io"
	"net/http"
	"path"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(BasicAttachmentService))
}

// BasicAttachmentService 附件上传
// @Description:
type BasicAttachmentService struct {
	sv          *repositoryBasic.BasicAttachmentRepository `autowire:"?"`
	FileService types.FileProvider                         `autowire:"?"`

	server configPg.Server `value:"${server}"`
}

// Upload 上传
func (c *BasicAttachmentService) Upload(ctx *gin.Context) (rt rg.Rs[modBasicAttachment.OkVo]) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return rt.ErrorMessage("上传文件不能为空")
	}

	f, err := file.Open()
	if err != nil {
		return rt.ErrorMessage("上传文件错误")
	}
	defer func() {
		// 打开的资源。一定要记住主动关闭
		_ = f.Close()
	}()

	log.Infof(ctx, log.TagAppDef, "file.Filename=%+v, file.Size = %+v", file.Filename, file.Size)
	atta, err := c.FileService.PutObject(f, modAttachment.PutFile(file.Filename, file.Size), nil)
	log.Infof(ctx, log.TagAppDef, "atta=%#v", atta)
	if err == nil {
		var vo modBasicAttachment.OkVo
		copier.Copy(&vo, &atta)
		vo.Domain = c.server.Domain
		vo.Url = vo.Domain + vo.Url
		return rg.OkData(vo)
	} else {
		return rt.ErrorMessage(err.Error())
	}
}

// UploadMore 多文件上传
func (c *BasicAttachmentService) UploadMore(ctx *gin.Context) (rt rg.Rs[map[int]modBasicAttachment.OkVo]) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return rt.ErrorMessage("上传文件不能为空")
	}
	headers := form.File["file"]
	if nil == headers || len(headers) <= 0 {
		return rt.ErrorMessage("上传文件不能为空")
	}
	data := make(map[int]modBasicAttachment.OkVo, 0)
	for i, file := range headers {
		f, err := file.Open()
		defer func() {
			// 打开的资源。一定要记住主动关闭
			_ = f.Close()
		}()
		if err != nil {
			data[i+1] = modBasicAttachment.OkVo{Error: file.Filename + " 上传文件错误:" + err.Error()}
			continue
		}
		log.Infof(ctx, log.TagAppDef, "file.Filename=%+v, file.Size = %+v", file.Filename, file.Size)
		atta, err := c.FileService.PutObject(f, modAttachment.PutFile(file.Filename, file.Size), nil)
		if err == nil {
			var vo modBasicAttachment.OkVo
			copier.Copy(&vo, &atta)
			vo.Domain = c.server.Domain
			vo.Url = vo.Domain + vo.Url
			data[i+1] = vo
		} else {
			data[i+1] = modBasicAttachment.OkVo{Error: file.Filename + " 上传文件错误: " + err.Error()}
		}
	}
	return rg.OkData(data)
}

// UploadLink 多url文件上传
func (c *BasicAttachmentService) UploadLink(ctx *gin.Context, ct modBasicAttachment.WebUrlCt) (rt rg.Rs[map[int]modBasicAttachment.OkVo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if nil == ct.Url || len(ct.Url) == 0 {
		return rt.ErrorMessage("上传文件不能为空")
	}
	data := make(map[int]modBasicAttachment.OkVo, 0)
	for i, url := range ct.Url {
		if "" == url {
			data[i+1] = modBasicAttachment.OkVo{Error: "地址错误：" + url}
			continue
		}
		res, err := http.Get(url)
		if err != nil {
			data[i+1] = modBasicAttachment.OkVo{Error: "远程下载错误：" + err.Error()}
			continue
		}
		// defer后的为延时操作，通常用来释放相关变量
		defer res.Body.Close()

		content, err := io.ReadAll(res.Body)
		if err != nil {
			data[i+1] = modBasicAttachment.OkVo{Error: "远程下载错误：" + err.Error()}
			continue
		}
		reader := bytes.NewReader(content)
		filename := path.Base(url)
		size := int64(len(content))
		log.Infof(ctx, log.TagAppDef, "file.Filename=%+v, file.Size = %+v", filename, size)
		atta, err := c.FileService.PutObject(reader, modAttachment.PutFile(filename, size), nil)
		if err == nil {
			var vo modBasicAttachment.OkVo
			copier.Copy(&vo, &atta)
			vo.Domain = c.server.Domain
			vo.Url = vo.Domain + vo.Url
			data[i+1] = vo
		} else {
			data[i+1] = modBasicAttachment.OkVo{Error: filename + " 上传文件错误: " + err.Error()}
		}
	}
	return rg.OkData(data)
}

// ListByOwner 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicAttachmentService) ListByOwner(ctx *gin.Context, ct modBasicAttachment.ListFileOwnerCt) (rt rg.Rs[modBasicAttachment.ListFIleOwnerVo]) {
	data := modBasicAttachment.ListFIleOwnerVo{}
	data.GroupData = make(map[string][]modBasicAttachment.Vo)
	data.Data = make([]modBasicAttachment.Vo, 0)
	//
	fileOwner := make([]string, 0)
	mapOwner := make(map[string]string)
	key := ""
	if nil != ct.GroupData {
		for _, item := range ct.GroupData {
			key = strings.TrimSpace(item.FileOwner)
			if strPg.IsNotBlank(key) && strPg.IsNotBlank(item.Group) {
				fileOwner = append(fileOwner, key)
				mapOwner[key] = strings.TrimSpace(item.Group)
			}
		}
	} else if nil != ct.FileOwner {
		for _, item := range ct.FileOwner {
			key = strings.TrimSpace(item)
			if strPg.IsNotBlank(key) {
				fileOwner = append(fileOwner, key)
			}
		}
	}
	if len(fileOwner) > 0 {
		r := c.sv
		var query entityBasic.BasicAttachmentEntity
		//
		query.State = enumStatePg.ENABLE.Index()
		//
		infos := r.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
			db = db.Order("create_at desc")
			db = db.Where("file_owner in ?", fileOwner)
			return db
		}))
		if nil != infos {
			//字段赋值
			for _, item := range infos {
				var vo modBasicAttachment.Vo
				copier.Copy(&vo, &item)
				vo.Url = item.Domain + item.File
				//分组
				if len(mapOwner) > 0 {
					if group, ok := mapOwner[item.FileOwner]; ok {
						if _, b := data.GroupData[group]; !b {
							data.GroupData[group] = make([]modBasicAttachment.Vo, 0)
						}
						data.GroupData[group] = append(data.GroupData[group], vo)
					}
				} else {
					data.Data = append(data.Data, vo)
				}
			}
		}
	}

	rt.Data = data
	return rt.Ok()
}

// DelByOwner 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicAttachmentService) DelByOwner(ctx *gin.Context, ct modBasicAttachment.DelFileOwnerCt) (rt rg.Rs[string]) {
	if strPg.IsBlank(ct.FileOwner) {
		return rt.ErrorMessage("文件拥有者参数错误")
	}
	if nil == ct.Nos || len(ct.Nos) < 1 {
		return rt.ErrorMessage("请选择要删除的文件")
	}
	nos := make([]string, 0)
	for _, item := range ct.Nos {
		item = strings.TrimSpace(item)
		if strPg.IsNotBlank(item) {
			nos = append(nos, item)
		}
	}
	if len(nos) < 1 {
		return rt.ErrorMessage("请选择要删除的文件")
	}
	if c.sv.Config().Data.Delete {
		c.sv.DeleteByIdAndFileOwner(ctx, nos, strings.TrimSpace(ct.FileOwner))
	} else {
		c.sv.UpdateByIdAndFileOwnerSetState13(ctx, nos, strings.TrimSpace(ct.FileOwner))
	}
	rt.Data = "删除成功"
	return rt.Ok()
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicAttachmentService) Query(ctx *gin.Context, ct modBasicAttachment.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicAttachment.Vo]]) {
	var query entityBasic.BasicAttachmentEntity
	copier.Copy(&query, &ct)
	r := c.sv
	slice := make([]modBasicAttachment.Vo, 0)
	rt.Data.Data = slice
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Order("create_at DESC")
		//自定义查询
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%").Where("source_name like ?", "%"+ct.Wd+"%")
		}
	}))
	if nil != err {
		return rt.Ok()
	}
	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicAttachment.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicAttachment.Vo
			copier.Copy(&vo, &item)
			vo.Url = item.Domain + item.File
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
		return rt.Ok()
	}
	return rt.Ok()
}

// MakeFileOwner 设置文件拥有者/file token/ 文件归属
func (c *BasicAttachmentService) MakeFileOwner(ctx *gin.Context, ct modBasicAttachment.MakeFileOwnerCt) (rt rg.Rs[modBasicAttachment.MakeFileOwner]) {
	vo := modBasicAttachment.MakeFileOwner{}
	prefix := ""
	if len(ct.Mark) > 0 {
		if ":auto" == ct.Mark {
			prefix = "default-"
		} else {
			str := strutil.Substring(ct.Mark, 0, 1)
			if str != ":" {
				prefix = strings.TrimSpace(ct.Mark) + "-"
			}
		}
	} else {
		prefix = "default-"
	}
	vo.FileOwner = noPg.MakeNo(prefix)
	rt.Data = vo
	return rt.Ok()
}

// MakeFileOwnerAll 设置文件拥有者/file token/ 文件归属
func (c *BasicAttachmentService) MakeFileOwnerAll(ctx *gin.Context, ct modBasicAttachment.MakeFileOwnerAllCt) (rt rg.Rs[[]modBasicAttachment.MakeFileOwner]) {
	data := make([]modBasicAttachment.MakeFileOwner, 0)
	if ct.Rule != nil && len(ct.Rule) > 0 {
		for _, item := range ct.Rule {
			vo := modBasicAttachment.MakeFileOwner{}
			prefix := ""
			if len(item.Mark) > 0 {
				if ":auto" == item.Mark {
					prefix = "default-"
				} else {
					str := strutil.Substring(item.Mark, 0, 1)
					if str != ":" {
						prefix = strings.TrimSpace(item.Mark) + "-"
					}
				}
			} else {
				prefix = "default-"
			}
			vo.FileOwner = noPg.MakeNo(prefix)
			vo.Mark = item.Mark
			//
			data = append(data, vo)
		}
	} else if ct.Num > 0 || ct.Num == 0 {
		num := ct.Num
		if num <= 0 {
			num = 1
		}
		prefix := ""
		for i := int32(0); i < num; i++ {
			vo := modBasicAttachment.MakeFileOwner{}
			prefix = ""
			if len(ct.Mark) > 0 {
				if ":auto" == ct.Mark {
					prefix = "default-"
				} else {
					str := strutil.Substring(ct.Mark, 0, 1)
					if str != ":" {
						prefix = strings.TrimSpace(ct.Mark) + "-"
					}
				}
			} else {
				prefix = "default-"
			}
			vo.FileOwner = noPg.MakeNo(prefix)
			vo.Mark = ct.Mark
			//
			data = append(data, vo)
		}
	}
	rt.Data = data
	return rt.Ok()
}

// UpdateByFileOwner 批量更新文件拥有者
func (c *BasicAttachmentService) UpdateByFileOwner(ctx *gin.Context, ct modBasicAttachment.UpdateByFileOwner) (rt rg.Rs[string]) {
	if ct.Data == nil || len(ct.Data) < 1 {
		return rt.ErrorMessage("请选择要更新的文件")
	}
	ids := make([]string, 0)
	mapOwner := make(map[string][]string)
	for key, item := range ct.Data {
		if nil != item && len(item) > 0 {
			mapOwner[key] = make([]string, 0)
			for _, str := range item {
				str = strings.TrimSpace(str)
				if strPg.IsNotBlank(str) {
					mapOwner[key] = append(mapOwner[key], str)
					ids = append(ids, str)
				}
			}
		}
	}
	//查询这个id，如果文件已经拥有者，那么复制一个新纪录
	if len(mapOwner) > 0 {
		data := make([]*entityBasic.BasicAttachmentEntity, 0)
		mapAtt := make(map[string]*entityBasic.BasicAttachmentEntity)
		if len(ids) > 0 {
			var query entityBasic.BasicAttachmentEntity
			query.State = enumStatePg.ENABLE.Index()
			//
			infos := c.sv.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
				db = db.Order("create_at desc")
				db = db.Where("id in ?", ids)
				return db
			}))
			if infos != nil && len(infos) > 0 {
				for _, item := range infos {
					//存在文件拥有者
					if strPg.IsNotBlank(item.FileOwner) {
						mapAtt[numberPg.Int64ToString(item.ID)] = item
					}
				}
			}
		}
		mapOwnerNew := make(map[string][]string)
		filter := make([]string, 0)
		for key, item := range mapOwner {
			if nil != item && len(item) > 0 {
				filter = make([]string, 0)
				for _, str := range item {
					if obj, ok := mapAtt[str]; ok {
						//存在文件拥有者
						obj.ID = 0
						obj.FileOwner = key
						obj.TypeData = "copy"
						data = append(data, obj)
					} else {
						//不存在文件拥有者
						filter = append(filter, str)
					}
				}
				if len(filter) > 0 {
					fmt.Printf("不存在文件拥有者:%+v\n", filter)
					mapOwnerNew[key] = filter
				}
			}
		}
		if len(data) > 0 {
			c.sv.DbModel().Create(data)
		}
		if len(mapOwnerNew) > 0 {
			for key, item := range mapOwnerNew {
				fmt.Printf("不存在文件拥有者:%+v\n", item)
				c.sv.UpdateByIdSetFileOwner(ctx, item, key)
			}
		}
		clear(filter)
		clear(data)
		clear(mapAtt)
		clear(mapOwnerNew)
	}
	return rt.Ok()
}

func (c *BasicAttachmentService) UpdateDetail(ctx *gin.Context, ct modBasicAttachment.DetailByFileOwnerCt) (rt rg.Rs[[]modBasicAttachment.Vo]) {
	if strPg.IsBlank(ct.FileOwner) {
		return rt.ErrorMessage("所有者key不能为空")
	}
	list := make([]modBasicAttachment.Vo, 0)
	entity := entityBasic.BasicAttachmentEntity{}
	entity.FileOwner = ct.FileOwner
	infos := c.sv.FindAll(ctx, entity)
	if infos != nil && len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicAttachment.Vo
			copier.Copy(&vo, &item)
			vo.Url = item.Domain + item.File
			list = append(list, vo)
		}
	}
	return rt.OkData(list)
}

// UpdateAddByFileOwner 添加拥有着
func (c *BasicAttachmentService) UpdateAddByFileOwner(ctx *gin.Context, ct modBasicAttachment.AddByFileOwnerCt) (rt rg.Rs[string]) {
	fileOwner := strings.TrimSpace(ct.FileOwner)
	if strPg.IsBlank(fileOwner) {
		return rt.ErrorMessage("所有者key不能为空")
	}
	if nil == ct.Nos || len(ct.Nos) < 1 {
		return rt.ErrorMessage("没有选择任何数据")
	}
	ids := make([]string, 0)
	for _, id := range ct.Nos {
		if strPg.IsNotBlank(id) {
			ids = append(ids, strings.TrimSpace(id))
		}
	}
	if len(ids) < 1 {
		return rt.ErrorMessage("没有选择任何数据")
	}
	list := make([]entityBasic.BasicAttachmentEntity, 0)
	infos := c.sv.FindAllOption(ctx, optionsPg.WithOption(func(db *optionsPg.OptionParams) {
		db.Db.Where("file_owner in ?", ids)
	}))
	if infos != nil && len(infos) > 0 {
		now := time.Now()
		for _, item := range infos {
			var vo entityBasic.BasicAttachmentEntity
			copier.Copy(&vo, &item)
			vo.Url = item.Domain + item.File
			vo.FileOwnerCopy = item.FileOwner
			vo.FileOwner = fileOwner
			vo.No = noPg.No()
			vo.ID = 0
			vo.CreateAt = &now
			vo.UpdateAt = &now
			if strPg.IsNotBlank(ct.FileOwnerSub) {
				vo.FileOwnerSub = strings.TrimSpace(ct.FileOwnerSub)
			} else {
				vo.FileOwnerSub = "main"
			}
			vo.TypeData = "copy"
			//
			list = append(list, vo)
		}
		err, _ := c.sv.CreateBatch(ctx, list)
		if err != nil {
			return rt.ErrorMessage("保存失败")
		}
	}
	return rt.Ok()
}

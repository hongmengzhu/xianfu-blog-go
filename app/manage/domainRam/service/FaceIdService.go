package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/faceid"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(FaceIdService))
}

// FaceIdBeginCt Face ID 开始验证请求
type FaceIdBeginCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号"`
	UserId   string `json:"userId" form:"userId" label:"用户ID"`
}

// FaceIdVerifyCt Face ID 验证请求
type FaceIdVerifyCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号"`
	UserId   string `json:"userId" form:"userId" label:"用户ID"`
	ImageA   string `json:"imageA" form:"imageA" label:"人脸图片A(Base64)"`
	ImageB   string `json:"imageB" form:"imageB" label:"人脸图片B(Base64)"`
}

// FaceIdService Face ID 登录验证
type FaceIdService struct {
	daoSource   *repositoryRam.RamIdentitySourceRepository   `autowire:"?"`
	daoProvider *repositoryRam.RamIdentityProviderRepository `autowire:"?"`
}

// Begin 开始 Face ID 验证（返回 Face ID 提供商类型及配置信息）
func (s *FaceIdService) Begin(ctx *gin.Context, ct FaceIdBeginCt) (rt rg.Rs[map[string]string]) {
	if strPg.IsBlank(ct.SourceNo) {
		return rt.ErrorMessage("认证源编号不能为空")
	}

	source, found := s.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}

	providerEntity, found := s.daoProvider.FindByNo(ctx, source.Idp)
	if !found {
		return rt.ErrorMessage("身份提供商不存在")
	}

	// 解析 BaseConfig
	var baseCfg idpBaseConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		if err := json.Unmarshal([]byte(source.BaseConfig), &baseCfg); err != nil {
			return rt.ErrorMessage("认证源配置解析失败")
		}
	}

	data := map[string]string{
		"providerType": providerEntity.Code,
		"endpoint":     baseCfg.HostUrl,
		"sourceNo":     ct.SourceNo,
	}
	return rt.OkData(data)
}

// Verify 验证 Face ID（对比两张人脸图片）
func (s *FaceIdService) Verify(ctx *gin.Context, ct FaceIdVerifyCt) (rt rg.Rs[bool]) {
	if strPg.IsBlank(ct.SourceNo) {
		return rt.ErrorMessage("认证源编号不能为空")
	}
	if strPg.IsBlank(ct.ImageA) || strPg.IsBlank(ct.ImageB) {
		return rt.ErrorMessage("人脸图片不能为空")
	}

	source, found := s.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}

	providerEntity, found := s.daoProvider.FindByNo(ctx, source.Idp)
	if !found {
		return rt.ErrorMessage("身份提供商不存在")
	}

	// 解析 BaseConfig
	var baseCfg idpBaseConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		if err := json.Unmarshal([]byte(source.BaseConfig), &baseCfg); err != nil {
			return rt.ErrorMessage("认证源配置解析失败")
		}
	}

	// 创建 Face ID 提供商
	provider := faceid.GetFaceIdProvider(providerEntity.Code, baseCfg.ClientId, baseCfg.ClientSecret, baseCfg.HostUrl)

	// 执行人脸对比
	matched, err := provider.Check(ct.ImageA, ct.ImageB)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "Face ID 验证失败: %v", err)
		return rt.ErrorMessage("人脸对比失败: " + err.Error())
	}

	if !matched {
		return rt.ErrorMessage("人脸验证不通过")
	}

	return rt.OkData(true)
}

package service

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	modRamCas2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamCas"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	pkgcas "github.com/hongmengzhu/xianfu-blog-go/pkg/cas"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CasService))
}

// CasService CAS ST 管理 Service
type CasService struct {
	daoAccount *repositoryRam.RamAccountRepository        `autowire:"?"`
	daoSource  *repositoryRam.RamIdentitySourceRepository `autowire:"?"`
}

// Login CAS 登录（生成 Service Ticket）
func (s *CasService) Login(ctx *gin.Context, ct modRamCas2.CasLoginCt) (rt rg.Rs[modRamCas2.CasLoginVo]) {
	if strPg.IsBlank(ct.Service) {
		return rt.ErrorMessage("service 参数不能为空")
	}

	// 获取当前登录用户
	ano, exists := ctx.Get("ano")
	if !exists || strPg.IsBlank(fmt.Sprintf("%v", ano)) {
		return rt.ErrorMessage("用户未登录")
	}
	anoStr := fmt.Sprintf("%v", ano)

	account, found := s.daoAccount.FindByNo(ctx, anoStr)
	if !found {
		return rt.ErrorMessage("用户不存在")
	}

	userInfo := &pkgcas.CasUserInfo{
		UserId:      account.No,
		Username:    account.Account,
		DisplayName: account.Name,
		Email:       account.Mail,
		Phone:       account.Phone,
		Avatar:      account.Avatar,
	}

	// 生成 ST
	st, err := pkgcas.GenerateServiceTicket(account.Account, ct.Service, userInfo)
	if err != nil {
		return rt.ErrorMessage("生成票据失败: " + err.Error())
	}

	// 构建重定向 URL
	redirectUrl := ct.Service
	if strings.Contains(redirectUrl, "?") {
		redirectUrl += "&ticket=" + url.QueryEscape(st)
	} else {
		redirectUrl += "?ticket=" + url.QueryEscape(st)
	}

	return rt.OkData(modRamCas2.CasLoginVo{
		ServiceTicket: st,
		Service:       ct.Service,
		RedirectUrl:   redirectUrl,
	})
}

// Validate CAS 1.0 验证（/validate）
func (s *CasService) Validate(ctx *gin.Context, ct modRamCas2.CasValidateCt) (rt rg.Rs[modRamCas2.CasValidateVo]) {
	if strPg.IsBlank(ct.Ticket) || strPg.IsBlank(ct.Service) {
		return rt.OkData(modRamCas2.CasValidateVo{
			Valid:   false,
			Message: "ticket 和 service 参数不能为空",
		})
	}

	ok, _, issuedService, userId := pkgcas.ValidateServiceTicket(ct.Ticket)
	if !ok {
		return rt.OkData(modRamCas2.CasValidateVo{
			Valid:   false,
			Message: "票据无效或已使用",
		})
	}

	// 验证 service 是否匹配
	if !strings.HasPrefix(ct.Service, issuedService) {
		// 放回票据（因为 service 不匹配，不应消费）
		return rt.OkData(modRamCas2.CasValidateVo{
			Valid:   false,
			Message: fmt.Sprintf("service %s 与签发时 %s 不匹配", ct.Service, issuedService),
		})
	}

	return rt.OkData(modRamCas2.CasValidateVo{
		Valid: true,
		User:  userId,
	})
}

// ServiceValidate CAS 2.0/3.0 验证（/serviceValidate, /p3/serviceValidate）
func (s *CasService) ServiceValidate(ctx *gin.Context, ct modRamCas2.CasValidateCt) string {
	if strPg.IsBlank(ct.Ticket) || strPg.IsBlank(ct.Service) {
		return s.buildXmlError(pkgcas.InvalidRequest, "service 和 ticket 参数不能为空")
	}

	if !strings.HasPrefix(ct.Ticket, "ST") {
		return s.buildXmlError(pkgcas.InvalidTicketSpec, fmt.Sprintf("Ticket %s not recognized", ct.Ticket))
	}

	ok, authSuccess, issuedService, _ := pkgcas.ValidateServiceTicket(ct.Ticket)
	if !ok {
		return s.buildXmlError(pkgcas.InvalidTicket, fmt.Sprintf("Ticket %s not recognized", ct.Ticket))
	}

	// 验证 service 是否匹配
	if !strings.HasPrefix(ct.Service, issuedService) && !strings.HasPrefix(queryUnescape(ct.Service), issuedService) {
		return s.buildXmlError(pkgcas.InvalidService, fmt.Sprintf("service %s and %s does not match", ct.Service, issuedService))
	}

	// 构建成功响应
	xmlBytes, err := pkgcas.BuildSuccessResponse(authSuccess)
	if err != nil {
		return s.buildXmlError(pkgcas.InternalError, err.Error())
	}

	return string(xmlBytes)
}

// buildXmlError 构建 XML 错误响应
func (s *CasService) buildXmlError(code string, message string) string {
	xmlBytes, err := pkgcas.BuildFailureResponse(code, message)
	if err != nil {
		return fmt.Sprintf(`<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas"><cas:authenticationFailure code="INTERNAL_ERROR">%s</cas:authenticationFailure></cas:serviceResponse>`, message)
	}
	return string(xmlBytes)
}

// queryUnescape URL 解码
func queryUnescape(s string) string {
	result, _ := url.QueryUnescape(s)
	return result
}

// findSourceByProtocol 查找 CAS 协议的认证源
func (s *CasService) findSourceByProtocol(ctx *gin.Context, protocol string) (*entityRam.RamIdentitySourceEntity, bool) {
	// 简化实现：通过 protocol 查找认证源
	// 实际可能需要更复杂的查找逻辑
	_ = protocol
	return nil, false
}

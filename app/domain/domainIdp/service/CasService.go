package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/idp/modIdpCas"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityIdp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryIdp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	pkgcas "github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/idpPg/cas"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IdpCasService))
}

// IdpCasService CAS ST 管理 Service
type IdpCasService struct {
	daoAccount *repositoryRam.RamAccountRepository        `autowire:"?"`
	daoSource  *repositoryIdp.IdpIdentitySourceRepository `autowire:"?"`
}

// Login CAS 登录（生成 Service Ticket）
func (s *IdpCasService) Login(ctx *gin.Context, ct modIdpCas.CasLoginCt) (rt rg.Rs[modIdpCas.CasLoginVo]) {
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
		return rt.ErrorMessage("生成票据失败")
	}

	// 构建重定向 URL
	redirectUrl := ct.Service
	if strings.Contains(redirectUrl, "?") {
		redirectUrl += "&ticket=" + url.QueryEscape(st)
	} else {
		redirectUrl += "?ticket=" + url.QueryEscape(st)
	}

	return rt.OkData(modIdpCas.CasLoginVo{
		ServiceTicket: st,
		Service:       ct.Service,
		RedirectUrl:   redirectUrl,
	})
}

// Validate CAS 1.0 验证（/validate）
func (s *IdpCasService) Validate(ctx *gin.Context, ct modIdpCas.CasValidateCt) (rt rg.Rs[modIdpCas.CasValidateVo]) {
	if strPg.IsBlank(ct.Ticket) || strPg.IsBlank(ct.Service) {
		return rt.OkData(modIdpCas.CasValidateVo{
			Valid:   false,
			Message: "ticket 和 service 参数不能为空",
		})
	}

	ok, _, issuedService, userId := pkgcas.ValidateServiceTicket(ct.Ticket)
	if !ok {
		return rt.OkData(modIdpCas.CasValidateVo{
			Valid:   false,
			Message: "票据无效或已使用",
		})
	}

	// 验证 service 必须精确匹配
	if !sameCasService(ct.Service, issuedService) {
		return rt.OkData(modIdpCas.CasValidateVo{
			Valid:   false,
			Message: "service 与签发时不匹配",
		})
	}

	return rt.OkData(modIdpCas.CasValidateVo{
		Valid: true,
		User:  userId,
	})
}

// ServiceValidate CAS 2.0/3.0 验证（/serviceValidate, /p3/serviceValidate）
func (s *IdpCasService) ServiceValidate(ctx *gin.Context, ct modIdpCas.CasValidateCt) string {
	if strPg.IsBlank(ct.Ticket) || strPg.IsBlank(ct.Service) {
		return s.buildXmlError(pkgcas.InvalidRequest, "service 和 ticket 参数不能为空")
	}

	if !validCasTicket(ct.Ticket) {
		return s.buildXmlError(pkgcas.InvalidTicketSpec, "INVALID_TICKET_SPEC")
	}

	ok, authSuccess, issuedService, _ := pkgcas.ValidateServiceTicket(ct.Ticket)
	if !ok {
		return s.buildXmlError(pkgcas.InvalidTicket, "INVALID_TICKET")
	}

	// 验证 service 必须精确匹配，避免相似域名/前缀绕过
	if !sameCasService(ct.Service, issuedService) {
		return s.buildXmlError(pkgcas.InvalidService, "INVALID_SERVICE")
	}

	// 构建成功响应
	xmlBytes, err := pkgcas.BuildSuccessResponse(authSuccess)
	if err != nil {
		return s.buildXmlError(pkgcas.InternalError, "internal error")
	}

	return string(xmlBytes)
}

// buildXmlError 构建 XML 错误响应
func (s *IdpCasService) buildXmlError(code string, message string) string {
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
func (s *IdpCasService) findSourceByProtocol(ctx *gin.Context, protocol string) (*entityIdp.IdpIdentitySourceEntity, bool) {
	// 简化实现：通过 protocol 查找认证源
	// 实际可能需要更复杂的查找逻辑
	_ = protocol
	return nil, false
}

var casTicketPattern = regexp.MustCompile(`^ST-[A-Za-z0-9]{8,64}$`)

func validCasTicket(ticket string) bool {
	return casTicketPattern.MatchString(strings.TrimSpace(ticket))
}

// sameCasService 对 service URL 做规范化后精确比较。
func sameCasService(a, b string) bool {
	na := normalizeCasService(a)
	nb := normalizeCasService(b)
	return na != "" && na == nb
}

func normalizeCasService(raw string) string {
	raw = queryUnescape(strings.TrimSpace(raw))
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	u.Fragment = ""
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	return strings.TrimSuffix(u.String(), "/")
}

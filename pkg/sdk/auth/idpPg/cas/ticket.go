package cas

import (
	"encoding/xml"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// CAS 协议常量
const (
	InvalidRequest    = "INVALID_REQUEST"
	InvalidTicketSpec = "INVALID_TICKET_SPEC"
	InvalidTicket     = "INVALID_TICKET"
	InvalidService    = "INVALID_SERVICE"
	InternalError     = "INTERNAL_ERROR"
)

// CasServiceResponse CAS 服务响应 XML 根结构
type CasServiceResponse struct {
	XMLName xml.Name                  `xml:"cas:serviceResponse"`
	Xmlns   string                    `xml:"xmlns:cas,attr"`
	Failure *CasAuthenticationFailure `xml:"cas:authenticationFailure,omitempty"`
	Success *CasAuthenticationSuccess `xml:"cas:authenticationSuccess,omitempty"`
}

// CasAuthenticationFailure CAS 认证失败
type CasAuthenticationFailure struct {
	XMLName xml.Name `xml:"cas:authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",innerxml"`
}

// CasAuthenticationSuccess CAS 认证成功
type CasAuthenticationSuccess struct {
	XMLName    xml.Name       `xml:"cas:authenticationSuccess"`
	User       string         `xml:"cas:user"`
	Attributes *CasAttributes `xml:"cas:attributes,omitempty"`
}

// CasAttributes CAS 用户属性
type CasAttributes struct {
	XMLName            xml.Name  `xml:"cas:attributes"`
	AuthenticationDate time.Time `xml:"cas:authenticationDate"`
	FirstName          string    `xml:"cas:firstName,omitempty"`
	LastName           string    `xml:"cas:lastName,omitempty"`
	Email              string    `xml:"cas:email,omitempty"`
	DisplayName        string    `xml:"cas:displayName,omitempty"`
	Phone              string    `xml:"cas:phone,omitempty"`
	Avatar             string    `xml:"cas:avatar,omitempty"`
	MemberOf           []string  `xml:"cas:memberOf,omitempty"`
}

// CasTicketWrapper CAS 票据包装（存储 ST 与用户/服务的映射）
type CasTicketWrapper struct {
	AuthenticationSuccess *CasAuthenticationSuccess
	Service               string
	UserId                string
	CreatedAt             time.Time
}

var (
	stToTicket sync.Map
)

// GenerateServiceTicket 生成 CAS Service Ticket (ST)
func GenerateServiceTicket(userId string, service string, user *CasUserInfo) (string, error) {
	if userId == "" {
		return "", fmt.Errorf("userId 不能为空")
	}
	if service == "" {
		return "", fmt.Errorf("service 不能为空")
	}

	authSuccess := &CasAuthenticationSuccess{
		User: userId,
		Attributes: &CasAttributes{
			AuthenticationDate: time.Now(),
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Email:              user.Email,
			DisplayName:        user.DisplayName,
			Phone:              user.Phone,
			Avatar:             user.Avatar,
		},
	}

	st := fmt.Sprintf("ST-%d", rand.Int63n(math.MaxInt64))
	stToTicket.Store(st, &CasTicketWrapper{
		AuthenticationSuccess: authSuccess,
		Service:               service,
		UserId:                userId,
		CreatedAt:             time.Now(),
	})

	return st, nil
}

// ValidateServiceTicket 验证 CAS Service Ticket
// 返回 (是否有效, 票据信息, 票据对应的 service, userId)
func ValidateServiceTicket(ticket string) (bool, *CasAuthenticationSuccess, string, string) {
	if wrapper, ok := stToTicket.LoadAndDelete(ticket); ok {
		tw := wrapper.(*CasTicketWrapper)
		return true, tw.AuthenticationSuccess, tw.Service, tw.UserId
	}
	return false, nil, "", ""
}

// GetServiceTicket 获取 ST 信息（不删除）
func GetServiceTicket(ticket string) (bool, *CasTicketWrapper) {
	if wrapper, ok := stToTicket.Load(ticket); ok {
		return true, wrapper.(*CasTicketWrapper)
	}
	return false, nil
}

// BuildSuccessResponse 构建认证成功 XML 响应
func BuildSuccessResponse(authSuccess *CasAuthenticationSuccess) ([]byte, error) {
	resp := CasServiceResponse{
		Xmlns:   "http://www.yale.edu/tp/cas",
		Success: authSuccess,
	}
	return xml.MarshalIndent(resp, "", "  ")
}

// BuildFailureResponse 构建认证失败 XML 响应
func BuildFailureResponse(code string, message string) ([]byte, error) {
	resp := CasServiceResponse{
		Xmlns: "http://www.yale.edu/tp/cas",
		Failure: &CasAuthenticationFailure{
			Code:    code,
			Message: message,
		},
	}
	return xml.MarshalIndent(resp, "", "  ")
}

// CasUserInfo 用户信息（用于生成 CAS 票据）
type CasUserInfo struct {
	UserId      string
	Username    string
	FirstName   string
	LastName    string
	Email       string
	DisplayName string
	Phone       string
	Avatar      string
	MemberOf    []string
}

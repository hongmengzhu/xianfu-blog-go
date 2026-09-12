package authTokenPg

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
)

// MakePaseToken
//
//	@Description: 生成 token
//	@param param
//	@param jwtC
//	@return rt
func MakePaseToken(param Param, jwtC pg.JwtConfig) (rt rg.Rs[Result]) {
	// Issuer: jwt签发者
	// Subject: jwt所面向的用户
	// Audience: 接收jwt的一方
	// Expiration: jwt的过期时间，这个过期时间必须要大于签发时间
	// NotBefore: 定义在什么时间之前，该jwt都是不可用的.
	// IssuedAt: jwt的签发时间
	// jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
	now := time.Now()
	// 默认2天
	minute := 60 * 24 * 2
	if jwtC.Expire > 0 {
		minute = jwtC.Expire
	}
	token := paseto.NewToken()
	token.SetIssuer(jwtC.Issuer)
	token.SetIssuedAt(now)
	token.SetSubject(param.LoginNo)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(time.Minute * time.Duration(minute)))
	token.SetJti(param.UniqueId)
	token.SetAudience(jwtC.Audience)
	//版本号
	if strPg.IsNotBlank(param.Version) {
		token.SetString(Version, param.Version)
	}
	//登录用户 名
	if strPg.IsNotBlank(param.LoginUserName) {
		token.SetString(LoginUserName, param.LoginUserName)
	}
	//登录用户 名
	if strPg.IsNotBlank(param.Name) {
		token.SetString(Name, param.Name)
	}
	//租户
	if strPg.IsNotBlank(param.TenantNo) {
		token.SetString(TenantNo, param.TenantNo)
	}
	//组织
	if strPg.IsNotBlank(param.OrgNo) {
		token.SetString(OrgNo, param.OrgNo)
	}
	//类型
	if strPg.IsNotBlank(param.Type) {
		token.SetString(Type, param.Type)
	}
	// 可选的 footer
	//token.SetFooter([]byte("foxwho.com"))
	if strPg.IsBlank(param.Result.PrivateKey) {
		return rt.ErrorMessage("私钥不能为空")
	}
	if strPg.IsBlank(param.Result.PublicKey) {
		return rt.ErrorMessage("公钥不能为空")
	}
	// 生成密钥对
	result := Result{PrivateKey: param.Result.PrivateKey, PublicKey: param.Result.PublicKey}
	hex, err := paseto.NewV4AsymmetricSecretKeyFromHex(param.Result.PrivateKey)
	if err != nil {

		return rt.ErrorMessage("生成密钥对错误")
	}
	// 签名 token
	signed := token.V4Sign(hex, nil)
	result.Token = AuthScheme + " " + signed
	return rt.OkData(result)
}

// VerifyByPublicKey
//
//	@Description: 验证
//	@param pubKey
//	@param signed
//	@return t
//	@return rt
func VerifyByPublicKey(pubKey, signed string) (t *paseto.Token, rt rg.Rs[string]) {
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(pubKey)
	if nil != err {
		return nil, rt.ErrorMessage("公钥转换失败")
	}
	parser := paseto.NewParserWithoutExpiryCheck()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	token, err := parser.ParseV4Public(publicKey, signed, nil)
	if nil != err {
		log.Errorf(context.Background(), log.TagAppDef, "验证失败= %+v", err)
		return nil, rt.ErrorMessage("验证失败")
	}
	return token, rt.Ok()
}

// MakePublicPrivateKey
//
//	@Description: 生成密钥对
//	@return result
func MakePublicPrivateKey() (result Result) {
	// 生成密钥对
	privateKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := privateKey.Public()
	result.PrivateKey = privateKey.ExportHex()
	result.PublicKey = publicKey.ExportHex()
	return result
}

func deconstructToken(token string) (string, bool) {
	parts := strings.Split(token, ".")
	partsLen := len(parts)
	if partsLen != 3 && partsLen != 4 {
		return "", false
	}
	return parts[2], true
}

var b64 = base64.RawURLEncoding.Strict()

// Encode Standard encoding for Paseto is URL safe base64 with no padding
func Encode(bytes []byte) string {
	return b64.EncodeToString(bytes)
}

// Decode Standard decoding for Paseto is URL safe base64 with no padding
func Decode(encoded string) ([]byte, bool) {
	// From: https://pkg.go.dev/encoding/base64#Encoding.Strict
	// Note that the input is still malleable, as new line characters (CR and LF) are still ignored.
	if strings.ContainsAny(encoded, "\n\r") {
		return nil, false
	}
	if b, err := b64.DecodeString(encoded); err != nil {
		return nil, false
	} else {
		return b, true
	}
}

// ParseUnverified
//
//	@Description: 获取 载荷内容，不验证签名
//	@param token
//	@return []byte
//	@return bool
func ParseUnverified(token string) ([]byte, bool) {
	str, ok := deconstructToken(token)
	if !ok {
		return nil, false
	}
	bytes, b := Decode(str)
	if !b {
		return nil, false
	}
	signatureOffset := len(bytes) - 64
	message := make([]byte, len(bytes)-64)
	copy(message, bytes[:signatureOffset])
	return message, true
}

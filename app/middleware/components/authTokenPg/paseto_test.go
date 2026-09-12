package authTokenPg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"go-spring.org/log"
)

func TestToken(t *testing.T) {
	privatePubKey := MakePublicPrivateKey()
	param := Param{}
	param.UniqueId = noPg.No()
	param.LoginNo = "1"
	param.LoginUserName = "1:name"
	param.No = param.LoginNo
	param.Version = "000001"
	param.TenantNo = "x00001"
	param.Type = "system"
	param.Result = privatePubKey
	//
	config := pg.JwtConfig{}
	config.Expire = 3600
	config.Issuer = "fox"
	config.Audience = "fox2222"
	//
	rt := MakePaseToken(param, config)
	if rt.SuccessIs() {
		t.Log(rt.Data)
		fmt.Printf("%+v\n", rt.Data)
		signed := rt.Data.Token

		signed = strings.Replace(signed, AuthScheme+" ", "", -1)

		publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(privatePubKey.PublicKey)
		if nil != err {
			fmt.Printf("err=%+v\n", err)
		}
		parser := paseto.NewParserWithoutExpiryCheck()
		parser.AddRule(paseto.NotExpired())
		parser.AddRule(paseto.ValidAt(time.Now()))
		token, err := parser.ParseV4Public(publicKey, signed, nil)
		if nil != err {
			log.Errorf(context.Background(), log.TagAppDef, "验证失败= %+v", err)
		} else {
			log.Infof(context.Background(), log.TagAppDef, "验证成功= %+v", token)
		}
	} else {
		t.Error(rt.Message)
	}
}

func TestTokenBody(t *testing.T) {
	tokenx := "v4.public.eyJhdWQiOiJmb3gyMjIyIiwiZXhwIjoiMjAyNS0wNy0xMVQwNDoyMjoxOCswODowMCIsImlhdCI6IjIwMjUtMDctMDhUMTY6MjI6MTgrMDg6MDAiLCJpc3MiOiJmb3giLCJqdGkiOiIyNTA3MDgxNjIyMTg1NjczODY0IiwibG9nVW4iOiIxOm5hbWUiLCJuYmYiOiIyMDI1LTA3LTA4VDE2OjIyOjE4KzA4OjAwIiwic3ViIjoiMSIsInRubyI6IngwMDAwMSIsInRwIjoic3lzdGVtIiwidmVyIjoiMDAwMDAxIn0m_e2O8ZzcRA_Uw9NeRm-3jQ59iES41LPJiefA_g51OVeWzhAiGeHyqghKHuXeBoSbaQwfsbKq376lJ-yakqIJ.Zm94d2hvLmNvbQ"
	fmt.Printf("tokenx=%+v\n", tokenx)
	unverified, b := ParseUnverified(tokenx)
	if b {
		fmt.Printf("unverified=%+v\n", unverified)
		fmt.Printf("unverified=%+v\n", string(unverified))
		// 解析为map便于查看
		var payload map[string]interface{}
		if err := json.Unmarshal(unverified, &payload); err != nil {
			fmt.Printf("解析载荷JSON失败：%v\n", err)
			return
		}
		//
		//// 输出结果
		fmt.Println("PASETO 载荷内容：")
		for key, value := range payload {
			fmt.Printf("%-8s: %v\n", key, value)
		}
	}
}

func TestTokenVerify(t *testing.T) {
	signed := "7278a98adb8beed32a3f4ebb00a1c0342f3a9156dfb20f14bc5561c9318461ba"
	str := "v4.public.eyJhdWQiOiJwY1dlYiIsImV4cCI6IjIwMjctMDYtMDdUMDI6MTA6NDArMDg6MDAiLCJpYXQiOiIyMDI1LTA3LTE4VDE0OjEwOjQwKzA4OjAwIiwiaXNzIjoiZm94d2hvLmNvbSIsImp0aSI6IjI1MDcxODE0MTA0MDM0NDM5ODk3OTYiLCJuYW1lIjoiZGVtbyIsIm5iZiI6IjIwMjUtMDctMThUMTQ6MTA6NDArMDg6MDAiLCJzdWIiOiJjdXN0b21lciIsInRubyI6IjEwMDAwIiwidHAiOiJjdXN0b21lciJ9V_-YMrzcAaB4fMgIoszN3xBOzQcHY8-dQ4JFVMU3JLd9BVgn0VMJY_KcJu0f7WGjKQ_CcEfKBxvxcfXXYpLhCA.Zm94d2hvLmNvbQ"
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(str)
	if nil != err {
		fmt.Printf("err=%+v\n", err)
	}
	parser := paseto.NewParserWithoutExpiryCheck()
	//parser.AddRule(paseto.NotExpired())
	//parser.AddRule(paseto.ValidAt(time.Now()))
	token, err := parser.ParseV4Public(publicKey, signed, nil)
	if nil != err {
		log.Errorf(context.Background(), log.TagAppDef, "验证失败= %+v", err)
	} else {
		log.Infof(context.Background(), log.TagAppDef, "验证成功= %+v", token)
	}
}

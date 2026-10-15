package cryptoPg

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"go-spring.org/log"
)

type SM2Util struct {
	PrivateHex string
	PublicHex  string
}

func Sm2GenerateKey() (rt rg.Rs[model.AuthPubPriveDto]) {
	dto := model.AuthPubPriveDto{}
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(context.Background(), log.TagAppDef, "生成密钥对失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("生成密钥对失败")
	}
	dto.PrivateKey = x509.WritePrivateKeyToHex(privateKey)
	dto.PublicKey = x509.WritePublicKeyToHex(&privateKey.PublicKey)
	return rt.OkData(dto)
}

func NewSm2Default() *SM2Util {
	return &SM2Util{}
}

func NewSm2(private, public string) *SM2Util {
	return &SM2Util{
		PublicHex:  public,
		PrivateHex: private,
	}
}
func NewSm2WithHex(private, public string) *SM2Util {
	return &SM2Util{
		PublicHex:  public,
		PrivateHex: private,
	}
}

// DecodeHex 解密
// 私钥解密
func (c *SM2Util) DecodeHex(text string) (rt rg.Rs[string]) {

	ciphertext, err2 := hex.DecodeString(text)
	if nil != err2 {
		log.Fatal(context.Background(), log.TagAppDef, "hex解析失败", log.Any("err", err2.Error()))
		return rt.ErrorMessage("hex解析失败")
	}
	privateKey, err := x509.ReadPrivateKeyFromHex(c.PrivateHex)
	if err != nil {
		log.Fatal(context.Background(), log.TagAppDef, "解析私钥失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("解析私钥失败")
	}

	// 4. 使用加载的私钥解密
	decryptedData, err := sm2.Decrypt(privateKey, ciphertext, 0)
	if err != nil {
		log.Fatal(context.Background(), log.TagAppDef, "解密失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("解密失败")
	}
	return rt.OkData(string(decryptedData))
}

// EncryptHex 加密
// 公钥加密
func (c *SM2Util) EncryptHex(plainText []byte) (rt rg.Rs[string]) {
	publicKey, err := x509.ReadPublicKeyFromHex(c.PublicHex)
	if err != nil {
		log.Fatal(context.Background(), log.TagAppDef, "解析公钥失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("解析公钥失败")
	}
	ciphertext, err := sm2.Encrypt(publicKey, plainText, rand.Reader, 0)
	if err != nil {
		log.Fatal(context.Background(), log.TagAppDef, "加密失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("加密失败")
	}
	return rt.OkData(hex.EncodeToString(ciphertext))
}

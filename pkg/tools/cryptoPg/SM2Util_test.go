package cryptoPg

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"go-spring.org/log"
)

func Test_SM2Util_001(t *testing.T) {
	// 1. 生成密钥对
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(context.Background(), log.TagBizDef, "%+v", err)
	}
	publicKey := &privateKey.PublicKey

	// 2. 转换为 Hex 字符串
	privHex := x509.WritePrivateKeyToHex(privateKey)
	pubHex := x509.WritePublicKeyToHex(publicKey)

	fmt.Println("私钥 (Hex):", privHex)
	fmt.Println("公钥 (Hex):", pubHex)

	// 3. 从 Hex 字符串恢复（演示）
	restoredPriv, _ := x509.ReadPrivateKeyFromHex(privHex)
	restoredPub, _ := x509.ReadPublicKeyFromHex(pubHex)

	// 验证恢复的密钥是否可用
	msg := []byte("test")
	cipher, _ := sm2.Encrypt(restoredPub, msg, rand.Reader, 0)
	plain, _ := sm2.Decrypt(restoredPriv, cipher, 0)
	fmt.Println("加解密验证成功:", string(plain))
}

func Test_SM2Util_002(t *testing.T) {
	// 1. 生成密钥对

	// 2. 转换为 Hex 字符串
	privHex := "00e53929453aec25ba93a6fc3ccc18ea0c92a0bb008d1c7b196d7cb8758a6b588a"
	pubHex := "0435053a7a91b2fc473405f475ac66cd12cf9201f82630d96cd983ce1e4200dfb31a50dca2f3428bbc55c8d2de96668c3dc82ee3427ff6d7b53310792b7585c85d"

	fmt.Println("私钥 (Hex):", privHex)
	fmt.Println("公钥 (Hex):", pubHex)

	// 3. 从 Hex 字符串恢复（演示）
	restoredPriv, _ := x509.ReadPrivateKeyFromHex(privHex)
	restoredPub, _ := x509.ReadPublicKeyFromHex(pubHex)

	// 验证恢复的密钥是否可用
	msg := []byte("test")
	cipher, _ := sm2.Encrypt(restoredPub, msg, rand.Reader, 0)
	plain, _ := sm2.Decrypt(restoredPriv, cipher, 0)
	cipherText := base64.StdEncoding.EncodeToString(cipher)
	fmt.Println("密文(Base64):", cipherText)
	toString := hex.EncodeToString(cipher)
	fmt.Println("密文(hex):", toString)
	fmt.Println("加解密验证成功:", string(plain))
	//
	str := "04b67a5adaafb533194e1ef3f8ea1664f23ff17744868e29b27fe6b6d235643e75ddf4a7ce8b311dfc102cad65b7e205d6aaaef5a2433d784830b874f616cdc46729a895cf4cf42ece4f942235594ee514b5bf3f84ccc65d9e78e907d356b53a410169f58570f0f91a8d9c"
	decodeString, err := hex.DecodeString(str)
	if nil != err {
		fmt.Println("失败:", err)
	}
	plain2, _ := sm2.Decrypt(restoredPriv, decodeString, 0)
	fmt.Println("加解密验证成功2:", string(plain2))
}

func Test_SM2Util_003(t *testing.T) {
	privHex := "00e53929453aec25ba93a6fc3ccc18ea0c92a0bb008d1c7b196d7cb8758a6b588a"
	pubHex := "0435053a7a91b2fc473405f475ac66cd12cf9201f82630d96cd983ce1e4200dfb31a50dca2f3428bbc55c8d2de96668c3dc82ee3427ff6d7b53310792b7585c85d"
	newSm2 := NewSm2(privHex, pubHex)
	msg := []byte("test")
	rt := newSm2.EncryptHex(msg)
	fmt.Println(rt)
	if rt.SuccessIs() {
		decodeHex := newSm2.DecodeHex(rt.Data)
		fmt.Println(decodeHex)
	}
}

func Test_SM2Util_004(t *testing.T) {
	privHex := "4e2d631271b4858786ea5afb2c068b2d79a9b9193cf2777aa25c4fcfe9424477"
	pubHex := "04adde4b0f15fd326d9db386a97c317289f8e60e571da0e3bc286b5c1f919baae6e94978ee12d330a106f0d841103dd11e0f309ece5001536fbbafd13f04eca036"
	newSm2 := NewSm2(privHex, pubHex)
	msg := "04e17517f4dfff10fa4a68e9213547a2ad0b9937e88081c82899278c832ae5d4525a02f95393560b708da578166b7159fd829549a407465b0e92e985361918b470823b40313f44c0c0aedb9627b0184879d0ab28b55bf18bddeaa8e9889331741debd1bbdce3edd7f6600e"
	decodeHex := newSm2.DecodeHex(msg)
	fmt.Println(decodeHex)
}

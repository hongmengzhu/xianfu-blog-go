package auditPg

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pangu-2/go-tools/tools/noPg"
)

// noPrefix 操作日志编号前缀：AL = Audit Log。
const noPrefix = "AL"

// nextNo 生成日志编号。
func nextNo() string {
	return noPg.MakeNo(noPrefix)
}

// newRequestID 生成随机请求 ID（16 字节 hex）。
func newRequestID() string {
	bs := make([]byte, 16)
	if _, err := rand.Read(bs); err != nil {
		return ""
	}
	return hex.EncodeToString(bs)
}

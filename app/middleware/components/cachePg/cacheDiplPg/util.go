package cacheDiplPg

import "github.com/pangu-2/go-tools/tools/cryptPg"

// HashSha
//
//	@Description: sha512
//	@param st1
//	@param st2
//	@return string
func HashSha(st1, st2 string) string {
	return cryptPg.Sha256(cryptPg.Sha512(st1+"::"+st2) + "::" + st2)
}

// HashShaVerify
//
//	@Description: 验证
//	@param st1
//	@param st2
//	@return string
func HashShaVerify(st1, st2 string, verify string) bool {
	return HashSha(st1, st2) == verify
}

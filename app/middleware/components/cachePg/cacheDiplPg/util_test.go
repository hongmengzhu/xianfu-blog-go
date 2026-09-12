package cacheDiplPg

import (
	"fmt"
	"testing"
)

func TestSha1(t *testing.T) {
	key := "85227155478337338126"
	secret := "YGJQDjCBy3VNMnKFQzll"
	fmt.Println(HashSha(key, secret))
}

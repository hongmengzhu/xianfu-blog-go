package typePg

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-spring.org/log"
)

type Uint64String uint64

func (i Uint64String) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, i)), nil
}

func (i *Uint64String) UnmarshalJSON(value []byte) error {

	str := string(value)
	//log.Info("UnmarshalJSON.%+v", value)
	//log.Info("UnmarshalJSON.%+v", str)
	if len(str) == 0 {
		*i = Uint64String(0)
		return nil
	}
	if str == `""` {
		*i = Uint64String(0)
		return nil
	}
	if strings.LastIndexAny(str, `"`) > -1 {
		str = string(value[1 : len(value)-1])
	}
	m, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "UnmarshalJSON")
		*i = Uint64String(0)
	} else {
		*i = Uint64String(m)
	}
	return nil
}

func (i Uint64String) ToString() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i Uint64String) ToInt64() int64 {
	return int64(i)
}
func (i Uint64String) ToUint64() uint64 {
	return uint64(i)
}

func (i *Uint64String) SetZero() uint64 {
	*i = Uint64String(0)
	return 0
}

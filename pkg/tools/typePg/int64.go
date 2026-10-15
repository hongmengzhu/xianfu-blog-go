package typePg

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-spring.org/log"
)

type Int64String int64

func (i Int64String) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, i)), nil
}

func (i *Int64String) UnmarshalJSON(value []byte) error {

	str := string(value)
	//log.Info("UnmarshalJSON.%+v", value)
	//log.Info("UnmarshalJSON.%+v", str)
	if len(str) == 0 {
		*i = Int64String(0)
		return nil
	}
	if str == `""` {
		*i = Int64String(0)
		return nil
	}
	if strings.LastIndexAny(str, `"`) > -1 {
		str = string(value[1 : len(value)-1])
	}
	m, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "UnmarshalJSON")
		*i = Int64String(0)
	} else {
		*i = Int64String(m)
	}
	return nil
}

func (i Int64String) ToString() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int64String) ToInt64() int64 {
	return int64(i)
}
func (i Int64String) ToUint64() int64 {
	return int64(i)
}

func (i *Int64String) SetZero() int64 {
	*i = Int64String(0)
	return 0
}

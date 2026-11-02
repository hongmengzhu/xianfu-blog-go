package enumDeletedPg

// 状态
const (
	DELETED StateEnum = 1 //删除
	NORMAL  StateEnum = 2 //正常
)

type StateEnum int8

func (this StateEnum) String() string {
	switch this {
	case 1:
		return "删除"
	case 2:
		return "正常"
	default:
		return "未知"
	}
}
func (this StateEnum) ToInt64() int64 {
	return int64(this)
}
func (this StateEnum) ToInt8() int8 {
	return int8(this)
}
func (this StateEnum) Index() int8 {
	return int8(this)
}

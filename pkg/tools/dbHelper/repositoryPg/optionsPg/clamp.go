package optionsPg

// Number 分页参数支持的整数类型
type Number interface {
	~int | ~int32 | ~int64
}

// 分页上限：防止 pageSize/pageNum 过大导致查询/内存放大（DoS）。
const (
	// DefaultPageSize 默认每页条数
	DefaultPageSize = 20
	// MaxPageSize 每页条数硬上限
	MaxPageSize = 100
	// MaxPageNum 页码硬上限
	MaxPageNum = 100000
)

// ClampPage 对分页参数做边界收敛：pageSize 走默认值与硬上限，pageNum 不为负且不超上限。
//
//	@Description: 统一分页上限，避免每个 service 各自实现且遗漏上限
//	@param pageNum
//	@param pageSize
//	@return T 收敛后的 pageNum
//	@return T 收敛后的 pageSize
func ClampPage[T Number](pageNum, pageSize T) (T, T) {
	if pageSize < 1 {
		pageSize = T(DefaultPageSize)
	}
	if pageSize > T(MaxPageSize) {
		pageSize = T(MaxPageSize)
	}
	if pageNum < 0 {
		pageNum = 0
	}
	if pageNum > T(MaxPageNum) {
		pageNum = T(MaxPageNum)
	}
	return pageNum, pageSize
}

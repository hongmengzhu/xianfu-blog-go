package opLogPg

// pathMatcher 操作日志排除规则匹配器。
// 规则按路径段匹配，模式耗尽即视为命中（前缀语义）：
//   - *  匹配任意单段
//   - ** 匹配零到多段
//
// 因此 "/xianfu/pub/*" 会排除 /xianfu/pub/login、/xianfu/pub/login/x 等，
// 但不会排除 /xianfu/public；"/xianfu/auth/**" 会排除 /xianfu/auth 及其下所有路径。
type pathMatcher struct {
	patterns [][]string
}

func newPathMatcher(patterns []string) *pathMatcher {
	m := &pathMatcher{}
	for _, pattern := range patterns {
		segs := splitPath(pattern)
		if len(segs) == 0 {
			continue
		}
		m.patterns = append(m.patterns, segs)
	}
	return m
}

// Match 判断路径是否命中任一排除规则。
func (m *pathMatcher) Match(path string) bool {
	if m == nil || len(m.patterns) == 0 {
		return false
	}
	segs := splitPath(path)
	for _, pattern := range m.patterns {
		if matchPathSegments(pattern, segs) {
			return true
		}
	}
	return false
}

func matchPathSegments(pattern, path []string) bool {
	return matchPathFrom(pattern, path, 0, 0)
}

func matchPathFrom(pattern, path []string, pi, si int) bool {
	for pi < len(pattern) {
		switch pattern[pi] {
		case "**":
			for k := si; k <= len(path); k++ {
				if matchPathFrom(pattern, path, pi+1, k) {
					return true
				}
			}
			return false
		case "*":
			if si >= len(path) {
				return false
			}
		default:
			if si >= len(path) || pattern[pi] != path[si] {
				return false
			}
		}
		pi++
		si++
	}
	// 模式已耗尽，路径剩余部分不再限制（前缀匹配）。
	return true
}

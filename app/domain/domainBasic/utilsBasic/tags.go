package utilsBasic

func TagCacheKey(prefix, tag string) string {
	return prefix + ":" + tag
}

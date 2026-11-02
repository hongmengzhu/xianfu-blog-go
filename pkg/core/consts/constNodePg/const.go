package constNodePg

// 默认根id
const ROOT = "_ROOT_"

func NoLinkDefault(no string) string {
	return "|" + no + "|"
}

func NoLinkAssemble(prefix, no string) string {
	return prefix + no + "|"
}

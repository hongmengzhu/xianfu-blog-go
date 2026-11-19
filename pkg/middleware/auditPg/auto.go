package auditPg

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constModulePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
)

// 审计元数据自动推导。
//
// handler 只需按约定暴露路由 /xianfu/{sys|manage}/<domain>/<resource>/<action>，
// 中间件按 <domain>/<resource> 查 constModulePg.ModuleMap 拿模块信息，
// 按 <action> 查 MatchActionMap 拿动作信息，因此不再需要每个 handler 写 *Audit.go。
//
// 约定之外的路由（模块或动作未登记）不会被审计；需要强制覆盖时仍可用 RegisterRoute。

const routeRootPrefix = "xianfu"

var (
	moduleIndexOnce sync.Once
	moduleIndex     map[string]constModulePg.Item
	routeMetaCache  sync.Map // 路由模板 -> autoMetaResult
)

type autoMetaResult struct {
	meta Meta
	ok   bool
}

// autoMetaForPath 返回按约定推导出的审计元数据；结果按路由模板缓存。
func autoMetaForPath(path string) (Meta, bool) {
	if v, ok := routeMetaCache.Load(path); ok {
		r := v.(autoMetaResult)
		return r.meta, r.ok
	}
	m, ok := resolveMeta(path)
	routeMetaCache.Store(path, autoMetaResult{meta: m, ok: ok})
	return m, ok
}

// resolveMeta 由路由模板解析审计元数据。
// 只处理 /xianfu/{sys|manage}/<domain>/<resource>/<action> 形式的路由。
func resolveMeta(path string) (Meta, bool) {
	moduleKey, action, ok := splitRoute(path)
	if !ok {
		return Meta{}, false
	}
	item, ok := lookupModule(moduleKey)
	if !ok {
		return Meta{}, false
	}
	match, ok := matchAction(action)
	if !ok {
		return Meta{}, false
	}
	return buildAutoMeta(item, match), true
}

// splitRoute 拆分路由模板：返回模块键 <domain>/<resource> 与动作名。
func splitRoute(path string) (moduleKey string, action string, ok bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 5 || parts[0] != routeRootPrefix {
		return "", "", false
	}
	// 仅接管后台两类路由；auth 等其它域不在审计范围。
	if parts[1] != "sys" && parts[1] != "manage" {
		return "", "", false
	}
	if parts[2] == "" || parts[3] == "" || parts[4] == "" {
		return "", "", false
	}
	return parts[2] + "/" + parts[3], parts[4], true
}

// lookupModule 按模块路径查表；路径比较忽略大小写、- 和 _，
// 使 /basic/tags-category 与 basic/tagsCategory 命中同一条目。
func lookupModule(moduleKey string) (constModulePg.Item, bool) {
	moduleIndexOnce.Do(func() {
		moduleIndex = buildModuleIndex()
	})
	item, ok := moduleIndex[normalizeModuleKey(moduleKey)]
	return item, ok
}

func buildModuleIndex() map[string]constModulePg.Item {
	index := make(map[string]constModulePg.Item, len(constModulePg.ModuleMap)*2)
	put := func(path string, item constModulePg.Item) {
		key := normalizeModuleKey(path)
		if key == "" {
			return
		}
		if prev, exists := index[key]; exists && (prev.Module != item.Module || prev.TargetType != item.TargetType) {
			panic(fmt.Sprintf("auditPg: 模块路径归一化后冲突: %q 已属于 %s，又映射到 %s", path, prev.Module, item.Module))
		}
		index[key] = item
	}
	for key, item := range constModulePg.ModuleMap {
		put(key, item)
		for _, mark := range item.PathMark {
			put(mark, item)
		}
	}
	return index
}

func normalizeModuleKey(path string) string {
	path = strings.ToLower(strings.TrimSpace(path))
	path = strings.ReplaceAll(path, "-", "")
	path = strings.ReplaceAll(path, "_", "")
	return strings.Trim(path, "/")
}

// buildAutoMeta 合并模块信息与动作信息。模块提供模块名、对象类型，动作提供类型与风险。
func buildAutoMeta(item constModulePg.Item, match Match) Meta {
	meta := Meta{
		Module:        item.Module,
		ModuleName:    item.Name,
		SubModule:     moduleDomain(item.Module),
		Action:        match.Action,
		ActionName:    match.ActionName,
		OperationType: match.OperationType,
		RiskLevel:     match.RiskLevel,
		TargetType:    item.TargetType,
		TargetFrom:    match.TargetFrom,
		TargetKey:     match.TargetKey,
	}
	// 导出/导入通常是二进制流：不保存请求体，响应体只在失败时保留。
	if match.OperationType == entityBasic.OperationLogTypeExport || match.OperationType == entityBasic.OperationLogTypeImport {
		meta.SkipRequestBody = true
		meta.RecordResponse = false
	}
	return meta.withDefaults()
}

// moduleDomain 取模块编码的顶级域，如 basic.area -> basic。
func moduleDomain(module string) string {
	if i := strings.IndexByte(module, '.'); i > 0 {
		return module[:i]
	}
	return module
}

// matchAction 解析动作名：先精确匹配 MatchActionMap 的 key，
// 再按 Match.Prefix 匹配别名（如 exportExcel -> export），多个命中时取最长前缀。
func matchAction(action string) (Match, bool) {
	action = strings.TrimSpace(action)
	if action == "" {
		return Match{}, false
	}
	if match, ok := MatchActionMap[action]; ok {
		return match, true
	}
	bestLen := 0
	var best Match
	found := false
	for _, match := range MatchActionMap {
		for _, prefix := range match.Prefix {
			prefix = strings.TrimSpace(prefix)
			if prefix == "" || len(prefix) <= bestLen {
				continue
			}
			if action == prefix || strings.HasPrefix(action, prefix) {
				best, bestLen, found = match, len(prefix), true
			}
		}
	}
	return best, found
}

package auditPg

import (
	"testing"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
)

// TestResolveMeta 校验路由 -> 模块 -> 动作的自动推导。
func TestResolveMeta(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantOK     bool
		wantModule string
		wantAction string
		wantType   string
		wantOpType int8
		wantRisk   int8
		wantKey    string
	}{
		{name: "系统端删除区域", path: "/xianfu/sys/basic/area/delete", wantOK: true,
			wantModule: "basic.area", wantAction: "delete", wantType: "basic_area",
			wantOpType: entityBasic.OperationLogTypeDelete, wantRisk: 3, wantKey: "ids"},
		{name: "kebab 资源名归一化", path: "/xianfu/manage/basic/tags-category/delete", wantOK: true,
			wantModule: "basic.tagsCategory", wantAction: "delete", wantType: "basic_tags_category",
			wantOpType: entityBasic.OperationLogTypeDelete, wantRisk: 3, wantKey: "ids"},
		{name: "数据字典 kebab", path: "/xianfu/manage/basic/data-dictionary/createUpdate", wantOK: true,
			wantModule: "basic.dataDictionary", wantAction: "createUpdate", wantType: "basic_data_dictionary",
			wantOpType: entityBasic.OperationLogTypeCreateUpdate, wantRisk: 2, wantKey: "id"},
		{name: "数据字典子项 PathMark", path: "/xianfu/sys/basic/data-dictionary-sub/state", wantOK: true,
			wantModule: "basic.dictionary.sub", wantAction: "state", wantType: "basic_data_dictionary",
			wantOpType: entityBasic.OperationLogTypeStateChange, wantRisk: 2, wantKey: "ids"},
		{name: "idp 路由挂在 ram-idp 下", path: "/xianfu/manage/ram-idp/identity-provider/delete", wantOK: true,
			wantModule: "idp.identityProvider", wantAction: "delete", wantType: "ram_identity_provider",
			wantOpType: entityBasic.OperationLogTypeDelete, wantRisk: 3, wantKey: "ids"},
		{name: "ram 多因素认证 PathMark", path: "/xianfu/sys/ram/account-mfa/delete", wantOK: true,
			wantModule: "ram.accountAuthMfa", wantAction: "delete", wantType: "ram_account_auth_mfa",
			wantOpType: entityBasic.OperationLogTypeDelete, wantRisk: 3, wantKey: "ids"},
		{name: "导出别名", path: "/xianfu/sys/basic/area/exportExcel", wantOK: true,
			wantModule: "basic.area", wantAction: "export", wantType: "basic_area",
			wantOpType: entityBasic.OperationLogTypeExport, wantRisk: 2},
		{name: "查询不审计", path: "/xianfu/sys/basic/area/query", wantOK: false},
		{name: "未知模块", path: "/xianfu/sys/public/login", wantOK: false},
		{name: "非后台路由", path: "/xianfu/auth/basic/area/delete", wantOK: false},
		{name: "空路径", path: "", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, ok := autoMetaForPath(tt.path)
			if ok != tt.wantOK {
				t.Fatalf("ok=%v, want %v (meta=%+v)", ok, tt.wantOK, meta)
			}
			if !ok {
				return
			}
			if err := meta.Validate(); err != nil {
				t.Fatalf("推导出的 Meta 未通过校验: %v", err)
			}
			if meta.Module != tt.wantModule {
				t.Errorf("Module=%q, want %q", meta.Module, tt.wantModule)
			}
			if meta.Action != tt.wantAction {
				t.Errorf("Action=%q, want %q", meta.Action, tt.wantAction)
			}
			if meta.TargetType != tt.wantType {
				t.Errorf("TargetType=%q, want %q", meta.TargetType, tt.wantType)
			}
			if meta.OperationType != tt.wantOpType {
				t.Errorf("OperationType=%d, want %d", meta.OperationType, tt.wantOpType)
			}
			if meta.RiskLevel != tt.wantRisk {
				t.Errorf("RiskLevel=%d, want %d", meta.RiskLevel, tt.wantRisk)
			}
			if meta.TargetKey != tt.wantKey {
				t.Errorf("TargetKey=%q, want %q", meta.TargetKey, tt.wantKey)
			}
			if meta.Action == "export" && (!meta.SkipRequestBody || meta.RecordResponse) {
				t.Errorf("导出应关闭请求体/响应体记录: %+v", meta)
			}
		})
	}
}

func TestMatchActionPrefix(t *testing.T) {
	if match, ok := matchAction("exportExcel"); !ok || match.Action != "export" {
		t.Fatalf("exportExcel 应匹配 export, got %+v ok=%v", match, ok)
	}
	if match, ok := matchAction("physicalDeletion"); !ok || match.Action != "physicalDeletion" {
		t.Fatalf("physicalDeletion 应精确匹配, got %+v ok=%v", match, ok)
	}
	if _, ok := matchAction("detail"); ok {
		t.Fatal("detail 不应匹配任何动作")
	}
}

func TestModuleIndexNoConflict(t *testing.T) {
	// 触发索引构建；归一化冲突会在构建期 panic。
	lookupModule("basic/area")
}

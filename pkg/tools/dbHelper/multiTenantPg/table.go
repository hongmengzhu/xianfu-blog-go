package multiTenantPg

import "github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"

// 系统
var multiTableSys = pg.MultiItem{
	Contain: []string{
		"ram_account",
		"ram_account_device",
		"ram_account_login_log",
		"ram_account_session",
		"ram_account_session_access_key",
		"ram_app",
		"ram_app_access_key",
		"ram_app_category",
		"ram_channel",
		"ram_department",
		"ram_favorites",
		"ram_group",
		"ram_level",
		"ram_position",
		"ram_post",
		"ram_role",
		"ram_team",
	},
	Not: make([]string, 0),
}

// 管理后台
var multiTableManage = pg.MultiItem{
	Contain: []string{
		"api_dipl",
		"api_dipl_access_key",
		"api_dipl_category",
		"ram_account",
		"ram_account_device",
		"ram_account_login_log",
		"ram_account_session",
		"ram_account_session_access_key",
		"ram_app",
		"ram_app_access_key",
		"ram_app_category",
		"ram_channel",
		"ram_department",
		"ram_favorites",
		"ram_group",
		"ram_level",
		"ram_position",
		"ram_post",
		"ram_role",
		"ram_team",
		"blog_article",
		"blog_article_category",
		"blog_attachment",
		"blog_topic",
		"blog_topic_category",
		"blog_topic_relation",
		"blog_collect",
		"blog_collect_category",
		"blog_bookmark",
		"blog_bookmark_category",
		"bm_bookmark",
		"bm_bookmark_category",
		"bm_tool",
		"bm_tool_category",
		"b72_config",
		"b72_config_category",
		"b72_config_event",
		"b72_config_event_fields",
		"b72_config_list",
		"b72_config_model",
		"b72_config_model_fields",
		"b72_config_model_rules",
		"basic_account_apply_deny_list",
		"basic_data_snapshot",
		"basic_explain",
		"basic_explain_category",
		"basic_model_rules",
		"basic_model_rules_template",
		"basic_module",
		"basic_operation_log",
		"basic_tags",
		"basic_tags_category",
		"basic_tags_relation",
		"ram_account_auth_mfa",
		"ram_account_auth_webauthn",
		"ram_account_authorization",
		"ram_account_deny_list",
		"ram_account_session_log",
		"ram_job_function",
		"ram_person",
		"ram_resource_authority",
		"ram_resource_group",
		"ram_resource_group_relation",
		"ram_resource_menu",
		"ram_resource_relation",
		"ram_identity_provider",
		"ram_identity_source",
		"ram_identity_source_callback",
		"ram_idp_binding",
		"ram_idp_credential",
		"ram_idp_metadata_cache",
		"tc_level",
		"tc_organization",
		"tc_tenant_domain",
	},
	Not: make([]string, 0),
}

// 管理后台
var multiTableCustomer = pg.MultiItem{
	Contain: []string{
		"api_dipl",
		"api_dipl_access_key",
		"api_dipl_category",
		"cc_ram_account",
		"cc_ram_account_device",
		"cc_ram_account_login_log",
		"cc_ram_account_session",
		"cc_ram_account_session_access_key",
		"cc_ram_app",
		"cc_ram_app_access_key",
		"cc_ram_app_category",
		"cc_ram_channel",
		"cc_ram_department",
		"cc_ram_favorites",
		"cc_ram_group",
		"cc_ram_level",
		"cc_ram_position",
		"cc_ram_post",
		"cc_ram_role",
		"cc_ram_team",
		"blog_article",
		"blog_article_category",
		"blog_attachment",
		"blog_topic",
		"blog_topic_category",
		"blog_topic_relation",
		"tc_tenant_domain",
		"lib_app_library",
		"mc_channel",
		"mc_logistics",
		"mc_store",
	},
	Not: make([]string, 0),
}

var multiTableOrg = pg.MultiItem{
	Contain: []string{
		"api_dipl",
		"api_dipl_access_key",
		"api_dipl_category",
		"cc_ram_account",
		"cc_ram_account_device",
		"cc_ram_account_login_log",
		"cc_ram_account_session",
		"cc_ram_account_session_access_key",
		"cc_ram_app",
		"cc_ram_app_access_key",
		"cc_ram_app_category",
		"cc_ram_channel",
		"cc_ram_department",
		"cc_ram_favorites",
		"cc_ram_group",
		"cc_ram_level",
		"cc_ram_position",
		"cc_ram_post",
		"cc_ram_role",
		"cc_ram_team",
		"blog_article",
		"blog_article_category",
		"blog_attachment",
		"blog_topic",
		"blog_topic_category",
		"blog_topic_relation",
		"tc_tenant_domain",
		"lib_app_library",
		"mc_channel",
		"mc_logistics",
		"mc_store",
		//
		"bc_group",
	},
	Not: make([]string, 0),
}

// GetMultiTableSys
//
//	@Description: 系统 表规则
//	@return pg.MultiItem
func GetMultiTableSys() pg.MultiItem {
	return multiTableSys
}

// GetMultiTableManage
//
//	@Description:  管理后台 表规则
//	@return pg.MultiItem
func GetMultiTableManage() pg.MultiItem {
	return multiTableManage
}

// GetMultiTableCustomer
//
//	@Description:  客户后台 表规则
//	@return pg.MultiItem
func GetMultiTableCustomer() pg.MultiItem {
	return multiTableCustomer
}

// GetMultiTableOrg
//
//	@Description:  客户后台 表规则
//	@return pg.MultiItem
func GetMultiTableOrg() pg.MultiItem {
	return multiTableOrg
}

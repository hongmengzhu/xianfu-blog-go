package constModulePg

// ModuleMap 模块映射
// key 为模块路径（与路由分组去掉 /xianfu/{sys|manage|auth} 前缀后的部分对齐），
// Module 为点分模块编码，ModuleRoot 为模块根（无上级模块时与 Module 一致），
// TargetType 为库表模型对应的表名（见 pkg/infrastructure/entity*）。
var ModuleMap = map[string]Item{
	// api 开放接口
	"api/dipl":          Item{Module: "api.dipl", Name: "外交馆", ModuleRoot: "api.dipl", TargetType: "api_dipl"},
	"api/diplAccessKey": Item{Module: "api.diplAccessKey", Name: "外交馆密钥", ModuleRoot: "api.diplAccessKey", TargetType: "api_dipl_access_key"},
	"api/diplCategory":  Item{Module: "api.diplCategory", Name: "外交馆分类", ModuleRoot: "api.diplCategory", TargetType: "api_dipl_category"},

	// b72 配置中心
	"b72/config":            Item{Module: "b72.config", Name: "配置", ModuleRoot: "b72.config", TargetType: "b72_config"},
	"b72/configCategory":    Item{Module: "b72.configCategory", Name: "配置分类", ModuleRoot: "b72.configCategory", TargetType: "b72_config_category"},
	"b72/configEvent":       Item{Module: "b72.configEvent", Name: "配置事件", ModuleRoot: "b72.configEvent", TargetType: "b72_config_event"},
	"b72/configEventFields": Item{Module: "b72.configEventFields", Name: "配置事件字段", ModuleRoot: "b72.configEventFields", TargetType: "b72_config_event_fields"},
	"b72/configList":        Item{Module: "b72.configList", Name: "配置列表", ModuleRoot: "b72.configList", TargetType: "b72_config_list"},
	"b72/configModel":       Item{Module: "b72.configModel", Name: "配置模型", ModuleRoot: "b72.configModel", TargetType: "b72_config_model"},
	"b72/configModelFields": Item{Module: "b72.configModelFields", Name: "配置模型字段", ModuleRoot: "b72.configModelFields", TargetType: "b72_config_model_fields"},
	"b72/configModelRules":  Item{Module: "b72.configModelRules", Name: "配置模型规则", ModuleRoot: "b72.configModelRules", TargetType: "b72_config_model_rules"},

	// basic 基础信息
	"basic/accountApplyDenyList": Item{Module: "basic.accountApplyDenyList", Name: "账号申请拒绝名单", ModuleRoot: "basic.accountApplyDenyList", TargetType: "basic_account_apply_deny_list"},
	"basic/area":                 Item{Module: "basic.area", Name: "区域", ModuleRoot: "basic.area", TargetType: "basic_area"},
	"basic/attachment":           Item{Module: "basic.attachment", Name: "附件", ModuleRoot: "basic.attachment", TargetType: "basic_attachment"},
	"basic/attachmentRelated":    Item{Module: "basic.attachmentRelated", Name: "附件关联", ModuleRoot: "basic.attachmentRelated", TargetType: "basic_attachment_related"},
	"basic/country":              Item{Module: "basic.country", Name: "国家", ModuleRoot: "basic.country", TargetType: "basic_country"},
	"basic/dataDictionary":       Item{Module: "basic.dataDictionary", Name: "数据字典", ModuleRoot: "basic.dataDictionary", TargetType: "basic_data_dictionary"},
	"basic/dataDictionarySub":    Item{Module: "basic.dictionary.sub", Name: "数据字典子项", ModuleRoot: "basic.dataDictionary", TargetType: "basic_data_dictionary"},
	"basic/dataSnapshot":         Item{Module: "basic.dataSnapshot", Name: "数据快照", ModuleRoot: "basic.dataSnapshot", TargetType: "basic_data_snapshot"},
	"basic/explain":              Item{Module: "basic.explain", Name: "说明", ModuleRoot: "basic.explain", TargetType: "basic_explain"},
	"basic/explainCategory":      Item{Module: "basic.explainCategory", Name: "说明分类", ModuleRoot: "basic.explainCategory", TargetType: "basic_explain_category"},
	"basic/modelRules":           Item{Module: "basic.modelRules", Name: "模型规则", ModuleRoot: "basic.modelRules", TargetType: "basic_model_rules"},
	"basic/modelRulesTemplate":   Item{Module: "basic.modelRulesTemplate", Name: "模型规则模板", ModuleRoot: "basic.modelRulesTemplate", TargetType: "basic_model_rules_template"},
	"basic/module":               Item{Module: "basic.module", Name: "模块", ModuleRoot: "basic.module", TargetType: "basic_module"},
	"basic/operationLog":         Item{Module: "basic.operationLog", Name: "操作日志", ModuleRoot: "basic.operationLog", TargetType: "basic_operation_log"},
	"basic/tags":                 Item{Module: "basic.tags", Name: "标签", ModuleRoot: "basic.tags", TargetType: "basic_tags"},
	"basic/tagsCategory":         Item{Module: "basic.tagsCategory", Name: "标签分类", ModuleRoot: "basic.tagsCategory", TargetType: "basic_tags_category"},
	"basic/tagsRelation":         Item{Module: "basic.tagsRelation", Name: "标签关系", ModuleRoot: "basic.tagsRelation", TargetType: "basic_tags_relation"},

	// blog 博客
	"blog/article":           Item{Module: "blog.article", Name: "文章", ModuleRoot: "blog.article", TargetType: "blog_article"},
	"blog/articleCategory":   Item{Module: "blog.articleCategory", Name: "文章分类", ModuleRoot: "blog.articleCategory", TargetType: "blog_article_category"},
	"blog/articleStatistics": Item{Module: "blog.articleStatistics", Name: "文章统计", ModuleRoot: "blog.articleStatistics", TargetType: "blog_article_statistics"},
	"blog/collect":           Item{Module: "blog.collect", Name: "收集", ModuleRoot: "blog.collect", TargetType: "blog_collect"},
	"blog/collectCategory":   Item{Module: "blog.collectCategory", Name: "收集分类", ModuleRoot: "blog.collectCategory", TargetType: "blog_collect_category"},
	"blog/topic":             Item{Module: "blog.topic", Name: "文章主题", ModuleRoot: "blog.topic", TargetType: "blog_topic"},
	"blog/topicCategory":     Item{Module: "blog.topicCategory", Name: "主题分类", ModuleRoot: "blog.topicCategory", TargetType: "blog_topic_category"},
	"blog/topicRelation":     Item{Module: "blog.topicRelation", Name: "主题关系", ModuleRoot: "blog.topicRelation", TargetType: "blog_topic_relation"},
	"blog/topicStatistics":   Item{Module: "blog.topicStatistics", Name: "主题统计", ModuleRoot: "blog.topicStatistics", TargetType: "blog_topic_statistics"},

	// bm 书签工具
	"bm/bookmark":         Item{Module: "bm.bookmark", Name: "书签", ModuleRoot: "bm.bookmark", TargetType: "bm_bookmark"},
	"bm/bookmarkCategory": Item{Module: "bm.bookmarkCategory", Name: "书签分类", ModuleRoot: "bm.bookmarkCategory", TargetType: "bm_bookmark_category"},
	"bm/tool":             Item{Module: "bm.tool", Name: "工具", ModuleRoot: "bm.tool", TargetType: "bm_tool"},
	"bm/toolCategory":     Item{Module: "bm.toolCategory", Name: "工具分类", ModuleRoot: "bm.toolCategory", TargetType: "bm_tool_category"},

	// idp 身份提供商（表名沿用 ram_ 前缀）
	"idp/identityProvider":       Item{Module: "idp.identityProvider", Name: "身份提供商", ModuleRoot: "idp.identityProvider", TargetType: "ram_identity_provider", PathMark: []string{"ram-idp/identity-provider"}},
	"idp/identitySource":         Item{Module: "idp.identitySource", Name: "身份认证源", ModuleRoot: "idp.identitySource", TargetType: "ram_identity_source", PathMark: []string{"ram-idp/identity-source"}},
	"idp/identitySourceCallback": Item{Module: "idp.identitySourceCallback", Name: "认证源回调白名单", ModuleRoot: "idp.identitySourceCallback", TargetType: "ram_identity_source_callback", PathMark: []string{"ram-idp/identity-source-callback"}},
	"idp/idpBinding":             Item{Module: "idp.idpBinding", Name: "身份绑定", ModuleRoot: "idp.idpBinding", TargetType: "ram_idp_binding", PathMark: []string{"ram-idp/idp-binding"}},
	"idp/idpCredential":          Item{Module: "idp.idpCredential", Name: "身份凭证", ModuleRoot: "idp.idpCredential", TargetType: "ram_idp_credential", PathMark: []string{"ram-idp/idp-credential"}},
	"idp/metadataCache":          Item{Module: "idp.metadataCache", Name: "身份元数据缓存", ModuleRoot: "idp.metadataCache", TargetType: "ram_idp_metadata_cache", PathMark: []string{"ram-idp/idp-metadata-cache"}},

	// ram 权限管理
	"ram/account":                 Item{Module: "ram.account", Name: "账号", ModuleRoot: "ram.account", TargetType: "ram_account"},
	"ram/accountAuthMfa":          Item{Module: "ram.accountAuthMfa", Name: "账号多因素认证", ModuleRoot: "ram.accountAuthMfa", TargetType: "ram_account_auth_mfa", PathMark: []string{"ram/account-mfa"}},
	"ram/accountAuthWebauthn":     Item{Module: "ram.accountAuthWebauthn", Name: "账号通行密钥", ModuleRoot: "ram.accountAuthWebauthn", TargetType: "ram_account_auth_webauthn"},
	"ram/accountAuthorization":    Item{Module: "ram.accountAuthorization", Name: "账号授权", ModuleRoot: "ram.accountAuthorization", TargetType: "ram_account_authorization"},
	"ram/accountDenyList":         Item{Module: "ram.accountDenyList", Name: "账号拒绝名单", ModuleRoot: "ram.accountDenyList", TargetType: "ram_account_deny_list"},
	"ram/accountDevice":           Item{Module: "ram.accountDevice", Name: "账号设备", ModuleRoot: "ram.accountDevice", TargetType: "ram_account_device"},
	"ram/accountLoginLog":         Item{Module: "ram.accountLoginLog", Name: "账号登录日志", ModuleRoot: "ram.accountLoginLog", TargetType: "ram_account_login_log"},
	"ram/accountSession":          Item{Module: "ram.accountSession", Name: "账号会话", ModuleRoot: "ram.accountSession", TargetType: "ram_account_session"},
	"ram/accountSessionAccessKey": Item{Module: "ram.accountSessionAccessKey", Name: "账号会话密钥", ModuleRoot: "ram.accountSessionAccessKey", TargetType: "ram_account_session_access_key"},
	"ram/accountSessionLog":       Item{Module: "ram.accountSessionLog", Name: "账号 SSO 登录日志", ModuleRoot: "ram.accountSessionLog", TargetType: "ram_account_session_log"},
	"ram/app":                     Item{Module: "ram.app", Name: "应用", ModuleRoot: "ram.app", TargetType: "ram_app"},
	"ram/appAccessKey":            Item{Module: "ram.appAccessKey", Name: "应用密钥", ModuleRoot: "ram.appAccessKey", TargetType: "ram_app_access_key"},
	"ram/appCategory":             Item{Module: "ram.appCategory", Name: "应用分类", ModuleRoot: "ram.appCategory", TargetType: "ram_app_category"},
	"ram/channel":                 Item{Module: "ram.channel", Name: "渠道", ModuleRoot: "ram.channel", TargetType: "ram_channel"},
	"ram/department":              Item{Module: "ram.department", Name: "部门", ModuleRoot: "ram.department", TargetType: "ram_department"},
	"ram/favorites":               Item{Module: "ram.favorites", Name: "收藏", ModuleRoot: "ram.favorites", TargetType: "ram_favorites"},
	"ram/group":                   Item{Module: "ram.group", Name: "用户组", ModuleRoot: "ram.group", TargetType: "ram_group"},
	"ram/jobFunction":             Item{Module: "ram.jobFunction", Name: "职务", ModuleRoot: "ram.jobFunction", TargetType: "ram_job_function"},
	"ram/level":                   Item{Module: "ram.level", Name: "级别", ModuleRoot: "ram.level", TargetType: "ram_level"},
	"ram/menu":                    Item{Module: "ram.menu", Name: "菜单", ModuleRoot: "ram.menu", TargetType: "ram_menu"},
	"ram/person":                  Item{Module: "ram.person", Name: "用户", ModuleRoot: "ram.person", TargetType: "ram_person"},
	"ram/position":                Item{Module: "ram.position", Name: "职位", ModuleRoot: "ram.position", TargetType: "ram_position"},
	"ram/post":                    Item{Module: "ram.post", Name: "岗位", ModuleRoot: "ram.post", TargetType: "ram_post"},
	"ram/resource":                Item{Module: "ram.resource", Name: "资源", ModuleRoot: "ram.resource", TargetType: "ram_resource"},
	"ram/resourceAuthority":       Item{Module: "ram.resourceAuthority", Name: "资源权限", ModuleRoot: "ram.resourceAuthority", TargetType: "ram_resource_authority"},
	"ram/resourceGroup":           Item{Module: "ram.resourceGroup", Name: "资源组", ModuleRoot: "ram.resourceGroup", TargetType: "ram_resource_group"},
	"ram/resourceGroupRelation":   Item{Module: "ram.resourceGroupRelation", Name: "资源组关系", ModuleRoot: "ram.resourceGroupRelation", TargetType: "ram_resource_group_relation"},
	"ram/resourceMenu":            Item{Module: "ram.resourceMenu", Name: "资源菜单", ModuleRoot: "ram.resourceMenu", TargetType: "ram_resource_menu"},
	"ram/resourceRelation":        Item{Module: "ram.resourceRelation", Name: "资源关系", ModuleRoot: "ram.resourceRelation", TargetType: "ram_resource_relation"},
	"ram/role":                    Item{Module: "ram.role", Name: "角色", ModuleRoot: "ram.role", TargetType: "ram_role"},
	"ram/team":                    Item{Module: "ram.team", Name: "团队", ModuleRoot: "ram.team", TargetType: "ram_team"},

	// tc 租户管理
	"tc/level":              Item{Module: "tc.level", Name: "租户级别", ModuleRoot: "tc.level", TargetType: "tc_level"},
	"tc/organization":       Item{Module: "tc.organization", Name: "组织", ModuleRoot: "tc.organization", TargetType: "tc_organization"},
	"tc/tenant":             Item{Module: "tc.tenant", Name: "租户", ModuleRoot: "tc.tenant", TargetType: "tc_tenant"},
	"tc/tenantDomain":       Item{Module: "tc.tenantDomain", Name: "租户域名", ModuleRoot: "tc.tenantDomain", TargetType: "tc_tenant_domain"},
	"tc/tenantDomainModule": Item{Module: "tc.tenantDomainModule", Name: "租户域名模块", ModuleRoot: "tc.tenantDomainModule", TargetType: "tc_tenant_domain_module"},

	// im 即时消息（黑名单表沿用 social_ 前缀）
	"im/message":        Item{Module: "im.message", Name: "消息", ModuleRoot: "im.message", TargetType: "im_message"},
	"im/messageOffline": Item{Module: "im.messageOffline", Name: "离线消息", ModuleRoot: "im.messageOffline", TargetType: "im_message_offline"},
	"im/blacklist":      Item{Module: "im.blacklist", Name: "黑名单", ModuleRoot: "im.blacklist", TargetType: "social_blacklist"},

	// social 社交
	"social/comment":        Item{Module: "social.comment", Name: "评论", ModuleRoot: "social.comment", TargetType: "social_comment"},
	"social/follower":       Item{Module: "social.follower", Name: "关注/粉丝", ModuleRoot: "social.follower", TargetType: "social_follower"},
	"social/friendRelation": Item{Module: "social.friendRelation", Name: "好友关系", ModuleRoot: "social.friendRelation", TargetType: "social_friend_relation"},
	"social/like":           Item{Module: "social.like", Name: "点赞", ModuleRoot: "social.like", TargetType: "social_like"},
}

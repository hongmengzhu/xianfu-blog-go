package entityBasic

import (
	"time"
)

// BasicOperationLogEntity 操作日志(审计日志)
// 追加写模型：只新增，不提供更新/删除接口；需要长期保留时把 State 置为冻结。
type BasicOperationLogEntity struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`

	// 身份与归属
	No         string `gorm:"column:no;type:varchar(80);not null;uniqueIndex:uk_basic_operation_log_no;comment:日志编号" json:"no" comment:"日志编号"`
	TenantNo   string `gorm:"column:tenant_no;type:varchar(80);index:idx_basic_operation_log_tenant_time,priority:1;not null;comment:租户编号" json:"tenant_no" comment:"租户编号"`
	OrgNo      string `gorm:"column:org_no;type:varchar(80);not null;comment:组织编号" json:"org_no" comment:"组织编号"`
	StoreNo    string `gorm:"column:store_no;type:varchar(80);not null;comment:门店编号" json:"store_no" comment:"门店编号"`
	TypeDomain string `gorm:"column:type_domain;type:varchar(40);not null;comment:域类型" json:"type_domain" comment:"域类型 manage/system/customer"`
	AppNo      string `gorm:"column:app_no;type:varchar(80);not null;comment:应用编号" json:"app_no" comment:"应用编号"`

	// 操作语义
	Module        string `gorm:"column:module;type:varchar(80);index:idx_basic_operation_log_module_time,priority:1;not null;comment:模块" json:"module" comment:"模块 code"`
	ModuleName    string `gorm:"column:module_name;type:varchar(120);not null;comment:模块名称" json:"module_name" comment:"模块名称快照"`
	SubModule     string `gorm:"column:sub_module;type:varchar(80);not null;comment:子模块" json:"sub_module" comment:"子模块/功能"`
	Action        string `gorm:"column:action;type:varchar(80);index:idx_basic_operation_log_module_time,priority:2;not null;comment:动作" json:"action" comment:"动作 code"`
	ActionName    string `gorm:"column:action_name;type:varchar(80);not null;comment:动作名称" json:"action_name" comment:"动作名称快照"`
	OperationType int8   `gorm:"column:operation_type;type:int2;not null;default:9;comment:操作类型" json:"operation_type" comment:"操作类型 1新增 2修改 3删除 4查询 5导出 6导入 7授权 8状态变更 9其他"`
	RiskLevel     int8   `gorm:"column:risk_level;type:int2;not null;default:1;comment:风险级别" json:"risk_level" comment:"风险级别 1低 2中 3高"`
	Description   string `gorm:"column:description;type:varchar(512);not null;comment:操作描述" json:"description" comment:"操作描述"`

	// 操作人
	OperatorNo      string `gorm:"column:operator_no;type:varchar(80);index:idx_basic_operation_log_operator_time,priority:1;not null;comment:操作人编号" json:"operator_no" comment:"操作人账号编号"`
	OperatorAccount string `gorm:"column:operator_account;type:varchar(120);not null;comment:操作人账号" json:"operator_account" comment:"登录账号"`
	OperatorName    string `gorm:"column:operator_name;type:varchar(120);not null;comment:操作人姓名" json:"operator_name" comment:"操作人姓名快照"`
	OperatorType    int8   `gorm:"column:operator_type;type:int2;not null;default:3;comment:操作人类型" json:"operator_type" comment:"1平台 2租户管理员 3普通用户 9系统任务"`
	RoleNo          string `gorm:"column:role_no;type:varchar(80);not null;comment:角色编号" json:"role_no" comment:"角色编号快照"`
	RoleName        string `gorm:"column:role_name;type:varchar(120);not null;comment:角色名称" json:"role_name" comment:"角色名称快照"`

	// 业务对象
	TargetType  string `gorm:"column:target_type;type:varchar(80);index:idx_basic_operation_log_target_time,priority:1;not null;comment:对象类型" json:"target_type" comment:"对象类型"`
	TargetID    string `gorm:"column:target_id;type:varchar(80);not null;comment:对象主键" json:"target_id" comment:"对象主键"`
	TargetNo    string `gorm:"column:target_no;type:varchar(80);index:idx_basic_operation_log_target_time,priority:2;not null;comment:对象编号" json:"target_no" comment:"对象业务编号"`
	TargetName  string `gorm:"column:target_name;type:varchar(255);not null;comment:对象名称" json:"target_name" comment:"对象名称快照"`
	TargetCount int    `gorm:"column:target_count;type:int4;not null;default:0;comment:影响条数" json:"target_count" comment:"影响条数"`

	// 变更快照
	BeforeData string `gorm:"column:before_data;type:text;comment:变更前快照JSON" json:"before_data" comment:"变更前快照"`
	AfterData  string `gorm:"column:after_data;type:text;comment:变更后快照JSON" json:"after_data" comment:"变更后快照"`
	DiffData   string `gorm:"column:diff_data;type:text;comment:字段差异JSON" json:"diff_data" comment:"字段差异"`

	// 请求上下文
	RequestMethod string `gorm:"column:request_method;type:varchar(16);not null;comment:请求方法" json:"request_method" comment:"HTTP 方法"`
	RequestPath   string `gorm:"column:request_path;type:varchar(512);not null;comment:请求路径" json:"request_path" comment:"请求路径"`
	RequestQuery  string `gorm:"column:request_query;type:text;comment:query参数" json:"request_query" comment:"query string"`
	RequestBody   string `gorm:"column:request_body;type:text;comment:请求体摘要" json:"request_body" comment:"请求体摘要"`
	ResponseBody  string `gorm:"column:response_body;type:text;comment:响应摘要" json:"response_body" comment:"响应体摘要"`

	// 结果
	Result       int8   `gorm:"column:result;type:int2;not null;index:idx_basic_operation_log_result_time,priority:1;default:0;comment:结果" json:"result" comment:"0未知 1成功 2失败 3部分成功"`
	ResultCode   string `gorm:"column:result_code;type:varchar(40);not null;comment:结果码" json:"result_code" comment:"业务码/HTTP 状态码"`
	ErrorMessage string `gorm:"column:error_message;type:text;comment:失败原因" json:"error_message" comment:"失败原因"`
	DurationMs   int64  `gorm:"column:duration_ms;type:int8;not null;default:0;comment:耗时毫秒" json:"duration_ms" comment:"耗时毫秒"`

	// 客户端
	IPAddress string `gorm:"column:ip_address;type:varchar(128);not null;comment:客户端IP" json:"ip_address" comment:"客户端 IP"`
	UserAgent string `gorm:"column:user_agent;type:text;comment:UA" json:"user_agent" comment:"User-Agent"`

	// 链路
	RequestID string `gorm:"column:request_id;type:varchar(80);index;not null;comment:请求ID" json:"request_id" comment:"请求唯一 ID"`
	TraceID   string `gorm:"column:trace_id;type:varchar(80);index;not null;comment:链路ID" json:"trace_id" comment:"链路 ID"`

	// 来源与生命周期
	Source    int8       `gorm:"column:source;type:int2;not null;default:1;comment:来源" json:"source" comment:"1接口 2事件 3定时任务 4补录"`
	State     int8       `gorm:"column:state;type:int2;not null;default:1;comment:状态" json:"state" comment:"1正常 2冻结保留"`
	OperateAt *time.Time `gorm:"column:operate_at;index:idx_basic_operation_log_tenant_time,priority:2;index:idx_basic_operation_log_module_time,priority:3;index:idx_basic_operation_log_operator_time,priority:2;index:idx_basic_operation_log_target_time,priority:3;index:idx_basic_operation_log_result_time,priority:2;comment:操作时间" json:"operate_at" comment:"操作发生时间"`
	CreateAt  *time.Time `gorm:"column:create_at;autoCreateTime;not null;default:current_timestamp;comment:创建时间" json:"create_at" comment:"入库时间"`

	ExtraData string `gorm:"column:extra_data;type:text;comment:扩展JSON" json:"extra_data" comment:"扩展数据"`
}

// TableName BasicOperationLogEntity's table name
func (*BasicOperationLogEntity) TableName() string {
	return "basic_operation_log"
}

// TableComment 表注释
func (*BasicOperationLogEntity) TableComment() string {
	return "操作日志(审计日志)"
}

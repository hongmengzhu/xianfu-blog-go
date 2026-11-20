package entityBasic

import "time"

// BasicOperationLogEntity 操作日志：只记录一次请求的基础信息，追加写模型。
type BasicOperationLogEntity struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	// 身份与归属
	No           string `gorm:"column:no;type:varchar(80);not null;uniqueIndex:uk_basic_operation_log_no;comment:日志编号" json:"no" comment:"日志编号"`
	TenantNo     string `gorm:"column:tenant_no;type:varchar(80);index:idx_basic_operation_log_tenant_time,priority:1;not null;comment:租户编号" json:"tenant_no" comment:"租户编号"`
	OrgNo        string `gorm:"column:org_no;type:varchar(80);not null;comment:组织编号" json:"org_no" comment:"组织编号"`
	StoreNo      string `gorm:"column:store_no;type:varchar(80);not null;comment:门店编号" json:"store_no" comment:"门店编号"`
	TypeDomain   string `gorm:"column:type_domain;type:varchar(40);not null;comment:域类型" json:"type_domain" comment:"域类型 manage/system/customer"`
	AppNo        string `gorm:"column:app_no;type:varchar(80);not null;comment:应用编号" json:"app_no" comment:"应用编号"`
	TerminalCode string `gorm:"column:terminal_code;type:varchar(80);not null;comment:终端编号" json:"terminal_code" comment:"终端编号"`
	Client       string `gorm:"column:client;type:varchar(80);not null;comment:客户端" json:"client" comment:"客户端"`
	DeviceNo     string `gorm:"column:device_no;type:varchar(80);not null;comment:设备编号" json:"device_no" comment:"设备编号"`
	// 操作语义
	Module     string `gorm:"column:module;type:varchar(80);index:idx_basic_operation_log_module_time,priority:1;not null;comment:模块" json:"module" comment:"模块 code"`
	ModuleName string `gorm:"column:module_name;type:varchar(120);not null;comment:模块名称" json:"module_name" comment:"模块名称快照"`
	ModuleRoot string `gorm:"column:module_root;type:varchar(80);not null;comment:模块根" json:"module_root" comment:"模块根"`
	Action     string `gorm:"column:action;type:varchar(80);index:idx_basic_operation_log_module_time,priority:2;not null;comment:动作" json:"action" comment:"动作 code"`
	ActionName string `gorm:"column:action_name;type:varchar(80);not null;comment:动作名称" json:"action_name" comment:"动作名称快照"`
	// 操作人
	OperatorNo      string   `gorm:"column:operator_no;type:varchar(80);index:idx_basic_operation_log_operator_time,priority:1;not null;comment:操作人编号" json:"operator_no" comment:"操作人账号编号"`
	OperatorAccount string   `gorm:"column:operator_account;type:varchar(120);not null;comment:操作人账号" json:"operator_account" comment:"登录账号"`
	OperatorName    string   `gorm:"column:operator_name;type:varchar(120);not null;comment:操作人姓名" json:"operator_name" comment:"操作人姓名快照"`
	Roles           []string `gorm:"column:roles;type:text;not null;default:'[]';serializer:json;comment:操作人角色" json:"roles" comment:"操作人角色快照"`
	//
	TraceID     string     `gorm:"column:trace_id;type:varchar(80);not null;default:'';index:idx_basic_operation_log_trace;comment:链路ID" json:"trace_id" comment:"链路 ID"`
	StatusCode  int        `gorm:"column:status_code;type:int4;not null;default:0;comment:HTTP状态码" json:"status_code" comment:"HTTP 状态码"`
	Method      string     `gorm:"column:method;type:varchar(16);not null;default:'';comment:请求方法" json:"method" comment:"HTTP 方法"`
	Path        string     `gorm:"column:path;type:varchar(512);not null;default:'';comment:请求路径" json:"path" comment:"请求路径"`
	IP          string     `gorm:"column:ip;type:varchar(128);not null;default:'';comment:客户端IP" json:"ip" comment:"客户端 IP"`
	Header      string     `gorm:"column:header;type:text;comment:请求头" json:"header" comment:"请求头"`
	Param       string     `gorm:"column:param;type:text;comment:query参数" json:"param" comment:"query string"`
	RequestBody string     `gorm:"column:request_body;type:text;comment:请求体" json:"request_body" comment:"请求体，文件流不记录"`
	CreateAt    *time.Time `gorm:"column:create_at;autoCreateTime;not null;index:idx_basic_operation_log_create_at;comment:创建时间" json:"create_at" comment:"创建时间"`
}

// TableName BasicOperationLogEntity's table name
func (*BasicOperationLogEntity) TableName() string {
	return "basic_operation_log"
}

// TableComment 表注释
func (*BasicOperationLogEntity) TableComment() string {
	return "操作日志"
}

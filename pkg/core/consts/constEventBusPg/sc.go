package constEventBusPg

const (
	// ScPurchaseIn 采购入库，审核通过入库
	ScPurchaseIn string = "scPurchaseIn" //采购入库，审核通过入库
	// ScInventoryChangesIn 库存变动，审核通过入库
	ScInventoryChangesIn string = "scInventoryChangesIn" //库存变动，审核通过入库

	// BasicAttachmentCreate 附件创建
	BasicAttachmentCreate = "BasicAttachmentCreate"

	// BasicCacheInit 缓存初始化
	BasicCacheInit = "BasicCacheInit"
	// BasicConfigModelCacheAllMake 生成缓存
	BasicConfigModelCacheAllMake = "BasicConfigModelCacheAllMake"
	// BasicUpdateCache 更新缓存
	BasicUpdateCache           = "BasicUpdateCache"
	BasicConfigEventCache      = "BasicConfigEventCache"
	BasicConfigEventFieldCache = "BasicConfigEventFieldCache"
	BasicModelRulesCache       = "BasicModelRulesCache"

	// GcGoodsSnapshot 商品 快照
	GcGoodsSnapshot = "GcGoodsSnapshot"
	GcSkuSnapshot   = "GcSkuSnapshot"
	// OmsOrderAction 订单动作
	OmsOrderAction = "OmsOrderAction"
	//OmsAsOrder 订单售后数据
	OmsAsOrder = "OmsAsOrder"
	// OcOrderSetStateOrder 订单状态设置
	OcOrderSetStateOrder = "OcOrderSetStateOrder"
	// OcOrderAutomaticallyCloseUnpaid 订单自动关闭未支付
	OcOrderAutomaticallyCloseUnpaid = "OcOrderAutomaticallyCloseUnpaid"
	// OcOrderPaymentInProgress 订单付款中
	OcOrderPaymentInProgress = "OcOrderPaymentInProgress"
	// QueueOcOrderQueryPaymentStatus 队列 订单支付状态
	QueueOcOrderQueryPaymentStatus = "QueueOcOrderQueryPaymentStatus"
	// RamAccountLoginLog 账号 登录日志
	RamAccountLoginLog   = "RamAccountLoginLog"
	CcRamAccountLoginLog = "CcRamAccountLoginLog"
	//BlogArticle 博客文章
	BlogArticle              = "BlogArticle"
	BlogArticleCategoryCache = "BlogArticleCategoryCache" // 博客文章分类缓存

)

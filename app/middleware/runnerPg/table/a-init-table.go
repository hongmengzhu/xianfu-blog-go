package table

import (
	"context"
	"fmt"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityApi"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbMakePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
	gormcore "go-spring.org/starter-gorm"
	"gorm.io/gorm"
)

// AInitTable 初始化创建表
type AInitTable struct {
	ser      configPg.Server   `value:"${server}"`
	database configPg.Database `value:"${database}"`
	DB       *gormcore.DB      `autowire:"postgres.primary"`
	db       *gorm.DB          `autowire:"?"`
}

func (b *AInitTable) Run(ctx context.Context) error {
	log.Infof(context.Background(), log.TagAppDef, "初始化表=>表不存在时,则进行初始化")
	entityData := make([]any, 0)
	{
		entityData = append(entityData, &entityApi.ApiDiplEntity{})
		entityData = append(entityData, &entityApi.ApiDiplAccessKeyEntity{})
		entityData = append(entityData, &entityApi.ApiDiplCategoryEntity{})
	}
	{
		entityData = append(entityData, &entityBasic.BasicAccountApplyDenyListEntity{})
		entityData = append(entityData, &entityBasic.BasicAreaEntity{})
		entityData = append(entityData, &entityBasic.BasicAttachmentEntity{})
		entityData = append(entityData, &entityBasic.BasicCountryEntity{})
		entityData = append(entityData, &entityBasic.BasicConfigEntity{})
		entityData = append(entityData, &entityBasic.BasicConfigEventEntity{})
		entityData = append(entityData, &entityBasic.BasicConfigEventFieldsEntity{})
		entityData = append(entityData, &entityBasic.BasicConfigListEntity{})
		entityData = append(entityData, &entityBasic.BasicConfigModelEntity{})
		entityData = append(entityData, &entityBasic.BasicConfigModelFieldsEntity{})
		entityData = append(entityData, &entityBasic.BasicDataDictionaryEntity{})
		entityData = append(entityData, &entityBasic.BasicDataSnapshotEntity{})
		entityData = append(entityData, &entityBasic.BasicModelRulesEntity{})
		entityData = append(entityData, &entityBasic.BasicModuleEntity{})
		entityData = append(entityData, &entityBasic.BasicTagsEntity{})
		entityData = append(entityData, &entityBasic.BasicTagsCategoryEntity{})
		entityData = append(entityData, &entityBasic.BasicTagsRelationEntity{})
	}
	{
		entityData = append(entityData, &entityBlog.BlogArticleEntity{})
		entityData = append(entityData, &entityBlog.BlogArticleCategoryEntity{})
		entityData = append(entityData, &entityBlog.BlogArticleStatisticsEntity{})
		entityData = append(entityData, &entityBlog.BlogBookmarkEntity{})
		entityData = append(entityData, &entityBlog.BlogBookmarkCategoryEntity{})
		entityData = append(entityData, &entityBlog.BlogCollectEntity{})
		entityData = append(entityData, &entityBlog.BlogCollectCategoryEntity{})
		entityData = append(entityData, &entityBlog.BlogTopicEntity{})
		entityData = append(entityData, &entityBlog.BlogTopicCategoryEntity{})
		entityData = append(entityData, &entityBlog.BlogTopicRelationEntity{})
		entityData = append(entityData, &entityBlog.BlogTopicStatisticsEntity{})
	}
	//
	{
		entityData = append(entityData, &entityRam.RamAccountEntity{})
		entityData = append(entityData, &entityRam.RamAccountAuthorizationEntity{})
		entityData = append(entityData, &entityRam.RamAccountDenyListEntity{})
		entityData = append(entityData, &entityRam.RamAccountDeviceEntity{})
		entityData = append(entityData, &entityRam.RamAccountLoginLogEntity{})
		entityData = append(entityData, &entityRam.RamAccountSessionEntity{})
		entityData = append(entityData, &entityRam.RamAccountSessionAccessKeyEntity{})
		entityData = append(entityData, &entityRam.RamAppEntity{})
		entityData = append(entityData, &entityRam.RamAppAccessKeyEntity{})
		entityData = append(entityData, &entityRam.RamAppCategoryEntity{})
		entityData = append(entityData, &entityRam.RamChannelEntity{})
		entityData = append(entityData, &entityRam.RamDepartmentEntity{})
		entityData = append(entityData, &entityRam.RamFavoritesEntity{})
		entityData = append(entityData, &entityRam.RamGroupEntity{})
		entityData = append(entityData, &entityRam.RamLevelEntity{})
		entityData = append(entityData, &entityRam.RamMenuEntity{})
		entityData = append(entityData, &entityRam.RamJobFunctionEntity{})
		entityData = append(entityData, &entityRam.RamPositionEntity{})
		entityData = append(entityData, &entityRam.RamPositionEntity{})
		entityData = append(entityData, &entityRam.RamPostEntity{})
		entityData = append(entityData, &entityRam.RamResourceEntity{})
		entityData = append(entityData, &entityRam.RamResourceAuthorityEntity{})
		entityData = append(entityData, &entityRam.RamResourceGroupEntity{})
		entityData = append(entityData, &entityRam.RamResourceGroupRelationEntity{})
		entityData = append(entityData, &entityRam.RamResourceMenuEntity{})
		entityData = append(entityData, &entityRam.RamResourceRelationEntity{})
		entityData = append(entityData, &entityRam.RamRoleEntity{})
		entityData = append(entityData, &entityRam.RamTeamEntity{})
	}
	//
	{
		entityData = append(entityData, &entityTc.TcLevelEntity{})
		entityData = append(entityData, &entityTc.TcTenantEntity{})
		entityData = append(entityData, &entityTc.TcTenantDomainEntity{})
	}

	//初始化创建表
	sv := &dbMakePg.CreateTable{
		Database: b.database,
		DB:       b.Db(),
	}
	//rt := sv.DbOpen()
	if true {
		sv.TableCreateAllByTransaction(entityData)
		//插入 基础数据
		b.sqlBasicDictionary()
		//
		log.Infof(context.Background(), log.TagAppDef, "[init].[主键序号保留].")
		b.seqEdit()
		log.Infof(context.Background(), log.TagAppDef, "初始化表 successfully")
	}
	sv = nil
	return nil
}
func (b *AInitTable) Db() *gorm.DB {
	return b.DB.DB
}

// seqEdit
//
//	@Description: 序号修改 初始值
//	@receiver b
func (b *AInitTable) seqEdit() {
	dialect := b.Db().Dialector.Name()
	if dialect == "sqlite" {
		log.Infof(context.Background(), log.TagAppDef, "[init][autoinc] sqlite skip auto increment init")
		return
	}
	// 判断是否已初始化过：以 api_dipl_id_seq 为标志序列，当前值已达初始值则跳过，避免重启重复执行
	ret, _ := IsAutoInited(b.Db(), "api_dipl", 100000)
	if ret {
		return
	}

	sql := make([]string, 0)
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "api_dipl", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "basic_account_apply_deny_list", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "basic_area", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "basic_country", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "basic_config_list", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "basic_data_dictionary", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "basic_tags", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "basic_tags_category", 100000))
	//
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_account", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_account_authorization", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_app", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_department", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_group", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_level", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_position", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_post", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "ram_team", 100000))
	//
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "tc_tenant", 100000))
	sql = append(sql, dbMakePg.MakeDataBaseAutoInc(dialect, "tc_tenant_domain", 100000))
	//
	// 事务批量执行
	tx := b.Db().Begin()
	if tx.Error != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()
	for _, raw := range sql {
		if strPg.IsBlank(raw) {
			continue
		}
		rs := tx.Exec(raw)
		if rs.Error != nil {
			log.Errorf(context.Background(), log.TagAppDef, "初始化序号异常:%+v", rs.Error)
			_ = tx.Rollback()
			return
		}
		log.Debugf(context.Background(), log.TagAppDef, "执行结果: %+v 行受影响", rs.RowsAffected)
		time.Sleep(time.Microsecond * 10)
	}
	err := tx.Commit().Error
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "创建表异常:%+v", err)
	}
}

func (b *AInitTable) sqlBasicDictionary() {
	data := []entityBasic.BasicDataDictionaryEntity{
		{ID: 1, No: "1", Name: "身份类型", Code: "typeIdentity", State: 1, TypeUniqueMd5: "7e3a158550d82e66c9f016977dcbec61", Value: "typeIdentity"},
		{ID: 2, No: "2", Name: "普通", Code: "general", State: 1, TypeUniqueMd5: "958153f1b8b96ec4c4eb2147429105d9", Value: "general", TypeCode: "typeIdentity"},
		{ID: 3, No: "3", Name: "经理", Code: "manager", State: 1, TypeUniqueMd5: "1d0258c2440a8d19e716292b231e3190", Value: "manager", TypeCode: "typeIdentity"},
		{ID: 4, No: "4", Name: "性别", Code: "sex", State: 1, TypeUniqueMd5: "3c3662bcb661d6de679c636744c66b62", Value: "sex"},
		{ID: 5, No: "5", Name: "副经理", Code: "assistant_manager", State: 1, TypeUniqueMd5: "adb153a9579779a6837a8205e07fa31f", Value: "assistant_manager", TypeCode: "typeIdentity"},
		{ID: 6, No: "6", Name: "男", Code: "male", State: 1, TypeUniqueMd5: "07cf4f8f5d8b76282917320715dda2ad", Value: "male", TypeCode: "sex"},
		{ID: 7, No: "7", Name: "女", Code: "female", State: 1, TypeUniqueMd5: "273b9ae535de53399c86a9b83148a8ed", Value: "female", TypeCode: "sex"},
		{ID: 8, No: "8", Name: "未知", Code: "unknown", State: 1, TypeUniqueMd5: "ad921d60486366258809553a3db49a4a", Value: "unknown", TypeCode: "sex"},
		{ID: 9, No: "9", Name: "菜单类型", Code: "typeMenu", State: 1, TypeUniqueMd5: "0e6b09acf0d718df71494094446eb4c0", Value: "typeMenu"},
		{ID: 10, No: "10", Name: "目录", Code: "catalog", State: 1, TypeUniqueMd5: "46f22f2a56ddd091f4b2b2c35c5ca989", Value: "catalog", TypeCode: "typeMenu"},
		{ID: 11, No: "11", Name: "菜单", Code: "menu", State: 1, TypeUniqueMd5: "8d6ab84ca2af9fccd4e4048694176ebf", Value: "menu", TypeCode: "typeMenu"},
		{ID: 12, No: "12", Name: "内嵌", Code: "embedded", State: 1, TypeUniqueMd5: "605abe26d014c72e3df9deb267e73756", Value: "embedded", TypeCode: "typeMenu"},
		{ID: 13, No: "13", Name: "按钮", Code: "button", State: 1, TypeUniqueMd5: "ce50a09343724eb82df11390e2c1de18", Value: "button", TypeCode: "typeMenu"},
		{ID: 14, No: "14", Name: "外链", Code: "link", State: 1, TypeUniqueMd5: "2a304a1348456ccd2234cd71a81bd338", Value: "link", TypeCode: "typeMenu"},
		{ID: 15, No: "15", Name: "短信服务商", Code: "spCodeSms", State: 1, TypeUniqueMd5: "5923f12d736703f67a352b8a20c14d61", Value: "spCodeSms"},
		{ID: 16, No: "16", Name: "邮件服务商", Code: "spCodeMail", State: 1, TypeUniqueMd5: "4ffa88124b350a0d4c731e2bdd04721c", Value: "spCodeMail"},
		{ID: 17, No: "17", Name: "阿里云", Code: "aliyun", State: 1, TypeUniqueMd5: "d41b8265807ce33038df2b3f7aa9fb10", Value: "aliyun", TypeCode: "spCodeSms"},
		{ID: 35, No: "35", Name: "终端类型", Code: "terminalCode", State: 1, TypeUniqueMd5: "2345a28058000459ae6173b7c66286bb", Value: "terminalCode"},
		{ID: 36, No: "36", Name: "系统", Code: "system", State: 1, TypeUniqueMd5: "54b53072540eeeb8f8e9343e71f28176", Value: "system", TypeCode: "terminalCode"},
		{ID: 37, No: "37", Name: "管理", Code: "manage", State: 1, TypeUniqueMd5: "70682896e24287b0476eff2a14c148f0", Value: "manage", TypeCode: "terminalCode"},
	}
	err := b.Db().Transaction(func(tx *gorm.DB) error {
		for i := range data {
			item := data[i]
			// 判断该主键数据是否已存在，存在则跳过，继续下一条
			var count int64
			if err := tx.Model(&entityBasic.BasicDataDictionaryEntity{}).Where("id = ?", item.ID).Count(&count).Error; err != nil {
				log.Errorf(context.Background(), log.TagAppDef, "查询数据字典主键 %d 异常:%+v", item.ID, err)
				return err
			}
			if count > 0 {
				log.Debugf(context.Background(), log.TagAppDef, "数据字典主键 %d 已存在,跳过插入", item.ID)
				continue
			}
			rs := tx.Create(&item)
			if rs.Error != nil {
				log.Errorf(context.Background(), log.TagAppDef, "初始化数据字典主键 %d 异常:%+v", item.ID, rs.Error)
				return rs.Error
			}
			log.Debugf(context.Background(), log.TagAppDef, "数据字典主键 %d 执行结果: %+v 行受影响", item.ID, rs.RowsAffected)
		}
		return nil
	})
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "创建数据字典数据异常:%+v", err)
	}
}

// IsAutoInited 判断是否已经初始化自增起始号（方言适配）
// startThreshold: 阈值 100000
func IsAutoInited(db *gorm.DB, table string, startThreshold int64) (bool, error) {
	dialect := db.Dialector.Name()
	switch dialect {
	case "postgres":
		var currentVal int64
		err := db.Raw(fmt.Sprintf("SELECT last_value FROM %s_id_seq", table)).Scan(&currentVal).Error
		if err != nil {
			return false, err
		}
		return currentVal >= startThreshold, nil
	case "mysql":
		var autoInc int64
		err := db.Raw(`SELECT AUTO_INCREMENT FROM information_schema.TABLES 
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`, table).Scan(&autoInc).Error
		if err != nil {
			return false, err
		}
		return autoInc >= startThreshold, nil
	case "sqlite":
		// SQLite无法可靠判断，直接返回true，跳过初始化
		return true, nil
	default:
		return false, fmt.Errorf("unsupported dialect: %s", dialect)
	}
}

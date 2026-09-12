# 知识
用接口进行 值 扩展传递的。用 指针形式
例如：StandardHolder.IHolderRule
admin.manage.go
```golang
FindByJwt(jwt *jwtMid.Jwt) (rt rg.Rs[holderPg.HolderPg])
// 120行
pg.Rule = &rule

//-----------------
        err = ctx.Set(constContext.CTX_RULE, rt.Data.Rule)
		if err != nil {
			knife.Delete(ctx.Context(), constContext.CTX_RULE)
			err = ctx.Set(constContext.CTX_RULE, rt.Data.Rule)
			if err != nil {
				f.log.Warnf("Set.constContext.CTX_RULE.err= %+v", err)
			}
		}
		
		
//---------------
func ScopeRulePgWhere(ctx *gin.Context, tableName string) func(db *gorm.DB) *gorm.DB {
    return func (db *gorm.DB) *gorm.DB {
        var multiRule *MultiRule
        var iHolderRule interfaces.IHolderRule
        var multi *pg.Multi
        //
        get := ctx.Get(constContext.CTX_RULE)
        if nil != get {
            log.Debugf(context.Background(), log.TagAppDef,"multiRule:%+v", get)
            iHolderRule = get.(interfaces.IHolderRule)
            multiRule = iHolderRule.(*MultiRule)
        }
    }
	......
}
```
package serviceCore

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/blog/key/blogKeyPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CoreArticleCategory))
}

type CoreArticleCategory struct {
	sv  *repositoryBlog.BlogArticleCategoryRepository `autowire:"?"`
	rdt *bigCachePg.BatchString                       `autowire:"?"`
}

// GetAllByKeysRetMap
//
//	@Description: 获取所有缓存
//	@receiver c
//	@param ctx
//	@param tenantNo
//	@return []*entityBlog.BlogArticleCategoryEntity
//	@return bool
func (c *CoreArticleCategory) GetAllByKeysRetMap(ctx *gin.Context, keys []string, tenantNo string) (map[string]*modBlogArticleCategory.Cache, bool) {
	data := make(map[string]*modBlogArticleCategory.Cache, 0)
	tmps := make([]string, 0)
	for _, item := range keys {
		if strPg.IsNotBlank(item) {
			tmps = append(tmps, blogKeyPg.ArticleCategoryTenantNo(tenantNo, item))
		}
	}
	if len(tmps) == 0 {
		return data, false
	}
	result, b := c.rdt.GetAllByKeys(ctx, tmps)
	//fmt.Printf("获取缓存结果: %v\n", keys)
	//fmt.Printf("获取缓存结果: %v\n", result)
	if b {
		for _, item := range result {
			// 解析值（JSON 字符串）
			valueStr, ok := item.(string)
			if !ok {
				fmt.Printf("键 %s 的值格式错误\n", item)
				continue
			}
			// 可选：将 JSON 字符串解析为结构体
			var category modBlogArticleCategory.Cache
			if err := json.Unmarshal([]byte(valueStr), &category); err != nil {
				fmt.Printf("解析键 JSON 失败: %v\n", err)
				continue
			}
			data[category.No] = &category
		}
		return data, true
	}
	return data, false
}

// GetAllByCache
//
//	@Description: 获取所有缓存
//	@receiver c
//	@param ctx
//	@param tenantNo
//	@return []*entityBlog.BlogArticleCategoryEntity
//	@return bool
func (c *CoreArticleCategory) GetAllByCache(ctx *gin.Context, tenantNo string) ([]*modBlogArticleCategory.Cache, bool) {
	data := make([]*modBlogArticleCategory.Cache, 0)
	// 原来是把集合 key 交给 Lua（SMEMBERS + MGET）；现在由 GetAllBySetKey 在进程内完成同样的两步。
	setKey := blogKeyPg.ArticleCategoryTenantNoKeys(tenantNo)
	result, b := c.rdt.GetAllBySetKey(ctx, setKey)
	//fmt.Printf("获取缓存结果: %v\n", result)
	if b {
		//从集合中获取的键值对
		for i := 0; i < len(result); i += 2 {
			// 解析键名
			key := result[i].(string)
			// 解析值（JSON 字符串）
			valueStr, ok := result[i+1].(string)
			if !ok {
				fmt.Printf("键 %s 的值格式错误\n", key)
				continue
			}

			// 可选：将 JSON 字符串解析为结构体
			var category modBlogArticleCategory.Cache
			if err := json.Unmarshal([]byte(valueStr), &category); err != nil {
				fmt.Printf("解析键 %s 的 JSON 失败: %v\n", key, err)
				continue
			}
			data = append(data, &category)
		}
		return data, true
	}
	return data, false
}

// FormatTree 格式化为树形
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param tenantNo
//	@return []model.BaseNodeTree
//	@return bool
func (c *CoreArticleCategory) FormatTree(ctx *gin.Context, tenantNo string) ([]model.BaseNodeTree, bool) {
	cache, b := c.GetAllByCache(ctx, tenantNo)
	if b {
		tree := c.ConvertToTree(cache)
		return tree, true
	}
	return make([]model.BaseNodeTree, 0), b
}

// ConvertToTree 将 数组转换为 BaseNodeTree 树形结构
//
//	[]BaseNodeTree: 树形结构的根节点列表
func (c *CoreArticleCategory) ConvertToTree(entities []*modBlogArticleCategory.Cache) []model.BaseNodeTree {
	// 初始化节点Map和根节点列表
	nodeMap := make(map[string]*model.BaseNodeTree)
	rootNodes := make([]model.BaseNodeTree, 0)

	// 先将所有实体转换为节点并存入Map
	for _, entity := range entities {
		if entity == nil {
			continue
		}

		// 创建基础节点
		node := model.BaseNodeTree{
			Id:       entity.No,
			ParentId: entity.ParentNo,
			Name:     entity.Name,
			// 将整个实体作为扩展数据
			Extend:    entity,
			ChildData: make([]model.BaseNodeTree, 0),
		}

		// 存入Map，方便后续查找
		nodeMap[entity.No] = &node
	}

	// 构建父子关系
	for _, node := range nodeMap {
		if node == nil {
			continue
		}

		// 父ID为空/空字符串，作为根节点
		if node.ParentId == "" || node.ParentId == "0" {
			rootNodes = append(rootNodes, *node)
			continue
		}

		// 查找父节点并添加子节点
		parentNode, exists := nodeMap[node.ParentId]
		if exists && parentNode != nil {
			parentNode.ChildData = append(parentNode.ChildData, *node)
		} else {
			// 找不到父节点的，也作为根节点（容错处理）
			rootNodes = append(rootNodes, *node)
		}
	}

	return rootNodes
}

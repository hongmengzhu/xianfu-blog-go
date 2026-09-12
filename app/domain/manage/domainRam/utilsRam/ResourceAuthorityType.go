package utilsRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/resourceTypeCategoryPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
)

func ResourceAuthorityMark(tp resourceTypeCategoryPg.ResourceTypeCategory, groupId, resourceId string) string {
	return tp.String() + "|" + groupId + "|" + resourceId
}

func ResourceAuthorityMarkByInt64(tp resourceTypeCategoryPg.ResourceTypeCategory, groupId, resourceId int64) string {
	return ResourceAuthorityMark(tp, numberPg.Int64ToString(groupId), numberPg.Int64ToString(resourceId))
}
func ResourceAuthorityMarkByUint64(tp resourceTypeCategoryPg.ResourceTypeCategory, groupId, resourceId int64) string {
	return ResourceAuthorityMark(tp, numberPg.Int64ToString(groupId), numberPg.Int64ToString(resourceId))
}

func ResourceRelationMark(tp resourceTypeCategoryPg.ResourceTypeCategory, typeValue, resourceId string) string {
	return tp.String() + "|" + typeValue + "|" + resourceId
}

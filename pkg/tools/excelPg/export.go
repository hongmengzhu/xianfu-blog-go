package excelPg

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"github.com/xuri/excelize/v2"
	"go-spring.org/log"
)

// Letter 遍历a-z
func Letter(length int) []string {
	var str []string
	for i := 0; i < length; i++ {
		str = append(str, string(rune('A'+i)))
	}
	return str
}

// ExportExcelByMap 导出excel 数据源为[]map
func ExportExcelByMap(c *gin.Context, titleList []string, data []map[string]interface{}, fileName, sheetName string) error {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)
	header := make([]string, 0)
	for _, v := range titleList {
		header = append(header, v)
	}
	//表格样式
	style := excelize.Style{
		Font: &excelize.Font{
			Color:  "#666666",
			Size:   13,
			Family: "arial",
		},
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
	}
	rowStyleID, _ := f.NewStyle(&style)
	_ = f.SetSheetRow(sheetName, "A1", &header)
	_ = f.SetRowHeight(sheetName, 1, 30)
	length := len(titleList)
	headStyle := Letter(length)
	var lastRow string
	var widthRow string
	for k, v := range headStyle {
		if k == length-1 {
			lastRow = fmt.Sprintf("%s1", v)
			widthRow = v
		}
	}
	if err := f.SetColWidth(sheetName, "A", widthRow, 30); err != nil {
		log.Errorf(context.Background(), log.TagAppDef, err, "")
	}
	rowNum := 1
	for _, value := range data {
		row := make([]any, 0)
		var dataSlice []string
		for key := range value {
			dataSlice = append(dataSlice, key)
		}
		sort.Strings(dataSlice)
		for _, v := range dataSlice {
			if val, ok := value[v]; ok {
				row = append(row, val)
			}
		}
		rowNum++
		if err := f.SetSheetRow(sheetName, fmt.Sprintf("A%d", rowNum), &row); err != nil {
			log.Errorf(context.Background(), log.TagAppDef, err, "", err)
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s", lastRow), rowStyleID); err != nil {
			log.Errorf(context.Background(), log.TagAppDef, err, "", err)
		}

	}
	disposition := fmt.Sprintf("attachment; filename=%s-%s.xlsx", url.QueryEscape(fileName), time.Now().Format(datetimePg.YMDHIS))
	c.Header("Content-TypeCategory", "application/octet-stream")
	c.Header("Content-Disposition", disposition)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Writer)
}

// ExportExcelByStruct excel导出(数据源为Struct)
func ExportExcelByStruct(c *gin.Context, titleList []string, data []interface{}, fileName string, sheetName string) error {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)
	header := make([]string, 0)
	for _, v := range titleList {
		header = append(header, v)
	}
	style := excelize.Style{
		Font: &excelize.Font{
			Color:  "#666666",
			Size:   13,
			Family: "arial",
		},
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
	}
	rowStyleID, _ := f.NewStyle(&style)
	_ = f.SetSheetRow(sheetName, "A1", &header)
	_ = f.SetRowHeight("Sheet1", 1, 30)
	length := len(titleList)
	headStyle := Letter(length)
	var lastRow string
	var widthRow string
	for k, v := range headStyle {
		if k == length-1 {
			lastRow = fmt.Sprintf("%s1", v)
			widthRow = v
		}
	}
	if err := f.SetColWidth(sheetName, "A", widthRow, 30); err != nil {
		log.Errorf(context.Background(), log.TagAppDef, err, "", err)
	}
	rowNum := 1
	for _, v := range data {
		t := reflect.TypeOf(v)
		value := reflect.ValueOf(v)
		row := make([]interface{}, 0)
		for l := 0; l < t.NumField(); l++ {
			val := value.Field(l).Interface()
			row = append(row, val)
		}
		rowNum++
		err := f.SetSheetRow(sheetName, "A"+strconv.Itoa(rowNum), &row)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s", lastRow), rowStyleID)
		if err != nil {
			return err
		}
	}
	// 根据指定路径保存文件
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
	disposition := fmt.Sprintf("attachment; filename=%s-%s.xlsx", url.QueryEscape(fileName), time.Now().Format(datetimePg.YMDHIS))
	c.Header("Content-TypeCategory", "application/octet-stream")
	c.Header("Content-Disposition", disposition)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Writer)
}

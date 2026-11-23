package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Form 表单Form control 控件
type Form string

const (
	FormButton        Form = "PgButton"        //按钮   off/on
	FormCheckbox      Form = "PgCheckboxGroup" //复选框组
	FormColorPicker   Form = "PgColorPicker"   //颜色的控件
	FormDatePicker    Form = "PgDatePicker"    //日期的控件
	FormImage         Form = "PgImage"         //图片
	FormInput         Form = "PgInput"         //文本框
	FormInputTextarea Form = "PgInputTextarea" //多行文本框
	FormInputNumber   Form = "PgInputNumber"   //数值框
	FormMarkdown      Form = "PgMarkdown"      //Markdown 编辑器
	FormRadioGroup    Form = "PgRadioGroup"    //多选框组
	FormRate          Form = "PgRate"          //评分
	FormSelect        Form = "PgSelect"        //下拉选择
	FormStrengthMeter Form = "PgStrengthMeter" //密码框
	FormSwitch        Form = "PgSwitch"        //开关
	FormTimePicker    Form = "PgTimePicker"    //时间选择
	FormTree          Form = "PgTree"          //时间选择
	FormTreeSelect    Form = "PgTreeSelect"    //树下拉选择
	FormUpload        Form = "PgUpload"        //上传
	FormUploadGroup   Form = "PgUploadGroup"   //上传组
)

// Name 名称
func (this Form) Name() string {
	switch this {
	case "PgInput":
		return "文本框"
	case "PgButton":
		return "按钮"
	case "PgCheckboxGroup":
		return "复选框组"
	default:
		return "未知"
	}
}

// 值
func (this Form) String() string {
	return string(this)
}

// 值
func (this Form) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Form) IsEqual(id string) bool {
	return string(this) == id
}

var FormMap = map[string]enumBasePg.EnumString{
	FormButton.String():        enumBasePg.EnumString{Id: FormButton.String(), Name: FormButton.Name()},
	FormCheckbox.String():      enumBasePg.EnumString{Id: FormCheckbox.String(), Name: FormCheckbox.Name()},
	FormColorPicker.String():   enumBasePg.EnumString{Id: FormColorPicker.String(), Name: FormColorPicker.Name()},
	FormDatePicker.String():    enumBasePg.EnumString{Id: FormDatePicker.String(), Name: FormDatePicker.Name()},
	FormTimePicker.String():    enumBasePg.EnumString{Id: FormTimePicker.String(), Name: FormTimePicker.Name()},
	FormInput.String():         enumBasePg.EnumString{Id: FormInput.String(), Name: FormInput.Name()},
	FormInputTextarea.String(): enumBasePg.EnumString{Id: FormInputTextarea.String(), Name: FormInputTextarea.Name()},
	FormInputNumber.String():   enumBasePg.EnumString{Id: FormInputNumber.String(), Name: FormInputNumber.Name()},
	FormMarkdown.String():      enumBasePg.EnumString{Id: FormMarkdown.String(), Name: FormMarkdown.Name()},
	FormRadioGroup.String():    enumBasePg.EnumString{Id: FormRadioGroup.String(), Name: FormRadioGroup.Name()},
	FormSelect.String():        enumBasePg.EnumString{Id: FormSelect.String(), Name: FormSelect.Name()},
	FormStrengthMeter.String(): enumBasePg.EnumString{Id: FormStrengthMeter.String(), Name: FormStrengthMeter.Name()},
	FormSwitch.String():        enumBasePg.EnumString{Id: FormSwitch.String(), Name: FormSwitch.Name()},
	FormTree.String():          enumBasePg.EnumString{Id: FormTree.String(), Name: FormTree.Name()},
	FormTreeSelect.String():    enumBasePg.EnumString{Id: FormTreeSelect.String(), Name: FormTreeSelect.Name()},
	FormUpload.String():        enumBasePg.EnumString{Id: FormUpload.String(), Name: FormUpload.Name()},
	FormUploadGroup.String():   enumBasePg.EnumString{Id: FormUploadGroup.String(), Name: FormUploadGroup.Name()},
}

func IsExistForm(id string) (Form, bool) {
	_, ok := FormMap[id]
	return Form(id), ok
}

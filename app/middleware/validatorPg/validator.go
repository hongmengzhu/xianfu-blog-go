package validatorPg

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// Validate 验证器
var (
	validate *validator.Validate
	trans    ut.Translator
	uni      *ut.UniversalTranslator
)

func init() {
	//注册翻译器
	zh2 := zh.New()
	uni = ut.New(zh2, zh2)
	//
	trans, _ = uni.GetTranslator("zh")
	//获取gin的校验器
	validate = binding.Validator.Engine().(*validator.Validate)

	//validate = validatorPg.New()
	//注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		label := fld.Tag.Get("label")
		if label == "" {
			return fld.Name
		}
		return label
	})
	//注册翻译器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0) //无法初始化验证器，退出应用
	}
}

// Translate 翻译工具
func Translate(err error, s interface{}) map[string]string {
	r := make(map[string]string)
	t := reflect.TypeOf(s).Elem()
	// 使用安全类型断言，err 可能是 JSON 反序列化错误而非校验错误
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		return r
	}
	for _, err := range errs {
		//使用反射方法获取struct种的json标签作为key
		var k string
		if field, ok := t.FieldByName(err.StructField()); ok {
			k = field.Tag.Get("json")
		}
		if k == "" {
			k = err.StructField()
		}
		if k == "" {
			k = err.Field()
		}
		r[k] = err.Translate(trans)
	}
	return r
}

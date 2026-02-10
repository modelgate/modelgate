package db

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// HasRecrods 当前页码是否有记录
func HasRecrods(total int64, page int64, pageSize int64) bool {
	if total == 0 {
		return false
	}
	if pageSize > 0 && page > 0 {
		return total > (page-1)*pageSize
	}
	return true
}

// EscapeField 字段名称加上双引号，防止字段名与SQL关键字相同导致查询错误。
func EscapeField(field string) string {
	if strings.Contains(field, "(") { // 字段是一个函数，不需要转义
		return field
	} else if strings.Contains(field, "->") { // json_field->>'$.name' 或 json_field->'$.name'
		return field
	} else if arr := strings.Split(field, "."); len(arr) == 2 { // 字段名是 a.id 类型
		return fmt.Sprintf("%s.`%s`", arr[0], arr[1])
	} else {
		return fmt.Sprintf("`%s`", field) // 默认直接加上双引号``
	}
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsDbError(err error) bool {
	return err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
}

package db

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 15
	MaxPageSize     = 10000
)

// Option 查询选项
type Option func(db *gorm.DB) *gorm.DB

// Apply 应用选项
func Apply(dbc *gorm.DB, opt Option, others ...Option) *gorm.DB {
	dbc = opt(dbc)
	for _, opt := range others {
		dbc = opt(dbc)
	}
	return dbc
}

// WithFilter 对结果进行过滤
func WithFilter(filter any) Option {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		where, args, err := BuildQuery(filter)
		if err != nil {
			db.AddError(err)
			return db
		}
		if where != "" {
			db = db.Where(where, args...)
		}
		return db
	}
}

// WithTable 指定表名
func WithTable(table string) Option {
	return func(db *gorm.DB) *gorm.DB {
		if table != "" {
			return db.Table(table)
		}
		return db
	}
}

// WithField 查询字段
func WithField(fields string) Option {
	return func(db *gorm.DB) *gorm.DB {
		if fields != "" {
			return db.Select(fields)
		}
		return db
	}
}

// WithPaging 对结果进行分页
func WithPaging(page, size int64) Option {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = DefaultPageSize
	} else if size > MaxPageSize {
		size = MaxPageSize
	}
	offset := (page - 1) * size
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(size))
	}
}

// WithOrder 对结果进行排序
func WithOrder(order string, allowFields []string) Option {
	var orders []string
	for v := range strings.SplitSeq(order, ",") {
		if v = strings.TrimSpace(v); v != "" {
			order := "ASC"
			switch v[0] {
			case '-':
				v = v[1:]
				order = "DESC"
			case '+':
				v = v[1:]
			}
			if len(allowFields) == 0 || lo.Contains(allowFields, v) {
				orders = append(orders, fmt.Sprintf("%s %s", EscapeField(v), order))
			}
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		if len(orders) > 0 {
			return db.Order(strings.Join(orders, ","))
		}
		return db
	}
}

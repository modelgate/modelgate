package db

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm/schema"
)

// OmitOption 跳过条件
type OmitOption[T any] func(*omitOptions[T])

type omitOptions[T any] struct {
	omitIf func(v T) bool
}

// OmitIfZero 跳过零值
func OmitIfZero[T any]() OmitOption[T] {
	return func(o *omitOptions[T]) {
		o.omitIf = func(v T) bool {
			return isZero(v)
		}
	}
}

// OmitIfEqual 跳过等于值
func OmitIfEqual[T any](v T) OmitOption[T] {
	return func(o *omitOptions[T]) {
		o.omitIf = func(x T) bool {
			return reflect.DeepEqual(x, v)
		}
	}
}

// OmitIfNil 跳过 nil 值
func OmitIfNil[T any]() OmitOption[T] {
	return func(o *omitOptions[T]) {
		o.omitIf = func(v T) bool {
			rv := reflect.ValueOf(v)
			return rv.Kind() == reflect.Pointer && rv.IsNil()
		}
	}
}

// OmitIf 跳过满足条件的值
func OmitIf[T any](fn func(T) bool) OmitOption[T] {
	return func(o *omitOptions[T]) {
		o.omitIf = fn
	}
}

func applyOmitOptions[T any](v T, opts []OmitOption[T]) bool {
	if len(opts) == 0 {
		return false
	}
	var o omitOptions[T]
	for _, opt := range opts {
		opt(&o)
	}
	return o.omitIf != nil && o.omitIf(v)
}

// BuildQuery 构建查询
func BuildQuery(f any) (string, []any, error) {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	var conditions []string
	var values []any
	ns := schema.NamingStrategy{}

	for i := range t.NumField() {
		field := t.Field(i)
		fieldVal := v.Field(i)

		ef, ok := fieldVal.Interface().(exprFilter)
		if !ok {
			continue
		}
		exprs := ef.GetExprs()
		if len(exprs) == 0 {
			continue
		}
		args := ef.GetArgs()
		expected := 0
		for _, expr := range exprs {
			expected += strings.Count(expr, "?")
		}
		if expected != len(args) {
			return "", nil, fmt.Errorf("db build query failed on field %s: expected %d args, got %d", field.Name, expected, len(args))
		}

		var columnName string
		if tag := field.Tag.Get("gorm"); tag != "" {
			for part := range strings.SplitSeq(tag, ";") {
				part = strings.TrimSpace(part)
				if v, ok := strings.CutPrefix(part, "column:"); ok {
					columnName = v
					break
				}
			}
		}
		if columnName == "" {
			columnName = ns.ColumnName("", field.Name)
		}
		var parts []string
		for _, expr := range exprs {
			parts = append(parts, fmt.Sprintf(expr, columnName))
		}
		if len(parts) > 1 {
			conditions = append(conditions, "("+strings.Join(parts, " OR ")+")")
		} else {
			conditions = append(conditions, parts[0])
		}
		values = append(values, args...)
	}
	if len(conditions) == 0 {
		return "", nil, nil
	}
	return strings.Join(conditions, " AND "), values, nil
}

type exprFilter interface {
	GetExprs() []string
	GetArgs() []any
}

type F[T any] struct {
	Exprs []string
	Args  []any
}

func (f F[T]) GetExprs() []string {
	return f.Exprs
}

func (f F[T]) GetArgs() []any {
	return f.Args
}

func NewF[T any](exprs []string, args ...T) F[T] {
	values := make([]any, len(args))
	for i, arg := range args {
		values[i] = arg
	}
	return F[T]{Exprs: exprs, Args: values}
}

func NewEmptyF[T any]() F[T] {
	return F[T]{Exprs: nil, Args: nil}
}

func Eq[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s = ?"}, v)
}

func NotEq[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s != ?"}, v)
}

func Gt[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s > ?"}, v)
}

func Gte[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s >= ?"}, v)
}

func Lt[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s < ?"}, v)
}

func Lte[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s <= ?"}, v)
}

func Like[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s LIKE ?"}, v)
}

func NotLike[T any](v T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s NOT LIKE ?"}, v)
}

func In[T any](v []T, opts ...OmitOption[[]T]) F[[]T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[[]T]()
	}
	return NewF([]string{"%s IN (?)"}, v)
}

func NotIn[T any](v []T, opts ...OmitOption[[]T]) F[[]T] {
	if applyOmitOptions(v, opts) {
		return NewEmptyF[[]T]()
	}
	return NewF([]string{"%s NOT IN (?)"}, v)
}

func Between[T any](v1 T, v2 T, opts ...OmitOption[T]) F[T] {
	if applyOmitOptions(v1, opts) || applyOmitOptions(v2, opts) {
		return NewEmptyF[T]()
	}
	return NewF([]string{"%s BETWEEN ? AND ?"}, v1, v2)
}

func EqOrNull[T any](v T) F[T] {
	return NewF([]string{"%s = ?", "%s IS NULL"}, v)
}

func NotEqOrNull[T any](v T) F[T] {
	return NewF([]string{"%s != ?", "%s IS NULL"}, v)
}

func GtOrNull[T any](v T) F[T] {
	return NewF([]string{"%s > ?", "%s IS NULL"}, v)
}

func GteOrNull[T any](v T) F[T] {
	return NewF([]string{"%s >= ?", "%s IS NULL"}, v)
}

func LtOrNull[T any](v T) F[T] {
	return NewF([]string{"%s < ?", "%s IS NULL"}, v)
}

func LteOrNull[T any](v T) F[T] {
	return NewF([]string{"%s <= ?", "%s IS NULL"}, v)
}

func LikeOrNull[T any](v T) F[T] {
	return NewF([]string{"%s LIKE ?", "%s IS NULL"}, v)
}

func NotLikeOrNull[T any](v T) F[T] {
	return NewF([]string{"%s NOT LIKE ?", "%s IS NULL"}, v)
}

func InOrNull[T any](v []T) F[[]T] {
	return NewF([]string{"%s IN (?)", "%s IS NULL"}, v)
}

func NotInOrNull[T any](v []T) F[[]T] {
	return NewF([]string{"%s NOT IN (?)", "%s IS NULL"}, v)
}

func isZero(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		return rv.IsNil()
	}
	if z, ok := v.(interface{ IsZero() bool }); ok {
		return z.IsZero()
	}
	return rv.IsZero()
}

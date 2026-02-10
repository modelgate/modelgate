package db

import (
	"testing"
)

func TestBuildQuery(t *testing.T) {
	type TestFindFilter struct {
		ID       F[int64]
		IDs      F[[]int64] `gorm:"column:id"`
		Name     F[string]
		Status   F[int]
		Statuses F[[]int] `gorm:"column:status"`
	}

	dest := &TestFindFilter{
		ID:       Eq(int64(1)),
		IDs:      In([]int64{1, 2, 3}),
		Name:     Like("test"),
		Status:   NotEq(1),
		Statuses: In([]int{1, 2, 3}),
	}
	query, args, err := BuildQuery(dest)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("query: %s, args: %v", query, args)
}

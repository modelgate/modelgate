package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func JSONStr(v any) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func Unique[T comparable](arr []T) []T {
	if len(arr) == 0 {
		return arr
	}
	m := make(map[T]struct{})
	ret := make([]T, 0, len(arr))
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
			m[v] = struct{}{}
		}
	}
	return ret
}

func Ptr[T any](v T) *T {
	return &v
}

// GenApiKey returns a random api key
func GenApiKey() string {
	return fmt.Sprintf("sk-mg-api01-%s", strings.ReplaceAll(uuid.New().String(), "-", ""))
}

// MaskApiKey returns the first 16 characters and the last 4 characters of the key
func MaskApiKey(key string) (string, string) {
	return key[:16], key[len(key)-4:]
}

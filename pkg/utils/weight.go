package utils

import "math/rand"

// Weighter 权重接口
type Weighter interface {
	GetWeight() int
}

// PickByWeight 按照权重随机选择一个
func PickByWeight[T Weighter](list []T) T {
	var total int
	for _, item := range list {
		total += item.GetWeight()
	}
	n := rand.Intn(total)
	for _, item := range list {
		n -= item.GetWeight()
		if n < 0 {
			return item
		}
	}
	return list[0]
}

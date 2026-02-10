package utils

import (
	"github.com/google/uuid"
)

type UUIDv7 []byte

func NewUUIDv7() UUIDv7 {
	id, _ := uuid.NewV7()
	return id[:]
}

func (u UUIDv7) String() string {
	id, _ := uuid.FromBytes(u)
	return id.String()
}

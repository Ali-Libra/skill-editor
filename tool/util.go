package tool

import "github.com/google/uuid"

// GenerateUniqueID 返回一个全局唯一的 ID
func GenerateUniqueID() string {
	return uuid.New().String()
}

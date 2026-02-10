package common

import "github.com/gin-gonic/gin"

const (
	AccountIdKey = "accountId"
	ApiKeyIdKey  = "apiKeyId"
)

func SetAccountId(c *gin.Context, accountId int64) {
	c.Set(AccountIdKey, accountId)
}

func GetAccountId(c *gin.Context) int64 {
	return c.GetInt64(AccountIdKey)
}

func SetApiKeyId(c *gin.Context, apiKeyId int64) {
	c.Set(ApiKeyIdKey, apiKeyId)
}

func GetApiKeyId(c *gin.Context) int64 {
	return c.GetInt64(ApiKeyIdKey)
}

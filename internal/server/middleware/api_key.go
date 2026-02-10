package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/common"
)

// CheckApiKey  校验apiKey
func CheckApiKey(i do.Injector) gin.HandlerFunc {
	relayService := do.MustInvoke[relay.Service](i)
	return func(c *gin.Context) {
		accountApiKey, err := checkApiKey(c, relayService)
		if err != nil {
			log.Warnf("failed to check api key: %v", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		common.SetAccountId(c, accountApiKey.AccountId)
		common.SetApiKeyId(c, accountApiKey.ID)
	}
}

func checkApiKey(c *gin.Context, relayService relay.Service) (apiKey *model.AccountApiKey, err error) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		err = errors.New("Authorization header is required")
		return
	}
	auths := strings.SplitN(auth, " ", 2)
	if len(auths) != 2 {
		err = errors.New("Authorization header is invalid")
		return
	}
	apiKey, err = relayService.GetAccountApiKey(c, auths[1])
	if err != nil {
		err = errors.New("failed to get account api key")
		return
	}
	return apiKey, nil
}

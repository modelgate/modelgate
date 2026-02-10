package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type ModelPricingDao struct {
	*db.BaseDAO[model.ModelPricing, model.ModelPricingFilter]
}

func NewModelPricingDao(i do.Injector) (relay.ModelPricingDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &ModelPricingDao{
		BaseDAO: db.NewBaseDAO[model.ModelPricing, model.ModelPricingFilter](dbConn),
	}, nil
}

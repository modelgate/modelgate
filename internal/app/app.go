package app

import (
	"github.com/samber/do/v2"

	admv1 "github.com/modelgate/modelgate/internal/app/admin/v1"
	apiv1 "github.com/modelgate/modelgate/internal/app/api/v1"
)

func Init(i do.Injector) {
	// Register GRPC/HTTP Service
	do.Provide(i, admv1.NewAuthService)
	do.Provide(i, admv1.NewSystemService)
	do.Provide(i, admv1.NewRelayService)

	// api v1 service
	do.Provide(i, apiv1.NewRelayService)
}

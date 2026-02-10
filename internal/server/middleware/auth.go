package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/system"
)

func Auth(i do.Injector) func(connectHandler http.Handler) http.Handler {
	systemService := do.MustInvoke[system.Service](i)

	return func(connectHandler http.Handler) http.Handler {
		middleware := authn.NewMiddleware(func(ctx context.Context, req *http.Request) (any, error) {
			token, ok := authn.BearerToken(req)
			if !ok {
				if isUnauthorizeAllowedMethod(req.URL.Path) {
					return nil, nil
				}
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid authorization"))
			}
			userId, err := systemService.Authenticate(req.Context(), token)
			if err != nil {
				if isUnauthorizeAllowedMethod(req.URL.Path) {
					return nil, nil
				}
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid authorization"))
			}
			// check api permission
			isSuperAdmin, apiPermList, err := systemService.GetUserApiPermissions(ctx, userId)
			if err != nil {
				return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no permission"))
			}
			if !isSuperAdmin {
				isAllowed := false
				urlPath := strings.TrimPrefix(req.URL.Path, "/")
				for _, item := range apiPermList {
					if item.Path == urlPath && item.Method == req.Method {
						isAllowed = true
						break
					}
				}
				if !isAllowed && !isNotCheckPermissionMethod(req.URL.Path) {
					return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no permission"))
				}
			}
			authn.SetInfo(req.Context(), userId)
			return userId, nil
		})
		return middleware.Wrap(connectHandler)
	}
}

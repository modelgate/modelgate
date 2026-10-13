package middleware

import (
	"net/http"
	"strings"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/system"
)

func Auth(i do.Injector) func(connectHandler http.Handler) http.Handler {
	systemService := do.MustInvoke[system.Service](i)
	errorWriter := connecthttp.NewErrorWriter()

	return func(connectHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			token, ok := bearerToken(req)
			if !ok {
				if isUnauthorizeAllowedMethod(req.URL.Path) {
					connectHandler.ServeHTTP(writer, req)
					return
				}
				_ = errorWriter.Write(writer, req, connect.NewError(connect.CodeUnauthenticated, "invalid authorization"))
				return
			}
			userID, err := systemService.Authenticate(ctx, token)
			if err != nil {
				if isUnauthorizeAllowedMethod(req.URL.Path) {
					connectHandler.ServeHTTP(writer, req)
					return
				}
				_ = errorWriter.Write(writer, req, connect.NewError(connect.CodeUnauthenticated, "invalid authorization"))
				return
			}
			// check api permission
			isSuperAdmin, apiPermList, err := systemService.GetUserApiPermissions(ctx, userID)
			if err != nil {
				_ = errorWriter.Write(writer, req, connect.NewError(connect.CodePermissionDenied, "no permission"))
				return
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
					_ = errorWriter.Write(writer, req, connect.NewError(connect.CodePermissionDenied, "no permission"))
					return
				}
			}
			ctx = system.ContextWithUserID(ctx, userID)
			connectHandler.ServeHTTP(writer, req.WithContext(ctx))
		})
	}
}

func bearerToken(req *http.Request) (string, bool) {
	const prefix = "Bearer "
	authorization := req.Header.Get("Authorization")
	if len(authorization) < len(prefix) || !strings.EqualFold(authorization[:len(prefix)], prefix) {
		return "", false
	}
	return authorization[len(prefix):], true
}

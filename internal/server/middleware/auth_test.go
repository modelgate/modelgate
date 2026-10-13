package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/do/v2"
	"go.uber.org/mock/gomock"

	"github.com/modelgate/modelgate/internal/system"
)

func TestAuthRejectsMissingTokenWithConnectV2Error(t *testing.T) {
	handler := newAuthTestHandler(t, func(http.ResponseWriter, *http.Request) {
		t.Fatal("protected handler was called")
	})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, newConnectRequest("/admin.v1.SystemService/GetUserList"))

	assertConnectError(t, recorder, http.StatusUnauthorized, "unauthenticated")
}

func TestAuthAddsUserIDToContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	systemService := system.NewMockService(ctrl)
	systemService.EXPECT().Authenticate(gomock.Any(), "test-token").Return(int64(42), nil)
	systemService.EXPECT().GetUserApiPermissions(gomock.Any(), int64(42)).Return(true, nil, nil)

	handler := authTestMiddleware(systemService)(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		userID, ok := system.UserIDFromContext(req.Context())
		if !ok || userID != 42 {
			t.Fatalf("UserIDFromContext() = (%d, %t), want (42, true)", userID, ok)
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	recorder := httptest.NewRecorder()
	req := newConnectRequest("/admin.v1.SystemService/GetUserList")
	req.Header.Set("Authorization", "Bearer test-token")

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestAuthRejectsMissingPermissionWithConnectV2Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	systemService := system.NewMockService(ctrl)
	systemService.EXPECT().Authenticate(gomock.Any(), "test-token").Return(int64(42), nil)
	systemService.EXPECT().GetUserApiPermissions(gomock.Any(), int64(42)).Return(false, nil, nil)

	handler := authTestMiddleware(systemService)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("protected handler was called")
	}))
	recorder := httptest.NewRecorder()
	req := newConnectRequest("/admin.v1.SystemService/GetUserList")
	req.Header.Set("Authorization", "Bearer test-token")

	handler.ServeHTTP(recorder, req)

	assertConnectError(t, recorder, http.StatusForbidden, "permission_denied")
}

func newAuthTestHandler(t *testing.T, handler http.HandlerFunc) http.Handler {
	t.Helper()
	ctrl := gomock.NewController(t)
	return authTestMiddleware(system.NewMockService(ctrl))(handler)
}

func authTestMiddleware(systemService system.Service) func(http.Handler) http.Handler {
	injector := do.New()
	do.ProvideValue[system.Service](injector, systemService)
	return Auth(injector)
}

func newConnectRequest(path string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, path, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertConnectError(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if recorder.Code != status {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, status, recorder.Body.String())
	}
	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != code {
		t.Fatalf("code = %q, want %q", body.Code, code)
	}
}

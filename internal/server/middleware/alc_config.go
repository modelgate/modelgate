package middleware

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

var authenticationAllowlistMethods = map[string]bool{
	"/admin.v1.AuthService/Login":               true,
	"/admin.v1.AuthService/RefreshToken":        true,
	"/admin.v1.AuthService/SignOut":             true,
	"/admin.v1.AuthService/SignUp":              true,
	"/admin.v1.SystemService/GetConstantRoutes": true,
	"/admin.v1.SystemService/GetVersion":        true,
}

// isUnauthorizeAllowedMethod returns whether the method is exempted from authentication.
func isUnauthorizeAllowedMethod(fullMethodName string) bool {
	log.Info("fullMethodName", fullMethodName)
	if strings.HasPrefix(fullMethodName, "/openapi.") {
		return true
	} else {
		return authenticationAllowlistMethods[fullMethodName]
	}
}

var notCheckPermissionMethods = map[string]bool{
	"/admin.v1.AuthService/GetUserInfo":     true,
	"/admin.v1.SystemService/GetUserRoutes": true,
}

func isNotCheckPermissionMethod(fullMethodName string) bool {
	if authenticationAllowlistMethods[fullMethodName] {
		return true
	}
	return notCheckPermissionMethods[fullMethodName]
}

var allowedMethodsOnlyForAdmin = map[string]bool{
	"/admin.v1.UserService/CreateUser": true,
}

// isOnlyForAdminAllowedMethod returns true if the method is allowed to be called only by admin.
func isOnlyForAdminAllowedMethod(methodName string) bool {
	return allowedMethodsOnlyForAdmin[methodName]
}

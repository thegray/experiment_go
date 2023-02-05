package middleware

import (
	"strings"

	"experiment_go/kafka/dyn_consumer/internal/pkg/errors"
	"experiment_go/kafka/dyn_consumer/internal/pkg/logger"
	"experiment_go/kafka/dyn_consumer/internal/service/auth"
	echo "github.com/labstack/echo/v4"
)

const bearerPrefix string = "Bearer "

func BearerAuthMiddleware(tokenParser auth.TokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errors.ErrAuthInvalidToken("Authorization header not found", nil)
			}

			// Get token
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return errors.ErrAuthInvalidToken("Authorization header must start with 'Bearer '", nil)
			}
			jwtToken := strings.TrimPrefix(authHeader, bearerPrefix)

			// Parse token
			payload, err := tokenParser.ParseToken(jwtToken)
			if err != nil {
				return err
			}

			// Set to context
			c.Set("email", payload.Email)
			c.Set("roleID", payload.RoleID)
			c.Set("permissions", payload.Permissions)

			return next(c)
		}
	}
}

func GetEmail(c echo.Context) string {
	email, ok := c.Get("email").(string)
	if !ok {
		ctx := GetContext(c)
		logger.ErrorCtx(ctx, "[GetEmail] context value email not set")
		return ""
	}

	return email
}

func GetRoleID(c echo.Context) uint {
	roleID, ok := c.Get("roleID").(uint)
	if !ok {
		ctx := GetContext(c)
		logger.ErrorCtx(ctx,
			"[Auth:GetRoleID]",
			"error", "context value roleID not set")
		return 0
	}

	return roleID
}

func GetPermissions(c echo.Context) []string {
	permissions, ok := c.Get("permissions").([]string)
	if !ok {
		ctx := GetContext(c)
		logger.ErrorCtx(ctx,
			"[Auth:GetPermissions]",
			"error", "context value permissions not set")
		return []string{}
	}

	return permissions
}

func PermissionMiddleware(requiredPermission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			permissions := GetPermissions(c)
			for _, permission := range permissions {
				if permission == requiredPermission {
					return next(c) // If permission match continue
				}
			}

			return errors.ErrAuthForbidden("User don't have required permission", nil)
		}
	}
}

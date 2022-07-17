package security

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type SessionHandler func(s *WebSecurity) gin.HandlerFunc

// Setup setup the security middleware.
// You may pass your custom additional session handler to this function.
func Setup(r *gin.Engine, w *WebSecurity, handlers ...SessionHandler) {
	autoMigrateSecurityPath()
	autoMigrateSecuritySession()
	r.Use(middleware(w))
	for _, handler := range handlers {
		r.Use(handler(w))
	}
}

// is public? yes -> pass, no -> next checking
// has session? yes -> pass, no -> fail
// (optional) match ip? yes -> pass, no -> fail
// has role? yes -> pass, no -> fail
func middleware(w *WebSecurity) gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path
		method := strings.ToUpper(c.Request.Method)

		potential := findPotentialPaths(method, path)
		if checkIsPublic(potential, method, path) {
			return
		}

		sessionId, err := c.Request.Cookie("sessionId")
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		session, err := findSessionById(sessionId.Value)
		if err != nil || session.SessionId == "" {
			c.AbortWithStatus(401)
			return
		}

		if session.IsExpired() {
			c.AbortWithStatus(401)
			return
		}

		if w.SecurityConfig.RestrictIP && session.ClientIP != w.GetRequestIP(c) {
			c.AbortWithStatus(401)
			return
		}

		if !checkRole(potential, method, path, strings.Split(session.Role, ",")) {
			c.AbortWithStatus(401)
			return
		}
	}
}

func checkIsPublic(potentialPaths []SecurityPath, method string, path string) bool {
	for _, potentialPath := range potentialPaths {
		match, _ := regexp.MatchString(potentialPath.PathRegex, path)
		if match && potentialPath.Method == method {
			return potentialPath.isPublic
		}
	}
	return false
}

func checkRole(potentialPaths []SecurityPath, method string, path string, roles []string) bool {
	for _, potentialPath := range potentialPaths {
		match, _ := regexp.MatchString(potentialPath.PathRegex, path)
		if match && potentialPath.Method == method {
			for _, role := range roles {
				if potentialPath.Role == "*" || potentialPath.Role == role {
					return true
				}
			}
		}
	}
	return false
}

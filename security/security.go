package security

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WebSecurity manages the security configs and the security table.
type WebSecurity struct {
	SecurityConfig *SecurityConfig
}

// BuildSecurityTable build the security table
func (w *WebSecurity) BuildSecurityTable(config *SecurityConfig) {
	deleteSecurityPathTable()
	autoMigrateSecurityPath()

	w.SecurityConfig = config

	// add public paths
	for _, path := range config.PublicPaths {
		createSecurityPath(SecurityPath{
			Method:    strings.ToUpper(path.Method),
			PathRegex: path.PathRegex,
			Role:      "*",
			isPublic:  true,
		})
	}

	// add secured paths
	for _, path := range config.SecuredPaths {
		createSecurityPath(SecurityPath{
			Method:    strings.ToUpper(path.Method),
			PathRegex: path.PathRegex,
			Role:      strings.Join(path.Roles, ","),
			isPublic:  false,
		})
	}
}

// GetRequestMethod get the request method.
func (w *WebSecurity) GetRequestMethod(c *gin.Context) string {
	return strings.ToUpper(c.Request.Method)
}

// GetRequestPath get the request path.
func (w *WebSecurity) GetRequestPath(c *gin.Context) string {
	return c.Request.URL.Path
}

// GetRequestIP get the request ip.
func (w *WebSecurity) GetRequestIP(c *gin.Context) string {
	return c.ClientIP()
}

// GetSessionId get the session id from the cookie.
func (w *WebSecurity) GetSessionId(c *gin.Context) string {
	cookie, err := c.Request.Cookie("sessionId")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// GetCurrentSession get the current session.
func (w *WebSecurity) GetCurrentSession(c *gin.Context) *SecuritySession {
	sessionId := w.GetSessionId(c)
	session, err := findSessionById(sessionId)
	if err != nil {
		return nil
	}
	return session
}

// GetRefreshToken get the refresh token.
func (w *WebSecurity) GetRefreshToken(c *gin.Context) string {
	return c.Request.Header.Get("Authorization")
}

// Login create a new session.
func (w *WebSecurity) Login(c *gin.Context, userId uint, role []string) *SecuritySession {
	if !w.SecurityConfig.MultipleSession {
		session, _ := findSessionByUserId(userId)
		if session != nil {
			deleteSessionById(session.SessionId)
		}
	}

	session := SecuritySession{
		SessionId: uuid.New().String(),
		UserId:    userId,
		Role:      strings.Join(role, ","),
		Expire:    time.Now().Add(w.SecurityConfig.SessionExpire),
	}

	if w.SecurityConfig.RefreshToken {
		session.RefreshToken = uuid.New().String()
		session.RefreshExpire = time.Now().Add(w.SecurityConfig.RefreshExpire)
	}

	if w.SecurityConfig.RestrictIP {
		session.ClientIP = c.ClientIP()
	}

	createSession(&session)
	w.setCookie(c, session.SessionId)

	return &session
}

// Logout delete the current session.
func (w *WebSecurity) Logout(c *gin.Context) {
	sessionId := w.GetSessionId(c)
	deleteSessionById(sessionId)
	c.SetCookie("sessionId", "", -1, "/", "", w.SecurityConfig.OnlyHttps, true)
}

// RefreshSession refresh the session.
// If the session is expired, it cannot be refreshed.
func (w *WebSecurity) RefreshSession(c *gin.Context) {
	session := w.GetCurrentSession(c)
	if session == nil {
		return
	}

	if session.Expire.Before(time.Now()) {
		return
	}

	if w.SecurityConfig.RestrictIP && session.ClientIP != w.GetRequestIP(c) {
		return
	}

	session.SessionId = uuid.New().String()
	session.Expire = time.Now().Add(w.SecurityConfig.SessionExpire)
	updateSession(session)

	w.setCookie(c, session.SessionId)
}

// RefreshWithToken refresh the session with token.
// Session can be expired or not.
func (w *WebSecurity) RefreshWithToken(c *gin.Context) {
	session := w.GetCurrentSession(c)
	if session == nil {
		return
	}

	if session.RefreshExpire.Before(time.Now()) {
		return
	}

	if w.GetRefreshToken(c) != session.RefreshToken {
		return
	}

	if w.SecurityConfig.RestrictIP && session.ClientIP != w.GetRequestIP(c) {
		return
	}

	session.SessionId = uuid.New().String()
	session.Expire = time.Now().Add(w.SecurityConfig.SessionExpire)
	session.RefreshExpire = time.Now().Add(w.SecurityConfig.RefreshExpire)
	session.RefreshToken = uuid.New().String()
	updateSession(session)

	w.setCookie(c, session.SessionId)
}

func (w *WebSecurity) setCookie(c *gin.Context, sessionId string) {
	c.SetCookie("sessionId", sessionId, int(w.SecurityConfig.SessionExpire.Seconds()), "/", "", w.SecurityConfig.OnlyHttps, true)
}

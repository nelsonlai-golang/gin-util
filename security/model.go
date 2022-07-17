package security

import (
	"time"

	"gorm.io/gorm"
)

// PublicPath is a path that is public.
type PublicPath struct {
	Method    string
	PathRegex string
}

// SecurityPath is a path that is secured.
type SecuredPath struct {
	Method    string
	PathRegex string
	Roles     []string
}

// SecurityConfig is the configuration of the security.
type SecurityConfig struct {
	PublicPaths     []PublicPath
	SecuredPaths    []SecuredPath
	SessionExpire   time.Duration
	RefreshExpire   time.Duration
	RefreshToken    bool
	MultipleSession bool
	RestrictIP      bool
	OnlyHttps       bool
}

// SecuritySession is the session of the security.
type SecuritySession struct {
	gorm.Model
	SessionId     string
	UserId        uint
	ClientIP      string
	Role          string
	Expire        time.Time
	RefreshToken  string
	RefreshExpire time.Time
}

func (s *SecuritySession) IsExpired() bool {
	return s.Expire.Before(time.Now())
}

// SecurityPath is the paths registered in the security.
type SecurityPath struct {
	gorm.Model
	Method    string
	PathRegex string
	Role      string
	isPublic  bool
}

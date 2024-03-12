package service

import (
	"sync"
	"time"
)

type Cookie struct {
	Token       string
	UserId      int64
	AddedTime   time.Time
	ExpiresTime time.Time
}

var CookieInfo = sync.Map{}

func IsValidCookie(cookie string) bool {
	v, ok := CookieInfo.Load(cookie)
	if !ok {
		return false
	}
	cookieInfo := v.(Cookie)
	if cookieInfo.ExpiresTime.Before(time.Now()) {
		return false
	}
	return true
}

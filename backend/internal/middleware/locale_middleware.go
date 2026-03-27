package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	localeKey        = "locale"
	LocaleCookieName = "final_whistle_locale"
	DefaultLocale    = "en"
)

func ResolveLocale() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := normalizeLocale(c.GetHeader("X-Final-Whistle-Locale"))
		if locale == "" {
			cookieValue, err := c.Cookie(LocaleCookieName)
			if err == nil {
				locale = normalizeLocale(cookieValue)
			}
		}
		if locale == "" {
			locale = DefaultLocale
		}

		c.Set(localeKey, locale)
		c.Next()
	}
}

func CurrentLocale(c *gin.Context) string {
	value, ok := c.Get(localeKey)
	if !ok {
		return DefaultLocale
	}

	locale, ok := value.(string)
	if !ok || locale == "" {
		return DefaultLocale
	}

	return locale
}

func normalizeLocale(locale string) string {
	switch strings.ToLower(strings.TrimSpace(locale)) {
	case "zh", "zh-cn", "zh-hans":
		return "zh"
	case "en", "en-us", "en-gb":
		return "en"
	default:
		return ""
	}
}

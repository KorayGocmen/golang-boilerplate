package context

import (
	gocontext "context"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
)

func NewFiberCtx(c *fiber.Ctx) (Ctx, gocontext.CancelFunc) {
	// Create a new context.
	ctx, cancel := WithCancel(Background())

	remoteIP := net.ParseIP(c.Context().RemoteIP().String())
	if !env.IsDev() && (remoteIP.IsPrivate() || remoteIP.IsLoopback() || remoteIP.IsUnspecified()) {
		remoteIP = net.ParseIP(c.Get(fiber.HeaderXForwardedFor))
	}

	ctx = WithValue(ctx, KeyRemoteIP, remoteIP.String())
	ctx = WithValue(ctx, KeyRequestID, c.Locals(string(KeyRequestID)))
	ctx = WithValue(ctx, KeyMethod, c.Method())
	ctx = WithValue(ctx, KeyPath, string(c.Context().Path()))
	ctx = WithValue(ctx, KeyQuery, string(c.Context().QueryArgs().QueryString()))

	ctx = WithValue(ctx, KeyLang, LangTR)
	if lang := c.Get(fiber.HeaderAcceptLanguage); len(lang) >= 2 {
		ctx = WithValue(ctx, KeyLang, ToLanguage(strings.TrimSpace(lang)[:2]))
	}

	return ctx, cancel
}

func FromFiberCtx(c *fiber.Ctx) Ctx {
	// Return context if already set.
	if ctx := c.Locals("ctx"); ctx != nil {
		return ctx.(Ctx)
	}

	// Technically this should never happen.
	ctx, cancel := NewFiberCtx(c)
	c.Locals("ctx", ctx)
	c.Locals("cancel", cancel)

	return ctx
}

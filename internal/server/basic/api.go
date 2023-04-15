package basic

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labi-le/server/internal"
	"github.com/labi-le/server/pkg/response"
	"net/http"
)

func RegisterHandlers(r fiber.Router, reply *response.Reply, link string) {
	res := &resource{
		reply: reply,
	}

	r.Get("/", res.HomePage)
	r.Get("version", res.Version)
	r.Get("discord", func(ctx *fiber.Ctx) error {
		return ctx.Redirect(link, http.StatusMovedPermanently)
	})
}

type resource struct {
	reply *response.Reply
}

func (r *resource) Version(ctx *fiber.Ctx) error {
	return r.reply.OK(ctx, internal.BuildVersion())
}

func (r *resource) HomePage(ctx *fiber.Ctx) error {
	return r.reply.OK(ctx, "privet")
}

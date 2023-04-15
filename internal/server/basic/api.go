package basic

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labi-le/server/internal"
	"github.com/labi-le/server/pkg/response"
)

func RegisterHandlers(r fiber.Router, reply *response.Reply) {
	res := &resource{
		reply: reply,
	}

	r.Get("/", res.HomePage)
	r.Get("version", res.Version)
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

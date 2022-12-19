package file

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"labile-me-serv/pkg/config"
	"labile-me-serv/pkg/log"
	"labile-me-serv/pkg/response"
	"net/http"
	"time"
)

var (
	ErrInvalidForm = errors.New("invalid form")
	ErrEmptyFile   = errors.New("file is empty")
	ErrInvalidKey  = errors.New("invalid key")
)

func RegisterHandlers(r fiber.Router, l log.Logger, s Service, c config.Config) {
	res := &resource{
		log:      l,
		s:        s,
		reply:    response.New(l),
		ownerKey: c.GetOwnerKey(),
	}

	r.Put("*", res.Upload)
	r.Get("*", res.Get)
}

type resource struct {
	log   log.Logger
	s     Service
	reply *response.Reply

	ownerKey string
}

func (r *resource) Upload(ctx *fiber.Ctx) error {
	customURL := ctx.Params("*")
	if customURL != "" {
		if !checkKey(ctx, r.ownerKey) {
			return r.reply.Unauthorized(ctx, ErrInvalidKey)
		}

	} else {
		customURL = Short(time.Now().Nanosecond())
	}

	// multipart form
	header, err := ctx.FormFile("file")
	if err != nil {
		return r.reply.BadRequest(ctx, ErrInvalidForm)
	}

	if header.Size == 0 {
		return r.reply.BadRequest(ctx, ErrEmptyFile)
	}

	mpFile, opErr := header.Open()
	if opErr != nil {
		return r.reply.BadRequest(ctx, ErrInvalidFile)
	}

	defer mpFile.Close()

	contentType, mimeErr := getContentType(mpFile)
	if mimeErr != nil {
		return r.reply.InternalServerError(ctx, ErrContentTypeAssertion)
	}

	filename := customURL + contentType.Extension()
	//filename := customURL + ".jpg"

	req := RequestFile{
		Name:        filename,
		ShortID:     customURL,
		ContentType: contentType.String(),
		//ContentType: "jpeg",
		Reader: mpFile,
	}

	add, sErr := r.s.Add(ctx.Context(), req)
	if errors.Is(sErr, ErrFileExists) {
		return r.reply.Conflict(ctx, fiber.Map{
			"short_id": add,
			"error":    sErr.Error(),
		})
	}

	if sErr != nil {
		return r.reply.InternalServerError(ctx, sErr)
	}

	return r.reply.Created(ctx, fiber.Map{"short_id": add})
}

func (r *resource) Get(ctx *fiber.Ctx) error {
	short := ctx.Params("*")
	if short == "" {
		return r.reply.BadRequest(ctx, ErrInvalidForm)
	}

	file, err := r.s.Get(ctx.Context(), short)
	if err != nil {
		return r.reply.NotFound(ctx, err)
	}

	ctx.Set("Content-Type", file.ContentType)

	return ctx.
		Status(http.StatusOK).
		SendStream(file)
}

func checkKey(ctx *fiber.Ctx, key string) bool {
	return ctx.Get("authorization") == key
}
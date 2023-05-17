package presentation

import (
	"shortner/domain"

	"github.com/gofiber/fiber/v2"
)

type httpRedirect struct {
	engine *fiber.App
	app    domain.Service
}

func NewHttpRedirect(s domain.Service) *httpRedirect {
	return &httpRedirect{
		engine: fiber.New(),
		app:    s,
	}
}

func (H *httpRedirect) Start(listen string) {
	H.engine.Get("/+", H.Redirect)
	H.engine.Listen(listen)

}
func (H *httpRedirect) Redirect(c *fiber.Ctx) error {
	id := c.Params("+")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "no redicrection found"})
	}
	data, err := H.app.Find(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"could not find " + id: " err > " + err.Error()})
	}
	return c.Redirect(data.Url, 301)
}

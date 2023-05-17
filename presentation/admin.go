package presentation

import (
	"shortner/domain"
	"shortner/libs"

	"github.com/gofiber/fiber/v2"
)

type httpApi struct {
	engine *fiber.App
	app    domain.Service
}

func NewAdmin(s domain.Service) *httpApi {
	return &httpApi{
		engine: fiber.New(),
		app:    s,
	}
}

func (H *httpApi) Start(listen string) {
	H.engine.Get("/find", H.Find)
	H.engine.Get("/searchurl", H.SearchUrl)
	H.engine.Post("/store", H.Store)
	H.engine.Delete("/delete", H.Delete)
	H.engine.Listen(listen)

}
func (H *httpApi) Find(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send id"})
	}
	data, err := H.app.Find(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"could not find " + id: " err > " + err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"result": libs.STM(data)})
}

func (H *httpApi) Store(c *fiber.Ctx) error {
	url := c.Query("url")
	if url == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send url"})
	}
	data, err := H.app.Store(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"could not store to redis: ": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"result": libs.STM(data)})
}

func (H *httpApi) SearchUrl(c *fiber.Ctx) error {
	url := c.Query("url")
	if url == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send url"})
	}
	data, err := H.app.SearchUrl(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"could not find url: " + url: " err > " + err.Error()})
	}
	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "not found"})
	}
	var output []map[string]interface{}
	for _, i := range data {
		d := libs.STM(i)
		output = append(output, d)
	}
	return c.Status(200).JSON(fiber.Map{"result": output})
}

func (H *httpApi) Delete(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send id"})
	}
	result, err := H.app.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"could not delete: " + id: err.Error()})
	}
	if !result {
		return c.Status(404).JSON(fiber.Map{"key not found to delete: " + id: err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"key deleted ": id})
}

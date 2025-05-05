package main

import (
	"fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type Link struct {
	OriginalUrl string `validate:"required,min=5"`
	ShortUrl    string
}

var validate = validator.New()

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {

		return c.Render("index", fiber.Map{})
	})

	app.Post("/shorten", storeUrl)

	app.Use(func(c *fiber.Ctx) error {
		log.Error("halaman tidak ditemukan")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"message": "halaman tidak ditemukan",
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func storeUrl(c *fiber.Ctx) error {
	var link Link
	fmt.Println("Received data:", link)
	if err := c.BodyParser(&link); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "invalid request",
		})
	}

	if err := validate.Struct(link); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"detail": err.Error(),
		})
	}

	link.ShortUrl = generateShortLink()

	// Return the result
	return c.JSON(fiber.Map{
		"message":   "URL shortened successfully",
		"original":  link.OriginalUrl,
		"short_url": link.ShortUrl,
	})
}

func generateShortLink() string {
	randomCode := generateRandomString(6)
	urlShort := fmt.Sprintf("www.urlty.link/%s", randomCode)
	return urlShort
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	fmt.Println(string(result))
	return string(result)
}

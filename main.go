package main

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

type Link struct {
	OriginalUrl string `validate:"required,min=5"`
	ShortUrl    string
	UniqueKey   string
}

var store = session.New()
var validate = validator.New()

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Use(logger.New())

	app.Get("/", index)

	app.Get("/:uniqueKey", RedirectHandler)

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

func RedirectHandler(c *fiber.Ctx) error {
	uniqueKey := c.Params("uniqueKey")

	originalURL, err := getOriginalURL(uniqueKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("URL not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	fmt.Println("visiting", originalURL)
	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
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

	link.ShortUrl, link.UniqueKey = generateShortLink()
	saveurltodb(link.OriginalUrl, link.ShortUrl, link.UniqueKey)
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	sess.Set("message", "URL shortened successfully")
	sess.Set("original", link.OriginalUrl)
	sess.Set("short_url", link.ShortUrl)
	sess.Save()

	// Redirect ke halaman utama
	return c.Redirect("/")

}

func index(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	data := fiber.Map{}

	if msg := sess.Get("message"); msg != nil {
		data["message"] = msg
		data["original"] = sess.Get("original")
		data["short_url"] = sess.Get("short_url")

		sess.Delete("message")
		sess.Delete("original")
		sess.Delete("short_url")
		sess.Save()
	}

	return c.Render("index", data)
}

func generateShortLink() (string, string) {
	randomCode := generateRandomString(5)
	urlShort := fmt.Sprintf("https://susut.ink/%s", randomCode)
	return urlShort, randomCode
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

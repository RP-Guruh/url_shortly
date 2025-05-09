package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

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
	app.Static("/images", "./images")

	app.Get("/", index)

	app.Get("/:uniqueKey", RedirectHandler)

	app.Post("/shorten", storeUrl)
	app.Post("/track/link", trackLink)
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
	fmt.Println("masuk redirect handler")

	uniqueKey := c.Params("uniqueKey")

	originalURL, err := getOriginalURL(uniqueKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("URL not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	ip := c.IP()
	info := getInformation(c, ip)

	fmt.Println("Visitor Info:", info)

	if _, err := storeVisitor(uniqueKey, info); err != nil {
		fmt.Println("Error storing visitor:", err)
		// Lanjutkan redirect meskipun gagal simpan
	}

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
	// Ambil session dari store
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	// Siapkan data untuk dikirim ke template
	data := fiber.Map{}

	// Ambil pesan dan data lainnya dari session jika ada
	if msg := sess.Get("message"); msg != nil {
		data["message"] = msg
		data["original"] = sess.Get("original")
		data["short_url"] = sess.Get("short_url")

		// Hapus data setelah digunakan
		sess.Delete("message")
		sess.Delete("original")
		sess.Delete("short_url")
		sess.Save()
	}

	// Ambil pesan tracking dan data lainnya dari session jika ada
	if msg_track := sess.Get("message_track"); msg_track != nil {
		data["message_track"] = msg_track
		data["visitorTotal"] = sess.Get("visitorTotal")

		// Hapus data setelah digunakan
		sess.Delete("message_track")
		sess.Delete("visitorTotal")
		sess.Save()
	}

	// Ambil nilai tab aktif dari query string (default "shorten")

	// Render halaman dengan data dan tab aktif
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

func trackLink(c *fiber.Ctx) error {
	// Ambil data "code" dari form

	fullUrl := c.FormValue("code")
	parsedUrl, err := url.Parse(fullUrl)
	if err != nil {
		return c.Status(400).SendString("Invalid URL")
	}

	// Ambil path URL setelah '/'
	unique := strings.TrimPrefix(parsedUrl.Path, "/")

	// Hitung jumlah visitor
	visitorCount, err := getMyVisitor(unique)
	if err != nil {

		return c.Status(500).SendString("Internal Server Error")
	}

	sess, err := store.Get(c)
	sess.Set("message_track", "Jumlah Visitor")
	sess.Set("visitorTotal", visitorCount)
	sess.Save()

	// Kirim response dengan visitorCount
	return c.Redirect("/")
}

package imagectrl

import (
	"image/jpeg"
	"io"
	"log"
	"os"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"github.com/nfnt/resize"
	"pentag.kr/dimimonster/config"
	"pentag.kr/dimimonster/middleware"
	"pentag.kr/dimimonster/models"
	"pentag.kr/dimimonster/utils/crypt"
)

func UploadImage(c *fiber.Ctx) error {
	userID := middleware.GetUserIDFromMiddleware(c)
	token := c.FormValue("token", "")
	if !crypt.RecaptchaCheck(token, "image_upload") {
		return c.Status(425).SendString("Recaptcha Failed")
	}
	description := c.FormValue("description", "")
	if utf8.RuneCountInString(description) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description is too long",
		})
	}

	location := c.FormValue("location", "")
	exist := func() bool {
		for _, l := range config.LocationList {
			if location == l {
				return true
			}
		}
		return false
	}
	if !exist() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Location is not valid",
		})
	}

	// download fileHeader as io.Reader
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	} else if fileHeader.Size > 1024*1024 { // check size is over 1024KB
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size is too big",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	defer file.Close()

	image, err := jpeg.Decode(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	x := image.Bounds().Dx()
	y := image.Bounds().Dy()
	if x > 1024 || y > 1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image size is too big",
		})
	} else if x != y {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image is not square",
		})
	}
	if x > 256 {
		image = resize.Resize(256, 256, image, resize.Lanczos3)
	}
	newImage := models.NewImage(userID, description, location)
	err = mgm.Coll(newImage).Create(newImage)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	fileID := newImage.ID.Hex()

	thumbnailFile, err := os.Create("./data/thumbnail/" + fileID + ".jpg")
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	defer thumbnailFile.Close()

	err = jpeg.Encode(thumbnailFile, image, &jpeg.Options{Quality: 70})
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	originalFile, err := os.Create("./data/original/" + fileID + ".jpg")
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	defer originalFile.Close()
	file.Seek(0, 0)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	_, err = originalFile.Write(fileBytes)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	return c.SendString("UploadImage")
}

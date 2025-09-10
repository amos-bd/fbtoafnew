package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func main() {
	_ = godotenv.Load()
	db := InitDB()
	defer db.Close()
	redisClient = InitRedis()

	app := fiber.New()

	app.Post("/event/:action", func(c *fiber.Ctx) error {
		action := c.Params("action")
		trackingID := c.FormValue("tracking_id")
		userID := c.FormValue("user_id")
		ip := c.IP()
		eventID := c.FormValue("event_id")
		if eventID == "" {
			eventID = GenerateEventID()
		}

		ok, err := CheckEventID(eventID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Redis error")
		}
		if !ok {
			return fiber.NewError(fiber.StatusConflict, "Duplicate event")
		}

		if err := SaveEvent(db, eventID, action, trackingID, userID, ip); err != nil {
			log.Println("DB error:", err)
		}
		if err := AffilikaPostback(action, trackingID); err != nil {
			log.Println("Affilika error:", err)
		}
		if err := FbConversionAPI(action, eventID, trackingID, ip); err != nil {
			log.Println("FB CAPI error:", err)
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"event_id": eventID,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
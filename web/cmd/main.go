package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"teamC/Global"
	"teamC/db/inMemory"
	"teamC/models"
	"teamC/web"
)

func main() {
	cfg := Global.Configuration{}

	err := web.LoadConfiguration(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	cfg.WebApp = fiber.New()
	app := cfg.WebApp

	app.Use(logger.New())

	if cfg.Production {
		productionConfiguration(&cfg)
	} else {
		nonProductionConfiguration(&cfg)
		// Setup Routes
		// static page to test out back and forth websocket connection
		app.Static("/", "./static/home.html")
	}
	// Authentication middleware
	app.Use(web.GetJwtMiddleware(&cfg))

	// Start the communication hub
	go web.RunHub() // on a separate goroutine|thread

	// TESTING
	myRepo := inMemory.NewInMemoryRepository()
	myRepo.SaveUser(&models.User{
		Id:        1,
		Username:  "1st-Username",
		SubId:     "1st-SubId",
		FirstName: "1st-FirstName",
		LastName:  "1st-LastName",
	})

	cfg.UserRepo.SaveUser(&models.User{
		Id:        2,
		Username:  "2nd-Username",
		SubId:     "2nd-SubId",
		FirstName: "2nd-FirstName",
		LastName:  "2nd-LastName",
	})

	// As tests they're small enough, but should be moved at a later date when functional to reduce *noise*
	{ // API routing
		api := app.Group("/api")

		{ // Deck API
			deckApi := api.Group("/deck")

			deckApi.Post("/", func(c *fiber.Ctx) error {
				//var users *[]models.User //users := make([]models.User)
				allUsers := make([]models.User, 0)
				cfg.UserRepo.GetAllUsers(&allUsers)
				myRepo.GetAllUsers(&allUsers)
				log.Println(allUsers)
				myRepo.PrintAllUsers()
				return c.SendString("POST CALLED")
			})

			deckApi.Get("/:id", func(c *fiber.Ctx) error {
				return c.SendString("GET CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
			})

			deckApi.Put("/:id?", func(c *fiber.Ctx) error {
				return c.SendString("PUT CALLED with ID: " + c.Params("id", "POSSIBLE_ID"))
			})

			deckApi.Patch("/:id", func(c *fiber.Ctx) error {
				return c.SendString("Patch CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
			})

			deckApi.Delete("/:id", func(c *fiber.Ctx) error {
				return c.SendString("DELETE CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
			})

			api.Head("/", func(c *fiber.Ctx) error {
				return c.SendStatus(200)

			})
		}

		{ // Question API
			questionApi := api.Group("questions")

			questionApi.Post("/:deck_id/", func(c *fiber.Ctx) error {
				return c.SendString("POST CALLED with deck_id" + c.Params("deck_id", "MUST_HAVE_ID"))
			})

			questionApi.Get("/:deck_id/:question_id", func(c *fiber.Ctx) error {
				return c.SendString("GET method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
					" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
			})

			questionApi.Put("/:deck_id/:question_id?", func(c *fiber.Ctx) error {
				return c.SendString("PUT CALLED with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
					" with question ID: " + c.Params("question_id", "POSSIBLE_ID"))
			})
			questionApi.Patch("/:deck_id/:question_id", func(c *fiber.Ctx) error {
				return c.SendString("PATCH method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
					" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
			})
			questionApi.Delete("/:deck_id/:question_id", func(c *fiber.Ctx) error {
				return c.SendString("DELETE method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
					" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
			})
			questionApi.Head("/", func(c *fiber.Ctx) error {
				return c.SendStatus(200)

			})
		}
	}

	// Websocket setup
	websockets := app.Group("/ws")
	websockets.Use(web.SetupWebsocketUpgrade())
	websockets.Get("/:id", web.WebsocketRoom())

	//app.Use(web.SetupWebsocketUpgrade())
	//app.Get("/ws/:id", web.WebsocketRoom())

	// Start the web server
	log.Fatalln(app.Listen(cfg.Port))
}

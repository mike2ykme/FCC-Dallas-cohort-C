package web

import (
	"github.com/gofiber/fiber/v2"
	"teamC/Global"
)

func SetupAPIRoutes(cfg *Global.Configuration) {
	app := cfg.WebApp

	api := app.Group("/api")

	deckApi := api.Group("/deck")
	deckApi.Post("/", deckPost())
	deckApi.Get("/:id", deckGet())
	deckApi.Put("/:id?", deckPut())
	deckApi.Patch("/:id", deckPatch())
	deckApi.Delete("/:id", deckDelete())

	questionApi := api.Group("questions")

	questionApi.Post("/:deck_id/", questionPost())
	questionApi.Get("/:deck_id/:question_id", questionGet())
	questionApi.Put("/:deck_id/:question_id?", questionPut())
	questionApi.Patch("/:deck_id/:question_id", questionPatch())
	questionApi.Delete("/:deck_id/:question_id", questionDelete())
	questionApi.Head("/", questionHead())

}

/*
	Deck API Functions
*/

func deckPost() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("POST CALLED")
	}
}
func deckGet() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("GET CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
	}
}

func deckPut() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("PUT CALLED with ID: " + c.Params("id", "POSSIBLE_ID"))
	}
}

func deckPatch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Patch CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
	}
}
func deckDelete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("DELETE CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
	}
}

/*
	Question API Functions
*/
func questionPost() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("POST CALLED with deck_id" + c.Params("deck_id", "MUST_HAVE_ID"))
	}
}

func questionGet() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("GET method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
	}
}
func questionPut() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("PUT CALLED with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "POSSIBLE_ID"))
	}
}
func questionPatch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("PATCH method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
	}
}

func questionDelete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("DELETE method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
	}
}

func questionHead() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	}
}

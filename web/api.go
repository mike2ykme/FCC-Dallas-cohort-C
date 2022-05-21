package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"teamC/Global"
	"teamC/models"
)

// These should all be authenticated
func SetupAPIRoutes(cfg *Global.Configuration) {
	app := cfg.WebApp

	// for all the API calls we're going to load the userID into locals for all calls
	api := app.Group("/api", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")     // "json"
		c.Accepts("application/json") // "application/json"
		if user, ok := c.Locals("user").(*jwt.Token); ok {
			if claims, ok := user.Claims.(jwt.MapClaims); ok {
				if userId, ok := claims["id"].(float64); ok {
					c.Locals("userId", uint(userId))
					return c.Next()
				}
			}
		}
		//if c.Method() != fiber.MethodDelete {
		//	var deck models.Deck
		//	if err := c.BodyParser(&deck); err == nil {
		//		deck.OwnerId =
		//		c.Locals("postedDeck", deck)
		//	} else {
		//		c.Locals("postedDeck", nil)
		//	}
		//
		//}
		return c.SendStatus(fiber.StatusInternalServerError)
	})

	deckApi := api.Group("/deck", func(c *fiber.Ctx) error {
		return c.Next()
	})
	deckApi.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("RESULTS OK %d", c.Locals("userId").(uint)))
	})
	deckApi.Post("/", deckPost(cfg))
	deckApi.Get("/:id", deckGet(cfg)).Name("deck.get")
	deckApi.Put("/:id?", deckPut(cfg))
	deckApi.Patch("/:id", deckPatch(cfg))
	deckApi.Delete("/:id", deckDelete(cfg))

	questionApi := api.Group("questions")

	questionApi.Post("/:deck_id/", questionPost(cfg))
	questionApi.Get("/:deck_id/:question_id", questionGet(cfg)).Name("question.get")
	questionApi.Put("/:deck_id/:question_id?", questionPut(cfg))
	questionApi.Patch("/:deck_id/:question_id", questionPatch(cfg))
	questionApi.Delete("/:deck_id/:question_id", questionDelete(cfg))
	questionApi.Head("/", questionHead(cfg))

}

/*
	Deck API Functions
*/

func deckPost(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var deck models.Deck
		err := c.BodyParser(&deck)
		if id, ok := c.Locals("userId").(uint); ok && err == nil {
			deck.Id = 0
			deck.OwnerId = id

			if deckID, err := cfg.DeckRepo.SaveDeck(&deck); err == nil {
				location, _ := c.GetRouteURL("deck.get", fiber.Map{"id": deckID})
				return c.Status(fiber.StatusCreated).SendString(location)
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func deckGet(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var decks []models.Deck
		id := c.Locals("userId").(uint)
		cfg.DeckRepo.GetDecksByUserId(&decks, id)
		if decks == nil {
			return c.Status(fiber.StatusCreated).SendString("[]")
		}
		return c.Status(fiber.StatusOK).JSON(decks)
	}
}

func deckPut(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("PUT CALLED with ID: " + c.Params("id", "POSSIBLE_ID"))
	}
}

func deckPatch(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Patch CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
	}
}
func deckDelete(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("DELETE CALLED with id: " + c.Params("id", "MUST_HAVE_ID"))
	}
}

/*
	Question API Functions
*/
func questionPost(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("POST CALLED with deck_id" + c.Params("deck_id", "MUST_HAVE_ID"))
	}
}

func questionGet(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("GET method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
	}
}
func questionPut(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("PUT CALLED with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "POSSIBLE_ID"))
	}
}
func questionPatch(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("PATCH method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
	}
}

func questionDelete(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("DELETE method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
	}
}

func questionHead(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	}
}

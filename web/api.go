package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"teamC/Global"
	"teamC/models"
)

const USER_ID = "userId"

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
					c.Locals(USER_ID, uint(userId))
					return c.Next()
				}
			}
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	})

	deckApi := api.Group("/deck", func(c *fiber.Ctx) error {
		return c.Next()
	})
	deckApi.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("RESULTS OK %d", c.Locals("userId").(uint)))
	})
	deckApi.Post("/", deckPost(cfg))
	deckApi.Get("/owner", deckGetByOwner(cfg))
	deckApi.Get("/:id", deckGetById(cfg)).Name("deck.get")
	deckApi.Put("/:id?", deckPut(cfg))
	deckApi.Patch("/:id", deckPut(cfg))
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

		if id, ok := c.Locals(USER_ID).(uint); ok && err == nil {
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

func deckGetByOwner(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		decks := make([]models.Deck, 0)
		userId := c.Locals(USER_ID).(uint)
		cfg.DeckRepo.GetDecksByUserId(&decks, userId)
		return c.Status(fiber.StatusOK).JSON(decks)
	}
}

func deckGetById(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var deck models.Deck

		if deckId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64); err == nil {
			cfg.DeckRepo.GetDeckById(&deck, uint(deckId))
			id := c.Locals(USER_ID).(uint)

			if deck.OwnerId == id {
				return c.Status(fiber.StatusOK).JSON(&deck)
			}
		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func deckPut(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var parsedDeck models.Deck

		repo := cfg.DeckRepo
		log := cfg.Logger
		if err := c.BodyParser(&parsedDeck); err != nil {
			log.Println(c.BaseURL(), err)
		}

		if deckId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64); err == nil {
			var oldDeck models.Deck
			repo.GetDeckById(&oldDeck, uint(deckId))
			userId := c.Locals(USER_ID).(uint)

			if oldDeck.OwnerId == userId {
				oldDeck.ReplaceFields(&parsedDeck)
				// The above method replaces all, we don't want to change Ids
				oldDeck.Id = uint(deckId)
				repo.SaveDeck(&oldDeck)
				location, _ := c.GetRouteURL("deck.get", fiber.Map{"id": deckId})
				return c.Status(fiber.StatusCreated).SendString(location)
			}
		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func deckDelete(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		repo := cfg.DeckRepo
		//log := cfg.Logger
		userId := c.Locals(USER_ID).(uint)

		if deckId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64); err == nil {
			var deck models.Deck
			if err := repo.GetDeckById(&deck, uint(deckId)); err == nil && deck.OwnerId == userId {
				err := repo.DeleteDeckById(uint(deckId))
				if err == nil {
					return c.SendStatus(fiber.StatusOK)
				}
			}
		}

		return c.SendStatus(fiber.StatusBadRequest)
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

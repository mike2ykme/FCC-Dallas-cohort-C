package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"teamC/Global"
	"teamC/models"
)

const USER_ID = "userId"
const FirstName = "firstName"

// These are all behind a JWT authentication layer so we can get a user's details
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

	flashcardAPI := api.Group("flashcard")

	flashcardAPI.Post("/", flashcardPost(cfg))
	flashcardAPI.Get("/deck/:deck_id", flashcardGetByDeck(cfg))
	flashcardAPI.Get("/:flashcard_id", flashcardGet(cfg)).Name("flashcard.get")
	flashcardAPI.Put("/:flashcard_id?", flashcardPut(cfg))
	flashcardAPI.Delete("/:flashcard_id", flashcardDelete(cfg))

	scoreAPI := api.Group("scores")

	scoreAPI.Get("/:room_id", getRoomScores(cfg))

	roomAPI := api.Group("room")
	roomAPI.Post("/create", postNewRoom(cfg))

	// Create an API call for a ROOM
	// Store the score by room ID
	// Create a leaderboard where everyone can see the scores for everyone, and you can narrow into a certain group

	// /scores/:room? if room then show for certain room
	// swith up rooms to usign UUID? this would allow

	// Andrew will want the decks to be shuffled in backend, or at least the option to do it
	// so shuffle the decks before sending them out

	// at the end of the game they will send the user score, correct # of answers only
	/*
			{

		    "action": "SCORE",
		    //"admin": true,
		    //"question": "",
		    //"answers": null
			SCORE: 5
		}

	*/

}

/*
	Room API Functions
*/

func postNewRoom(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {

		newUUID := uuid.New()
		cfg.Logger.Printf("new UUID room created: %s", newUUID.String())
		cfg.Logger.Printf("Have a user ID of %d\n", c.Locals(USER_ID))
		adminId := c.Locals(USER_ID).(uint)
		newRoom <- models.RoomCreation{
			AdminId:   adminId,
			NewRoomID: newUUID,
			Logger:    cfg.Logger,
		}

		return c.SendString(newUUID.String())
	}
}

/*
	Score API Functions
*/
func getRoomScores(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

/*
	Deck API Functions
*/

func deckPost(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var deck models.Deck
		err := c.BodyParser(&deck)
		if err != nil {
			cfg.Logger.Printf("there was an error trying to parse the deck %s\n", err.Error())
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if id, ok := c.Locals(USER_ID).(uint); ok {
			deck.ID = 0
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
			if deckErr := cfg.DeckRepo.GetDeckById(&deck, uint(deckId)); deckErr == nil {
				id := c.Locals(USER_ID).(uint)

				shouldShuffle := c.Query("shuffle", "false")

				if strings.ToLower(shouldShuffle) == "true" {
					deck.Shuffle()
				}

				if deck.OwnerId == id {
					return c.Status(fiber.StatusOK).JSON(&deck)
				}

			} else {
				cfg.Logger.Printf("there was an error getting the deck: %s\n", deckErr.Error())
			}
		} else {
			cfg.Logger.Printf("there was an error parsing the deckID: %s\n", err.Error())
		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func deckPut(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var parsedDeck models.Deck

		repo := cfg.DeckRepo
		logger := cfg.Logger
		if err := c.BodyParser(&parsedDeck); err != nil {
			logger.Println(c.BaseURL(), err)

		} else {
			if deckId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64); err == nil {
				var oldDeck models.Deck
				var dbErr error = nil

				// if the value is 0 then we won't have something to look for in the DB
				if deckId != 0 {
					dbErr = repo.GetDeckById(&oldDeck, uint(deckId))
				}

				if dbErr != nil {
					logger.Printf("there was a problem getting the deck by ID: %s\n", dbErr.Error())
				} else {
					userId := c.Locals(USER_ID).(uint)

					if oldDeck.OwnerId == userId {
						oldDeck.ReplaceFields(&parsedDeck)
						// The above method replaces all, we don't want to change Ids
						oldDeck.ID = uint(deckId)
						if id, saveErr := repo.SaveDeck(&oldDeck); saveErr == nil {
							location, _ := c.GetRouteURL("deck.get", fiber.Map{"id": id})
							return c.Status(fiber.StatusCreated).SendString(location)

						} else {
							logger.Println("there was an error trying to save %s\n", saveErr.Error())
						}
					}
				}
			}
		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func deckDelete(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		repo := cfg.DeckRepo
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
func flashcardPost(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var flashcard models.FlashCard

		if err := c.BodyParser(&flashcard); err == nil {

			if flashcard.DeckId > 0 {

				if cardID, saveErr := cfg.FlashcardRepo.SaveFlashcard(&flashcard); saveErr == nil {
					location, _ := c.GetRouteURL("flashcard.get", fiber.Map{
						"flashcard_id": cardID,
					})

					return c.Status(fiber.StatusCreated).SendString(location)

				} else {
					cfg.Logger.Printf("There was an error saving the card: %s\n", err.Error())

				}
			} else {
				cfg.Logger.Println("the flashcard did not contain a deck ID")

			}
		} else {
			cfg.Logger.Printf("there was an error trying to parse the flashcard %s\n", err.Error())

		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func flashcardGet(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {

		flashcardID, qErr := strconv.ParseUint(c.Params("flashcard_id", "0"), 10, 64)
		if qErr == nil {
			var flashcard models.FlashCard
			if err := cfg.FlashcardRepo.GetFlashcardById(&flashcard, uint(flashcardID)); err == nil {
				return c.Status(fiber.StatusOK).JSON(flashcard)

			} else {
				cfg.Logger.Printf("there was an error getting the flashcard by ID: %s\n", err.Error())
			}

		} else {
			cfg.Logger.Printf("there was an error during conversion -> %s\n", qErr.Error())
		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}
func flashcardGetByDeck(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		deckId, err := strconv.ParseUint(c.Params("deck_id", "0"), 10, 64)
		if err == nil {
			var allCards = make([]models.FlashCard, 0)
			if err := cfg.FlashcardRepo.GetAllFlashcardByDeckId(&allCards, uint(deckId)); err == nil {
				return c.Status(fiber.StatusOK).JSON(allCards)

			}
		} else {
			cfg.Logger.Printf("there was an error during conversion -> %s\n", err.Error())

		}
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func flashcardPut(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var parsedFlashcard models.FlashCard

		repo := cfg.FlashcardRepo
		logger := cfg.Logger

		if parseErr := c.BodyParser(&parsedFlashcard); parseErr != nil {
			logger.Printf("there was an error parsing flashcard %s\n", parseErr)

		} else {
			if cardId, convertErr := strconv.ParseUint(c.Params("flashcard_id", "0"), 10, 64); convertErr != nil {
				logger.Printf("there was an error parsing flashcard_id value %s\n", convertErr.Error())
			} else {
				var oldCard models.FlashCard
				var getByIdErr error = nil

				// We're only going to look in DB if this is not 0
				if cardId != 0 {
					getByIdErr = repo.GetFlashcardById(&oldCard, uint(cardId))
				}

				if getByIdErr != nil {
					logger.Printf("there was an error getting from DB")
				} else {
					userId, ok := c.Locals(USER_ID).(uint)
					if !ok {
						logger.Println("there was an error converting and getting the user ID")
					} else {
						var oldDeck models.Deck

						if getDeckErr := cfg.DeckRepo.GetDeckById(&oldDeck, parsedFlashcard.DeckId); getDeckErr == nil {
							if oldDeck.OwnerId == userId {
								oldCard.ReplaceFields(&parsedFlashcard)
								oldCard.ID = uint(cardId)

								if id, saveErr := repo.SaveFlashcard(&oldCard); saveErr == nil {
									location, _ := c.GetRouteURL("flashcard.get", fiber.Map{
										"flashcard_id": id,
									})
									return c.Status(fiber.StatusCreated).SendString(location)
								}
							} else {
								logger.Printf("this user (%d)does not own this deck (%d)\n", userId, oldDeck.ID)
							}
						} else {
							logger.Printf("there was an error getting the deck by ID: %s\n", getDeckErr.Error())
						}

					}
				}
			}
		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}
func flashcardPatch(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("PATCH method with Deck: " + c.Params("deck_id", "MUST_HAVE_ID") +
			" with question ID: " + c.Params("question_id", "MUST_HAVE_ID"))
	}
}

func flashcardDelete(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userId := c.Locals(USER_ID).(uint)

		if flashcardID, err := strconv.ParseUint(c.Params("flashcard_id", "0"), 10, 64); err == nil {
			var card models.FlashCard
			if err = cfg.FlashcardRepo.GetFlashcardById(&card, uint(flashcardID)); err == nil {
				deckId := card.DeckId
				var deck models.Deck
				if err := cfg.DeckRepo.GetDeckById(&deck, uint(deckId)); err == nil && deck.OwnerId == userId {

					//err := cfg.DeckRepo.DeleteDeckById(uint(deckId))
					err := cfg.FlashcardRepo.DeleteFlashcardById(uint(card.ID))
					if err == nil {
						return c.SendStatus(fiber.StatusOK)
					}
				}

			} else {
				cfg.Logger.Printf("there was an error trying to get the flashcard by ID %s\n", err.Error())
			}

		} else {
			cfg.Logger.Printf("there was an error parsing uint %s\n", err.Error())
		}

		return c.SendStatus(fiber.StatusBadRequest)
	}
}

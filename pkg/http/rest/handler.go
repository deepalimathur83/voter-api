package rest

import (
	"strconv"

	"drexel.edu/voter-api/pkg/create"
	"drexel.edu/voter-api/pkg/delete"
	"drexel.edu/voter-api/pkg/read"
	"drexel.edu/voter-api/pkg/update"
	"github.com/gofiber/fiber/v2"
)

func Handler(port int, createAdapter create.Adapter, updateAdapter update.Adapter, readAdapter read.Adapter, deleteAdapter delete.Adapter) *fiber.App {

	router := fiber.New()

	router.Post("/voters/:id", func(c *fiber.Ctx) error {

		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {

			c.Status(fiber.StatusInternalServerError)
			return err
		}

		newVoter := create.Voter{
			Id: voterId,
		}

		if err := c.BodyParser(&newVoter); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		if err = createAdapter.CreateVoter(newVoter); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		c.Status(fiber.StatusOK)

		return c.SendString("New Voter got created ")
	})

	router.Post("/voters/:voterId/polls/:pollId", func(c *fiber.Ctx) error {

		voterId, err := strconv.Atoi(c.Params("voterId"))

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		pollId, err := strconv.Atoi(c.Params("pollId"))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		newHistory := create.VoterHistory{
			PollId: pollId,
		}

		if err := c.BodyParser(&newHistory); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		err = createAdapter.CreateVoterHistory(voterId, newHistory)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		c.Status(fiber.StatusOK)

		return c.SendString("New Voter History got created")
	})

	//Note : Adding PUT/ Update

	router.Put("/voters/:id", func(c *fiber.Ctx) error {

		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {

			c.Status(fiber.StatusInternalServerError)
			return err
		}

		newVoter := update.Voter{
			Id: voterId,
		}

		if err := c.BodyParser(&newVoter); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		if err = updateAdapter.UpdateVoter(newVoter); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		c.Status(fiber.StatusOK)

		return c.SendString("Voter got updated")
	})

	router.Put("/voters/:voterId/polls/:pollId", func(c *fiber.Ctx) error {

		voterId, err := strconv.Atoi(c.Params("voterId"))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		pollId, err := strconv.Atoi(c.Params("pollId"))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		newHistory := update.VoterHistory{
			PollId: pollId,
		}

		if err := c.BodyParser(&newHistory); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		err = updateAdapter.UpdateVoterHistory(voterId, newHistory)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		c.Status(fiber.StatusOK)

		return c.SendString("Voter History got Updated")
	})

	// Adding Read

	// GET voter by ID

	router.Get("/voters", func(c *fiber.Ctx) error {
		voters, err := readAdapter.ReadAllVoter()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.Status(fiber.StatusOK)
		return c.JSON(voters)
	})

	// GET voter by : ID

	router.Get("/voters/:id", func(c *fiber.Ctx) error {
		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		voter, err := readAdapter.ReadVoter(voterId)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.Status(fiber.StatusOK)

		return c.JSON(voter)
	})

	//GET voter history by : voter ID and poll ID

	router.Get("/voters/:voterId/polls/:pollId", func(c *fiber.Ctx) error {
		voterId, err := strconv.Atoi(c.Params("voterId"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		pollId, err := strconv.Atoi(c.Params("pollId"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		voterHistory, err := readAdapter.ReadVoterHistory(voterId, pollId)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		return c.JSON(voterHistory)
	})

	// GET all voters

	router.Get("/voters", func(c *fiber.Ctx) error {
		voters, err := readAdapter.ReadAllVoter()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		return c.JSON(voters)
	})

	// GET all voter history for a specific voter

	router.Get("/voters/:voterId/polls", func(c *fiber.Ctx) error {
		voterId, err := strconv.Atoi(c.Params("voterId"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		voterHistories, err := readAdapter.ReadAllVoterHistory(voterId)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		return c.JSON(voterHistories)
	})
	//Delete voter
	router.Delete("/voters/:id", func(c *fiber.Ctx) error {
		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		if err := deleteAdapter.DeleteVoter(voterId); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.Status(fiber.StatusOK)

		return c.SendString("Voter got deleted")
	})

	// Delete voterhistory by ID and PollID
	router.Delete("/voters/:voterId/polls/:pollId", func(c *fiber.Ctx) error {

		voterId, err := strconv.Atoi(c.Params("voterId"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		pollId, err := strconv.Atoi(c.Params("pollId"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		if err := deleteAdapter.DeleteVoterHistory(voterId, pollId); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		c.Status(fiber.StatusOK)

		return c.SendString("Voter History got deleted")
	})

	return router

}

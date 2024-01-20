package main

import (
	"AnimeBirthday/birthday"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type ReturnMsg struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data []birthday.AnimePerson `json:"data"`
}

func main() {

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/animebirthday", func(c fiber.Ctx) error {
		params := c.Queries()
		month, err1 := strconv.Atoi(params["month"])
		day, err2 := strconv.Atoi(params["day"])

		returnMeg := ReturnMsg{
			Code: 200,
			Msg:  "success",
			Data: nil,
		}

		if err1 != nil || err2 != nil {
			returnMeg.Code = 400
			returnMeg.Msg = "Error: month and day must be integer"

			return c.JSON(returnMeg)
		}

		persons, err := birthday.GetAnimePersonBirthday(month, day)

		if err != nil {
			returnMeg.Code = 500
			returnMeg.Msg = "Error: " + err.Error()

			return c.JSON(returnMeg)
		}

		birthday.PrintPersons(persons)

		returnMeg.Data = persons

		return c.JSON(returnMeg)

	})

	app.Listen(":22400")
}

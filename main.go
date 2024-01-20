package main

import (
	"AnimeBirthday/birthday"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/pkg/errors"
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

		persons, err := get_anime_person_birthday(month, day)

		if err != nil {
			returnMeg.Code = 500
			returnMeg.Msg = "Error: " + err.Error()

			log.Error(err)

			return c.JSON(returnMeg)
		}

		// birthday.PrintPersons(persons)

		returnMeg.Data = persons

		log.Infof("get month:%s, day:%s\n", params["month"], params["day"])

		return c.JSON(returnMeg)

	})

	app.Listen(":22400")
}

func get_anime_person_birthday(month, day int) ([]birthday.AnimePerson, error) {
	persons, err := birthday.GetAnimePersonBirthdayFromDatabase(month, day)

	if err != nil {
		return nil, err
	}

	if len(persons) > 0 {
		return persons, nil
	} else {
		persons, err := birthday.GetAnimePersonBirthdayFromWeb(month, day)

		if err != nil {
			return nil, err
		}

		err_insert := birthday.InsertAnimePersonBirthdayToDatabase(month, day, persons)

		if err_insert != nil {
			return nil, errors.Wrapf(err_insert, "insert anime person birthday to database error with month: %d, day: %d", month, day)
		}

		return persons, nil
	}
}

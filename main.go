package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"AnimeBirthday/birthday"
	"AnimeBirthday/config"
	"AnimeBirthday/tasks"
)

type ReturnMsg struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data []birthday.AnimePerson `json:"data"`
}

func main() {

	config := config.LoadConfig()

	log.Info("Server start...")
	log.Infof("config load: %v\n", config)

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/animebirthday", getAnimeBirthday)

	app.Get("/checkweekcache", checkWeekCache)

	app.Get("/checkallcache", checkAllCache)

	// 后台任务：每天定时更新缓存
	go tasks.ScheduleDailyCacheUpdate()

	app.Listen(fmt.Sprintf(":%d", config.ServerConfig.Port))
}

func getAnimeBirthday(c fiber.Ctx) error {
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

		log.Error(err)

		return c.JSON(returnMeg)
	}

	// birthday.PrintPersons(persons)

	returnMeg.Data = persons

	log.Infof("get month:%s, day:%s\n", params["month"], params["day"])

	return c.JSON(returnMeg)
}

func checkWeekCache(c fiber.Ctx) error {
	err := tasks.CheckWeekBirthdayCache()

	if err != nil {
		return c.SendString("Error: " + err.Error())
	}

	return c.SendString("Check week cache success")
}

func checkAllCache(c fiber.Ctx) error {
	tasks.CheckAllBirthdayCacheNeedUpdate()

	return c.SendString("Check all cache success")
}

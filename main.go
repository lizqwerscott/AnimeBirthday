package main

import (
	"AnimeBirthday/birthday"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
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

	app.Get("/checkweekcache", func(c fiber.Ctx) error {
		err := checkWeekBirthdayCache()

		if err != nil {
			return c.SendString("Error: " + err.Error())
		}

		return c.SendString("Check week cache success")
	})

	app.Get("/checkallcache", func(c fiber.Ctx) error {
		checkAllBirthdayCacheNeedUpdate()

		return c.SendString("Check all cache success")
	})

	// 后台任务：每天定时更新缓存
	go scheduleDailyCacheUpdate()

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
		persons, err := birthday.GetAnimePersonBirthdayFromWebSlow(month, day)

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

func getAfterDays(days int) []birthday.Birthday {
	now := time.Now()
	result := make([]birthday.Birthday, 0, days)

	for i := 0; i < days; i++ {
		afterDay := now.AddDate(0, 0, i)
		month := afterDay.Month()
		day := afterDay.Day()
		afterDayBirthday := birthday.Birthday{Month: int(month), Day: day}

		result = append(result, afterDayBirthday)
	}

	return result
}

func checkWeekBirthdayCache() error {
	afterWeekDays := getAfterDays(7)
	// 设置种子，以确保每次运行程序时都会生成不同的随机数序列
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for _, day := range afterWeekDays {
		_, err := get_anime_person_birthday(day.Month, day.Day)

		if err != nil {
			return errors.Wrapf(err, "check week birthday cache error with month: %d, day: %d", day.Month, day.Day)
		}

		// 生成在区间[5, 10]之间的随机整数
		randomNumber := rng.Intn(6) + 5

		time.Sleep(time.Duration(randomNumber) * time.Second)
	}

	return nil
}

func checkAllBirthdayCacheNeedUpdate() {

	allAnimeBirthday, err := birthday.GetAllAnimePersonBirthdayFromDatabase()

	if err != nil {
		log.Error(errors.Wrap(err, "check all birthday cache error"))
	}

	for _, animeBirthday := range allAnimeBirthday {
		animeBirthdayLastUpdate := time.Unix(animeBirthday.LastUpdate, 0)
		now := time.Now()
		if now.Sub(animeBirthdayLastUpdate).Hours() > 100*24 {
			month := animeBirthday.Birthday.Month
			day := animeBirthday.Birthday.Day

			persons, err := birthday.GetAnimePersonBirthdayFromWebSlow(month, day)

			if err != nil {
				log.Error(errors.Wrapf(err, "check all birthday cache error with month: %d, day: %d", month, day))
				continue
			}

			err_update := birthday.UpdateAnimePersonBirthdayToDatabase(month, day, persons)

			if err_update != nil {
				log.Error(errors.Wrapf(err_update, "check all birthday cache error with month: %d, day: %d", month, day))
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func scheduleDailyCacheUpdate() {
	c := cron.New()
	// 每天的凌晨 2 点执行更新缓存的任务
	spec := "0 2 * * *"
	c.AddFunc(spec, func() {
		log.Info("Updating cache...")
		// 检查今天和一周之内的生日是否已经缓存
		err := checkWeekBirthdayCache()

		if err != nil {
			log.Error(err)
		}

		// 检查信息是否过期
		checkAllBirthdayCacheNeedUpdate()
	})

	c.Start()

	select {}
}

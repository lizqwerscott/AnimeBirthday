package tasks

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"AnimeBirthday/birthday"
)

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

func CheckWeekBirthdayCache() error {
	log.Info("Checking week birthday cache...")
	afterWeekDays := getAfterDays(7)
	// 设置种子，以确保每次运行程序时都会生成不同的随机数序列
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for _, day := range afterWeekDays {
		log.Infof("Checking month: %d, day: %d\n", day.Month, day.Day)
		_, err := birthday.GetAnimePersonBirthday(day.Month, day.Day)

		if err != nil {
			err_get := errors.Wrapf(err, "check week birthday cache error with month: %d, day: %d", day.Month, day.Day)

			log.Errorf("%v", err_get)
			log.Infof("retry: %d, %d", day.Month, day.Day)

			_, err2 := birthday.GetAnimePersonBirthday(day.Month, day.Day)
			if err2 != nil {
				log.Errorf("retry error: %v", err2)
			}
		}

		// 生成在区间 [10, 20] 之间的随机整数
		randomNumber := rng.Intn(11) + 10

		log.Infof("Sleep %d seconds\n", randomNumber)
		time.Sleep(time.Duration(randomNumber) * time.Second)
	}

	return nil
}

func CheckAllBirthdayCacheNeedUpdate() {

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

func ScheduleDailyCacheUpdate() {
	c := cron.New()
	// 每天的凌晨 2 点执行更新缓存的任务
	spec := "0 2 * * *"
	c.AddFunc(spec, func() {
		log.Info("Updating cache...")
		// 检查今天和一周之内的生日是否已经缓存
		err := CheckWeekBirthdayCache()

		if err != nil {
			log.Error(err)
		}

		// 检查信息是否过期
		CheckAllBirthdayCacheNeedUpdate()
	})

	c.Start()

	select {}
}

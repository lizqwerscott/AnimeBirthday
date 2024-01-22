package birthday

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type AnimeBirthdayModel struct {
	Birthday   Birthday      `bson:"birthday"`
	Persons    []AnimePerson `bson:"persons"`
	LastUpdate int64         `bson:"last_update"`
}

func initMongoDB() (*qmgo.QmgoClient, context.Context, error) {
	ctx := context.Background()

	cred := qmgo.Credential{
		Username:   "animebirthday",
		Password:   "12138",
		AuthSource: "AnimeBirthday",
	}

	conf := qmgo.Config{Uri: "mongodb://192.168.3.14:27017", Database: "AnimeBirthday", Coll: "animebirthdays", Auth: &cred}

	cli, err := qmgo.Open(ctx, &conf)

	return cli, ctx, errors.Wrap(err, "init mongodb error")
}

func GetAnimePersonBirthdayFromDatabase(month, day int) ([]AnimePerson, error) {
	cli, ctx, err := initMongoDB()

	if err != nil {
		return nil, errors.Wrapf(err, "Get anime person birthday from database with month: %d, day: %d", month, day)
	}

	anime_persons := AnimeBirthdayModel{}

	cli.Find(ctx, bson.M{"birthday": Birthday{Month: month, Day: day}}).One(&anime_persons)

	// fmt.Printf("get from database: %v\n", anime_persons)

	return anime_persons.Persons, nil
}

func GetAllAnimePersonBirthdayFromDatabase() ([]AnimeBirthdayModel, error) {
	cli, ctx, err := initMongoDB()

	if err != nil {
		return nil, errors.Wrap(err, "Get all anime person birthday from database error")
	}

	anime_birthdays := []AnimeBirthdayModel{}

	cli.Find(ctx, bson.M{}).All(&anime_birthdays)

	// fmt.Printf("get from database: %v\n", anime_persons)

	return anime_birthdays, nil
}

func InsertAnimePersonBirthdayToDatabase(month, day int, persons []AnimePerson) error {
	cli, ctx, err := initMongoDB()

	if err != nil {
		return errors.Wrapf(err, "insert anime person birthday to database error with month: %d, day: %d", month, day)
	}

	anime_persons := AnimeBirthdayModel{Birthday: Birthday{Month: month, Day: day}, Persons: persons, LastUpdate: time.Now().Unix()}

	_, errInsert := cli.InsertOne(ctx, anime_persons)

	if errInsert != nil {
		return errors.Wrapf(errInsert, "insert anime person birthday to database error with month: %d, day: %d", month, day)
	}

	return nil
}

func UpdateAnimePersonBirthdayToDatabase(month, day int, persons []AnimePerson) error {
	cli, ctx, err := initMongoDB()

	if err != nil {
		return errors.Wrapf(err, "update anime person birthday to database error with month: %d, day: %d", month, day)
	}

	errUpdate := cli.UpdateOne(ctx, bson.M{"birthday": Birthday{Month: month, Day: day}}, bson.M{"$set": bson.M{"persons": persons, "last_update": time.Now().Unix()}})

	if errUpdate != nil {
		return errors.Wrapf(errUpdate, "update anime person birthday to database error with month: %d, day: %d", month, day)
	}

	return nil
}

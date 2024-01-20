package birthday

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type AnimeBirthdayModel struct {
	Birthday string        `bson:"birthday"`
	Persons  []AnimePerson `bson:"persons"`
}

func initMongoDB() (*qmgo.QmgoClient, context.Context) {
	ctx := context.Background()

	cred := qmgo.Credential{
		Username:   "animebirthday",
		Password:   "12138",
		AuthSource: "AnimeBirthday",
	}

	conf := qmgo.Config{Uri: "mongodb://192.168.3.14:27017", Database: "AnimeBirthday", Coll: "animebirthdays", Auth: &cred}

	cli, err := qmgo.Open(ctx, &conf)

	// defer func() {
	// 	if err = cli.Close(ctx); err != nil {
	// 		log.Fatal(err)
	// 		panic(err)
	// 	}
	// }()

	if err != nil {
		log.Fatal(err)
	}

	return cli, ctx
}

func GetAnimePersonBirthdayFromDatabase(month, day int) ([]AnimePerson, error) {
	cli, ctx := initMongoDB()

	birthday := fmt.Sprintf("%d-%d", month, day)

	anime_persons := AnimeBirthdayModel{}

	cli.Find(ctx, bson.M{"birthday": birthday}).One(&anime_persons)

	// fmt.Printf("get from database: %v\n", anime_persons)

	return anime_persons.Persons, nil
}

func InsertAnimePersonBirthdayToDatabase(month, day int, persons []AnimePerson) error {
	cli, ctx := initMongoDB()

	birthday := fmt.Sprintf("%d-%d", month, day)

	anime_persons := AnimeBirthdayModel{Birthday: birthday, Persons: persons}

	// fmt.Printf("anime_persons: %v\n", anime_persons)

	_, err := cli.InsertOne(ctx, anime_persons)

	// fmt.Printf("insert to database: %v\n", result)

	return errors.Wrapf(err, "insert anime person birthday to database error with month: %d, day: %d", month, day)
}

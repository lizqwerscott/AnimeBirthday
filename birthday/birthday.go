package birthday

import (
	"github.com/pkg/errors"
)

func GetAnimePersonBirthday(month, day int) ([]AnimePerson, error) {
	persons, err := GetAnimePersonBirthdayFromDatabase(month, day)

	if err != nil {
		return nil, err
	}

	if len(persons) > 0 {
		return persons, nil
	} else {
		persons, err := GetAnimePersonBirthdayFromWebSlow(month, day)

		if err != nil {
			return nil, err
		}

		err_insert := InsertAnimePersonBirthdayToDatabase(month, day, persons)

		if err_insert != nil {
			return nil, errors.Wrapf(err_insert, "insert anime person birthday to database error with month: %d, day: %d", month, day)
		}

		return persons, nil
	}
}

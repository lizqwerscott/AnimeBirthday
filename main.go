package main

import "AnimeBirthday/birthday"

func main() {

	persons := birthday.GetAnimePersonBirthday(1, 19)

	birthday.PrintPersons(persons)
}

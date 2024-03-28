package main

import "goddd/internal/orm"

func main() {
	db, err := orm.ConnectMySQL()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(orm.UserSchema{})
}

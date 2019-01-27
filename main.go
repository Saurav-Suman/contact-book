package main

import "os"

func main() {
	//a := new(Apps)
	var a Apps

	a.Initialize(os.Getenv("USERNSME"), os.Getenv("PASSWORD"), os.Getenv("DBSTRING"))

	a.Run(":" + os.Getenv("PORT"))
}

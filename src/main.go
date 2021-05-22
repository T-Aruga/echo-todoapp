package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()
	initDB()

	e.GET("/", index)
	e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

func index(c echo.Context) error {
	users := getAlUsers()
	data, _ := json.Marshal(&users)
	result := string(data)
	return c.String(http.StatusOK, "Hello, GO World!"+result)
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	user := findUser(userId)
	resStr := "User:  " + user.Name + "  Email:  " + user.Email
	return c.String(http.StatusOK, resStr)
}

func saveUser(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	save(name, email)
	return c.String(http.StatusOK, "Name:"+name+", Email:"+email)
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	name := c.FormValue("name")
	email := c.FormValue("email")
	update(userId, name, email)
	return c.String(http.StatusOK, "Updated!!")
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	delete(userId)
	return c.String(http.StatusOK, "Deleted!!")
}

type User struct {
	gorm.Model
	Name  string
	Email string
}

func initDB() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベースの初期化に失敗 at initDB")
	}
	db.AutoMigrate(&User{})
	defer db.Close()
}

func findUser(id int) User {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベースの初期化に失敗 at findUser")
	}
	var user User
	db.First(&user, id)
	db.Close()
	return user
}

func getAlUsers() []User {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベースの初期化に失敗 at getAllUsers")
	}
	var users []User
	db.Order("created_at desc").Find(&users)
	db.Close()
	return users
}

func save(name, email string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベースの初期化に失敗 at save")
	}
	db.Create(&User{Name: name, Email: email})
	defer db.Close()
}

func update(id int, name, email string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベースの初期化に失敗 at update")
	}
	var user User
	db.First(&user, id)
	user.Name = name
	user.Email = email
	db.Save(&user)
	db.Close()
}

func delete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベースの初期化に失敗 at delete")
	}
	var user User
	db.First(&user, id)
	db.Delete(&user)
	db.Close()
}

package main

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

var users []User
var locker sync.Mutex

const TRIAL_KISS_COUNT = 3

func main() {
	locker = sync.Mutex{}
	router := gin.Default()
	router.Use(cors.Default())
	router.StaticFile("/", "www/index.html")
	router.StaticFile("/autokiss.zip", "autokiss.zip")
	router.GET("/autokiss/js", jsHandler)
	router.POST("/autokiss/who", whoHandler)
	router.GET("/autokiss/all", allHandler)
	router.GET("/autokiss/auth/:id", authHandler)
	router.GET("/autokiss/init/:id", initHandler)
	//router.Run(":8080")
	router.RunTLS(":443", "../certs/cert.crt", "../certs/pk.key")

}

func initHandler(c *gin.Context) {
	strID := c.Param("id")
	id, _ := strconv.Atoi(strID)

	addUser(id)

	c.JSON(http.StatusOK, nil)
}

func jsHandler(c *gin.Context) {

	file, _ := ioutil.ReadFile("in.js")

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
		"src":    string(file),
	})
}

func whoHandler(context *gin.Context) {

	users := getUsers()

	type Response struct {
		Code  int           `json:"code"`
		Data  []interface{} `json:"data"`
		Delay int           `json:"delay"`
	}

	body := context.Request.Body
	data, err := ioutil.ReadAll(body)

	if err != nil {
		context.JSON(http.StatusOK, &Response{Code: 0})
		return
	}

	if len(data) < 6 {
		context.JSON(http.StatusOK, &Response{Code: 0})
		return
	}

	pType := binary.LittleEndian.Uint16(data[4:6])

	log.Printf("type: %v data: %v", pType, data)

	res := &Response{}
	switch pType {
	case 29:

		if len(data) < 14 {
			context.JSON(http.StatusOK, &Response{Code: -1})
			return
		}

		leaderID := binary.LittleEndian.Uint32(data[6:10])
		rolledID := binary.LittleEndian.Uint32(data[10:14])
		log.Printf("leaderID: %v rolledID: %v", leaderID, rolledID)

		for i, v := range users {
			if v.UserID == int(leaderID) || v.UserID == int(rolledID) {
				end := stepUpUser(&users[i])

				if end {
					context.JSON(403, nil)
					return
				}

				res.Code = 29
				res.Data = []interface{}{1}
				res.Delay = 7000
				break
			}
		}

	case 28:
		if len(data) < 10 {
			context.JSON(http.StatusOK, &Response{Code: -1})
			return
		}

		leaderID := binary.LittleEndian.Uint32(data[6:10])

		log.Printf("leaderID: %v", leaderID)

		for _, v := range users {
			if v.UserID == int(leaderID) {
				res.Code = 28
				res.Data = []interface{}{0}
				res.Delay = 5000
				break
			}
		}
	}

	context.JSON(http.StatusOK, res)
}

func authHandler(context *gin.Context) {

	strID := context.Param("id")
	id, _ := strconv.Atoi(strID)

	addUser(id)
	u := getUser(id)
	u.IsTrial = false
	updateUser(u)
	context.JSON(http.StatusOK, gin.H{"user": u})
}

func allHandler(context *gin.Context) {
	users := getUsers()
	context.JSON(http.StatusOK, gin.H{"users": users})
}

//Увеличиваем kissCount++
//Возвращаем true, если кончился триал
func stepUpUser(user *User) bool {

	if user == nil {
		return false
	}

	user.KissCount++

	if user.KissCount > TRIAL_KISS_COUNT && user.IsTrial {
		return true
	}

	updateUser(user)
	return false
}

//Обновляем user в базе
func updateUser(user *User) {

	if user == nil {
		return
	}

	u := getUser(user.UserID)

	if u == nil {
		return
	}

	u = user

	log.Printf("Update user %v\n", u)
	saveJSON()
}

//Добавляем user в базу
func addUser(userID int) {

	u := getUser(userID)

	if u != nil {
		log.Println("User yes do")
		return
	}

	user := User{
		UserID:  userID,
		IsTrial: true,
	}

	users = append(users, user)
	log.Printf("New user added %v\n", user)
	saveJSON()
}

//Сохраняем Users в `users.json`
func saveJSON() {

	locker.Lock()
	defer locker.Unlock()

	file, _ := os.OpenFile("users.json", os.O_CREATE, os.ModePerm)
	defer file.Close()

	encoder := json.NewEncoder(file)
	users := getUsers()
	encoder.Encode(&users)
}

// Загружаем из `users.json` в Users
func loadJSON() {
	locker.Lock()
	defer locker.Unlock()
	file, _ := ioutil.ReadFile("users.json")
	err := json.Unmarshal(file, &users)
	if err != nil {
		log.Println(err.Error())
	}
}

// Ищим user в users
func getUser(userID int) *User {

	for i, v := range getUsers() {
		if v.UserID == userID {
			return &users[i]
		}
	}

	return nil
}

// Получаем список всех Users
func getUsers() []User {

	if users == nil {
		users = make([]User, 0)
		loadJSON()
	}

	return users
}

// User ...
type User struct {
	UserID    int  `json:"user_id"`
	IsTrial   bool `json:"is_trial"`
	KissCount int  `json:"kiss_count"`
}

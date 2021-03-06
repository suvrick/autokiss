package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
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
var states []State

var locker sync.Mutex
var isDebug = true

const TRIAL_KISS_COUNT = 30

func main() {

	states = make([]State, 0)

	locker = sync.Mutex{}
	router := gin.Default()
	router.Use(cors.Default())
	router.StaticFile("/", "www/index.html")
	router.GET("/autokiss", downloadZipHandler)
	router.GET("/autokiss/zip", zipHandler)
	router.GET("/autokiss/js", jsHandler)
	router.POST("/autokiss/who", whoHandler)
	router.GET("/autokiss/all", allHandler)
	router.GET("/autokiss/auth/:id", authHandler)
	router.GET("/autokiss/init/:id", initHandler)

	if isDebug {
		router.Run(":8080")
	} else {
		router.RunTLS(":443", "../certs/cert.crt", "../certs/pk.key")
	}
}

func initHandler(c *gin.Context) {
	strID := c.Param("id")
	id, _ := strconv.Atoi(strID)

	addUser(id)

	c.JSON(http.StatusOK, nil)
}

func jsHandler(c *gin.Context) {

	file, _ := ioutil.ReadFile("in.js")
	file2, _ := ioutil.ReadFile("style.css")

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
		"js":     string(file),
		"css":    string(file2),
	})
}

func zipHandler(c *gin.Context) {
	loadStateJSON()
	fmt.Println(states)
	c.JSON(200, gin.H{
		"count":     len(states),
		"downloads": &states,
	})
}

func downloadZipHandler(c *gin.Context) {

	loadStateJSON()
	state := State{
		IP: c.ClientIP(),
	}

	fmt.Println(state)

	for i, v := range states {
		if state.IP == v.IP {
			states[i].Count++
			saveStateJSON()
			fmt.Println("save state.json")
			c.File("autokiss.zip")
			return
		}
	}

	state.Count++
	states = append(states, state)
	saveStateJSON()
	c.File("autokiss.zip")
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

	case 308:
		if len(data) < 14 {
			context.JSON(http.StatusOK, &Response{Code: -1})
			return
		}

		kickID := binary.LittleEndian.Uint32(data[6:10])
		whoKickID := binary.LittleEndian.Uint32(data[10:14])
		log.Printf("KICK >> kickID: %v whoKickID: %v", kickID, whoKickID)

		u := getUser(int(kickID))

		if u.KissCount > TRIAL_KISS_COUNT && u.IsTrial {
			context.JSON(403, nil)
			return
		}

		if u != nil {
			res.Code = 30
			res.Data = []interface{}{kickID}
			res.Delay = 4000
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

//?????????????????????? kissCount++
//???????????????????? true, ???????? ???????????????? ??????????
func stepUpUser(user *User) bool {

	if user == nil {
		return false
	}

	user.KissCount++
	log.Printf("user[%d], kiss: %d, isTrial: %v\n", user.UserID, user.KissCount, user.IsTrial)
	if user.KissCount > TRIAL_KISS_COUNT && user.IsTrial {
		return true
	}

	updateUser(user)
	return false
}

//?????????????????? user ?? ????????
func updateUser(user *User) {

	if user == nil {
		return
	}
	users = getUsers()
	for i, v := range users {
		if user.UserID == v.UserID {
			users[i].KissCount = user.KissCount
			users[i].IsTrial = user.IsTrial

			log.Printf("Update user %v\n", user)
			saveJSON()
			break
		}
	}
}

//?????????????????? user ?? ????????
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

//?????????????????? Users ?? `users.json`
func saveJSON() {

	locker.Lock()
	defer locker.Unlock()

	file, _ := os.OpenFile("users.json", os.O_CREATE|os.O_RDWR, 777)
	defer file.Close()

	encoder := json.NewEncoder(file)
	users := getUsers()
	if err := encoder.Encode(&users); err != nil {
		log.Println(err.Error())
	}
}

// ?????????????????? ???? `users.json` ?? Users
func loadJSON() {
	locker.Lock()
	defer locker.Unlock()
	file, _ := ioutil.ReadFile("users.json")
	err := json.Unmarshal(file, &users)
	if err != nil {
		log.Println(err.Error())
	}
}

// ???????? user ?? users
func getUser(userID int) *User {

	for i, v := range getUsers() {
		if v.UserID == userID {
			return &users[i]
		}
	}

	return nil
}

// ???????????????? ???????????? ???????? Users
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

//State ...
type State struct {
	IP    string `json:"ip"`
	Count int    `json:"count"`
}

//?????????????????? Users ?? `state.json`
func saveStateJSON() {

	locker.Lock()
	defer locker.Unlock()

	file, _ := os.OpenFile("state.json", os.O_CREATE|os.O_RDWR, 777)
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&states); err != nil {
		log.Println(err.Error())
	}
}

// ?????????????????? ???? `state.json` ?? Stats
func loadStateJSON() {
	locker.Lock()
	defer locker.Unlock()
	file, _ := ioutil.ReadFile("state.json")
	states = make([]State, 0)
	err := json.Unmarshal(file, &states)
	if err != nil {
		log.Println(err.Error())
	}
}

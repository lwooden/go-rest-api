package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// album represents data about a record album.
// struct tags such as json:"artist" specify what a field’s name should be when the struct’s contents are serialized into JSON
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type verse struct {
	ID      string `json:"id"`
	Book    string `book:"book"`
	Chapter string `json:"chapter"`
	Verse   string `json:"verse"`
	Text    string `json:"text"`
}

type remoteverses struct {
	Book_ID   string `json:"book_id"`
	Book_Name string `json:"book_name"`
	Chapter   int    `json:"chapter"`
	Verse     int    `json:"verse"`
	Text      string `json:"text"`
}

type remoteverse struct {
	Reference        string         `json:"reference"`
	Verses           []remoteverses `json:"verses"`
	Text             string         `json:"text"`
	Translation_ID   string         `json:"translation_id"`
	Translation_Name string         `json:"translation_name"`
	Translation_Note string         `json:"translation_note"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var verses = []verse{
	{ID: uuid.NewString(), Book: "John", Chapter: "3", Verse: "16", Text: "For God so loved the world, that He gave His only Begotten Son..."},
	{ID: uuid.NewString(), Book: "Acts", Chapter: "24", Verse: "16", Text: "For this being so, I myself, always strive to have a clear consious before God and men!"},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getVerses(c *gin.Context) {
	println("Got a request!")
	c.IndentedJSON(http.StatusOK, verses)

}

// Fetching from "https://bible-api.com/John 3:16"
func fetchRemoteVerse(c *gin.Context) {

	println("Preparing HTTP Client")
	url := "https://bible-api.com/John 3:16"

	// Set up the client
	client := http.Client{}

	// Prep the request
	println("Preparing the request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		println(err)
	}

	// Invoke the request
	println("Sending the request")
	res, getErr := client.Do(req)

	if getErr != nil {
		println(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	// Response comes back in the form of an address to a location in memory -- Response => 0xc000506000
	println("Response =>", res)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	// Same for the body -- Body => (0x14502f0,0xc000474d00)
	println("Body =>", res.Body)

	payload := remoteverse{}

	jsonErr := json.Unmarshal(body, &payload)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(payload.Reference)

	c.IndentedJSON(http.StatusOK, "Done!")

}

func contains(vId string) ([]verse, bool) {

	var tempSlice []verse

	for _, item := range verses {

		if item.ID == vId {
			tempSlice = append(tempSlice, item)
			return tempSlice, true
		}
	}

	return tempSlice, false
}

// func getVerseById(c *gin.Context) {
// 	print("Got a request!")
// 	id := c.Param("id")

// 	var result, doesExist = contains(id)

// 	if !doesExist {
// 		c.String(http.StatusNotFound, "Resource %v does not exists!", id)
// 	} else {
// 		c.IndentedJSON(http.StatusOK, result[0])
// 	}
// }

func getVerseById(c *gin.Context) {
	print("Got a request!")
	id := c.Param("id")

	for _, v := range verses {
		if v.ID == id {
			c.IndentedJSON(http.StatusOK, v)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "verse not found"})
}

func main() {
	router := gin.Default()          // initialize a Gin router
	router.GET("/albums", getAlbums) // register method, route, and handler
	router.GET("/api/verse", getVerses)
	router.GET("/api/verse/:id", getVerseById)
	router.GET("/api/remoteverse", fetchRemoteVerse)
	router.Run(":8080") // start the webserver

}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/Icorp/petProject/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var posts map[string]*models.Post
var templates *template.Template

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("../web/templates/*")

	// Read static files from web/assets
	readStaticFiles(router)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.Run(":8080")
}
func readStaticFiles(router *gin.Engine) {
	router.Static("/css/", "../web/assets/css")
	router.Static("/js/", "../web/assets/js")
	router.Static("/images/", "../web/assets/images")
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := getData()
	templates.ExecuteTemplate(w, "index", data)
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../templates/index.html", "../templates/header.html", "../templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "index", posts)
}

func getData() []models.Post {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("blog").Collection("posts")
	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var posts []models.Post
	if err = cur.All(ctx, &posts); err != nil {
		panic(err)
	}
	return posts
}

func sendData(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("../templates/addpost.html", "../templates/header.html", "../templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		t.ExecuteTemplate(w, "index", nil)
	case "POST":
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		collection := client.Database("blog").Collection("posts")
		r.ParseForm()
		post := models.Post{Title: r.FormValue("title"), Body: r.FormValue("body")}
		insertResult, err := collection.InsertOne(ctx, post)
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult.InsertedID)
	}
}

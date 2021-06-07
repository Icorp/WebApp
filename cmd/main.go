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

var templates *template.Template
var posts map[string]*models.Post

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("../web/templates/*")
	// Read static files from web/assets
	readStaticFiles(router)
	router.GET("/", indexHandler)
	router.Run(":8080")
}
func readStaticFiles(router *gin.Engine) {
	router.Static("/css/", "../web/assets/css")
	router.Static("/js/", "../web/assets/js")
	router.Static("/images/", "../web/assets/images")
}
func indexHandler(c *gin.Context) {
	data := getData()
	c.HTML(http.StatusOK, "index.html", data)
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

func sendExampleData() {
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
	insertResult, err := collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: 0},
		{Key: "title", Value: "There’s a Cool New Way for Men to Wear Socks and Sandals"},
		{Key: "body", Value: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Praesentium nam quas inventore, ut iure iste modi eos adipisci ad ea itaque labore earum autem nobis et numquam, minima eius. Nam eius, non unde ut aut sunt eveniet rerum repellendus porro."},
		{Key: "authorName", Value: "Colorlib"},
		{Key: "imagePost", Value: "images/img_5.jpg"},
		{Key: "imageAuthor", Value: "images/person_1.jpg"},
		{Key: "postDate", Value: "March 15, 2018 "},
		{Key: "numOfComments", Value: 3},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertResult)
}

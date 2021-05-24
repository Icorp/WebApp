package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/Icorp/petProject/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var posts map[string]*models.Post

func main() {
	posts = make(map[string]*models.Post, 0)
	fmt.Println("Start Application ...")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/book", bookHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../templates/index.html", "../templates/header.html", "../templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	data := getData()
	t.ExecuteTemplate(w, "index", data)
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

package models

type Post struct {
	Id            int    `bson:"id"`
	Title         string `bson:"title"`
	Body          string `bson:"body"`
	ImagePost     string `bson:"imagePost"`
	ImageAuthor   string `bson:"imageAuthor"`
	AuthorName    string `bson:"authorName"`
	PostDate      string `bson:"postDate"`
	NumOfComments int    `bson:"numOfComments"`
}

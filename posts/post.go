package posts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

// PostDir holds the path of the directory containing the post files
var PostDir string

// Post holds info about posts appearing on the site
type Post struct {
	ID         string    `json:"ID"`
	Title      string    `json:"Title"`
	Author     string    `json:"Author"`
	DatePosted string    `json:"DatePosted"`
	Content    string    `json:"Content"`
	Category   string    `json:"Category"`
	Image      string    `json:"Image"`
	Comments   []Comment `json:"Comments"`
}

// Comment holds info about comments submited to a post
type Comment struct {
	ID       string    `json:"ID"`
	UserID   string    `json:"UserID"`
	Text     string    `json:"Text"`
	IsReply  bool      `json:"IsReply"`
	Comments []Comment `json:"Comments"`
}

// PostPath returns the file path of a post
func PostPath(postID string) string {
	if PostDir == "" {
		PostDir = "./posts"
	}
	return fmt.Sprintf("%s/%s.json", PostDir, postID)
}

// CreatePost creates a new post struct
func CreatePost(title, author, text, category, image string) Post {
	return Post{
		uuid.NewV4().String(),
		title,
		author,
		time.Now().Local().Format("2006-01-02"),
		text,
		category,
		image,
		make([]Comment, 0),
	}
}

// ReadPost reads a post file with the specified id
func ReadPost(id string) (Post, error) {
	filepath := PostPath(id)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Post{}, err
	}
	var post Post
	err = json.Unmarshal(data, &post)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

// ReadPosts reads all the posts on the post directory
func ReadPosts() ([]Post, error) {
	posts := make([]Post, 0)
	if PostDir == "" {
		PostDir = "./posts"
	}
	files, err := ioutil.ReadDir(PostDir)
	if err != nil {
		log.Fatal(err)
		return posts, err
	}
	for _, f := range files {
		id := strings.Split(f.Name(), ".")[0]
		post, err := ReadPost(id)
		if err != nil {
			log.Fatal(err)
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, err
}

// SavePost saves a new post
func (post *Post) SavePost() bool {
	data, err := json.Marshal(*post)
	if err != nil {
		return false
	}
	filepath := PostPath(post.ID)
	if _, err := os.Stat(PostDir); os.IsNotExist(err) {
		os.Mkdir(PostDir, os.ModeDir)
	}
	if err = ioutil.WriteFile(filepath, data, 0644); err != nil {
		return false
	}
	return true
}

// DeletePost deletes a saved post based on id
func (post *Post) DeletePost() bool {
	if err := os.Remove(PostPath(post.ID)); err != nil {
		return false
	}
	return true
}

// EditPost edits the title or text of an existing post
func EditPost(id, title, content, category, image string) bool {
	post, err := ReadPost(id)
	if err != nil {
		return false
	}
	post.Title = title
	post.Content = content
	saved := post.SavePost()
	return saved
}

// CreateComment creates a new comment
func CreateComment(userID string, text string, isReply bool) Comment {
	return Comment{
		uuid.NewV4().String(),
		userID,
		text,
		isReply,
		make([]Comment, 0),
	}
}

// AddComment to a centain post
func (post *Post) AddComment(comment Comment) bool {
	post.Comments = append(post.Comments, comment)
	saved := post.SavePost()
	return saved
}

// removes a comment from a comment array
func remove(slice []Comment, s int) []Comment {
	return append(slice[:s], slice[s+1:]...)
}

// DeleteComment deletes a saved comment
func DeleteComment(postID, commentID string) bool {
	post, err := ReadPost(postID)
	if err != nil {
		return false
	}
	for i := 0; i < len(post.Comments); i++ {
		if post.Comments[i].ID == commentID {
			post.Comments = remove(post.Comments, i)
		}
	}
	return post.SavePost()
}

// EditComment edits a saved comment
func EditComment(postID, commentID, text string) bool {
	post, err := ReadPost(postID)
	if err != nil {
		return false
	}
	for i := 0; i < len(post.Comments); i++ {
		if post.Comments[i].ID == commentID {
			post.Comments[i].Text = text
		}
	}
	return post.SavePost()
}

// ReplyComment adds a reply (Comment) to a comment
func ReplyComment(postID, commentID string, comment Comment) bool {
	if comment.IsReply != true {
		return false
	}
	post, err := ReadPost(postID)
	if err != nil {
		return false
	}
	for i := 0; i < len(post.Comments); i++ {
		if post.Comments[i].ID == commentID {
			post.Comments[i].Comments = append(post.Comments[i].Comments, comment)
		}
	}
	saved := post.SavePost()
	return saved
}

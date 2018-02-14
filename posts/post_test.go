package posts

import "testing"

var PostID string
var CommentID string
var Title string
var Author string
var Content string

func TestCreatePost(t *testing.T) {
	title := "Test"
	content := "This is a test post"
	author := "Nikos Raftogiannis"
	post := CreatePost(title, author, content, "Web", "http://www.techtechnik.com/wp-content/uploads/2014/11/hacking.jpg")
	if post.ID == "" || post.Author == "" || post.Content == "" {
		t.Errorf("CreatePost(%s, %s, %s) == %v", title, author, content, post)
	}
}

func TestSavePost(t *testing.T) {
	post := CreatePost("Test", "Nikos Raftogiannis", "This is a test post", "Web", "http://www.techtechnik.com/wp-content/uploads/2014/11/hacking.jpg")
	PostID = post.ID
	Title = post.Title
	Author = post.Author
	Content = post.Content
	saved := post.SavePost()
	if saved != true {
		t.Errorf("SavePost() failed to save post = %v", post)
	}
}

func TestPostPath(t *testing.T) {
	postID := PostID
	postPath := PostPath(postID)
	postPathWant := "./posts/" + postID + ".json"
	if postPath != postPathWant {
		t.Errorf("PostPath(%s) == %s, want %s", postID, postPath, postPathWant)
	}
}

func TestReadPost(t *testing.T) {
	post, err := ReadPost(PostID)
	if err != nil {
		t.Errorf("ReadPost() failed to read post = %v, with error %v", post, err)
	}
}

func TestReadPosts(t *testing.T) {
	posts, err := ReadPosts()
	if err != nil {
		t.Errorf("ReadPosts() failed to read posts = %v, with error %v", posts, err)
	}
}

func TestEditPost(t *testing.T) {
	post, _ := ReadPost(PostID)
	saved := EditPost(post.ID, "TestEdited", post.Content, "Systems", "http://www.techtechnik.com/wp-content/uploads/2014/11/hacking.jpg")
	post, _ = ReadPost(PostID)
	if saved != true || post.Title == Title {
		t.Errorf("EditPost(%s,%s,%s) failed", post.ID, post.Title, post.Content)
	}
}

func TestAddComment(t *testing.T) {
	post, _ := ReadPost(PostID)
	userID := "NikosID12345"
	text := "I say go"
	isReply := false
	comment := CreateComment(userID, text, isReply)
	CommentID = comment.ID
	post.AddComment(comment)
	post, _ = ReadPost(PostID)
	if len(post.Comments) <= 0 {
		t.Errorf("AddComment(%s,%s,%v) failed to add comment to post %v", userID, text, isReply, post)
	}
}

func TestEditComment(t *testing.T) {
	text := "I say go edited"
	saved := EditComment(PostID, CommentID, text)
	if saved != true {
		t.Errorf("EditComment(%s, %s, %s) failed", PostID, CommentID, text)
	}
}

func TestReplyComment(t *testing.T) {
	saved := ReplyComment(PostID, CommentID, CreateComment("userId1", "The is a reply comment", true))
	if saved != true {
		t.Errorf("ReplyComment() failed")
	}
}

func TestDeleteComment(t *testing.T) {
	deleted := DeleteComment(PostID, CommentID)
	if deleted != true {
		t.Errorf("DeleteComment(%s,%s) failed", PostID, CommentID)
	}
}

func TestDeletePost(t *testing.T) {
	post, _ := ReadPost(PostID)
	post.DeletePost()
	if _, err := ReadPost(PostID); err == nil {
		t.Errorf("DeletePost() failed to delete post = %v", post)
	}
}

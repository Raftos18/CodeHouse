package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Raftos18/authentication"
	"github.com/Raftos18/posts"
)

func main() {
	var role authentication.Role
	(&role).SetPermission("Admin")
	(&role).SetLevel(5)
	fmt.Printf("Role: %s\r\n", role.GetPermission())
	fmt.Printf("Level: %d\r\n", role.GetLevel())

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/posts", servePost)
	http.HandleFunc("/", servePage)

	log.Println("Listening...")
	http.ListenAndServe("192.168.1.98:80", nil)
}

// servePost responds with a post specified in the
func servePost(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", "post.html")
	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	post, err := posts.ReadPost(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
	if err = tmpl.Execute(w, post); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

}

func servePage(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path)) + ".html"

	if r.URL.RequestURI() == "/home" {
		tmpl, err := template.ParseFiles(lp, fp)
		if err != nil {
			// Log the detailed error
			log.Println(err.Error())
			// Return a generic "Internal Server Error" message
			http.Error(w, http.StatusText(500), 500)
			return
		}
		posts, err := posts.ReadPosts()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		if err = tmpl.Execute(w, posts); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	} else {
		fmt.Printf("Url Query= %v\r\n", r.URL.Query()["id"])

		// Return a 404 if the template doesn't exist
		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
		}

		// Return a 404 if the request is for a directory
		if info.IsDir() {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(lp, fp)
		if err != nil {
			// Log the detailed error
			log.Println(err.Error())
			// Return a generic "Internal Server Error" message
			http.Error(w, http.StatusText(500), 500)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}

}

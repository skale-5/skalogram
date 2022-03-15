package templates

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/skale-5/skalogram/web"
)

//go:embed posts.html
var postsFS embed.FS

type RenderPostsArgs struct {
	Posts          []web.Post
	PostsAsciiHTML []template.HTML
}

func RenderPosts(w http.ResponseWriter, args RenderPostsArgs) error {
	tpl, err := template.ParseFS(postsFS, "posts.html")
	if err != nil {
		log.Fatalf("failed to load posts.html template: %s", err)
	}
	return tpl.Execute(w, args)
}

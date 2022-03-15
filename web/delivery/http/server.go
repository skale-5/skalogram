package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/skale-5/skalogram/web"
	"github.com/skale-5/skalogram/web/delivery/http/templates"
)

type Server struct {
	listenAddr          string
	postDatabaseService *web.PostDatabaseService
	postCacheService    *web.PostCacheService
}

type NewServerArgs struct {
	ListenAddr          string
	PostDatabaseService *web.PostDatabaseService
	PostCacheService    *web.PostCacheService
}

func NewServer(args NewServerArgs) *Server {
	return &Server{
		listenAddr:          args.ListenAddr,
		postDatabaseService: args.PostDatabaseService,
		postCacheService:    args.PostCacheService,
	}
}

func httpError(w http.ResponseWriter, code int, message string, err error) {
	w.WriteHeader(code)
	log.Printf("[ERROR][%d] %s: %s", code, message, err)
	fmt.Fprint(w, message)
}

func (s *Server) voidHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) postsUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids) < 1 {
		httpError(w, http.StatusBadRequest, "id params is missing", fmt.Errorf("id param is missing"))
		return
	}
	id := ids[0]

	uid, err := uuid.Parse(id)
	if err != nil {
		httpError(w, http.StatusBadRequest, "malformed id params", fmt.Errorf("malformed id params"))
		return

	}
	err = s.postDatabaseService.UpvotePost(r.Context(), uid)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "server error", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) postsDownvoteHandler(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids) < 1 {
		httpError(w, http.StatusBadRequest, "id params is missing", fmt.Errorf("id param is missing"))
		return
	}
	id := ids[0]

	uid, err := uuid.Parse(id)
	if err != nil {
		httpError(w, http.StatusBadRequest, "malformed id params", fmt.Errorf("malformed id params"))
		return

	}
	err = s.postDatabaseService.DownvotePost(r.Context(), uid)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "server error", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.postDatabaseService.ListPosts(r.Context())
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to list posts", err)
		return
	}

	postsAscii := make([]string, len(posts))
	for i, post := range posts {
		cachedAscii, err := s.postCacheService.GetPost(r.Context(), post.ID)
		if err != nil && err != web.ErrPostCacheNotFound {
			log.Println("[WARNING] failed to retreive ascii in cache")
		}
		if len(cachedAscii) == 0 {
			postsAscii[i], err = post.GenerateAscii()
			if err != nil {
				httpError(w, http.StatusInternalServerError, "failed to generate post ascii", err)
				return
			}
			_, err = s.postCacheService.CachePost(r.Context(), post.ID, postsAscii[i], time.Second*60)
			if err != nil {
				log.Println("[WARNING] failed to cache ascii")
			}
			continue
		}
		postsAscii[i] = cachedAscii
	}
	postsAsciiHTML := make([]template.HTML, len(postsAscii))
	for i, postAscii := range postsAscii {
		postsAsciiHTML[i] = template.HTML(postAscii)
	}
	err = templates.RenderPosts(w, templates.RenderPostsArgs{
		Posts:          posts,
		PostsAsciiHTML: postsAsciiHTML,
	})
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to render posts", err)
		return
	}
}

func (s *Server) Run() {
	http.HandleFunc("/", s.postsHandler)
	http.HandleFunc("/upvote", s.postsUpvoteHandler)
	http.HandleFunc("/downvote", s.postsDownvoteHandler)

	http.HandleFunc("/favicon.ico", s.voidHandler)

	log.Printf("HTTP Server running on %s...\n", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

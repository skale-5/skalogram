package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/skale-5/skalogram/web"
	"github.com/skale-5/skalogram/web/config"
	"github.com/skale-5/skalogram/web/delivery/http/templates"
)

type Server struct {
	listenAddr          string
	postDatabaseService *web.PostDatabaseService
	postCacheService    *web.PostCacheService
	postStorageService  *web.PostStorageService
}

type NewServerArgs struct {
	ListenAddr          string
	PostDatabaseService *web.PostDatabaseService
	PostCacheService    *web.PostCacheService
	PostStorageService  *web.PostStorageService
}

func NewServer(args NewServerArgs) *Server {
	return &Server{
		listenAddr:          args.ListenAddr,
		postDatabaseService: args.PostDatabaseService,
		postCacheService:    args.PostCacheService,
		postStorageService:  args.PostStorageService,
	}
}

func httpError(w http.ResponseWriter, code int, message string, err error) {
	w.WriteHeader(code)
	log.Printf("[ERROR][%d] %s: %s", code, message, err)
	fmt.Fprint(w, message)
}

func (s *Server) voidHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) postsUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		httpError(w, http.StatusBadRequest, "failed to parse multipart form", err)
		return
	}

	f, h, err := r.FormFile("postImg")
	if err != nil {
		httpError(w, http.StatusBadRequest, "failed to retreive file from request form", err)
		return
	}
	defer f.Close()

	id := uuid.New()

	allowedContentType := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
	}

	if !allowedContentType[h.Header.Get("Content-Type")] {
		err = fmt.Errorf("unauthorized file format: %s", h.Header.Get("Content-Type"))
		httpError(w, http.StatusBadRequest, "file format not allowed", err)
		return
	}

	fullObjectPath := fmt.Sprintf("%s://%s/%s",
		config.Env().Get("STORAGE_TYPE"),
		config.Env().Get("STORAGE_BUCKET"),
		id.String(),
	)
	object, err := web.NewObjectPath(fullObjectPath)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to create object path", err)
		return
	}
	err = s.postStorageService.Write(r.Context(), object, f)
	if err != nil {
		httpError(w, http.StatusBadRequest, "failed to upload object", err)
		return
	}

	err = s.postDatabaseService.CreatePost(r.Context(), web.CreatePostParams{
		ID:     id,
		ImgUrl: object.URL(),
	})
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to create post", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

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
			obj, err := web.NewObjectPath(post.ImgUrl)
			if err != nil {
				httpError(w, http.StatusInternalServerError, "invalid object path", err)
				return
			}
			fileReader, err := s.postStorageService.Get(r.Context(), obj)
			if err != nil {
				httpError(w, http.StatusInternalServerError, "failed to get object", err)
				return
			}
			postsAscii[i], err = web.GenerateAscii(fileReader)
			if err != nil {
				httpError(w, http.StatusInternalServerError, "failed to generate post ascii", err)
				return
			}

			ttlConfig := config.Env().Get("CACHE_TTL")
			ttl, err := time.ParseDuration(ttlConfig)
			if err != nil {
				httpError(w, http.StatusInternalServerError, "invalid CACHE_TTL duration format", err)
				return
			}
			_, err = s.postCacheService.CachePost(r.Context(), post.ID, postsAscii[i], ttl)
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

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "200 OK")
	if err != nil {
		httpError(w, http.StatusInternalServerError, "failed to render posts", err)
		return
	}
	return
}

func (s *Server) Run() {
	http.HandleFunc("/", s.postsHandler)
	http.HandleFunc("/upvote", s.postsUpvoteHandler)
	http.HandleFunc("/downvote", s.postsDownvoteHandler)
	http.HandleFunc("/upload", s.postsUploadHandler)
	http.HandleFunc("/healthz", s.healthzHandler)

	http.HandleFunc("/favicon.ico", s.voidHandler)

	log.Printf("HTTP Server running on %s...\n", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}

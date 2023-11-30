package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework3/internal/pkg/repository"
	"homework3/internal/pkg/repository/postgresql"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	ArticleRepo *postgresql.ArticleRepo
	CommentRepo *postgresql.CommentRepo
}

type articleRequest struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Rating int64  `json:"rating"`
}

type commentRequest struct {
	ArticleID int64  `json:"article_id"`
	Text      string `json:"text"`
}

func mapArticleRequest(article articleRequest) *repository.Article {
	return &repository.Article{
		ID:     article.ID,
		Name:   article.Name,
		Rating: article.Rating,
	}
}

func mapCommentRequest(comment commentRequest) *repository.Comment {
	return &repository.Comment{
		ArticleID: comment.ArticleID,
		Text:      comment.Text,
	}
}

func CreateRouter(server Server) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/article", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			server.Create(w, req)
		case http.MethodPut:
			server.Update(w, req)
		case http.MethodGet:
			server.Get(w, req)
		case http.MethodDelete:
			server.Delete(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc("/comment", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			server.CreateComment(w, req)
		default:
			fmt.Println("error")
		}
	})
	return router
}

func (server *Server) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var articleReq articleRequest
	if err = json.Unmarshal(body, &articleReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	article := mapArticleRequest(articleReq)
	article, err = server.ArticleRepo.Add(req.Context(), article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(articleJson)
}

func (server *Server) Get(w http.ResponseWriter, req *http.Request) {
	articleID, err := strconv.ParseInt(req.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	article, err := server.ArticleRepo.GetByID(req.Context(), articleID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Article json: "))
	w.Write(articleJson)
	w.Write([]byte("\n"))

	comments, err := server.CommentRepo.GetCommentsForArticle(req.Context(), articleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Comments: \n"))
	if len(comments) == 0 {
		w.Write([]byte("no comments\n"))
	} else {
		for _, comment := range comments {
			commentJson, err := json.Marshal(&comment)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(commentJson)
			w.Write([]byte("\n"))
		}
	}

}

func (server *Server) Delete(w http.ResponseWriter, req *http.Request) {
	articleID, err := strconv.ParseInt(req.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := server.ArticleRepo.DeleteByID(req.Context(), articleID); err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func (server *Server) Update(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updateData articleRequest
	if err := json.Unmarshal(body, &updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	articleRepo := mapArticleRequest(updateData)

	if err := server.ArticleRepo.Update(req.Context(), articleRepo); err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func (server *Server) CreateComment(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var comment commentRequest
	if err = json.Unmarshal(body, &comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	commentRepo := mapCommentRequest(comment)
	commentRepo, err = server.CommentRepo.AddComment(req.Context(), commentRepo)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	commentJson, err := json.Marshal(commentRepo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(commentJson)
}

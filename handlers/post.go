package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"luis/go/rest-ws/models"
	"luis/go/rest-ws/repository"
	"luis/go/rest-ws/server"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	PostContent string `json:"postContent"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"postContent"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = InsertPostRequest{}
			err := json.NewDecoder(r.Body).Decode(&postRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post := models.Post{
				Id:          id.String(),
				PostContent: postRequest.PostContent,
				UserId:      claims.UserId,
			}
			err = repository.InsertPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id:          post.Id,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func GetPostByIDHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostByID(r.Context(), params["postId"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)

	}
}

// func DeletePostByIdHandler(s server.Server) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		params := mux.Vars(r)
// 		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
// 		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(s.Config().JWTSecret), nil
// 		})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return
// 		}
// 		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
// 			err = repository.DeletePost(r.Context(), params["postId"], claims.UserId)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(PostDeletedResponse{
// 				Message: "Post deleted",
// 			})
// 		} else {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

// func UpdatePostByIdHandler(s server.Server) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		params := mux.Vars(r)
// 		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
// 		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(s.Config().JWTSecret), nil
// 		})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return
// 		}
// 		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
// 			var postRequest = UpsertPostRequest{}
// 			err := json.NewDecoder(r.Body).Decode(&postRequest)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusBadRequest)
// 				return
// 			}
// 			post := models.Post{
// 				PostContent: postRequest.PostContent,
// 				Id:          params["postId"],
// 			}
// 			err = repository.UpdatePost(r.Context(), &post, claims.UserId)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(PostDeletedResponse{
// 				Message: "Post Update",
// 			})
// 		} else {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	b "github.com/devfullcycle/service-weaver-demo/components/book"
	"github.com/devfullcycle/service-weaver-demo/components/comment"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()
	root := weaver.Init(ctx)
	opts := weaver.ListenerOptions{LocalAddress: "localhost:12345"}
	lis, err := root.Listener("books", opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on %v\n", lis)

	booksComponent, err := weaver.Get[b.BookComponent](root)
	if err != nil {
		panic(err)
	}

	commentsComponent, err := weaver.Get[comment.CommentComponent](root)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Post("/books", func(w http.ResponseWriter, r *http.Request) {
		var book b.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := booksComponent.Add(ctx, book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		result, err := json.Marshal(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
	})

	r.Get("/books", func(w http.ResponseWriter, r *http.Request) {
		books, err := booksComponent.List(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonRes, err := json.Marshal(books)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonRes)
	})

	r.Post("/books/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "id")
		var comment comment.Comment
		comment.BookID = bookID
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := commentsComponent.Post(ctx, comment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		result, err := json.Marshal(comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
	})
	otelHandler := otelhttp.NewHandler(r, "books")
	http.Serve(lis, otelHandler)
}

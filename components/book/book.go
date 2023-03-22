package book

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/devfullcycle/service-weaver-demo/components/comment"
)

type BookComponent interface {
	Add(ctx context.Context, book Book) error
	List(ctx context.Context) ([]Book, error)
}

type Book struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Author   string            `json:"author"`
	Comments []comment.Comment `json:"comments,omitempty"`
	weaver.AutoMarshal
}

type bookComponent struct {
	weaver.Implements[BookComponent]
	Books []Book  `json:"books"`
	db    *sql.DB `json:"-"`
}

func (b *bookComponent) Init(ctx context.Context) error {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	b.db = db
	return nil
}

func (b *bookComponent) Add(_ context.Context, book Book) error {
	time.Sleep(time.Second)
	_, err := b.db.Exec("INSERT INTO books (id, title, author) VALUES (?, ?, ?)", book.ID, book.Title, book.Author)
	if err != nil {
		return err
	}
	b.Books = append(b.Books, book)
	return nil
}

func (b *bookComponent) List(_ context.Context) ([]Book, error) {
	time.Sleep(time.Second)
	var books []Book
	res, err := b.db.Query("SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var book Book
		err := res.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

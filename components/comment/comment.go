package comment

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
)

type CommentComponent interface {
	Post(ctx context.Context, comment Comment) error
	GetByBook(ctx context.Context, author string) ([]Comment, error)
}

type Comment struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Comment string `json:"comment"`
	BookID  string `json:"book_id"`
	weaver.AutoMarshal
}

type commentComponent struct {
	weaver.Implements[CommentComponent]
	Comments []Comment `json:"comments"`
	db       *sql.DB   `json:"-"`
}

func (b *commentComponent) Init(ctx context.Context) error {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	b.db = db
	return nil
}

func (c *commentComponent) Post(_ context.Context, comment Comment) error {
	time.Sleep(time.Second)
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO comments (id, author, comment, book_id) VALUES (?, ?, ?, ?)", id, comment.Author, comment.Comment, comment.BookID)
	if err != nil {
		return err
	}
	return nil
}

func (c *commentComponent) GetByBook(_ context.Context, bookID string) ([]Comment, error) {
	time.Sleep(time.Second)
	var comments []Comment
	res, err := c.db.Query("SELECT id, author, comment, book_id FROM comments WHERE book_id = ?", bookID)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var comment Comment
		err := res.Scan(&comment.ID, &comment.Author, &comment.Comment, &comment.BookID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

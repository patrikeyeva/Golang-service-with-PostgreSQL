package postgresql

import (
	"context"

	"homework3/internal/pkg/db"
	"homework3/internal/pkg/repository"
)

type CommentRepo struct {
	db *db.Database
}

func NewComments(database *db.Database) *CommentRepo {
	return &CommentRepo{db: database}
}

func (r *CommentRepo) AddComment(ctx context.Context, comment *repository.Comment) (*repository.Comment, error) {
	var exists bool
	err := r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM articles WHERE id = $1)", comment.ArticleID)
	if err != nil {
		return comment, err
	}
	if !exists {
		return comment, repository.ErrObjectNotFound
	}

	error := r.db.ExecQueryRow(ctx, "INSERT INTO comments(article_id, text) VALUES($1, $2) RETURNING id,created_at;", comment.ArticleID, comment.Text).Scan(&comment.ID, &comment.CreatedAt)
	return comment, error
}

func (r *CommentRepo) GetCommentsForArticle(ctx context.Context, articleID int64) ([]repository.Comment, error) {
	var comments []repository.Comment
	err := r.db.Select(ctx, &comments, "SELECT id, text, created_at FROM comments WHERE article_id = $1", articleID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

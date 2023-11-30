package postgresql

import (
	"context"

	"homework3/internal/pkg/db"
	"homework3/internal/pkg/repository"

	"github.com/georgysavva/scany/pgxscan"
)

type ArticleRepo struct {
	db *db.Database
}

func NewArticles(database *db.Database) *ArticleRepo {
	return &ArticleRepo{db: database}
}

func (r *ArticleRepo) Add(ctx context.Context, article *repository.Article) (*repository.Article, error) {
	err := r.db.ExecQueryRow(ctx, `INSERT INTO articles(name,rating) VALUES($1,$2) RETURNING id, created_at;`, article.Name, article.Rating).Scan(&article.ID, &article.CreatedAt)
	return article, err
}

func (r *ArticleRepo) GetByID(ctx context.Context, id int64) (*repository.Article, error) {
	articleData := &repository.Article{}
	err := r.db.Get(ctx, articleData, "SELECT id,name,rating,created_at FROM articles WHERE id=$1", id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return articleData, nil
}

func (r *ArticleRepo) DeleteByID(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, "DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return repository.ErrObjectNotFound
	}
	return nil
}

func (r *ArticleRepo) Update(ctx context.Context, article *repository.Article) error {
	res, err := r.db.Exec(ctx, "UPDATE articles SET name=$1, rating=$2 WHERE id=$3", article.Name, article.Rating, article.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return repository.ErrObjectNotFound
	}
	return nil
}

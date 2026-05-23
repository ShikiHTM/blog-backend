package repository

import (
	"database/sql"
	"fmt"

	"github.com/shikihtm/blog-backend/internal/model"
)

type BlogRepository interface {
	Get(uuid string) (*model.PostStats, error)
	GetAll() ([]model.PostStats, error)
	IncreaseView(slug string) (*model.PostStats, error)
	IncreaseLike(slug string) (*model.PostStats, error)
	SyncPost(slug string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(dbConn *sql.DB) BlogRepository {
	return &repository{
		db: dbConn,
	}
}

func (r *repository) Get(slug string) (*model.PostStats, error) {
	var stats model.PostStats

	query := "SELECT uuid, slug, views, likes, created_at, updated_at FROM post_stats WHERE slug = ?"

	err := r.db.QueryRow(query, slug).Scan(&stats.UUID, &stats.Slug, &stats.Views, &stats.Likes, &stats.CreatedAt, &stats.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no post found using uuid: %s", slug)
		}

		return nil, err
	}

	mdx, err := readFull(slug)
	if err != nil {
		return nil, err
	}
	stats.Title = mdx.Title
	stats.Subtitle = mdx.Subtitle
	stats.Topic = mdx.Topic
	stats.Cover = mdx.Cover
	stats.Content = mdx.Content

	return &stats, nil
}

func (r *repository) GetAll() ([]model.PostStats, error) {
	var stats []model.PostStats

	query := "SELECT uuid, slug, views, likes, created_at, updated_at FROM post_stats"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.PostStats

		err = rows.Scan(&item.UUID, &item.Slug, &item.Views, &item.Likes, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}

		meta, err := readMeta(item.Slug)
		if err != nil {
			return nil, err
		}
		item.Title = meta.Title
		item.Subtitle = meta.Subtitle
		item.Topic = meta.Topic
		item.Cover = meta.Cover

		stats = append(stats, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, err
}

func (r *repository) IncreaseView(slug string) (*model.PostStats, error) {
	var post model.PostStats

	query := `
		UPDATE post_stats
		SET views = views + 1
		WHERE slug = ?
		RETURNING uuid, slug, cover, views, likes, created_at, updated_at
	`

	err := r.db.QueryRow(query, slug).Scan(&post.UUID, &post.Slug, &post.Cover, &post.Views, &post.Likes, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot update likes to DB: %w", err)
	}

	return &post, nil
}

func (r *repository) IncreaseLike(slug string) (*model.PostStats, error) {
	var post model.PostStats

	query := `
		UPDATE post_stats
		SET likes = likes + 1
		WHERE slug = ?
		RETURNING uuid, slug, cover, views, likes, created_at, updated_at
	`

	err := r.db.QueryRow(query, slug).Scan(&post.UUID, &post.Slug, &post.Cover, &post.Views, &post.Likes, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot update likes to DB: %w", err)
	}

	return &post, nil
}

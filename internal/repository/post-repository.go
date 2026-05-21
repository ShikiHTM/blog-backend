package repository

func (r *repository) SyncPost(slug string) error {
	query := `
		INSERT INTO post_stats (uuid, slug, views, likes, created_at, updated_at)
		VALUES (lower(hex(randomblob(16))), ?, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(slug) DO NOTHING;
	`

	_, err := r.db.Exec(query, slug)
	return err
}

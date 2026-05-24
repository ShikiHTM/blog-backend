package repository

func (r *repository) SyncPost(slug string) error {
	meta, err := readMeta(slug)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO post_stats (uuid, slug, cover, author, views, likes, created_at, updated_at)
		VALUES (lower(hex(randomblob(16))), ?, ?, ?, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(slug) DO UPDATE SET cover = excluded.cover, updated_at = CURRENT_TIMESTAMP;
	`

	_, err = r.db.Exec(query, slug, meta.Cover, meta.Author)
	return err
}

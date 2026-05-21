package repository

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/frontmatter"
	"github.com/shikihtm/blog-backend/internal/model"
)

const postsDir = "./posts"

func mdxPath(slug string) string {
	return filepath.Join(postsDir, slug+".mdx")
}

func readMeta(slug string) (model.PostStats, error) {
	var meta model.PostStats

	f, err := os.Open(mdxPath(slug))
	if err != nil {
		return meta, fmt.Errorf("open mdx %s: %w", slug, err)
	}
	defer f.Close()

	if _, err := frontmatter.Parse(f, &meta); err != nil {
		return meta, fmt.Errorf("parse frontmatter %s: %w", slug, err)
	}
	return meta, nil
}

func readFull(slug string) (model.PostStats, error) {
	var meta model.PostStats

	raw, err := os.ReadFile(mdxPath(slug))
	if err != nil {
		return meta, fmt.Errorf("read mdx %s: %w", slug, err)
	}

	body, err := frontmatter.Parse(bytes.NewReader(raw), &meta)
	if err != nil {
		return meta, fmt.Errorf("parse frontmatter %s: %w", slug, err)
	}
	meta.Content = string(body)
	return meta, nil
}

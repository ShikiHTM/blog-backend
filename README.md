# blog-backend

The backend service for Shiki's personal blog. It serves blog posts authored as MDX files and tracks per-post engagement (views, likes) in a local SQLite database.

## Overview

- **Content source:** MDX files in `./posts`, each with YAML frontmatter (`title`, `topic`, `cover`). `subtitle` is still accepted but deprecated and will be removed in a future release.
- **Persistence:** SQLite at `./data/blog.db`. Only engagement metadata (UUID, slug, views, likes, timestamps) lives in the DB — post bodies are read from disk on request.
- **Hot reload:** an fsnotify watcher monitors `./posts` and automatically syncs new MDX files into the database.

## Stack

- Go 1.26.3
- [Gin](https://github.com/gin-gonic/gin) — HTTP router
- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite driver
- [fsnotify](https://github.com/fsnotify/fsnotify) — filesystem watcher
- [adrg/frontmatter](https://github.com/adrg/frontmatter) — MDX frontmatter parser

## Project layout

```
.
├── cmd/api/            # Application entrypoint (main.go)
├── internal/
│   ├── database/       # SQLite connection and schema bootstrap
│   ├── handler/        # Gin HTTP handlers and route registration
│   ├── logger/         # Colored log helpers
│   ├── model/          # Domain types (PostStats)
│   ├── repository/     # DB access, MDX reading, fs watcher
│   └── service/        # (reserved)
├── middleware/         # (reserved)
├── posts/              # MDX content
└── data/               # SQLite database file (gitignored)
```

## Getting started

### Prerequisites

- Go 1.26 or newer
- A C toolchain (required by `go-sqlite3` via cgo)

### Run

```sh
go mod download
go run ./cmd/api
```

The server listens on `:3050`. The `./data` directory and SQLite file are created automatically on first run.

### Build

```sh
go build -o blog-backend ./cmd/api
./blog-backend
```

## Authoring posts

Drop an `.mdx` file into `./posts`. The filename (without extension) becomes the post `slug`. Frontmatter fields are required for rendering metadata:

```mdx
---
title: Hello World
topic: General
cover: /hello-world/cover.png
---

Your post body in MDX goes here.
```

> `subtitle` is deprecated. Existing posts using it will continue to work, but new posts should omit it.

The watcher inserts a new row in `post_stats` on file create/write. Existing rows are preserved (slug is unique, conflicts are ignored), so view/like counters survive edits.

## HTTP API

Base path: `/api/v1`

| Method | Path                      | Description                                  |
| ------ | ------------------------- | -------------------------------------------- |
| GET    | `/posts`                  | List all posts with metadata and stats (no content). |
| GET    | `/posts/:slug`            | Fetch one post including the MDX body.       |
| POST   | `/posts/:slug/view`       | Increment the view counter. Returns 204.     |
| POST   | `/posts/:slug/like`       | Increment the like counter. Returns 204.     |

### Response shape

```json
{
  "id": "…uuid…",
  "slug": "hello-world",
  "title": "Hello World",
  "topic": "General",
  "cover": "/hello-world/cover.png",
  "subtitle": "",
  "views": 0,
  "likes": 0,
  "content": "…mdx body (single-post endpoint only)…",
  "created_at": "2026-05-22T00:00:00Z",
  "updated_at": "2026-05-22T00:00:00Z"
}
```

## License

[MIT](./LICENSE) © Phạm Nguyễn Khánh Đăng

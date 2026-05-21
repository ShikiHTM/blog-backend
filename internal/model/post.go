package model

import "time"

type PostStats struct {
	UUID      string    `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title,omitempty" yaml:"title"`
	Subtitle  string    `json:"subtitle,omitempty" yaml:"subtitle"`
	Topic     string    `json:"topic,omitempty" yaml:"topic"`
	Views     int       `json:"views"`
	Likes     int       `json:"likes"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

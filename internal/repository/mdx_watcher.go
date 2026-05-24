package repository

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/shikihtm/blog-backend/internal/logger"
)

func SyncAll(repo BlogRepository) {
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		log.Printf(logger.Error("Failed to read posts directory: %v"), err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".mdx" && ext != ".md" {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), ext)
		if err := repo.SyncPost(slug); err != nil {
			log.Printf("[SYNC] [ERROR] Failed to sync post %s to DB: %v\n", slug, err)
		} else {
			log.Printf("[SYNC] [INFO] Successfully synced database for: %s\n", slug)
		}
	}
}

func Watch(repo BlogRepository) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf(logger.Error("Failed to initialize file watcher: %v"), err)
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				ext := filepath.Ext(event.Name)

				if ext != ".mdx" && ext != ".md" {
					continue
				}

				slug := strings.TrimSuffix(filepath.Base(event.Name), ext)

				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
					err := repo.SyncPost(slug)
					if err != nil {
						log.Printf("[WATCHER] [ERROR] Failed to sync post %s to DB: %v\n", event.Name, err)
					} else {
						log.Printf("[WATCHER] [INFO] Successfully synced database for: %s\n", event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("[WATCHER] [ERROR] Runtime watcher error: %v\n", err)
			}
		}
	}()

	err = watcher.Add("./posts")
	if err != nil {
		log.Printf(logger.Error("Runtime watcher error: %v"), err)
	}

	log.Println(logger.System("Watcher service successfully started. Monitoring: ./posts"))
}

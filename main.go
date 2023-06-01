package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/aleitner/wallpaper-generator/config"
	"github.com/aleitner/wallpaper-generator/operations"
)

var logger *log.Logger

func init() {
	logDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "WallpaperGenerator", "Logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Cannot create log directory: %v", err)
	}

	logFileName := filepath.Join(logDir, "logfile.log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Cannot create or open log file: %v", err)
	}

	logger = log.New(logFile, "", log.LstdFlags)
}

func main() {

	config, err := config.ReadConfig()
	if err != nil {
		logger.Printf("Error reading config file: %s", err)
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	width, height, err := operations.GetDesktopResolution()
	if err != nil {
		logger.Printf("Error getting desktop resolution: %s", err)
		panic(err)
	}

	logger.Printf("Detected Desktop Resolution: %dx%d\n", width, height)

	userProfile, _ := os.UserHomeDir()
	downloadsPath := filepath.Join(userProfile, "Downloads")
	wallpapersPath := filepath.Join(downloadsPath, "wallpapers")

	if _, err := os.Stat(wallpapersPath); os.IsNotExist(err) {
		os.Mkdir(wallpapersPath, 0755)
	}

	for {
		for _, keyword := range config.Keywords {
			logger.Printf("Searching for keyword Image from %s\n", keyword)

			imageItems, err := operations.GoogleImageSearch(config.APIKey, config.CX, keyword)
			if err != nil {
				logger.Printf("Error searching image on Google: %s", err)
			}

			for _, imageItem := range imageItems {
				if int(imageItem.Image.Width) < width ||
					int(imageItem.Image.Height) < height ||
					int(imageItem.Image.Width/imageItem.Image.Height) != width/height ||
					imageItem.Image.Height > imageItem.Image.Width {
					continue
				}
				imageName := filepath.Base(imageItem.Link)
				filePath := filepath.Join(wallpapersPath, imageName)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					err = operations.DownloadImage(imageItem.Link, filePath)
					if err != nil {
						logger.Printf("Error (%s) downloading the image from %s", err, imageItem.Link)
					}

					logger.Printf("Downloaded Image from %s\n", imageItem.Link)
				}
			}
		}

		files, err := ioutil.ReadDir(wallpapersPath)
		if err != nil {
			logger.Printf("Could not read directory: %s", err)
			panic(fmt.Sprintf("Could not read directory: %s", err))
		}

		if len(files) > 0 {
			randomIndex := rand.Intn(len(files))
			randFile := files[randomIndex]

			wallpaperPath := filepath.Join(wallpapersPath, randFile.Name())

			err = operations.SetWallpaper(wallpaperPath)
			if err != nil {
				logger.Printf("Error setting wallpaper: %s", err)
				continue
			}

			logger.Printf("Set %s as wallpaper\n", randFile.Name())
		}

		time.Sleep(24 * time.Hour)
	}
}

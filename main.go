package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//channel_stats := getChannelStats()

	startdate := "2025-05-26"
	//startdate := "2024-07-20"
	enddate := "2025-05-30"
	today := time.Now().Format("2006-01-02")

	video_list := gatherVideoStats(startdate, enddate, today)
	getherThumbnailImages(video_list)

	m, _ := json.MarshalIndent(video_list, "", "    ")
	fmt.Println(string(m))

	// 通常レポート(yearly=false)
	render_report(video_list, today, false)

	// 年間レポート(yearly=true)
	//render_report(video_list, today, true)
}

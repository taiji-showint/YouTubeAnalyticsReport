package main

import (
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

	channel_stats := getChannelStats()
	fmt.Printf("Channel Stats: Title=%s, Subscribers=%s\n", channel_stats.Channel_title, channel_stats.Subscribers)

	startdate := "2026-02-23"
	enddate := "2026-04-28"
	//enddate := "2026-03-01"
	today := time.Now().Format("2006-01-02")

	// Get new subscribers during the period (from startdate to today)
	newSubs := getNewSubscribers(startdate, today)
	fmt.Printf("New Subscribers: %d\n", newSubs)

	// Get video list and basic stats from YouTube Data API
	video_list := gatherVideoStats(startdate, enddate, today)
	getherThumbnailImages(video_list)

	// Load impressions and CTR data from YouTube Studio CSV export
	// Expected CSV path: analytics_csv/{YYYYMMDD}/表データ.csv
	csvPath := "analytics_csv/" + today + "/表データ.csv"
	analyticsData, err := loadAnalyticsFromCSV(csvPath)
	if err != nil {
		log.Fatalf("Error loading analytics CSV: %v", err)
	}

	// Merge analytics data into video list
	video_list = mergeAnalyticsData(video_list, analyticsData)

	//m, _ := json.MarshalIndent(video_list, "", "    ")
	//fmt.Println(string(m))

	// 通常レポート(yearly=false)
	render_report(video_list, today, false, &channel_stats, newSubs)

	// 年間レポート(yearly=true)
	//render_report(video_list, today, true, &channel_stats, newSubs)
}

package main

import (
	"log"
	"fmt"
	//"time"
	"encoding/json"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//channel_stats := getChannelStats()
	
	//startdate := "2024-09-10T00:00:00Z"
	//enddate := "2024-09-20T23:00:00Z"

	//today := time.Now().Format("2006-01-02")
	//video_list := getVideoStats(startdate, enddate)
	
	startdate := "2024-06-03"
	//enddate := "2024-08-05"
	enddate := "2024-06-15"
	
	getVideoList(startdate, enddate)
	video_list := gatherVideoStats(startdate, enddate)
	
	m, _ := json.MarshalIndent(video_list,"","    ")
	fmt.Println(string(m))

	// 通常レポート(yearly=false)
	//render_report(video_list, today, false)

	// 年間レポート(yearly=true)
	//render_report(video_list, today, true)
}

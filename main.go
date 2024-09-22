package main

import (
	"log"
	"fmt"
//	"time"

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

	//today := time.Now().Format("20060102")
	//video_list := getVideoStats(startdate, enddate)
	
	startdate := "2024-09-03"
	enddate := "2024-09-20"
	//video_analytics := getVideoTrafficSource("ePJIVZZhzRg")
	//getVideoTrafficSource("ePJIVZZhzRg")
	//fmt.println(video_analytics)

	//getVideoBasicStats(startdate, enddate)
	//fmt.Println(getVideoList(startdate, enddate))
	//getVideoList(startdate, enddate)
	video_list := gatherVideoStats(startdate, enddate)
	fmt.Println(video_list)

	// 通常レポート(yearly=false)
	//render_report(video_list, today, false)

	// 年間レポート(yearly=true)
	//render_report(video_list, today, true)
}

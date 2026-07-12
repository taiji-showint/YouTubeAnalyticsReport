package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"strings"
	"sort"
	"math"
)

// formatViewCount formats view counts: 10000+ as "1.2万", <10000 as "123"
func formatViewCount(views float64) string {
	if views >= 10000 {
		return fmt.Sprintf("%.1f万", views/10000)
	}
	return fmt.Sprintf("%.0f", views)
}

func aggregateChannelStats(video_list []Video, channelStats *ChannelStats, newSubscribers int64) []Video {
	if len(video_list) == 0 {
		return video_list
	}

	// Calculate totals
	var totalViews, totalImpressions float64
	var totalWatchHours float64
	var totalAgeCount [7]float64
	var totalGenderCount [2]float64

	for _, video := range video_list {
		totalViews += video.View_counts
		totalImpressions += video.Impressions
		// Calculate total watch hours (Duration is in days, convert to hours)
		totalWatchHours += float64(video.Duration) * video.View_counts / 24.0

		// Aggregate age
		totalAgeCount[0] += video.Age_percentage.AGE13_17 * video.View_counts
		totalAgeCount[1] += video.Age_percentage.AGE18_24 * video.View_counts
		totalAgeCount[2] += video.Age_percentage.AGE25_34 * video.View_counts
		totalAgeCount[3] += video.Age_percentage.AGE35_44 * video.View_counts
		totalAgeCount[4] += video.Age_percentage.AGE45_54 * video.View_counts
		totalAgeCount[5] += video.Age_percentage.AGE55_64 * video.View_counts
		totalAgeCount[6] += video.Age_percentage.AGE65_ * video.View_counts

		// Aggregate gender
		totalGenderCount[0] += video.Gender_percentage.MALE * video.View_counts
		totalGenderCount[1] += video.Gender_percentage.FEMEL * video.View_counts
	}

	// Calculate percentages (Age_percentage and Gender_percentage are already in 0-100 scale)
	channelAge := Age_percentage{}
	channelGender := Gender_percentage{}

	if totalViews > 0 {
		channelAge.AGE13_17 = totalAgeCount[0] / totalViews
		channelAge.AGE18_24 = totalAgeCount[1] / totalViews
		channelAge.AGE25_34 = totalAgeCount[2] / totalViews
		channelAge.AGE35_44 = totalAgeCount[3] / totalViews
		channelAge.AGE45_54 = totalAgeCount[4] / totalViews
		channelAge.AGE55_64 = totalAgeCount[5] / totalViews
		channelAge.AGE65_ = totalAgeCount[6] / totalViews

		channelGender.MALE = totalGenderCount[0] / totalViews
		channelGender.FEMEL = totalGenderCount[1] / totalViews
	}

	// Calculate channel CTR
	var channelCTR float64
	if totalImpressions > 0 {
		channelCTR = (totalViews / totalImpressions) * 100
	}

	// Sort by views and get top 10
	sortedVideos := make([]Video, len(video_list))
	copy(sortedVideos, video_list)
	sort.Slice(sortedVideos, func(i, j int) bool {
		return sortedVideos[i].View_counts > sortedVideos[j].View_counts
	})

	var topTenVideos []Video
	if len(sortedVideos) > 10 {
		topTenVideos = sortedVideos[:10]
	} else {
		topTenVideos = sortedVideos
	}

	// Attach aggregated data to all videos
	subscribersStr := "0"
	if channelStats != nil && channelStats.Subscribers != "" {
		subscribersStr = channelStats.Subscribers
	}
	fmt.Fprintf(os.Stderr, "DEBUG: channelStats ptr=%p, channelStats=%+v, subscribersStr=%s\n", channelStats, channelStats, subscribersStr)

	for i := range video_list {
		video_list[i].TotalViews = formatViewCount(totalViews)
		video_list[i].TotalWatchHours = fmt.Sprintf("%.0f", totalWatchHours)
		video_list[i].TotalImpressions = fmt.Sprintf("%.1f万", math.Max(1, totalImpressions/10000))
		video_list[i].ChannelCTR = channelCTR
		video_list[i].ChannelGender = channelGender
		video_list[i].ChannelAge = channelAge
		video_list[i].TopTenVideos = topTenVideos
		video_list[i].TotalSubscribers = subscribersStr
		video_list[i].NewSubscribers = fmt.Sprintf("%d", newSubscribers)
	}

	return video_list
}

func render_report(video_list []Video, date string, yearly bool, channelStats *ChannelStats, newSubscribers int64) {
	fmt.Fprintf(os.Stderr, "render_report called: channelStats=%+v, newSubscribers=%d\n", channelStats, newSubscribers)

	// Aggregate channel-wide statistics
	video_list = aggregateChannelStats(video_list, channelStats, newSubscribers)

	var template_filename string
	if yearly {
		template_filename = "templates/tmpl_report_yearly.md"
	} else {
		template_filename = "templates/tmpl_report.md"
	}

	tmpl, err := template.New(filepath.Base(template_filename)).Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).ParseFiles(template_filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Before template: video_list[0].TotalSubscribers=%s, NewSubscribers=%s\n", video_list[0].TotalSubscribers, video_list[0].NewSubscribers)
	if err := tmpl.Execute(os.Stdout, video_list); err != nil {
		panic(err)
	}

	var report_filename string
	date_clean := strings.ReplaceAll(date, "-", "")

	if yearly {
		report_filename = "reports/showint_report_yearly_" + date_clean + ".md"
	} else {
		report_filename = "reports/showint_report_" + date_clean + ".md"
	}

	f, err := os.Create(report_filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, video_list)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"strings"
	"math"
	
	//"encoding/json"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/api/youtubeanalytics/v2"
	
)

type ChannelStats struct {
	Channel_title string
	Subscribers   string
	Channel_id    string
}

// These values are defined on sightTrafficSourceType API reference
// insightTrafficSourceType
// https://developers.google.com/youtube/analytics/dimensions#Traffic_source_Dimensions
type Traffic_source struct {
	SUBSCRIBER		float64
	RELATED_VIDEO	float64
	EXT_URL_ratio 	float64
	EXT_URL_count 	float64
	NO_LINK_OTHER	float64
	YT_CHANNEL		float64
	YT_OTHER_PAGE	float64
	YT_SEARCH		float64
	PLAYLIST		float64
	NOTIFICATION	float64
}

type Age_percentage struct {
	AGE13_17	float64
	AGE18_24	float64
	AGE25_34	float64
	AGE35_44 	float64
	AGE45_54 	float64
	AGE55_64	float64
	AGE65_		float64
}

type Gender_percentage struct {
	MALE			float64
	FEMEL			float64
	USER_SPECIFIED	float64
}

type Video struct {
	Video_title    string
	Video_id       string
	Updated_date   string
	View_counts    float64
	Like_counts    float64
	Dislike_counts float64
	AnnotationImpressions float64
	AnnotationClickThroughRate float64
	Thumbnail_url  string
	Traffic_source Traffic_source
	External_sites []map[string]float64
	Age_percentage Age_percentage
	Gender_percentage Gender_percentage
	Today string
	Duration int64
}

const (
	CHANNEL_ID_QUERY =  "channel==MINE"
	STARTDATE_SHOWINT = "2019-01-01"
)

func reverseVideoList(v []Video) []Video {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
	return v
}

func downloadImage(url string, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Println("Image downloaded successfully : " + filename)
	return nil

}

func getChannelStats() ChannelStats {
	API_KEY := os.Getenv("API_KEY")
	CHANNEL_ID := os.Getenv("CHANNEL_ID")

	var channel_stats ChannelStats

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(API_KEY))

	if err != nil {
		log.Fatalf("Error creating new YouTube service: %v", err)
	}

	call_channel := service.Channels.List([]string{"snippet", "statistics"}).Id(CHANNEL_ID)
	response_channel, err := call_channel.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}
	for _, ch := range response_channel.Items {
		channel_stats.Channel_title = ch.Snippet.Title
		channel_stats.Subscribers = string(ch.Statistics.SubscriberCount)
		//channel_stats.Channel_id = ch.Id
		//fmt.Println(strings.Repeat("=", 50))
		//fmt.Println("Channel Title: ", channel.Snippet.Title)
		//fmt.Println("Channel ID: ", channel.Id)
		//fmt.Printf("Suscribers: %v\n", channel.Statistics.SubscriberCount)
		//fmt.Println(strings.Repeat("=", 50))
	}
	return channel_stats
}

// YouTube Data API
// Official Reference: https://developers.google.com/youtube/v3/getting-started?hl=ja
// Official Sample: https://developers.google.com/youtube/v3/code_samples/go
// Library Source Code: https://github.com/googleapis/google-api-go-client/blob/main/youtube/v3/youtube-gen.go
func callYTDataAPI(
	startdate string,
	enddate string,
	maxresult int64,
) *youtube.SearchListResponse {
	API_KEY := os.Getenv("API_KEY")
	CHANNEL_ID := os.Getenv("CHANNEL_ID")
	
	// New() は非推奨になってた。Deprecated: please use NewService instead.とのこと
	// https://github.com/googleapis/google-api-go-client/blob/d0e0dc31cd30ec9b5e71541ad905236401b56d96/youtube/v3/youtube-gen.go#L159-L163
	//service, err := youtube.New(client)

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		log.Fatalf("Error creating new YouTube service: %v", err)
	}

	startdate += "T00:00:00Z"
	enddate += "T23:59:59Z"

	call := service.Search.List([]string{"id", "snippet"}).
		ChannelId(CHANNEL_ID).
		MaxResults(maxresult).
		Order("date").
		PublishedAfter(startdate).
		PublishedBefore(enddate)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making YouTube Data API call : %v", err)
	}

	return response
}

// YouTube Analytics API
// Official Reference: https://developers.google.com/youtube/reporting
// Official Sample: https://github.com/youtube/api-samples/tree/07263305b59a7c3275bc7e925f9ce6cabf774022/python
// Sample API call : https://developers.google.com/youtube/analytics/sample-requests#channel-geographic-reports
// Good blog: https://blog.codecamp.jp/youtube-analytics-python
// Google標準APIライブラリでは、 YouTube Data API のみ対応していて、AnalyticsAPIは呼び出し対応していなそうだった。
// 別のライブラリ google.golang.org/api/youtubeanalytics/v2 であれば、YouTube AnalyticsAPIを呼び出せそうだ。
// https://pkg.go.dev/google.golang.org/api/youtubeanalytics/v2#pkg-functions
func callYTAnalyticsAPI(
	dimentions string,
	metrics string,
	filters string, 
	startdate string,
	enddate string,
	sort string,
	maxresult int64,
) *youtubeanalytics.QueryResponse {
	client := getClient(youtubeanalytics.YtAnalyticsReadonlyScope)
	service, err := youtubeanalytics.New(client)

	call := service.Reports.Query().
		Ids(CHANNEL_ID_QUERY).
		Dimensions(dimentions).
		Metrics(metrics).
		StartDate(startdate).
		EndDate(enddate).
		MaxResults(maxresult).
		Sort(sort).
		Filters(filters)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making YouTube Analytics API call: %v", err)
	}
	return response
}

func getVideoList(startdate string, enddate string) []Video {
	response := callYTDataAPI(
		startdate, // startdate string
		enddate, // enddate string
		15, // maxresult int64
	)
	
	var video Video
	var video_list []Video

	for _, item := range response.Items {
		video.Video_title = item.Snippet.Title
		video.Video_id = item.Id.VideoId
		video.Updated_date = strings.Split(item.Snippet.PublishedAt, "T")[0]
		video.Thumbnail_url = item.Snippet.Thumbnails.High.Url
		video_list = append(video_list, video)
	}
	return reverseVideoList(video_list)
}

func calcDuration(video *Video, today string){
	video.Today = today

	var format = "2006-01-02"
	time_start, _ := time.Parse(format, video.Updated_date)
	time_end, _   := time.Parse(format, today)
	time_diff_hours := time_end.Sub(time_start).Hours()
	video.Duration = int64(time_diff_hours) / 24
}

func updateVideoCount(video *Video) {
	enddate_today := time.Now().Format("2006-01-02")
	filter_query := fmt.Sprintf("video==%s", video.Video_id)
	
	response := callYTAnalyticsAPI(
		"video", // dimentions
		"views,likes,dislikes", // metrics
		filter_query, // filters
		STARTDATE_SHOWINT, // startdate
		enddate_today, // enddate
		"", // sort
		5, // maxresult
	)
	video.View_counts = response.Rows[0][1].(float64)
	video.Like_counts = response.Rows[0][2].(float64)
	video.Dislike_counts = response.Rows[0][3].(float64)
}

// TODO : うまくAnnoutation関連の値が取得できなかった。すべて結果が0のままだった。
func updateAnnoutationImplession(video *Video) {
	enddate_today := time.Now().Format("2006-01-02")
	filter_query := fmt.Sprintf("video==%s", video.Video_id)
	response := callYTAnalyticsAPI(
		"video", // dimentions
		"annotationImpressions,annotationClickThroughRate,", // metrics
		//filter_query, // filters
		filter_query,
		STARTDATE_SHOWINT, // startdate
		enddate_today, // enddateå
		"", // sort
		5, // maxresult
	)
	//m, _ := json.MarshalIndent(response,"","    ")
	//fmt.Println(string(m))

	for _, row := range response.Rows {
		video.AnnotationImpressions = row[1].(float64)
		video.AnnotationClickThroughRate = row[2].(float64)
	}
}

func truncFloat(num float64) float64 {
	return math.Floor(num * 10) / 10
}

func updateVideoTrafficSourceType(video *Video) {
	enddate_today := time.Now().Format("2006-01-02")
	filter_query := fmt.Sprintf("video==%s", video.Video_id)

	response := callYTAnalyticsAPI(
		"insightTrafficSourceType", // dimentions
		"views", // metrics
		filter_query, // filters
		STARTDATE_SHOWINT, // startdate
		enddate_today, // enddate
		"-views", // sort
		15, // maxresult
	)

	for _, row := range response.Rows {
		ratio := truncFloat( row[1].(float64) / video.View_counts * 100 )
		switch row[0] {
		case "SUBSCRIBER":
			video.Traffic_source.SUBSCRIBER 	= ratio
		case "RELATED_VIDEO":
			video.Traffic_source.RELATED_VIDEO 	= ratio
		case "EXT_URL":
			video.Traffic_source.EXT_URL_ratio	= ratio
			video.Traffic_source.EXT_URL_count  = row[1].(float64)
		case "NO_LINK_OTHER":
			video.Traffic_source.NO_LINK_OTHER 	= ratio
		case "YT_CHANNEL":
			video.Traffic_source.YT_CHANNEL 	= ratio
		case "YT_OTHER_PAGE":
			video.Traffic_source.YT_OTHER_PAGE 	= ratio
		case "YT_SEARCH":
			video.Traffic_source.YT_SEARCH 		= ratio
		case "PLAYLIST":
			video.Traffic_source.PLAYLIST 		= ratio
		case "NOTIFICATION":
			video.Traffic_source.NOTIFICATION 	= ratio
		}
	}
	//m, _ := json.MarshalIndent(video,"","    ")
	//fmt.Println(string(m))	
}

func updateVideoExternalSites(video *Video) {
	enddate_today := time.Now().Format("2006-01-02")
	filter_query := fmt.Sprintf("video==%s;insightTrafficSourceType==EXT_URL", video.Video_id)

	response := callYTAnalyticsAPI(
		"insightTrafficSourceDetail", // dimentions
		"views", // metrics
		filter_query, // filters
		STARTDATE_SHOWINT, // startdate
		enddate_today, // enddate
		"-views", // sort
		10, // maxresult
	)

	//m, _ := json.MarshalIndent(response,"","    ")
	//fmt.Println(string(m))

	for _, row := range response.Rows {
		site_name := row[0].(string)
		site_count := row[1].(float64)
		site_ratio := truncFloat( site_count / video.Traffic_source.EXT_URL_count * 100 )
		site_list := map[string]float64{ site_name : site_ratio }
		video.External_sites = append(video.External_sites, site_list)
	}

}

func updateAgePercentage(video *Video) {
	enddate_today := time.Now().Format("2006-01-02")
	filter_query := fmt.Sprintf("video==%s", video.Video_id)

	response := callYTAnalyticsAPI(
		"ageGroup", // dimentions
		"viewerPercentage", // metrics
		filter_query, // filters
		STARTDATE_SHOWINT, // startdate
		enddate_today, // enddate
		"", // sort
		10, // maxresult
	)

	//m, _ := json.MarshalIndent(response,"","    ")
	//fmt.Println(string(m))

	for _, row := range response.Rows {
		switch row[0] {
		case "age13-17":
			video.Age_percentage.AGE13_17 = row[1].(float64)
		case "age18-24":
			video.Age_percentage.AGE18_24 = row[1].(float64)
		case "age25-34":
			video.Age_percentage.AGE25_34 = row[1].(float64)
		case "age35-44":
			video.Age_percentage.AGE35_44 = row[1].(float64)
		case "age45-54":
			video.Age_percentage.AGE45_54 = row[1].(float64)
		case "age55-64":
			video.Age_percentage.AGE55_64 = row[1].(float64)
		case "age65-":
			video.Age_percentage.AGE65_ = row[1].(float64)
		}
	}
}

func updateGenderPercentage(video *Video) {
	enddate_today := time.Now().Format("2006-01-02")
	filter_query := fmt.Sprintf("video==%s", video.Video_id)

	response := callYTAnalyticsAPI(
		"gender", // dimentions
		"viewerPercentage", // metrics
		filter_query, // filters
		STARTDATE_SHOWINT, // startdate
		enddate_today, // enddate
		"", // sort
		10, // maxresult
	)

	//m, _ := json.MarshalIndent(response,"","    ")
	//fmt.Println(string(m))

	for _, row := range response.Rows {
		switch row[0] {
		case "male":
			video.Gender_percentage.MALE = row[1].(float64)
		case "female":
			video.Gender_percentage.FEMEL = row[1].(float64)
		case "user_specified":
			video.Gender_percentage.USER_SPECIFIED = row[1].(float64)
		}
	}
}


func gatherVideoStats(startdate string, enddate string, today string) []Video {
	var video_list_final []Video

	video_list_init := getVideoList(startdate, enddate)
	for _, video := range video_list_init {
		fmt.Printf("Gathering Video Stats : %s\n", video.Video_title )
		calcDuration(&video, today)
		updateVideoCount(&video)
		updateVideoTrafficSourceType(&video)
		updateVideoExternalSites(&video)
		updateAgePercentage(&video)
		updateGenderPercentage(&video)
		// TODO : 正常にAnnoutation関連が取得できないので、一旦停止
		updateAnnoutationImplession(&video)
		

		video_list_final = append(video_list_final, video)
	}
	return video_list_final
}

func getherThumbnailImages(video_list []Video) {
	for _, video := range video_list {
		image_name := "reports/images/thumbnail_" + video.Video_id + ".jpg"
		_, err := os.Stat(image_name)
		if os.IsNotExist(err) {
			fmt.Printf("Downloading Thumbnail iamge : %s", image_name)
			downloadImage(video.Thumbnail_url, image_name)
			trim_YT_Thumbnail(image_name)
		}else {
			fmt.Printf("%s is already exist. Skipping...\n", image_name)
		}
	}
}


/*
func getVideoStats(startdate string, enddate string) []Video {

	API_KEY := os.Getenv("API_KEY")
	CHANNEL_ID := os.Getenv("CHANNEL_ID")

	var video Video
	var video_list []Video

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		log.Fatalf("Error creating new YouTube service: %v", err)
	}

	//This is draft code for YouTubeAnalitycs API with Oauth2.0
	//client := getClient(youtube.YoutubeReadonlyScope)
	//service, err := youtube.New(client)

	// MAX result は 50個以下でしか動作できない
	call_search := service.Search.List([]string{"id", "snippet"}).ChannelId(CHANNEL_ID).Order("date").MaxResults(50).PublishedAfter(startdate).PublishedBefore(enddate)
	response_search, err_search := call_search.Do()
	if err_search != nil {
		log.Fatalf("Error making search API call: %v", err_search)
	}

	for _, search := range response_search.Items {
		video.Video_title = search.Snippet.Title
		video.Video_id = search.Id.VideoId
		video.Updated_date = strings.Split(search.Snippet.PublishedAt, "T")[0]
		video.Thumbnail_url = search.Snippet.Thumbnails.High.Url

		call_video := service.Videos.List([]string{"snippet", "contentDetails", "statistics"}).Id(search.Id.VideoId)
		response_video, err := call_video.Do()
		if err != nil {
			log.Fatalf("Error making Video API call: %v", err)
		}
		for _, video_stats := range response_video.Items {
			video.View_counts = video_stats.Statistics.ViewCount
			video.Like_counts = video_stats.Statistics.LikeCount
			video.Dislike_counts = video_stats.Statistics.DislikeCount
		}

		// 公式ページには、Go言語でYouTube Analytics APIを利用できるサンプルが無かった。
		// インプレッション数や視聴者の性別、トラフィック流入元などのデータを取るには Analytics APIへのアクセスが必須。
		// ふと見つけた google.golang.org/api/youtubeanalytics/v2 が利用可能かもしれない。
		// https://pkg.go.dev/google.golang.org/api/youtubeanalytics/v2#pkg-functions
		// ただし公開日が2023/1/23。まだReadyじゃない？いちおうGoogleの正式なプロダクトに見える。
		// まだ正しい使い方がわからないので、要検討。だめなら、Pythonで書き直す必要があるかも。
		// https://developers.google.com/youtube/reporting/v1/code_samples/python
		// 実現できる手段が無いわけではないが、Go言語で実現できるかは微妙。要検証。
		// この資料、読んで見る価値あるかも
		// https://blog.codecamp.jp/youtube-analytics-python

		//TODO: This is draft code for YouTubeAnalitycs API with Oauth2.0
		//client_analytics := getClient(youtubeanalytics.YtAnalyticsReadonlyScope)
		//service_analytics, err := youtubeanalytics.New(client_analytics)
		//if err != nil {
		//	log.Fatalf("Error creating new YouTube service: %v", err)
		//}
		// 動作OK
		//call_analytics := service_analytics.Reports.Query().Ids("channel==MINE").Dimensions("ageGroup").Metrics("viewerPercentage").StartDate("2023-01-01").EndDate("2023-02-11")

		// 動作OK
		//call_analytics := service_analytics.Reports.Query().Ids("channel==MINE").Dimensions("channel").Metrics("views,likes").StartDate("2023-01-01").EndDate("2023-02-11")

		// 動作OK
		//call_analytics := service_analytics.Reports.Query().Ids("channel==MINE").Dimensions("day").Metrics("views,likes").StartDate("2023-01-01").EndDate("2023-02-11")

		// Ids("channel==MINE") は必須。
		// Dimensions("video") がなぜが使えない。
		// Dimensions("day")やDimensions("ageGroup")、Dimensions("channel")は使えてる。
		// 2023/03/03 18:41:27 Error making Analytics API call: googleapi: Error 400: The query is not supported. Check the documentation at https://developers.google.com/youtube/analytics/v2/available_reports for a list of supported queries., badRequest

		// 動いた！ .Sort("-views")が必須ぽい。
		// annotationImpressions は結果が0になる。これはAPI server側ががまだ対応していないのかもしれない。なんかいろいろ制約があるようにみえるな。
		// Note: YouTube Analytics API reports only return data for the annotationClickThroughRate and annotationCloseRate metrics as of June 10, 2012. In addition, YouTube Analytics API reports only return data for the remaining annotation metrics as of July 16, 2013.
		// 動作OK

		//TODO: This is draft code for YouTubeAnalitycs API with Oauth2.0
		//call_analytics := service_analytics.Reports.Query().Ids("channel==MINE").Dimensions("video").Metrics("views,likes,//////annotationImpressions").StartDate("2022-01-01").EndDate("2023-02-11").MaxResults(10).Sort("-views")

		//response_anlytics, err := call_analytics.Do()
		//if err != nil {
		//	log.Fatalf("Error making Analytics API call: %v", err)
		//}

		video_list = append(video_list, video)

		image_name := "reports/images/thumbnail_" + video.Video_id + ".jpg"
		downloadImage(video.Thumbnail_url, image_name)
		trim_YT_Thumbnail(image_name)
	}
	video_list = reverseVideoList(video_list)
	return video_list
}
*/

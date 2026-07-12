# YouTubeAnalyticsReport
The script will get videos for an arbitrarily designated time period　on your YouTube Channel.

# Generated report sample

See [ Report Sample](./reports/showint_report_20231021.md)

# How to run this script

Download Impresion data from YouTube Studio.

YouTube Studio ->　左メニュー「アナリティクス」-> 「コンテンツ」タブ -> グラフ左下の「詳細」をクリック

グラフが出たら、左メニュー「内訳」をクリックして「コンテンツ」を選ぶ。-> 右上アイコンからエクスポートで「カンマ区切り形式 .csv」を選択

Example
https://studio.youtube.com/channel/UCpO3RcIrPaDJJ0Q3cbZrvZA/analytics/tab-content/period-default/explore?entity_type=CHANNEL&entity_id=UCpO3RcIrPaDJJ0Q3cbZrvZA&ur_dimensions=CREATOR_CONTENT_TYPE&ur_values=%27VIDEO_ON_DEMAND%27&ur_inclusive_starts=&ur_exclusive_ends=&time_period=1771833600000%2C1783753200000&explore_type=TABLE_AND_CHART&metrics_computation_type=DELTA&metric=EXTERNAL_VIEWS&granularity=DAY&t_metrics=EXTERNAL_VIEWS&t_metrics=VIDEO_THUMBNAIL_IMPRESSIONS&t_metrics=VIDEO_THUMBNAIL_IMPRESSIONS_VTR&t_metrics=AVERAGE_WATCH_TIME&v_metrics=EXTERNAL_VIEWS&v_metrics=VIDEO_THUMBNAIL_IMPRESSIONS&v_metrics=VIDEO_THUMBNAIL_IMPRESSIONS_VTR&v_metrics=AVERAGE_WATCH_TIME&dimension=VIDEO&o_column=EXTERNAL_VIEWS&o_direction=ANALYTICS_ORDER_DIRECTION_DESC



Make `.env` file with your YouTube API KEY(YouTube Data API v3) & Channel ID.
```
vi .env

API_KEY="XXXXXXXX"
CHANNEL_ID="XXXXXXXX"
```

Update `startdate` & `enddate` on [main.go](main.go) like this

```
	startdate := "2023-05-15T00:00:00Z"
	enddate := "2023-07-24T23:00:00Z"
```

Run the scprits like this. 

```
$ go run main.go oauth2.go report.go youtube.go image.go
```
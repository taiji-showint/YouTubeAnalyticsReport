# show int レポート

## 動画名

{{ range . }}
1. [{{ .Video_title }}](https://www.youtube.com/watch?v={{ .Video_id }})  
 ( {{ .Updated_date }} 公開)
{{ end }}

{{ range . }}
|項目|内容|
|---|---|
|動画名|{{ .Video_title }}|
|動画URL|https://www.youtube.com/watch?v={{ .Video_id }}|
|動画公開日|{{ .Updated_date }}|
|集計期間|{{ .Updated_date }} ~ {{ .Today }} ( {{ .Duration }} 日間 ) |
|サムネイル|<img src="images/thumbnail_{{ .Video_id }}_trim.jpg" width="560" height="315">|
|再生回数|{{ .View_counts }} 回|
|グッド回数|{{ .Like_counts }}|
|バッド回数|{{ .Dislike_counts }}|
|インプレッション数|{{ printf "%.0f" .Impressions }} 回|
|インプレッションからのクリック率|{{ printf "%.2f" .CTR }} %|
|視聴者の年齢と性別| 男性: {{ .Gender_percentage.MALE }} %  女性: {{ .Gender_percentage.FEMEL  }}%<br>13～17 歳 {{ .Age_percentage.AGE13_17 }}%        25～34 歳 {{ .Age_percentage.AGE25_34 }}%<br>18～24 歳 {{ .Age_percentage.AGE18_24 }}%        35～44 歳 {{ .Age_percentage.AGE35_44 }}%<br>45～54 歳 {{ .Age_percentage.AGE45_54 }}%        55～64 歳 {{ .Age_percentage.AGE55_64 }}%        65 歳以上 {{ .Age_percentage.AGE65_ }}% |
|トラフィック流入元|show int 登録者へのおすすめ : {{ .Traffic_source.SUBSCRIBER }}% <br> show int チャンネルページ : {{ .Traffic_source.YT_CHANNEL }}% YouTube関連動画 : {{ .Traffic_source.RELATED_VIDEO }}%  <br> YouTube検索 : {{ .Traffic_source.YT_SEARCH }}% <br> 外部サイトからの流入 : {{ .Traffic_source.EXT_URL_ratio }}%|
|外部サイトからの流入の内訳|{{- range $idx, $ext_site_list := .External_sites }}{{ if gt $idx 0 }}<br>{{ end }}{{ range $ext_site_name, $ext_site_ratio := $ext_site_list }}{{ $ext_site_name }} : {{ $ext_site_ratio }}%{{ end }}{{ end }}|

<div style="page-break-before:always"></div>
{{ end }}

---

## 応援いただいた期間の show int チャンネル全体の統計情報

| 項目 | 内容 |
|---|---|
| 応援期間 | {{ (index . 0).Updated_date }} ~ {{ (index . 0).Today }} |
| 視聴回数 | {{ (index . 0).TotalViews }} 回 |
| 総再生時間 | {{ (index . 0).TotalWatchHours }} 時間 |
| インプレッション数 | {{ (index . 0).TotalImpressions }} 回 |
| インプレッションからのクリック率 | {{ printf "%.1f" (index . 0).ChannelCTR }} % |
| 総チャンネル登録者数 | {{ (index . 0).TotalSubscribers }} 名 |
| 新規チャンネル登録者数 | +{{ (index . 0).NewSubscribers }} 名 |

### 視聴者の年齢と性別（チャンネル全体）

| 項目 | 割合 |
|---|---|
| 男性 | {{ printf "%.1f" (index . 0).ChannelGender.MALE }} % |
| 女性 | {{ printf "%.1f" (index . 0).ChannelGender.FEMEL }} % |
| 13～17 歳 | {{ printf "%.1f" (index . 0).ChannelAge.AGE13_17 }} % |
| 18～24 歳 | {{ printf "%.1f" (index . 0).ChannelAge.AGE18_24 }} % |
| 25～34 歳 | {{ printf "%.1f" (index . 0).ChannelAge.AGE25_34 }} % |
| 35～44 歳 | {{ printf "%.1f" (index . 0).ChannelAge.AGE35_44 }} % |
| 45～54 歳 | {{ printf "%.1f" (index . 0).ChannelAge.AGE45_54 }} % |
| 55～64 歳 | {{ printf "%.1f" (index . 0).ChannelAge.AGE55_64 }} % |
| 65 歳以上 | {{ printf "%.1f" (index . 0).ChannelAge.AGE65_ }} % |

### 期間内のコンテンツ別 視聴回数上位10位

{{ range $idx, $video := (index . 0).TopTenVideos -}}
{{ add $idx 1 }}. [{{ $video.Video_title }}](https://www.youtube.com/watch?v={{ $video.Video_id }}) (視聴回数: {{ printf "%.0f" $video.View_counts }} 回)

<img src="images/thumbnail_{{ $video.Video_id }}_trim.jpg" width="320" height="180">

---

{{ end }}

---

## 謝辞

この度は show int の活動にご理解とご支援をいただきまして誠にありがとうございました。

引き続きどうぞお付き合いいただけますと幸いです。

{{ (index . 0).Today }}

show int 土屋 太二

taiji@showint.net

![show int](images/showint_logo.png)

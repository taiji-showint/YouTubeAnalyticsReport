# show int レポート

## 動画名

{{ range . }}
1. [{{ .Video_title }}](https://www.youtube.com/watch?v={{ .Video_id }})  
 ( {{ .Updated_date }} 公開)
{{ end }}

{{ range . }}
|---|---|
|動画名|{{ .Video_title }}|
|動画URL|https://www.youtube.com/watch?v={{ .Video_id }}|
|動画公開日|{{ .Updated_date }}|
|集計期間|{{ .Updated_date }} ~ TODAY |
|サムネイル|<img src="images/thumbnail_{{ .Video_id }}_trim.jpg">|
|再生回数|{{ .View_counts }} 回|
|グッド回数|{{ .Like_counts }}|
|バッド回数|{{ .Dislike_counts }}|
|インプレッション数|  回|
|インプレッションからのクリック率|  %|
|視聴者の年齢と性別| 男性: {{ .Age_percentage.male }} %  女性: {{ .Age_percentage.female  }} |
|| 13～17 歳 {{ .Age_percentage.age13-17 }}  18～24  歳 {{ .Age_percentage.age18-24 }}   |
|トラフィック流入元||
<div style="page-break-before:always"></div>
{{ end }}

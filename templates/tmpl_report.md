# show int レポート

## 動画名

{{ range . }}
1. [{{ .Video_title }}](https://www.youtube.com/watch?v={{ .Video_id }})  
 ( {{ .Updated_date }} 公開)
{{ end }}

{{ range . }}
|||
|---|---|
|動画名|{{ .Video_title }}|
|動画URL|https://www.youtube.com/watch?v={{ .Video_id }}|
|動画公開日|{{ .Updated_date }}|
|集計期間|{{ .Updated_date }} ~ {{ .Today }} ( {{ .Duration }} 日間 ) |
|サムネイル|<img src="images/thumbnail_{{ .Video_id }}_trim.jpg">|
|再生回数|{{ .View_counts }} 回|
|グッド回数|{{ .Like_counts }}|
|バッド回数|{{ .Dislike_counts }}|
|インプレッション数| {{ .AnnotationImpressions }} 回|
|インプレッションからのクリック率| {{ .AnnotationClickThroughRate }} %|
|視聴者の年齢と性別| 男性: {{ .Gender_percentage.MALE }} %  女性: {{ .Gender_percentage.FEMEL  }}%  
13～17 歳 {{ .Age_percentage.AGE13_17 }}%
18～24 歳 {{ .Age_percentage.AGE18_24 }}%   
25～34 歳 {{ .Age_percentage.AGE25_34 }}%
35～44 歳 {{ .Age_percentage.AGE35_44 }}% 
44～54 歳 {{ .Age_percentage.AGE45_54 }}%  
55～64 歳 {{ .Age_percentage.AGE55_64 }}% 
65 歳以上 {{ .Age_percentage.AGE65_ }}% 
|
|トラフィック流入元|
show int 登録者へのおすすめ : {{ .Traffic_source.SUBSCRIBER }}% 
show int チャンネルページ : {{ .Traffic_source.YT_CHANNEL }}% 
show int 再生リスト : {{ .Traffic_source.PLAYLIST }}% 
show int 登録者への通知機能 : {{ .Traffic_source.NOTIFICATION }}%
YouTube関連動画 : {{ .Traffic_source.RELATED_VIDEO }}% 
YouTube検索 : {{ .Traffic_source.YT_SEARCH }}% 
その他のYouTube機能 : {{ .Traffic_source.YT_OTHER_PAGE }}% 
その他 : {{ .Traffic_source.NO_LINK_OTHER }}% 
外部サイトからの流入 : {{ .Traffic_source.EXT_URL_ratio }}% 
外部サイトからの流入の内訳 :
{{- range $_, $ext_site_list := .External_sites }}
  {{- range $ext_site_name, $ext_site_ratio := $ext_site_list }}
    {{ $ext_site_name }} : {{ $ext_site_ratio }}%
  {{- end}}    
{{- end}}
|
<div style="page-break-before:always"></div>
{{ end }}

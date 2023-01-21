# YouTubeAnalyticsReport
YouTube Analytics Report for Show Int Sponsors

Add your [YouTube API key](https://developers.google.com/youtube/registering_an_application) and your YouTube Channel ID into `.env` file.

```
touch .env
vi .env

API_KEY="XXXXXXXX"
CHANNEL_ID="XXXXXXXX"
```

Run `main.go`. This script will get the most recent 15 videos on your YouTube Channel.

```
$ go run main.go
==================================================
Channel Title:  show int インターネットの裏側解説
Channel ID:  UCpO3RcIrPaDJJ0Q3cbZrvZA
Suscribers: 3780
==================================================
----------------------------------------------------------------------------------------------------
Video Title:  show int の2022年振り返り &amp; 2023年に挑戦したいこと
Video ID:  &{ youtube#video  vn_1gqEh8hw [] []}
Uploaded Date:  2023-01-16T09:00:26Z
----------------------------------------------------------------------------------------------------
Video Title:  【開催直前】JANOGスタッフがおすすめするJANOG51 Meeting みどころ紹介【最新技術から現地観光情報まで】
Video ID:  &{ youtube#video  NFNNBOflnkQ [] []}
Uploaded Date:  2023-01-09T09:00:18Z
----------------------------------------------------------------------------------------------------
Video Title:  フルリモートワークで働く現役エンジニアが地方移住して一年半経って見えてきた実際のところ
Video ID:  &{ youtube#video  EnqjEN_7ngU [] []}
Uploaded Date:  2023-01-02T03:00:06Z
----------------------------------------------------------------------------------------------------
Video Title:  Twitter・Facebookなどのグローバル企業で人員削減が相次いだ2022年。「大規模人員削減時代」にエンジニアがいま考えておくべきこと
Video ID:  &{ youtube#video  MDqkONWiu1k [] []}
Uploaded Date:  2022-12-26T09:00:17Z
----------------------------------------------------------------------------------------------------
Video Title:  ワールドカップ2022 の裏側でネットワークエンジニア達がヒヤヒヤしていた事情を解説します【ABEMAおよび各社エンジニアのみなさまお疲れさまでした】
Video ID:  &{ youtube#video  iQg49y7KSVU [] []}
Uploaded Date:  2022-12-19T09:00:17Z
----------------------------------------------------------------------------------------------------
Video Title:  トップレベルのエンジニアに求められるコミュニケーション力・リーダーシップ力
Video ID:  &{ youtube#video  FkAUTYip1n8 [] []}
Uploaded Date:  2022-12-13T10:03:19Z
----------------------------------------------------------------------------------------------------
Video Title:  アジア各国を飛び回るスターエンジニア maz 松崎 吉伸さん【エンジニア対談】
Video ID:  &{ youtube#video  YkvU8JmF2qw [] []}
Uploaded Date:  2022-12-05T09:00:23Z
----------------------------------------------------------------------------------------------------
Video Title:  「ネットワークエンジニア未経験者歓迎！」は本当に未経験でもOKなのか
Video ID:  &{ youtube#video  f4Rb4HIuBJI [] []}
Uploaded Date:  2022-11-28T09:00:31Z
----------------------------------------------------------------------------------------------------
Video Title:  複数の会社から転職オファーをもらったときの選び方【視聴者さんからのお悩み相談】
Video ID:  &{ youtube#video  hB4CBLmK0lw [] []}
Uploaded Date:  2022-11-21T09:00:24Z
----------------------------------------------------------------------------------------------------
Video Title:  衛星インターネット Starlink が爆速な理由 を現役ネットワークエンジニアが解説してみた
Video ID:  &{ youtube#video  Kfnd_iA_rm8 [] []}
Uploaded Date:  2022-11-14T09:00:09Z
----------------------------------------------------------------------------------------------------
Video Title:  現役ネットワークエンジニアが実践する「モチベーションに頼らない成長戦略」
Video ID:  &{ youtube#video  SGJVDqkquT4 [] []}
Uploaded Date:  2022-11-07T09:00:18Z
----------------------------------------------------------------------------------------------------
Video Title:  日本人が陥りがちな外国人エンジニアとのコミュニケーション失敗例
Video ID:  &{ youtube#video  fHZpgFVRlNY [] []}
Uploaded Date:  2022-10-31T09:00:26Z
----------------------------------------------------------------------------------------------------
Video Title:  StackStorm で実現するネットワーク完全自動化の世界
Video ID:  &{ youtube#video  Z8e5vso5SRU [] []}
Uploaded Date:  2022-10-24T09:00:14Z
----------------------------------------------------------------------------------------------------
Video Title:  【ストレス爆発寸前！】子育てエンジニアのストレス・悩みを育児経験者にぶつけてみた
Video ID:  &{ youtube#video  EFMTGgQEwTc [] []}
Uploaded Date:  2022-10-17T09:00:07Z
----------------------------------------------------------------------------------------------------
Video Title:  ロシア・ウクライナ軍事侵攻によるインターネット分断の危機【スプリンターネット】【JANOG50プレゼン解説】
Video ID:  &{ youtube#video  XC5J_Qbv7eE [] []}
Uploaded Date:  2022-10-03T09:00:30Z
```
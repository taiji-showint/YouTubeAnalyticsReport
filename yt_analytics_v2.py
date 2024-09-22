# -*- coding: utf-8 -*-

# Documents
# Official Reference: https://developers.google.com/youtube/reporting
# Official Sample: https://github.com/youtube/api-samples/tree/07263305b59a7c3275bc7e925f9ce6cabf774022/python
# Sample API call : https://developers.google.com/youtube/analytics/sample-requests#channel-geographic-reports
# Good blog: https://blog.codecamp.jp/youtube-analytics-python

import os
from pprint import pprint

import google.oauth2.credentials
import google_auth_oauthlib.flow
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from google_auth_oauthlib.flow import InstalledAppFlow

SCOPES = ['https://www.googleapis.com/auth/yt-analytics.readonly']

API_SERVICE_NAME = 'youtubeAnalytics'
API_VERSION = 'v2'
CLIENT_SECRETS_FILE = 'client_secret.json'

def get_service():
  flow = InstalledAppFlow.from_client_secrets_file(CLIENT_SECRETS_FILE, SCOPES)
  credentials = flow.run_local_server(port=0)
  return build(API_SERVICE_NAME, API_VERSION, credentials = credentials)

def execute_api_request(client_library_function, **kwargs):
  response = client_library_function(
    **kwargs
  ).execute()

  pprint(response)

if __name__ == '__main__':
  # Disable OAuthlib's HTTPs verification when running locally.
  # *DO NOT* leave this option enabled when running in production.
  os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = '1'

  youtubeAnalytics = get_service()
  video_id = "ePJIVZZhzRg"
  execute_api_request(
      youtubeAnalytics.reports().query,
      ids='channel==MINE',
      maxResults='10',
      startDate='2024-09-10',
      endDate='2024-09-20',

      # 動画ごとの視聴回数
      #dimensions='video',
      #metrics='views',
      #sort='-views',

      # 1日あたりのトラフィックソースごとの視聴回数
      #dimensions='day,insightTrafficSourceType',
      #metrics='views,estimatedMinutesWatched',
      #sort='day',


      # トップ 10 – 動画のトラフィックが最も多い外部ウェブサイト
      dimensions='insightTrafficSourceDetail',
      metrics='estimatedMinutesWatched,views',
      filters='video==ePJIVZZhzRg;insightTrafficSourceType==EXT_URL',
      sort='-estimatedMinutesWatched',
  )
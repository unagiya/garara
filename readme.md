# garara
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## APIs

### v1

| API                    | Description                                                      | Implements |
|:-----------------------|:-----------------------------------------------------------------|-----------:|
| 配信予約API<br>（キューイングモード） | 配信予約及び即時配信指示を行います                                                |       100% |
| 配信予約削除API              | 配信予約中の配信データを削除することができる API です。<br>配信予約データ削除機能、重複処理抑止機能があります。     |       100% |
| 配信ステータス取得API           | 配信や削除処理などの、処理結果を取得するAPIです。                                       |       100% |
| 配信結果リスト取得API           | 実際に配信処理を行ったメールアドレスの配信結果を取得するAPIです。                               |       100% |
| エラーリスト取得API            | 配信を行った結果、エラーとなったアドレスについての情報を取得するAPIです。                           |       100% |
| クリックカウントログ取得API        | クリックカウント・開封確認の詳細ログ取得APIとなります。                                    |       100% |
| エラーフィルターリスト取得API       | 管理画面のエラーリスト管理、配信除外リスト管理に登録されていメールアドレスを取得するAPIです。登録済みの全データを取得します。 |       100% |

配信結果リスト取得APIの実装でRequestID, 期間での結果取得を行うと、詳細リストが取得できない `※1` ので実装見送り。
※1 : 私の環境では取得できない状態だった。

追加メソッド
- SimpleV1SendListBuilder
  - メールアドレスのstring sliceからSendListを作成する。
  - Addressのdevice type は0（NONE）を付与する。
- SimpleV1SendListDataBuilder
  - id, deviceType, メールアドレス, intText, exttext, extImageの各string slice, keyFieldを引数に取り、SendListのDataを一つ作成する
  - intText, exttext, extImageの各string sliceはIDを0から自動採番する。
  - device typeは存在しない番号を入力するとエラーを返す
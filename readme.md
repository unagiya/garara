# garara
gararaは[アララメッセージ](https://am.arara.com/)を利用するためのサードパーティライブラリです。

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://opensource.org/licenses/mit-license.php)

## APIs

### v1

| API                    | Description                                                      | Implements | Method                                                             |
|:-----------------------|:-----------------------------------------------------------------|-----------:|:-------------------------------------------------------------------|
| 配信予約API<br>（キューイングモード） | 配信予約及び即時配信指示を行います                                                |       100% | SendQueueMode                                                      |
| 配信予約削除API              | 配信予約中の配信データを削除することができる API です。<br>配信予約データ削除機能、重複処理抑止機能があります。     |       100% | DeleteReservationByDeliverIDs<br>DeleteReservationByRequestIDs     |
| 配信ステータス取得API           | 配信や削除処理などの、処理結果を取得するAPIです。                                       |       100% | GetStatusByDeliverIDs<br>GetStatusByRequestIDs<br>GetStatusByTerms |
| 配信結果リスト取得API           | 実際に配信処理を行ったメールアドレスの配信結果を取得するAPIです。                               |       100% | GetResultListByDeliverIDs                                          |
| エラーリスト取得API            | 配信を行った結果、エラーとなったアドレスについての情報を取得するAPIです。                           |       100% | GetErrorListByDeliverIDs<br>GetErrorListByTerms                    |
| クリックカウントログ取得API        | クリックカウント・開封確認の詳細ログ取得APIとなります。                                    |       100% | GetClickLog                                                        |
| エラーフィルターリスト取得API       | 管理画面のエラーリスト管理、配信除外リスト管理に登録されていメールアドレスを取得するAPIです。登録済みの全データを取得します。 |       100% | GetErrorFilter                                                     |

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
- SimpleV1MailContentsBuilder
  - メール配信のためのContentsを作成します。
  - imagesはContent.Imageにid,0から順番にデータを作成します。
  - textsはContent.Textにid,0から順番にデータを作成します。
  - filesがattachMaxFileSize（現仕様では3）よりデータ量が大きい場合エラーを返却します。
- SimpleV1SettingBuilder
  - メール配信のためのsetting structを作成します。
  - S/MIME, Openedの指定が範囲外の場合エラーを返却します。
- V1MailDeliveryBuilder
  - メール配信のためのdelivery structを作成します。
- AttrIdCDataBuilder, AttrIdStringsBuilder
  - []stringから[]AttrIdCdata, []AttrIdStringを作成します。
  - idは0から順に振られます。

#### v1 api usage

メール送信
``` go
  c := garara.NewDefaultV1Client(
    [API User], 
    [API Password], 
    [ManageSiteUser], 
    [ManageSitePassword], 
    [SiteID], 
    [ServiceID],
    )
  
  var r garara.V1MailRequest
  u, _ := uuid.NewRandom()
  s, _ := garara.SimpleV1SettingBuilder(
    "now",
    "From Name",
    "from-address@example.com",
    "",
    0,
    garara.UNUSE,
    garara.USE,
    garara.Option{},
  )
  
  subject := "test message"
  body := `
<html>
  <head></head>
  <body>
    Hello Message!!
  </body>
</html>
`
  cont, _ := garara.SimpleV1MailContentsBuilder(subject, body, garara.HTML, nil, nil, nil)
  l := garara.SimpleV1SendListBuilder([]string{"send-mail-1@example.com", "send-mail-2@example.com"})
  r.Delivery = garara.V1MailDeliveryBuilder(1, u.String(), s, cont, l)
  
  ctx := context.Background()
  res, _ := c.SendQueueMode(ctx, r, [RequestHostName])
```

※他のAPIはある程度直感的な使い方ができるかと思います。
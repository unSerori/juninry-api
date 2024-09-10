# juninry

[RepositoryURL](https://github.com/unSerori/juninry-api)

## 概要

[docker-juninry](https://github.com/unSerori/docker-juninry)のGo APIサーバ

### 環境

macOS Sonoma version: 14.6.1
Visual Studio Code: 1.88.1  
Docker Desktop: 4.34.0: Engine: 27.2.0  
image: golang:1.22.2-bullseye  
image: mysql:latest

## 環境構築

[docker-juninry](https://github.com/unSerori/docker-juninry)を使ってDokcerコンテナーで開発・デプロイする形を想定している  
構築手順は[docker-juninryの環境構築項目](https://github.com/unSerori/docker-juninry/blob/main/README.md#環境構築)に記載  
cloneしてスクリプト実行で、自動的にコンテナー作成と開発環境（: またはデプロイ）を行う  
開発環境へのアタッチはVS CodeのDocker, DevContainer拡張機能の「Attach Visual Studio Code」を用いて、VS Codeの機能をそのまま使うことを想定している  
詳しくは[docker-juninryのREADME.md](https://github.com/unSerori/docker-juninry/blob/main/README.md)に記載

### ./vscode-ext.txt

Goのデバッグ用VS Code拡張機能や便利な拡張機能のリスト  
VS Codeアタッチしたコンテナー内で、以下のコマンド実行で一括インストールできる

```bash
cat vscode-ext.txt | while read line; do code --install-extension $line; done
```

### 自前でのローカル環境構築

想定はしていないが、ローカル環境にインストールすることも可能

1. [Goのインストール](https://go.dev/doc/install)
2. このリポジトリをclone

    ```bash
    git clone https://github.com/unSerori/juninry-api
    ```

3. [.env](#env)ファイルをもらうか作成
4. 必要なパッケージの依存関係など

    ```bash
    go mod tidy
    ```

5. プロジェクトを起動

    ```bash
    # 実行(VSCodeならF5keyで実行)
    go run .

    # ワンファイルにビルド
    go build -o output 

    # ビルドで生成されたファイルを実行
    ./output
    ```

#### おまけ: Goでプロジェクト作成時のコマンド

```bash
# goモジュールの初期化
go mod init ddd

# ginのインストール
go get -u github.com/gin-gonic/gin

# main.goの作成
echo package main > main.go
```

## ディレクトリ構成

<details>
  <summary>`tree -LFa 3 --dirsfirst`に加筆修正</summary>

```txt
./juninry-api
|-- .git/
|-- application/
|-- asset/
|-- common/
|   `-- logging/
|       |-- init.go
|       |-- log.go
|       `-- server.log
|-- controller/
|-- domain/
|-- infrastructure/
|   |-- model/
|   `--/
|-- middleware/
|-- model/
|-- presentation/
|-- route/
|   |-- dig.go
|   `-- router.go
|-- service/
|-- upload/
|   |-- homework/
|   |-- t_material/
|-- utility/
|   |-- auth/
|   |-- batch/
|   |-- config/
|   |-- custom/
|   |-- dip/
|   |-- scheduler/
|   |-- security/
|   `-- utility.go
|-- view/
|   `-- views/
|       |-- scripts/
|       |   `-- common.js
|       |-- styles/
|       |   `-- common.css
|       `-- index.html
|-- .env
|-- .gitignore
|-- README.md
|-- go.mod
|-- go.sum
|-- init.go
|-- main.go
|-- request.rest
`-- vscode-ext.txt
```

</details>

現状アーキテクチャが一貫されていない  
TODO: DDDに統合していく予定

### 主要なディレクトリの説明

- presentation, application, domain, infrastructure: DDDの4パッケージ
- view: テスト用ページの静的ファイル
- middleware: ミドルウェアを置くが、この中でDDD形式などに分割すべきかを検討中
- route: ルーティングや付随する初期設定
- utility: 再利用性の高い単体の処理群
- common: utilityの中でもより一般性の高い処理群
- asset: サーバー自体が最初から持つリソースや画像送信テストなどで使うリソースを置いておく
- upload: アップロードされたファイル
- controller, service, model: 旧アーキテクチャの3パッケージで、DDD形式に変換したい  
  modelはテーブルモデル定義ファイルとしてinfrastructure/modelに移動したい

## API仕様書

エンドポイント、リクエストレスポンスの形式、その他情報のAPIの仕様書

### エンドポインツ

<details>
  <summary>サーバー側で疎通確認するエンドポイント</summary>

- **URL:** `/v1/test/cfmreq`
- **メソッド:** GET
- **説明:** 鯖側でリクエストが受け取れたか確認できる。グループを作ったときの疎通を確かめたりする野に使う。
- **リクエスト:**
  - ヘッダー:
  - ボディ:

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResCode": "OK",
        "srvResData": {
          "message": "hello go server!"
        }
      }      
      ```

</details>

</details>

<details>
  <summary>ユーザを作成するエンドポイント</summary>

- **URL:** `/v1/users/register`
- **メソッド:** POST
- **説明:** 新規ユーザを登録。
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json
  - ボディ:

    ```json
    {
      "userName": "test teacher",
      "userTypeId": 1,
      "mailAddress": "test-teacher@gmail.com",
      "password": "C@tt"
    }
    ```

- **レスポンス:**
  - ステータスコード: 201 Created
    - ボディ:

      ```json
      {
        "srvResMsg":  "Created",
        "srvResData": {
          "authenticationToken": "token@h",
        },
      }
      ```

</details>

<details>
  <summary>ログインするエンドポイント</summary>

- **URL:** `/v1/users/login`
- **メソッド:** POST
- **説明:** メアドとパスワードでログインし、トークンを取得する
- **リクエスト:**
- ヘッダー:
  - `Content-Type`: application/json
  - ボディ:

    ```json
    {
      "mailAddress": "test-pupil@gmail.com",
      "password": "C@tp"
    }
    ```

- **レスポンス:**
  - ステータスコード: 200 OK
  - ボディ:

    ```json
    {
      "srvResMsg":  "OK",
      "srvResData": {
        "authenticationToken": "token@hogeta"
      },
    }
    ```

</details>

<details>
  <summary>ユーザー情報を取得するエンドポイント</summary>

- **URL:** `/v1/auth/auth/users/user`
- **メソッド:** GET
- **説明:** jwtから取得したidからユーザーを検索して情報を返す
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json
    - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: ＊ステータスコード ステータス＊
    - ボディ:
      ＊さまざまな形式のレスポンスデータ（基本はJSON）＊

      ```json
      {
        "srvResMsg":  "Successful user get.",
        "srvResData": {
          "userData": {
            "userUUID": "3cac1684-c1e0-47ae-92fd-6d7959759224",
            "userName": "test pupil",
            "userTypeId": 2,
            "mailAddress": "test-pupil@gmail.com",
            "password": "$2a$10$8hJGyU235UMV8NjkozB7aeHtgxh39wg/ocuRXW9jN2JDdO/MRz.fW",
            "jwtUUID": "14dea318-8581-4cab-b233-995ce8e1a948",
            "ouchiUUID": null
          }
        }
      }
      ```

</details>

<details>
  <summary>新規宿題を登録するエンドポイント</summary>

- **URL:** `/v1/auth/users/homeworks/register`
- **メソッド:** POST
- **説明:** 教師権限を持つユーザーがクラスに対して宿題を登録する
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json
    - `Authorization`: (string) 認証トークン
  - ボディ:

    ```json
    {
      "homeworkLimit": "2024-08-2T23:59:59Z",
      "classUUID": "09eba495-fe09-4f54-a856-9bea9536b661",
      "homeworkNote": "がんばってくださ～い＾＾",
      "teachingMaterialUUID": "978f9835-5a16-4ac0-8581-7af8fac06b4e",
      "startPage": 2,
      "pageCount": 8
    }
    ```

- **レスポンス:**
  - ステータスコード: 201 Created
    - ボディ:

      ```json
      {
      "srvResMsg":  "201",
      "srvResData": {
          "homeworkUUID": "6e8ad122-2ca9-453b-92ba-65edaf786ec2"
        },
      }
      ```

</details>

<details>
  <summary>特定の月の課題の提出状況を取得するエンドポイント</summary>

- **URL** `/v1/auth/users/homeworks/record?targetMonth=2025-01-01 00:00:00.000Z`
- **メソッド** GET
- **説明** 送られてきた特定の月の各日に設定されている課題の数と提出状況を返す
- **リクエスト**
  - ヘッダー:
    - Authorization: (string) 認証トークン

- **レスポンス**:
  - ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResData": [
          {
            "limitDate": "2025-01-21T00:00:00Z",
            "submissionCount": 0,
            "homeworkCount": 2
          },,,
        ],
        "srvResMsg": "OK"
      }
      ```

  - ステータスコード: 403 Forbidden
    - ボディ:

      ```json
      {
        "srvResData": {},
        "srvResMsg": "Forbidden"
      }
      ```

</details>

<details>
  <summary>クラスの課題情報一覧を取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/homeworks/upcoming`
- **メソッド:** GET
- **説明:** 自分が所属するクラスの期限が先のものを取得
- **リクエスト:**
  - ヘッダー:
    - Authorization: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
        {
          "srvResMsg": "OK",
          "srvResData": [
            {
              "homeworkLimit": "0001-01-01T00:00:00Z",
              "homeworkData": [
                {
                  "homeworkUUID": "a3579e71-3be5-4b4d-a0df-1f05859a7104",
                  "startPage": 24,
                  "pageCount": 2,
                  "homeworkNote": "がんばってくださ～い＾＾",
                  "teachingMaterialName": "漢字ドリル3",
                  "subjectId": 1,
                  "subjectName": "国語",
                  "teachingMaterialImageUUID": "a575f18c-d639-4b6d-ad57-a9d7a7f84575",
                  "className": "3-2 ふたば学級",
                  "submitFlag": 1  // 提出フラグ 1 提出 0 未提出
                },,,
              ]
            },,,
          ]
        }
      ```

</details>

<details>
  <summary>特定の宿題に対する任意のユーザーの提出状況と宿題の詳細情報を取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/homeworks/{homework_uuid}`
- **メソッド:** GET
- **説明:** 特定の宿題の詳細情報を取得する。生徒はクエパラなしで自分の提出状況を、教師はクエパラ設定で特定生徒を、保護者は家庭内特定児童をクエパラで設定すると提出状況を見られる。
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
  - パラメーター
    - パスパラメーター:
      - `homework_uuid`: 宿題リソースを指定するパラメーター
    - クエリパラメーター
      - `user_uuid`: どの自動ユーザーの宿題状況を確認するかのクエパラ
    - パラメーター例

      ```url
      /v1/auth/users/homeworks/a3579e71-3be5-4b4d-a0df-1f05859a7104?user_uuid=3cac1684-c1e0-47ae-92fd-6d7959759224
      ```

  リクエスト例

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResMsg":  "OK",
        "srvResData": {
          "teachingMaterialUUID": "978f9835-5a16-4ac0-8581-7af8fac06b4e",
          "teachingMaterialName": "漢字ドリル3",
          "subjectId": 1,
          "startPage": 2,
          "pageCount": 8,
          "isSubmitted": true,  // or false
          "images": ["bbbbbbbb-a6ad-4059-809c-6df866e7c5e6.jpg, gggggggg-176f-4dea-bec0-21464f192869.jpg, rrrrrrrr-bb84-4565-9666-d53dfcb59dd3.jpg"]
        },
      }
      ```

</details>

<details>
  <summary>特定の提出済み宿題の画像を取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/homeworks/{homework_uuid}/images/{image_file_name}`
- **メソッド:** GET
- **説明:** 特定の提出済み宿題に紐づいている画像を取得する。一枚取得なのでそれぞれの画像に対してGETすべき
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      画像

</details>

<details>
  <summary>終了した宿題を提出するエンドポイント</summary>

- **URL:** `/v1/auth/users/homeworks/submit`
- **メソッド:** POST
- **説明:** 宿題を提出する
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: multipart/form-data
    - `Authorization`: (string) 認証トークン
  - ボディ: Form
    - Form Fields - 宿題のID
      - homeworkUUID: a3579e71-3be5-4b4d-a0df-1f05859a7104,
    - Files - 提出する宿題の画像
      - images: page_67.jpg
      - images: page_68.png

- **レスポンス:**
  - ステータスコード: 201 Created
    - ボディ:

      ```json
      {
        "srvResMsg":  "Created",
        "srvResData": {
        },
      }
      ```

</details>

<details>
  <summary>おてがみ情報一覧を取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/notice/notices`
- **メソッド:** GET
- **説明:** 自分が所属するクラスのおてがみ情報一覧取得
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json
    - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResData": {
          "notices": [
            {
              "noticeUUID": "51e6807b-9528-4a4b-bbe2-d59e9118a70d",
              "noticeTitle": "【持ち物】おべんとうとぞうきん",
              "noticeDate": "2024-07-27T10:53:22Z",
              "userName": "test teacher",
              "classUUID": "09eba495-fe09-4f54-a856-9bea9536b661",
              "className": "3-2 ふたば学級",
              "readStatus": 0 // 未読: 0, 既読: 1, 対象外: null
            },,,
          ]
        },
        "srvResMsg": "OK"
      }
      ```

  - ステータスコード: 403 Forbidden
    - ボディ:

      ```json
      {
        "srvResData": {},
        "srvResMsg": "Forbidden"
      }
      ```

  - ステータスコード: 404
    - ボディ:

      ```json
      {
        "srvResData": {},
        "srvResMsg": "Not Found"
      }
      ```

</details>

<details>
  <summary>おてがみの詳細情報を取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/notices/{notice_uuid}`
- **メソッド:** GET
- **説明:** パスパラメーターで指定したおしらせの詳細情報を取得する
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResData": {
          "noticeTitle": "【持ち物】習字道具必要です",
          "noticeExplanatory": "国語授業で習字を行いますので持たせていただくようお願いします",
          "noticeDate": "2024-07-16T00:45:47Z",
          "userName": "test teacher",
          "className": "3-2 ふたば学級",
          "classUUID": "09eba495-fe09-4f54-a856-9bea9536b661",
          "quotedNoticeUUID": "2097a7bb-5140-460d-807e-7173a51672bd",
          "readStatus": 0   // 未読: 0, 既読: 1, 対象外: null
        },
        "srvResMsg": "OK"
      }
      ```

</details>

<details>
  <summary>お知らせの新規登録するエンドポイント</summary>

- **URL:** `/v1/auth/users/notice/register`
- **メソッド:** POST
- **説明:** お知らせの新規登録をする
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json
    - `Authorization`: (string) 認証トークン
  - ボディ:

    ```json
      {
        "srvResData": {
          "notices": {
            "NoticeTitle": "【持ち物】習字道具必要です",
            "NoticeDate": "2024-06-11T03:23:39Z",
            "NoticeExplanatory": "国語授業で習字を行いますので持たせていただくようお願いします",
            "UserUuid": "9efeb117-1a34-4012-b57c-7f1a4033adb9",
            "ClassUui": "817f600e-3109-47d7-ad8c-18b9d7dbdf8b",
        }},
      }
    ```

- **レスポンス:**
  - ステータスコード: 200 Created
    - ボディ:

      ```json
      {
        "srvResData": {
          "authenticationToken": "トークン",
          "srvResMsg": "OK"
        },
      }
      ```

  - ステータスコード: 403 Forbidden
    - ボディ:

      ```json
      {
        "srvResData": {},
        "srvResMsg": "Forbidden"
      }
      ```

</details>

<details>
  <summary>お知らせ既読処理をするエンドポイント</summary>

- **URL:** `/v1/auth/users/notices/read/{notice_uuid}`
- **メソッド:** POST
- **説明:** notice_read_statusにデータを追加する
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:
  
      ```json
      {
        "srvResData": {},
        "srvResMsg": "OK" 
      }
      ```

<details>
  <summary>クラスメイトを取得するエンドポイント</summary>

- **URL:** `/＊エンドポイントパス＊`
- **メソッド:** GET
- **説明:** 自分のクラスのクラスメイトを取得
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
    - `Content-Type`: application/json

- **レスポンス:**
  - ステータスコード: ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResData": [
            {
              "className": "3-2 ふたば学級",
              "juniorData": [
                {
                  "userUUID": "3cac1684-c1e0-47ae-92fd-6d7959759224",
                  "userName": "test pupil",
                  "genderId": 1,
                  "studentNumber": null // 数字 or null
                }
              ]
            }
        ],
        "srvResMsg": "OK"
      }
      ```

</details>

</details>

<details>
  <summary>特定のお知らせ既読一覧取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/notices/status/{notice_uuid}`
- **メソッド:** GET
- **説明:** 先生が特定のお知らせの生徒の既読情報を取得する
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200　OK
    - ボディ:

      ```json
      {
        "srvResData": [
        {
          "StudentNo": 0,   //定義がないので0デス
          "UserName": "test oooo",
          "GenderCode": null, //定義がないのでnullです 
          "ReadStatus": 0
        }
        ],
        "srvResMsg": "Successful noticeStatus get."
      }
      ```

</details>

<details>
  <summary>所属クラス一覧を取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/classes/affiliations`
- **メソッド:** GET
- **説明:** 子供、教師は自身の所属するクラスを、親は子供たちの所属するクラスの一覧を取得
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json

- **レスポンス:**
  - ステータスコード: 200
    - ボディ:

    ```json
    {
      "srvResData": {
        "classes": [
          {
            "classUUID": "09eba495-fe09-4f54-a856-9bea9536b661",
            "className": "3-2 ふたば学級"
          },,,
        ]
      },
      "srvResMsg": "OK"
    }
    ```

  - ステータスコード: 404
    - ボディ:

    ```json
    {
      "srvResData": {},
      "srvResMsg": "Not Found"
    }
    ```

</details>

<details>
  <summary>クラスを新規作成するエンドポイント</summary>

- **URL:** `/v1/auth/users/classes/register`
- **メソッド:** POST
- **説明:** クラスを新規作成し、招待コードを発行する。新規作成を行なったユーザーはクラスに所属する。
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
    - `Content-Type`: application/json
  - ボディ:

    ```json
    {
      "className": "クラスを立てる"
    }
    ```

- **レスポンス:**
  - ステータスコード: 201 OK
    - ボディ:

      ```json
      {
      "srvResData": {
        "ouchiUUID": "fe9462d6-bd7e-4b04-8b6a-785e9231b4d5",
        "ouchiName": "テスト家",
        "inviteCode": "009574",
        "validUntil": "2024-07-16T13:44:02.603671112Z"
      },
      "srvResMsg": "Created"
      }
      ```

  - ステータスコード: 403 Forbidden
    - ボディ:

      ```json
      {
        "srvResMsg": "Forbidden",
        "srvResData": {}
      }
      ```

</details>

<details>
  <summary>クラスの招待IDを更新するエンドポイント</summary>

- **URL:** `v1/auth/users/classes/refresh/{class_uuid}`
- **メソッド:** PUT
- **説明:** クラスの招待IDを更新する
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResData": {
          "classUUID": "53faea61-ae69-45e9-8b66-73481f9ca879",
          "className": "最新のクラス",
          "inviteCode": "7895",
          "validUntil": "2024-07-04T03:15:25Z"
        },
        "srvResMsg": "Created"
      }
      ```

  - ステータスコード: 403 Forbidden
    - ボディ:

      ```json
      {
        "srvResMsg": "Forbidden",
        "srvResData": {}
      }
      ```

  - ステータスコード: 404 Not Found
    - ボディ:

      ```json
        {
          "srvResData": {},
          "srvResMsg": "Not Found"
        }
        ```

</details>

<details>
  <summary>クラスに所属するエンドポイント</summary>

- **URL:** `/v1/auth/users/classes/join/:invite_code`
- **メソッド:** POST
- **説明:** クラスに生徒、職員を所属させる。
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
    - `Content-Type`: application/json
  - ボディ: ※任意

    ```json
      "studentNumber": 20
    ```

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

    ```json
    {
      "srvResData": {
        "className": "ゆるふわ"
      },
      "srvResMsg": "OK"
    }
    ```

  - ステータスコード: 409 Conflict
    - ボディ:

    ```json
    {
      "srvResData": {},
      "srvResMsg": "Conflict"
    }
    ```

  - ステータスコード: 403 Forbidden
    - ボディ:

    ```json
    {
      "srvResData": {},
      "srvResMsg": "Forbidden"
    }
    ```

</details>

<details>
  <summary>おうちを新規登録するエンドポイント</summary>

- **URL:** `/v1/auth/users/ouchies/register`
- **メソッド:** POST
- **説明:** おうちを新規作成し、招待コードを発行する。新規作成を行なったユーザーはおうちに所属する。
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
    - `Content-Type`: application/json
  - ボディ:

    ```json
    {
      "ouchiName": "おうちを立てる"
    }
    ```

</details>

<details>
  <summary>おうち招待コード更新するエンドポイント</summary>

- **URL:** `/v1/auth/users/ouchies/refresh/{ouchi_uuid}`
- **メソッド:** PUT
- **説明:** おうち招待コードの更新
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
  
- **レスポンス:**
  - ステータスコード: 201 Created
    - ボディ:

      ```json
      {
        "srvResData": {
          "ouchiUUID": "6fd7caf3-9ec9-4487-917e-f0fa75fb5ad2",
          "ouchiName": "テスト3家",
          "inviteCode": "007019",
          "validUntil": "2024-07-17T05:31:39.384195368Z"
        },
        "srvResMsg": "Created"
      }
      ```

</details>

<details>
  <summary>おうちに所属するエンドポイント</summary>

- **URL:** `/v1/auth/users/ouchies/join/{invite_code}`
- **メソッド:** POST
- **説明:** ユーザにouchiUuidを付与する
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
        {
          "srvResData": {
            "ouchiName": "テスト3家"
          },
          "srvResMsg": "OK"
        }      
      ```

</details>

<details>
  <summary>おうち情報取得するエンドポイント</summary>

- **URL:** `/v1/auth/users/ouchies/info`
- **メソッド:** GET
- **説明:** おうち名、おうちに所属している人全員を取得する
- **リクエスト:**
  - `Authorization`: (string) 認証トークン

- **レスポンス:**
  - ステータスコード: 200　OK
    - ボディ:

      ```json
      {
        "srvResData": {
          "ouchiUUID": "2e17a448-985b-421d-9b9f-62e5a4f28c49",
          "ouchiName": "piyonaka家",
          "ouchiMembers": [
            {
              "userUUID": "868c0804-cf1b-43e2-abef-08f7ef58fcd0",
              "userName": "test parent",
              "userTypeId": 3,
              "genderId": 0
            },
            {
              "userUUID": "3cac1684-c1e0-47ae-92fd-6d7959759224",
              "userName": "test pupil",
              "userTypeId": 2,
              "genderId": 1
            }
          ]
        },
        "srvResMsg": "OK"
      }      
      ```

</details>

<details>
  <summary>教材を登録するエンドポイント</summary>

- **URL:** `/v1/auth/users/t_materials/register`
- **メソッド:** POST
- **説明:** 教師ユーザーが教科をもとに教材をクラスに登録する
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
    - `Content-Type`: multipart/form-data
  - ボディ: Form
    - Form Fields - 教材の情報
      - teachingMaterialName: リピート2
      - subjectId: 4
      - classUUID: 09eba495-fe09-4f54-a856-9bea9536b661
    - Files - 教材の画像
      - images: repeat_2.jpg

- **レスポンス:**
  - ステータスコード: 201 Created
    - ボディ:

      ```json
      {
        "srvResMsg":  "Created.",
        "srvResData": {
          "teachingMaterialUuid": "95af0199-3692-40af-b68f-a76e46cfad95"
        },
      }
      ```

</details>

### API仕様書てんぷれ

<details>
  <summary>＊○○＊するエンドポイント</summary>

- **URL:** `/＊エンドポイントパス＊`
- **メソッド:** ＊HTTPメソッド名＊
- **説明:** ＊○○＊
- **リクエスト:**
  - ヘッダー:
    - `＊HTTPヘッダー名＊`: ＊HTTPヘッダー値＊
  - ボディ:
    ＊さまざまな形式のボディ値＊

- **レスポンス:**
  - ステータスコード: ＊ステータスコード ステータスメッセージ＊
    - ボディ:
      ＊さまざまな形式のレスポンスデータ（基本はJSON）＊

      ```json
      {
        "srvResMsg":  "レスポンスステータスメッセージ",
        "srvResData": {
        
        },
      }
      ```

</details>

## エラー処理

APIがエラーを返す場合、詳細なエラーメッセージが含まれます。~~エラーに関する情報は[サーバーエラー]を参照してください。~~

## .ENV

.evnファイルの各項目と説明

- `./.env`: 実行時に必要だが、環境によって変わったり、リポジトリに含めたくない値

```env:.env
MYSQL_USER=DBに接続する際のログインユーザ名: juninry_user
MYSQL_PASSWORD=パスワード: juninry_pass
MYSQL_HOST=ログイン先のDBホスト名（dockerだとサービス名）: mysql-db-srv
MYSQL_PORT=ポート番号（dockerだとコンテナのポート）: 3306
MYSQL_DATABASE=使用するdatabase名: juninry_db
JWT_SECRET_KEY="openssl rand -base64 32"で作ったJWTトークン作成用のキー
JWT_TOKEN_LIFETIME=JWTトークンの有効期限: 315360000
MULTIPART_IMAGE_MAX_SIZE=Multipart/form-dataの画像の制限サイズ（10MBなら10485760）: 10485760
REQ_BODY_MAX_SIZE=リクエストボディのマックスサイズ（50MBなら52428800）: 52428800
```

## TODO

- 三層アーキテクチャなエンドポイントをDDDにリファクタリング（現状はmodel層として使っていたものがinfrastructure層外に置き去りにされている）

## 開発者

- Author:[]
- Mail:[]

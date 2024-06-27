# juninry

## 概要

juninryのGo APIサーバ。

### 環境

Visual Studio Code: 1.88.1  
golang.Go: v0.41.4  
image Golang: go version go1.22.2 linux/amd64  

## 環境構築

1. 以下のDocker環境を作成  
[リポジトリURL](https://github.com/unSerori/docker-juninry)  
SSH URL:  

    ```SSH:SSH URL
    git@github.com:unSerori/juninry-api.git
    ```

2. ここまでが1の内容（フォルダーをVScodeで開きgo_serverをVScodeアタッチ。）
3. リポジトリをクローン

    ```bash
    # カレントディレクトリにリポジトリの中身を展開
    git clone git@github.com:unSerori/juninry-api.git .
    
    # developブランチに移動
    git switch develop
    ```

4. shareディレクトリ内で以下のコマンド。

    ```bash:Build an environment
    # vscode 拡張機能を追加　vscode-ext.txtにはプロジェクトごとに必要なものを追記している。  
    cat vscode-ext.txt | while read line; do code --install-extension $line; done
    ```

5. .envファイルをもらうか作成。[.envファイルの説明](#env)

## API仕様書

エンドポイント、リクエストレスポンスの形式、その他情報のAPIの仕様書。

### エンドポインツ

#### ユーザを作成するエンドポイント

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

#### クラスの課題情報一覧を取得するエンドポイント

- **URL:** `/v1/auth/users/homework/upcoming`
- **メソッド:** GET
- **説明:** 自分が所属するクラスの期限が先のものを取得
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
        "srvResMsg":  "レスポンスステータスメッセージ",
        "srvResData": {
        
        },
      }
      ```

#### クラスのおてがみ情報一覧を取得するエンドポイント

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
          "notices": {
            "NoticeTitle": "【持ち物】おべんと",
            "NoticeDate": "2024-06-11T03:23:39Z",
            "UserName": "test teacher",
            "ClassName": "3-2 ふたば学級",
            "ReadStatus": 0
        }},
      }
      ```

#### おてがみの詳細情報を取得するエンドポイント

- **URL:** `/v1/auth/users/notice/{notice_uuid}`
- **メソッド:** GET
- **説明:** パスパラメーターで指定したおしらせの詳細情報を取得する
- **リクエスト:**
  - ヘッダー:
    - `＊HTTPヘッダー名＊`: ＊HTTPヘッダー値＊
  - ボディ:
    ＊さまざまな形式のボディ値＊

- **レスポンス:**
  - ステータスコード: ＊ステータスコード ステータス＊
    - ボディ:
      ＊さまざまな形式のレスポンスデータ（基本はJSON）＊

      ```json
      {
        "srvResMsg":  "レスポンスステータスメッセージ",
        "srvResData": {
        
        },
      }
      ```

#### ユーザー情報を取得するエンドポイント

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

#### クラスを新規登録するエンドポイント

- **URL:** `/v1/auth/users/classes/register`
- **メソッド:** POST
- **説明:** ＊クラスを新規作成し、招待コードを発行する。新規作成を行なったユーザーはクラスに所属する。＊
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
          "srvResMsg": "Created",
          "srvResData": {
            "classUUID": "19ea35a6-1e43-4cdd-bc2e-f6c790f0858e",
            "className": "クラスを立てる",
            "inviteCode": "1385",
            "validUntil": "2024-07-04T00:49:41.462371507Z"
          }
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

#### クラスの招待IDを更新するエンドポイント

- **URL:** `v1/auth/users/classes/refresh/{class_uuid}`
- **メソッド:** PUT
- **説明:** ＊クラスの招待IDを更新する＊
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

#### ログインするエンドポイント

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

#### 終了した宿題を提出するエンドポイント

- **URL:** `/v1/auth`
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

### API仕様書てんぷれ

#### ＊○○＊するエンドポイント

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

## エラー処理

APIがエラーを返す場合、詳細なエラーメッセージが含まれます。~~エラーに関する情報は[サーバーエラー]を参照してください。~~

## SERVER ERROR MESSAGE

サーバーレスポンスメッセージとして"srvResMsg"キーでメッセージを返す。  
サーバーレスポンスステータスコードと合わせてデバックする。

## .ENV

.evnファイルの各項目と説明

```env:.env
MYSQL_USER=DBに接続する際のログインユーザ名
MYSQL_PASSWORD=パスワード
MYSQL_HOST=ログイン先のDBホスト名。dockerだとサービス名。
MYSQL_PORT=ポート番号。dockerだとコンテナのポート。
MYSQL_DATABASE=使用するdatabase名
JWT_SECRET_KEY="openssl rand -base64 32"で作ったJWTトークン作成用のキー。
JWT_TOKEN_LIFETIME=JWTトークンの有効期限
```

## 開発者

- Author:[]
- Mail:[]

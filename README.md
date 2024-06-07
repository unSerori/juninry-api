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
        "srvResMsg":  "Successful user registration.",
        "srvResData": {
          "authenticationToken": "token@h",
        },
      }
      ```

#### クラスの課題情報一覧を取得するエンドポイント

- **URL:** `/v1/auth/class/homework/upcoming`
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
          "srvResData": [
            {
              "homework_limit": "0001-01-01T00:00:00Z",
              "homework_data": [
                {
                  "homework_uuid": "a3579e71-3be5-4b4d-a0df-1f05859a7104",
                  "start_page": 24,
                  "page_count": 2,
                  "homework_note": "がんばってくださ～い＾＾",
                  "teaching_material_name": "漢字ドリル3",
                  "subject_id": 1,
                  "subject_name": "国語",
                  "teaching_material_image_uuid": "a575f18c-d639-4b6d-ad57-a9d7a7f84575",
                  "class_name": "3-2 ふたば学級",
                  "submit_flag": 1
                },,,
              ]
            },,,
          ]
        }
      ```

#### クラスのおてがみ情報一覧を取得するエンドポイント

- **URL:** `/v1/auth/class/notice`
- **メソッド:** GET
- **説明:** 自分が所属するクラスのおてがみ情報一覧取得
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
        "srvResMsg":  "",
        "srvResData": {

        },
      }
      ```

#### おてがみの詳細情報を取得するエンドポイント

- **URL:** `/v1/auth/class/notice/{notice_uuid}`
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
        "srvResMsg":  "",
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
  - ステータスコード: ＊ステータスコード ステータス＊
    - ボディ:
      ＊さまざまな形式のレスポンスデータ（基本はJSON）＊

      ```json
      {
        "srvResMsg":  "",
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
TOKEN_LIFETIME=JWTトークンの有効期限
```

## 開発者

- Author:[]
- Mail:[]

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

- **URL:** `/api/v1/users/register`
- **メソッド:** POST
- **説明:** 新規ユーザを登録。
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json
  - ボディ:

    ```json
    {
        "name": "hogeta piyonaka",
        "type": "teacher", // pupil, parents  // 数値のフラグでもよいかも。
        "mailAddress": "hogeta@gmail.com",
        "password": "C@h",
    }
    ```

- **レスポンス:**
  - ステータスコード: 201 Created
    - ボディ:

      ```json
      {
        "srvResCode": 1001,
        "srvResMsg":  "Successful user registration.",
        "srvResData": {
          "authenticationToken": "token@h",
        },
      }
      ```

#### ユーザ情報の取得エンドポイント

- **URL:** `/api/v1/auth/users/user`
- **メソッド:** GET
~ - **説明:** トークンからidを取得、そのユーザの詳細情報を返す。
- **リクエスト:**
  - ヘッダー:
    - `Authorization`: (string) 認証トークン
- **レスポンス:**
  - ステータスコード: 200 OK
    - ボディ:

      ```json
      {
        "srvResCode": 1002,
        "srvResMsg":  "Successful acquisition of user information.",
        "srvResData": {
          "userInfo": {
            "name": "hogeta piyonaka",
            "type": "teacher", // pupil, parents  // 数値のフラグでもよいかも。
            "mailAddress": "hogeta@gmail.com",
            "classes": [
              "d025523f-bb80-44a5-4bdb-5c8628b4d080",
              "a5801d4d-e00f-37f2-fa67-1ab58534696b",
            ],
            "home": "bf0bcb96-4527-6b4f-9077-a32d69af316f",
          }
        },
      }
      ```

#### クラスを作成するコード

- **URL:** `/api/v1/auth/class`
- **メソッド:** POST
- **説明:** 新規クラスを登録。
- **リクエスト:**
  - ヘッダー:
    - `Content-Type`: application/json
    - `Authorization`: (string) 認証トークン
  - ボディ:

    ```json
    {
      "name": "IE2A",
    }
    ```

- **レスポンス:**
  - ステータスコード: 201 Created
    - ボディ:

      ```json
      {
        "srvResCode": 1003,
        "srvResMsg":  "Successful class registration.",
        "srvResData": {
          "classUUID": "bf7a1768-8458-4469-5047-48b072c27aa4",
          "entryCode": "7777",
        },
      }
      ```



・参加ID更新
・くらす参加
・くらすじょうほう（課題一覧）取得
・くらすじょうほう（おてがみ一覧）取得
・特定の課題情報を取得
・特定のおてがみ情報を取得
・課題を付与
・おてがみを付与
・課題完了。(せんせいに通知と課題の画像)
・提出された課題を取得

・おうち情報取得

## エラー処理

APIがエラーを返す場合、詳細なエラーメッセージが含まれます。エラーに関する情報は[サーバーエラーコード](#server-error-code)を参照してください。　　

## SERVER ERROR CODE

サーバーレスポンスコードとして"srvResCode"キーで数値を返す。  
以下に意味を羅列。  

- 成功関連
  - 1000: Successful authentication.
  - 1001: Successful user registration.
  - 1002: Successful acquisition of user information.
  - 1003: Successful class registration.

- エラー関連
  - 7000: Authentication unsuccessful.  

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

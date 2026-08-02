# iknow-tools

iKnow! の学習統計を定期的に取得し、グラフと差分レポートを公開するツールです。

## Docker Composeで実行する

1. 環境変数ファイルを作成し、公開プロフィールのユーザーIDを設定します。

   ```shell
   cp .env.example .env
   $EDITOR .env
   ```

2. GitHub Container Registryからイメージを取得して起動します。サーバー上でのビルドは不要です。

   ```shell
   docker compose pull
   docker compose up -d
   ```

3. `http://localhost:8080/` を開きます。

収集処理は従来のcrontab設定と同じく毎時59分に `cron.sh` を実行します。SQLiteデータベースと公開ファイルは `iknow-data` Docker volumeに保存され、nginxから読み取り専用で公開されます。

実行ログは次のコマンドで確認できます。

```shell
docker compose logs -f collector
```

ポートを変更する場合は `.env` の `HTTP_PORT` を変更してください。
使用するイメージのバージョンを固定する場合は、`IKNOW_IMAGE` をリリースタグ（例: `ghcr.io/blue-jam/iknow-tools:1.2.3`）へ変更してください。

### イメージの公開

`main` ブランチへのpush時に、GitHub Actionsがlinux/amd64・linux/arm64用イメージをビルドし、`ghcr.io/blue-jam/iknow-tools:latest` として公開します。`v1.2.3` のようなタグをpushすると、`1.2.3`、`1.2`、`1` タグも公開されます。

初回公開後、認証なしでサーバーからpullする場合はGitHubのPackage settingsでパッケージのvisibilityをPublicにしてください。Privateのまま使用する場合は、サーバーでpackagesのread権限を持つトークンを使って先にログインします。

```shell
echo "$GHCR_TOKEN" | docker login ghcr.io -u <github-user> --password-stdin
```

## ホスト上で直接実行する

`IKNOW_USER_ID` を環境変数に設定するか、第1引数にユーザーIDを指定します。

```shell
IKNOW_USER_ID=<user-id> ./cron.sh
```

crontabから毎時59分に実行する例:

```cron
59 * * * * cd <path-to-iknow-tools> && IKNOW_USER_ID=<user-id> ./cron.sh
```

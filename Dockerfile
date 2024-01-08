# リリース用のビルドを行うステージ

# debian:bullseye-slimをベースにしているGoイメージと言える
FROM golang:1.18.2-bullseye as deploy-builder

WORKDIR /app

# 依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 不要な情報は排除しながら、ビルド（成果物バイナリはapp）
RUN go build -trimpath -ldflags "-w -s" -o app

# --------------------------------------------

# デプロイ用のコンテナ

# バイナリを実行するだけなのでlinuxでいい（実行環境によく使われるslimを採用）
FROM debian:bullseye-slim as deploy

RUN apt-get update

# ビルドステージからバイナリだけを取得
COPY --from=deploy-builder /app/app .

CMD [ "./app" ]

# --------------------------------------------

# ローカル開発環境で使用する ホットリロード環境
FROM golang:1.18.2 as dev

WORKDIR /app

# ホットリロードするツール air を入れる
# ファイルを更新するたびにgo buildを実行して実行中のGoプログラムを再起動する
RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
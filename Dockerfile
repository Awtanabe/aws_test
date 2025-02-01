# 💡 `build` ステージ（コンパイル用）
FROM golang:1.23 AS build

WORKDIR /app

# 依存関係をキャッシュするため、go.modとgo.sumを先にコピー
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ソースコードをコピー
COPY . .

# 実行ファイルをビルド
RUN go build -v -o main main.go

# 💡 `base` ステージ（ランタイム環境）
FROM golang:1.23 AS base
WORKDIR /app

COPY --from=build /app/main /app/main

CMD ["/app/main"]
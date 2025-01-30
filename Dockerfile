FROM golang:1.23

WORKDIR /app

# パッケージマネージャー
COPY go.mod go.sum ./
## ビルド時のコマンドrun, コンテナの時は cmd
RUN go mod download && go mod verify
COPY . .

RUN go build -v -o app/main main.go

CMD ["./app/main"]
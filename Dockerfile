FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X 'github.com/SETTER2000/prove/internal/app.OSString=`go env GOOS`' -X 'github.com/SETTER2000/prove/internal/app.archString=`go env GOARCH`' -X 'github.com/SETTER2000/prove/internal/app.dateString=`date`' -X 'github.com/SETTER2000/prove/internal/app.versionString=`git describe --tags`' -X 'github.com/SETTER2000/prove/internal/app.commitString=`git rev-parse HEAD`'" -o /prove ./cmd/prove/main.go


FROM debian:buster-slim

WORKDIR /
COPY --from=builder /prove ./
COPY --from=builder ./app/config/ ./config/

CMD ["./prove"]
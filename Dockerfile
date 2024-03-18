FROM golang:latest

ENV TZ /usr/share/zoneinfo/Asia/Tokyo

WORKDIR /app

COPY /app/* ./

RUN go mod download

RUN go install github.com/cosmtrek/air@latest

EXPOSE 5050

CMD ["air", "-c", ".air.toml"]
FROM golang:latest

ENV TZ /usr/share/zoneinfo/Asia/Tokyo

COPY ./app /app

WORKDIR /app

RUN go mod download

RUN go install github.com/cosmtrek/air@latest

EXPOSE 5050

RUN go build -o main .

# CMD ["air", "-c", ".air.toml"]

CMD ["./main"]
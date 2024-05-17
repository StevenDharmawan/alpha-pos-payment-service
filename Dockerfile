FROM golang:1.22-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

COPY .env /app/.env

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

EXPOSE 8080

CMD [ "./main" ]
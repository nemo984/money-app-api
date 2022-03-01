FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /money-app-api

EXPOSE 8080

CMD [ "/money-app-api" ]
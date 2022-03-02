FROM golang:1.17-alpine AS builder
WORKDIR /go/src/github.com/nemo984/money-app-api/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM scratch
WORKDIR /src
COPY --from=builder /go/src/github.com/nemo984/money-app-api/app .
COPY images ./images
COPY docs ./docs
EXPOSE 8080
CMD [ "/src/app" ]
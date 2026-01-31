# Use latest Go image (>= 1.25.6)
FROM golang:1.25.6

WORKDIR /app

COPY . .

RUN go build -o feedback-app

EXPOSE 8080

CMD ["./feedback-app"]

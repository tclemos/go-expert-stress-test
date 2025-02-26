FROM golang:1.23.5
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o tb ./cmd/tb/
ENTRYPOINT ["./tb"]
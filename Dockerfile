FROM golang:1.18-alpine as build
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o shortner .
FROM alpine:3.16 
#RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/shortner .
CMD ["./shortner"]

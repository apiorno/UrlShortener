FROM golang:latest as build-env

RUN mkdir /urlshortener

WORKDIR /urlshortener

COPY service-acc-key.json .
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o go/bin/urlshortener

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=build-env urlshortener .

EXPOSE 8080
CMD ["./go/bin/urlshortener"]
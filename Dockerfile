FROM golang:1.21

WORKDIR /usr/src/app

# Pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Specify the main package to build
RUN go build -v -o /usr/local/bin/app 
CMD ["app"]



#docker build -t serialization . 

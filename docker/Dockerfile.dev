FROM golang:1.17 as development


WORKDIR /app

COPY . /app

RUN go mod download

# Install Reflex for development
RUN go install github.com/cespare/reflex@latest

EXPOSE 8004

CMD reflex -r '\.go$$' -s go run *.go
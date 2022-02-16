FROM golang:1.17 as builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 go build -o /server/myapp server/*.go

FROM alpine:3.15

WORKDIR /app

COPY --from=builder /server/myapp /server/myapp

COPY --from=builder app/access-private.pem access-private.pem
COPY --from=builder app/access-public.pem access-public.pem
COPY --from=builder app/refresh-private.pem refresh-private.pem
COPY --from=builder app/refresh-public.pem refresh-public.pem
COPY --from=builder app/swagger.json swagger.json

EXPOSE 8080

CMD ["/server/myapp"]

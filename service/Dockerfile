FROM golang:1.20-alpine AS build

WORKDIR /
COPY . /

RUN go build -o ./service/cmd/main ./service/cmd/main.go

FROM scratch

WORKDIR /app
COPY --from=build /service/cmd/main .

EXPOSE 9010
ENTRYPOINT ["/app/main"]


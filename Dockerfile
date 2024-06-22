FROM --platform=linux/amd64 golang:1.21 as build

WORKDIR /go/src/app

COPY go.mod go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM --platform=linux/amd64 gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/.env /

EXPOSE 8080
ENV PORT 8080

CMD ["/app"]
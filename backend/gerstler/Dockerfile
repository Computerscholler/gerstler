FROM golang:1.17-alpine AS build-backend

WORKDIR /app
RUN apk add g++
ADD go.mod go.sum ./
RUN go mod download
ADD server.go ./
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o gerstler ./cmd/main.go

FROM alpine

WORKDIR /app
COPY --from=build-backend /app/gerstler .
RUN chmod +x ./gerstler

COPY data_integration/wait/wait-for /
RUN apk add --no-cache ca-certificates
EXPOSE 5000

RUN mkdir /app/data
RUN mkdir /app/secrets
RUN mkdir /data/config

VOLUME [ "/app/data", "/app/secrets", "/app/config" ]

ENV GIN_MODE=release

CMD ["/wait-for", "elasticsearch:9200", "--", "/app/gerstler"]
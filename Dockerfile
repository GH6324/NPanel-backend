FROM golang:1.26-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=v1.0.9

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" \
	-o /out/npanel ./cmd/npanel

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates \
		netbase \
		&& rm -rf /var/lib/apt/lists/ \
		&& apt-get autoremove -y \
		&& apt-get autoclean -y

COPY --from=builder /out/npanel /app/npanel

WORKDIR /app

EXPOSE 8081 9012
VOLUME ["/data/conf"]

CMD ["/app/npanel", "-conf", "/data/conf"]

FROM golang:1.11-alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

# RUN apk add --no-cache git
WORKDIR /app/
COPY . .
RUN go build -a -installsuffix cgo -ldflags "-s -w" -mod=vendor -o hcm cmd/hcm/main.go

FROM gruebel/upx:latest as upx
COPY --from=builder /app/hcm /hcm.org
RUN upx --best --lzma -o /hcm /hcm.org

FROM daixijun1990/scratch
COPY --from=upx /hcm /hcm
ENTRYPOINT [ "/hcm" ]
CMD ["-c", "/config.yaml"]

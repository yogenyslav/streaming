FROM golang:1.22.2 as builder
WORKDIR /builder

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/server/main.go

FROM scratch as prod
WORKDIR /orchestrator

COPY config/config.yaml ./config/
COPY --from=builder /builder/main .

CMD [ "./main" ]

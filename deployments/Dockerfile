FROM golang:1.19.3-alpine AS builder
WORKDIR /app

RUN --mount=type=secret,id=git-credentials cp /run/secrets/git-credentials ~/.git-credentials

# todo: remove this after repositories are public (+ git-credentials)
RUN apk add git
RUN git config --global credential.helper store

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN GOOS=linux CGO_ENABLED=0 go build -o santos ./cmd/santos/santos.go

FROM scratch
WORKDIR /app

COPY --from=builder /app /app/

CMD [ "/app/santos" ]
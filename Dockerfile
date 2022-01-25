FROM golang:1.17 as builder

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./scripts ./scripts

COPY ./pkg ./pkg

RUN ./scripts/build.sh

FROM redhat/ubi8:latest

COPY --from=builder /app/build/bin/kaleido-api /
COPY ./tools/solc-static-linux /usr/local/bin/solc
COPY ./contracts ./contracts

CMD ["/kaleido-api"]

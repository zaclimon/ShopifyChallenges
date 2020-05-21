# Use a multi-stage build for not including unused dependencies
FROM golang:1.14 as build

WORKDIR /tmp/utsuru

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./UtsuruConcept UtsuruConcept

# Use Alpine for the smallest image possible
FROM alpine:3

COPY --from=build /tmp/utsuru/UtsuruConcept /app/

CMD ["/app/UtsuruConcept"]

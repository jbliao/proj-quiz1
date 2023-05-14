FROM golang:1.18 AS build

WORKDIR /proj-quiz1
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

FROM ubuntu:22.04 AS run
COPY --from=build /server /server
USER daemon:daemon

ENV QUIZ1_APP_PORT=8080 \
    QUIZ1_DB_USER=root \
    QUIZ1_DB_HOST=database \
    QUIZ1_DB_NAME=quiz1

ENTRYPOINT [ "/server" ]
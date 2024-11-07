FROM golang:1.23.2-alpine3.20
ARG maintainer="binarch"
LABEL "maintainer"=$maintainer
USER root
ENV APP=/app
WORKDIR $APP
COPY . .
RUN go build -o server cmd/server/main.go
CMD ["./server"]
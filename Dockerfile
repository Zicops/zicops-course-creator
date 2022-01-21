FROM golang:1.17 as builder
RUN apt-get update -qq

ARG GO_MODULES_TOKEN
ENV GO111MODULE=on
ENV GO_MODULES_TOKEN=$GO_MODULES_TOKEN
WORKDIR /go/src/app
RUN git config --global url."https://${GO_MODULES_TOKEN}:x-oauth-basic@github.com/zicops/".insteadOf "https://github.com/zicops/"
COPY go.mod .
COPY go.sum .
# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o zicops-course-creator .

FROM golang:latest
LABEL maintainer="Puneet Saraswat <puneet.saraswat10074@gmail.com>"

RUN apt-get update -y -qq

COPY --from=builder /go/src/app/zicops-course-creator /usr/bin/
EXPOSE 8090

ENTRYPOINT ["/usr/bin/zicops-course-creator"]

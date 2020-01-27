FROM golang:1.13-alpine AS build 
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk add git make openssl

WORKDIR /go/src/github.com/grubastik/kubernetes-admission-control
ADD . .
RUN go build -v -o addLabelWebhook app/main.go

FROM alpine
RUN apk --no-cache add ca-certificates && mkdir -p /app
WORKDIR /app
COPY --from=build /go/src/github.com/grubastik/kubernetes-admission-control/addLabelWebhook .
COPY --from=build /go/src/github.com/grubastik/kubernetes-admission-control/certs certs
CMD ["/app/addLabelWebhook"]
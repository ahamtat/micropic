# Build docker image for application
FROM alpine:latest
LABEL maintainer="artem.orlov@gmail.com"

RUN apk update && apk upgrade

# Set Moscow timezone
RUN apk add --update tzdata
ENV TZ=Europe/Moscow

# Clean APK cache
RUN rm -rf /var/cache/apk/*

ENV GOPATH /go
COPY --from=builder $GOPATH/src/github.com/AcroManiac/micropic/bin/cache /bin
EXPOSE 50051
ENTRYPOINT ["/bin/cache"]

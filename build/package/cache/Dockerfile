# Build docker image for application
FROM alpine:latest

RUN apk update && apk upgrade

# Set Moscow timezone
RUN apk add --update tzdata
ENV TZ=Europe/Moscow

# Clean APK cache
RUN rm -rf /var/cache/apk/*

COPY ./build/bin/cache /bin/cache
EXPOSE 50051
ENTRYPOINT ["/bin/cache"]

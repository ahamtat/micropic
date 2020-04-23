# Build docker image for application
FROM alpine:latest

RUN apk update && apk upgrade \
    && apk add ca-certificates && update-ca-certificates \
    && apk add openssl

# Add dokerize utility
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# Set Moscow timezone
RUN apk add --update tzdata
ENV TZ=Europe/Moscow

# Clean APK cache
RUN rm -rf /var/cache/apk/*

COPY ./build/bin/proxy /bin/proxy
CMD ["dockerize", "-wait", "tcp://rabbitmq:15672", "-timeout", "10s", "-wait-retry-interval", "5s", "/bin/proxy"]

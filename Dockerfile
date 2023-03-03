FROM golang:1.20-bullseye as builder

ENV CGO_ENABLED=0

RUN apt-get update -y
RUN apt-get install -y webp libwebp-dev

RUN mkdir /content-app
WORKDIR /content-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.16.2

#we need timezone database
RUN apk --no-cache add tzdata

RUN apk update && \
    apk upgrade -U && \
    apk add ca-certificates ffmpeg libwebp libwebp-tools libwebp-dev

COPY --from=builder /content-app/bin/content /
COPY --from=builder /content-app/docs/swagger.yaml /docs/swagger.yaml

COPY --from=builder /content-app/driver/web/authorization_model.conf /driver/web/authorization_model.conf
COPY --from=builder /content-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv

COPY --from=builder /etc/passwd /etc/passwd

#we need timezone database
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo 

ENTRYPOINT ["/content"]

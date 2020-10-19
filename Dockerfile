FROM alpine
MAINTAINER Park, Jinhong <jinhong0719@naver.com>

COPY ./club-service ./club-service
ENTRYPOINT [ "/club-service" ]

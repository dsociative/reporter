FROM alpine:3.4

RUN apk update
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN echo "Europe/Moscow" >  /etc/timezone

EXPOSE 9443 9092

ADD ./build/reporter /usr/local/bin

ENTRYPOINT reporter

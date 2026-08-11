FROM alpine:3.20
COPY agc /usr/local/bin/agc
ENTRYPOINT ["agc"]


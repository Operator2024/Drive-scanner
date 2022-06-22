FROM golang:1.17-alpine3.15
# LABEL maintainer="Operator2024 <work.pwnz+github@gmail.com>"
LABEL version="0.2.0-build.1"
ENV VER="0.2.0-build.1"
ENV TZ=Asia/Yekaterinburg
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \ 
  && apk add --no-cache smartmontools make && mkdir /workdir
COPY main.go /workdir
COPY Makefile /workdir
COPY go.mod /workdir
WORKDIR "/workdir"
RUN make -f Makefile build && chmod +x drive_scanner
CMD ["./drive_scanner"]
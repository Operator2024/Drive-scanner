FROM golang:1.17-alpine3.15
# LABEL maintainer="Operator2024 <work.pwnz+github@gmail.com>"
LABEL version="0.1.0-beta.3"
ENV VER="0.1.0-beta.3"
ENV TZ=Asia/Yekaterinburg
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \ 
  && apk add --no-cache smartmontools make && mkdir /workdir
COPY main.go /workdir
COPY Makefile /workdir
COPY go.mod /workdir
WORKDIR "/workdir"
RUN make -f Makefile build && chmod +x drive_scanner
CMD ["./drive_scanner", "-V"]
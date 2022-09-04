FROM golang:1.17-alpine3.16 as build
ENV DS_VER="0.2.1"
ENV TZ=Asia/Yekaterinburg
WORKDIR /go/src/git.selectel.org/belomestnykh.v/drive-scanner/
COPY /main.go /Makefile /go.mod  ./
RUN apk add --no-cache make && make -f Makefile build

FROM alpine:3.16
# LABEL maintainer="Operator2024 <work.pwnz+github@gmail.com>"
LABEL version="1.0.0"
ENV TZ=Asia/Yekaterinburg
WORKDIR "/workdir"
COPY --from=build /go/src/git.selectel.org/belomestnykh.v/drive-scanner/drive-scanner /workdir
COPY entrypoint.sh /workdir
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \ 
  && apk add --no-cache smartmontools jq && chmod +x drive-scanner && chmod +x entrypoint.sh
ENTRYPOINT [ "./entrypoint.sh" ]
CMD ["$1"]

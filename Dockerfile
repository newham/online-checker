FROM ubuntu:16.04
ADD pages/ public/ app.conf online-checker /online-checker/
EXPOSE 8082
WORKDIR /online-checker
ENTRYPOINT ["./online-checker"]
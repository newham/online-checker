FROM ubuntu:16.04
ADD pages /online-checker/pages
ADD public /online-checker/public
ADD app.conf /online-checker
ADD online-checker /online-checker
EXPOSE 8082
WORKDIR /online-checker
ENTRYPOINT ["./online-checker"]
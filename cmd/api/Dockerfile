FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get upgrade -y

ADD api /bin/api

CMD ["/bin/api"]

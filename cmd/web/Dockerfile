FROM --platform=$TARGETPLATFORM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get upgrade -y

ARG TARGETOS
ARG TARGETARCH
ADD web-${TARGETOS}-${TARGETARCH} /bin/web

CMD ["/bin/web"]

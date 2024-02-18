FROM --platform=$TARGETPLATFORM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get upgrade -y

ARG TARGETOS
ARG TARGETARCH
ADD crawler-${TARGETOS}-${TARGETARCH} /bin/crawler

CMD ["/bin/crawler"]

FROM --platform=$TARGETPLATFORM debian:stable-slim

ARG TARGETOS
ARG TARGETARCH
ADD testram-${TARGETOS}-${TARGETARCH} /bin/testram

CMD ["/bin/testram"]

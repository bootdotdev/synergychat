FROM --platform=$TARGETPLATFORM debian:stable-slim

ARG TARGETOS
ARG TARGETARCH
ADD testcpu-${TARGETOS}-${TARGETARCH} /bin/testcpu

CMD ["/bin/testcpu"]

FROM --platform=$TARGETPLATFORM scratch

ARG TARGETOS
ARG TARGETARCH
ADD testcpu-${TARGETOS}-${TARGETARCH} /bin/testcpu

CMD ["/bin/testcpu"]

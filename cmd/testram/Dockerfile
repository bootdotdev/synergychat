FROM --platform=$TARGETPLATFORM scratch

ARG TARGETOS
ARG TARGETARCH
ADD testram-${TARGETOS}-${TARGETARCH} /bin/testram

CMD ["/bin/testram"]

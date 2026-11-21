# Static, non-root (UID 65532) base with CA certificates; the binary is built by goreleaser
# with CGO_ENABLED=0. goreleaser's dockers_v2 lays binaries out by $TARGETPLATFORM.
FROM gcr.io/distroless/static-debian13:nonroot
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/argonap /usr/bin/argonap
ENTRYPOINT ["/usr/bin/argonap"]

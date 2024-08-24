FROM debian:bookworm-slim

ENTRYPOINT ["/bin/argonap"]
COPY argonap /bin

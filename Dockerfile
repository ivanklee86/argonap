FROM debian:bookworm-slim

ENTRYPOINT ["/bin/argonap"]
COPY /bin/argonap /bin

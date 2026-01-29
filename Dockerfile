FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
COPY bin/ghupdater-linux-amd64 /ghupdater
WORKDIR /out
ENTRYPOINT ["/ghupdater"]

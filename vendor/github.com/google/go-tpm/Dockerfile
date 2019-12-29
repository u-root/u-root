FROM golang:latest
# We need OpenSSL headers to build the simulator
RUN apt-get update && apt-get install -y \
    libssl-dev \
 && rm -rf /var/lib/apt/lists/*

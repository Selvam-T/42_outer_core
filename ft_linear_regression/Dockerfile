FROM golang:1.25

WORKDIR /usr/src/app

# Install useful tools
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        nano \
        curl \
        git \
        postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Expose app port
EXPOSE 8080

# Default: keep container alive and allow interactive use
CMD ["tail", "-f", "/dev/null"]


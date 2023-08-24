FROM golang:1.20 as base

ARG USER
ARG UID

# Create Container user
RUN useradd -ms /bin/bash -G www-data,root -u $UID $USER

FROM base as dev

# Add Air for live reloading with Go
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Install golang-migrate
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash && \
    apt-get update && \
    apt-get install -y migrate

WORKDIR /app

COPY ../../ .

USER $USER

CMD ["air"]

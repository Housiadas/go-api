FROM golang:1.20 as base

ARG USER=localdev
ARG UID

# Create Container user
RUN useradd -ms /bin/bash -G www-data,root -u $UID $USER
USER $USER

FROM base as dev

# Add Air for live reloading with Go
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /app
COPY ../../ .

USER $USER

CMD ["air"]
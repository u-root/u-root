FROM ubuntu:22.04

RUN apt-get update -y
RUN apt-get dist-upgrade -y

# create user "builder" with sudo privileges
ARG GID
ARG UID
ARG USER=builder
RUN groupadd --gid ${GID} $USER
RUN useradd --uid ${UID} --gid $USER --shell /bin/bash --home-dir /home/$USER --create-home $USER
RUN apt-get install -y sudo
RUN echo "$USER ALL=(ALL) NOPASSWD: ALL" | tee -a /etc/sudoers

ENV TZ=Europe/Rome
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get install -y cmake gcc gcc-arm-none-eabi gcc-mingw-w64 git gzip \
                       libsodium-dev libusb-1.0-0 make pkg-config protobuf-compiler \
                       u-boot-tools vim wget

# install tamago-go
ENV TAMAGO_VERSION="1.20.6"
ENV TAMAGO_CHECKSUM="6319b1778e93695b62bb63946c5dd28c4d8f3c1ac3c4bf28e49cb967d570dfd5"
RUN wget -O tamago-go.tgz https://github.com/usbarmory/tamago-go/releases/download/tamago-go${TAMAGO_VERSION}/tamago-go${TAMAGO_VERSION}.linux-amd64.tar.gz
RUN echo "${TAMAGO_CHECKSUM} tamago-go.tgz" | sha256sum --strict --check -
RUN tar -C / -xzf tamago-go.tgz && rm tamago-go.tgz

ENV GOPATH "/home/${USER}/go"
ENV USBARMORY_GIT "/home/${USER}/usbarmory"
ENV TAMAGO "/usr/local/tamago-go/bin/go"
ENV PATH "${PATH}:${GOPATH}/bin:/usr/local/tamago-go/bin"

USER $USER
WORKDIR /build

# for u-root circleci
RUN sudo mkdir -p /home/circleci
# IDs are still in flux.
RUN sudo chown $USER:$USER /home/circleci

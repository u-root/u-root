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
RUN echo "builder ALL=(ALL) NOPASSWD: ALL" | tee -a /etc/sudoers

ENV TZ=Europe/Rome
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get install -y cmake gcc gcc-arm-none-eabi gcc-mingw-w64 git gzip \
                       libsodium-dev libusb-1.0-0 make pkg-config protobuf-compiler \
                       u-boot-tools vim wget

# install tamago-go
ENV TAMAGO_VERSION="1.19.3"
ENV TAMAGO_CHECKSUM="b614e5266e0933c67f00d77987d49e7fd290704c86482f97c8bb2c2e26e46e9f"
RUN wget -O tamago-go.tgz https://github.com/usbarmory/tamago-go/releases/download/tamago-go${TAMAGO_VERSION}/tamago-go${TAMAGO_VERSION}.linux-amd64.tar.gz
RUN echo "${TAMAGO_CHECKSUM} tamago-go.tgz" | sha256sum --strict --check -
RUN tar -C / -xzf tamago-go.tgz && rm tamago-go.tgz

# install minisign
ENV MINISIGN_VERSION="0.10"
ENV MINISIGN_CHECKSUM="9fe40c2bd899a91f6b62a6ff3d469ece670f155307df50c2482ddd31337ab6da"
RUN wget -O minisign.tgz https://github.com/jedisct1/minisign/archive/refs/tags/${MINISIGN_VERSION}.tar.gz
RUN echo "${MINISIGN_CHECKSUM} minisign.tgz" | sha256sum --strict --check -
RUN tar -xzf minisign.tgz
RUN cd minisign-${MINISIGN_VERSION}; mkdir build; cd build; cmake ..; make; make install
RUN rm -r minisign*

# install habtool
ENV HABTOOL_VERSION="v2022.05.25"
RUN su - $USER -c "/usr/local/tamago-go/bin/go install github.com/usbarmory/crucible/cmd/habtool@${HABTOOL_VERSION}"

ENV GOPATH "/home/${USER}/go"
ENV USBARMORY_GIT "/home/${USER}/usbarmory"
ENV TAMAGO "/usr/local/tamago-go/bin/go"
ENV PATH "${PATH}:${GOPATH}/bin:/usr/local/tamago-go/bin"

USER $USER
WORKDIR /build

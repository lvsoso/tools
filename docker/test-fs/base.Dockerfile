FROM ubuntu:18.04

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    make \
    autoconf \
    automake \
    linux-headers-$(uname -r) \
    zlib1g-dev \
    vim \
    wget \
    gettext \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /code

CMD ["/bin/bash"]

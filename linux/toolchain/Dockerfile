FROM gcc:15.1-bookworm

RUN apt update && \
    apt install -y \
        autoconf \
        automake \
        autotools-dev \
        bc \
        bison \
        bzip2 \
        ccache \
        cmake \
        cpio \
        curl \
        device-tree-compiler \
        flex \
        gawk \
        git \
        gperf \
        libelf-dev \
        libexpat-dev \
        libglib2.0-dev \
        libgmp-dev \
        libmpc-dev \
        libmpfr-dev \
        libslirp-dev \
        libssl-dev \
        libtool \
        make \
        ninja-build \
        patchutils \
        python3 \
        python3-pip \
        python3-tomli \
        texinfo \
        wget \
        xz-utils \
        zlib1g-dev \
    ;


RUN groupadd -r user && \
    useradd -m -r -g user user

USER user
WORKDIR /home/user

RUN git clone --depth 1 --single-branch --branch 2025.07.16 \
    https://github.com/riscv/riscv-gnu-toolchain toolchain-source && \
    cd toolchain-source && \
    ./configure --prefix=/home/user/toolchain --enable-multilib && \
    make --jobs=$(nproc) linux && \
    cd .. && \
    rm -rf toolchain-source

USER root
RUN apt install -y gcc-riscv64-linux-gnu
USER user

ENV PATH="$PATH:/home/user/toolchain/bin"
WORKDIR /home/user/work
ENV CCACHE_DIR=/home/user/work/_ccache

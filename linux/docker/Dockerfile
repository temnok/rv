FROM debian:bookworm

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y \
    ccache gcc-riscv64-linux-gnu bzip2 libssl-dev make wget xz-utils \
    # required for creating initramfs:
    cpio \
    # required for kernel build:
    bc bison flex libelf-dev \
    # required for opensbi build:
    device-tree-compiler python3

RUN groupadd -r user && \
    useradd -m -r -g user user

USER user

WORKDIR /home/user/work

ENV CCACHE_DIR=/home/user/work/_ccache

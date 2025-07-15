FROM debian:bookworm

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y \
    device-tree-compiler gcc-riscv64-linux-gnu libssl-dev make wget xz-utils \
    # required for opensbi build:
    python3 \
    # required for kernel build:
    bc bison flex libelf-dev \
    # required for creating initramfs:
    cpio \
    # caching build artefacts for speed:
    ccache

RUN groupadd -r user && \
    useradd -m -r -g user user

USER user

WORKDIR /home/user/work

ENV CCACHE_DIR=/home/user/work/_ccache

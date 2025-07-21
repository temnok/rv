FROM gcc:15.1-bookworm

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y \
    autoconf automake autotools-dev curl python3 python3-pip python3-tomli \
    libmpc-dev libmpfr-dev libgmp-dev gawk bison flex texinfo \
    gperf libtool patchutils bc zlib1g-dev libexpat-dev ninja-build git cmake \
    libglib2.0-dev libslirp-dev

RUN groupadd -r user && \
    useradd -m -r -g user user

USER user
WORKDIR /home/user

RUN git clone --depth 1 --single-branch --branch 2025.07.16 \
    https://github.com/riscv/riscv-gnu-toolchain toolchain-source

WORKDIR /home/user/toolchain-source

RUN ./configure --prefix=/home/user/toolchain --enable-multilib && \
    make --jobs=$(nproc) || true

RUN make --jobs=$(nproc) || true

ENV CCACHE_DIR=/home/user/work/_ccache
ENV PATH="$PATH:/home/user/toolchain/bin"

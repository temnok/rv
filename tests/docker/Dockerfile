FROM gcc:14.3

RUN apt update && apt install -y \
    autoconf automake autotools-dev curl python3 python3-pip python3-tomli libmpc-dev \
    libmpfr-dev libgmp-dev gawk build-essential bison flex texinfo gperf libtool patchutils \
    bc zlib1g-dev libexpat-dev ninja-build git cmake libglib2.0-dev libslirp-dev

RUN groupadd -r user && \
    useradd -m -r -g user user

USER user

RUN cd ~ && \
    git clone --branch=2025.06.13 https://github.com/riscv/riscv-gnu-toolchain && \
    cd riscv-gnu-toolchain && \
    ./configure --prefix=/home/user/riscv && \
    make

ENV PATH="$PATH:/home/user/riscv/bin"

RUN cd ~ && \
    git clone https://github.com/riscv/riscv-tests && \
    cd riscv-tests && \
    git submodule update --init --recursive && \
    autoconf && \
    ./configure --prefix=/home/user/riscv-target && \
    make && \
    make install

WORKDIR /home/user

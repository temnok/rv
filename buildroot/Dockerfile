FROM gcc:14.3

RUN apt update && \
    apt install -y bc bison build-essential cpio dosfstools file flex git \
      libncurses5-dev libssl-dev mtools python3 python3-yaml rsync u-boot-tools unzip wget zip

RUN groupadd -r user && \
    useradd -m -r -g user user

USER user

RUN cd ~ && git clone \
    --single-branch --branch 2024.02.13 \
    https://github.com/buildroot/buildroot

WORKDIR /home/user/buildroot

COPY --chown=user external /home/user/external

RUN make BR2_EXTERNAL=~/external rv32_defconfig && make source
RUN make && make linux-rebuild opensbi-rebuild all
RUN mkdir ~/output && \
    gzip -n -c output/images/fw_payload.bin > ~/output/rv32.kernel.gz && \
    cp output/images/rv32.dtb ~/output/

RUN make BR2_EXTERNAL=~/external clean rv64_defconfig && make source
RUN make && make linux-rebuild opensbi-rebuild all
RUN gzip -n -c output/images/fw_payload.bin > ~/output/rv64.kernel.gz && \
    cp output/images/rv64.dtb ~/output/

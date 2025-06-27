FROM gcc:14.3

RUN apt update && \
    apt install -y bc bison build-essential cpio dosfstools file flex git \
      libncurses5-dev libssl-dev mtools python3 python3-yaml rsync u-boot-tools unzip wget zip

RUN groupadd -r user && \
    useradd -m -r -g user user

USER user

RUN cd ~ && git clone https://github.com/buildroot/buildroot

WORKDIR /home/user/buildroot

FROM gcc:14.3.0

RUN apt update && \
    apt install -y bc bison build-essential cpio dosfstools file flex git \
      libncurses5-dev libssl-dev mtools python3 python3-yaml rsync u-boot-tools unzip wget zip

RUN groupadd -r user && \
    useradd -m -r -g user user

USER user
WORKDIR /home/user

RUN curl -L https://github.com/buildroot/buildroot/archive/refs/tags/2023.11.1.tar.gz -o buildroot.tar.gz && \
  tar -xf buildroot.tar.gz && \
  rm buildroot.tar.gz && \
  mv buildroot*/ buildroot/

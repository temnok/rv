FROM gcc:15

RUN apt update && apt install -y cpio rsync bc

RUN cd /opt && \
  curl -L https://github.com/buildroot/buildroot/archive/refs/tags/2025.05.tar.gz -o buildroot.tar.gz && \
  tar -xf buildroot.tar.gz && \
  rm buildroot.tar.gz && \
  mv buildroot*/ buildroot/

ENV FORCE_UNSAFE_CONFIGURE=1
COPY extern/ /opt/extern/

WORKDIR /opt

#RUN make -C buildroot BR2_EXTERNAL=/opt/extern rv_defconfig && \
#    make -C buildroot

#RUN make -C buildroot linux-rebuild opensbi-rebuild all

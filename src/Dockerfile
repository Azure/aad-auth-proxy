FROM mcr.microsoft.com/oss/busybox/busybox:1.33.1 as builder

FROM mcr.microsoft.com/cbl-mariner/distroless/minimal:2.0

# Copy static shell into base image.
COPY --from=builder /bin/sh /bin/sh

# Copy mkdir executable
COPY --from=builder /bin/mkdir /bin/mkdir
COPY --from=builder /bin/tar /bin/tar


RUN mkdir /build
WORKDIR /build

# Copy the tar file to image and extract it
COPY ./aad-auth-proxy.tar .
RUN tar xvf aad-auth-proxy.tar

# execute command when docker launches / run
CMD ["/build/main"]
FROM golang:1.14 as gobuilder

WORKDIR /root

ENV GOOS=linux\
    GOARCH=amd64

COPY / /root/
COPY /common/* /root/out/
COPY /common/* /root/in/

COPY /out_gdb/* /root/out
COPY /in_gdb/* /root/in

RUN go mod edit -replace github.com/fluent/fluent-bit-go=github.com/fluent/fluent-bit-go@master
RUN go mod download & make -C out all
RUN go mod download & make -C in all

FROM fluent/fluent-bit:2.2

COPY --from=gobuilder /root/out/out_gdb.so /fluent-bit/bin/
COPY --from=gobuilder /root/in/in_gdb.so /fluent-bit/bin/
COPY --from=gobuilder /root/fluent-bit.conf /fluent-bit/etc/
COPY --from=gobuilder /root/plugins.conf /fluent-bit/etc/

EXPOSE 2020

CMD ["/fluent-bit/bin/fluent-bit", "--config", "/fluent-bit/etc/fluent-bit.conf"]

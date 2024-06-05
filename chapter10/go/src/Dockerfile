FROM golang:1.21-bullseye AS gobuilder

WORKDIR /root

ENV GOOS=linux\
    GOARCH=amd64\
    CGO_ENABLED=1

COPY / /root/
COPY /common/* /root/out/
COPY /common/* /root/in/

COPY /out_gdb/* /root/out
COPY /in_gdb/* /root/in

RUN go mod edit -replace github.com/fluent/fluent-bit-go=github.com/fluent/fluent-bit-go@master
# RUN go mod download & make -C out all
# RUN go mod download & make -C in all

RUN go mod tidy
RUN go build -C out -buildmode=c-shared  -a -gcflags=all="-C -l -B" -ldflags="-w -s" -trimpath  -tags netgo,osusergo -o /root/out/out_gdb.so
RUN go build -C in -buildmode=c-shared  -a -gcflags=all="-C -l -B" -ldflags="-w -s" -trimpath  -tags netgo,osusergo -o /root/in/in_gdb.so


FROM fluent/fluent-bit:3.0.6-debug

COPY --from=gobuilder /root/out/out_gdb.so /fluent-bit/bin/
COPY --from=gobuilder /root/in/in_gdb.so /fluent-bit/bin/
COPY --from=gobuilder /root/fluent-bit.conf /fluent-bit/etc/
COPY --from=gobuilder /root/plugins.conf /fluent-bit/etc/

EXPOSE 2020

#CMD ["ldd", "--version"]
#CMD ["ldd", "/fluent-bit/bin/out_gdb.so"]
CMD ["/fluent-bit/bin/fluent-bit", "--config", "/fluent-bit/etc/fluent-bit.conf"]

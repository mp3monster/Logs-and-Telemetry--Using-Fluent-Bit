docker build . -t fluent-bit-gdb -f Dockerfile
# docker builder prune -f
docker run -it --rm fluent-bit-gdb
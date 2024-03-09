docker build . -t fluent-bit-gdb -f Dockerfile
REM docker builder prune -f
docker run -it --rm fluent-bit-gdb
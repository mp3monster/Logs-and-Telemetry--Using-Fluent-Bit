all:
	go mod tidy
	go build -buildmode=c-shared -a -gcflags=all="-C -l -B" -ldflags="-w -s" -trimpath  -tags netgo,osusergo -o out_gdb.so

clean:
	rm -rf out_gdb.so *.h *~

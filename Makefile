build-lib-current:
	go build -o libbw.so -buildmode=c-shared ./cbindings
	patchelf --remove-rpath libbw.so

build-lib-386:
	GOOS=linux GOARCH=386 CGO_ENABLED=1 CC=$$CC_386 go build -buildmode=c-shared -o libbw.so ./cbindings
	patchelf --remove-rpath libbw.so

build-lib-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=$$CC_ARM64 go build -buildmode=c-shared -o libbw.so ./cbindings
	patchelf --remove-rpath libbw.so

build-lib-arm7:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=$$CC_ARMV7 go build -buildmode=c-shared -o libbw.so ./cbindings
	patchelf --remove-rpath libbw.so

release-lib-current: build-lib-current
	mkdir -p out/current-os
	mv libbw.so out/current-os/
	mv libbw.h out/current-os/
	cp cbindings/bw_*.h out/current-os/

release-lib-386: build-lib-386
	mkdir -p out/386
	mv libbw.so out/386/
	mv libbw.h out/386/
	cp cbindings/bw_*.h out/386/

release-lib-arm7: build-lib-arm7
	mkdir -p out/arm7
	mv libbw.so out/arm7/
	mv libbw.h out/arm7/
	cp cbindings/bw_*.h out/arm7/

release-lib-arm64: build-lib-arm64
	mkdir -p out/arm64
	mv libbw.so out/arm64/
	mv libbw.h out/arm64/
	cp cbindings/bw_*.h out/arm64/

release-all: release-lib-current release-lib-arm7 release-lib-386 release-lib-arm64
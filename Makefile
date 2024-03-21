BINARY_NAME = "jail"

init_rootfs:
	sh scripts/init.sh	

build:
	mkdir -p output rootfs_upper rootfs_work rootfs_merged
	sudo mount -t overlay overlay -o lowerdir=rootfs,upperdir=rootfs_upper,workdir=rootfs_work ./rootfs_merged
	cd src && go mod tidy -compat=1.19 && goimports -l -w . && gofmt -s -w . && go build -o ../output/${BINARY_NAME}

clean:
	sudo umount proc
	sudo umount rootfs_merged
	sudo rm -rf output rootfs_upper rootfs_work rootfs_merged

run:
	sudo ./output/${BINARY_NAME}

all:
	make build
	make run

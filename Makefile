CURRENT_DIR	= $(shell pwd)
OUTPUT = ${CURRENT_DIR}/output

# container
CONTAINER_BINARY_NAME	= container
CONTAINER_DIR = ${CURRENT_DIR}/src/container
ROOTFS_SCRIPT = ${CURRENT_DIR}/scripts/init.sh
ROOTFS = ${CURRENT_DIR}/rootfs
ROOTFS_UPPER = ${CURRENT_DIR}/rootfs_upper
ROOTFS_WORK = ${CURRENT_DIR}/rootfs_work
ROOTFS_MERGED = ${CURRENT_DIR}/rootfs_merged

# server 
SERVER_BINARY_NAME = server
SERVER_DIR = ${CURRENT_DIR}/src/server

rootfs:
	bash ${ROOTFS_SCRIPT}	

init: rootfs

build_init:
	mkdir -p ${OUTPUT}
	cp -r scripts ${OUTPUT}
	cp -r rootfs ${OUTPUT}

build_container:
	# mkdir -p ${ROOTFS_UPPER} ${ROOTFS_WORK} ${ROOTFS_MERGED}
	# sudo mount -t overlay overlay \
	# 	-o lowerdir=${ROOTFS},upperdir=${ROOTFS_UPPER},workdir=${ROOTFS_WORK} ${ROOTFS_MERGED}
	cd ${CONTAINER_DIR} && \
		go mod tidy -compat=1.19 && \
		goimports -l -w . && \
		gofmt -s -w . && \
		go build -o ${OUTPUT}/${CONTAINER_BINARY_NAME}
	echo "container compiled successfully"

build_server:
	cd ${SERVER_DIR} && \
		go mod tidy -compat=1.19 && \
		goimports -l -w . && \
		gofmt -s -w . && \
		go build -o ${OUTPUT}/${SERVER_BINARY_NAME}
	echo "server compiled successfully"

build: build_init build_container build_server

clean:
	sudo rm -rf ${OUTPUT}

run_container:
	sudo ${OUTPUT}/${CONTAINER_BINARY_NAME}

run_server:
	sudo ${OUTPUT}/${SERVER_BINARY_NAME}

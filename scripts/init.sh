#!/bin/sh -e

ROOT_NAME=rootfs
ATTACHMENT_DIR=attachments

mkdir ${ROOT_NAME}
mkdir -p ${ROOT_NAME}/bin ${ROOT_NAME}/sbin ${ROOT_NAME}/etc ${ROOT_NAME}/proc ${ROOT_NAME}/sys ${ROOT_NAME}/usr/bin ${ROOT_NAME}/usr/sbin

cp ${ATTACHMENT_DIR}/busybox ${ROOT_NAME}/bin/busybox
cp ${ATTACHMENT_DIR}/bash ${ROOT_NAME}/bin/bash
cp ${ATTACHMENT_DIR}/port-scanner ${ROOT_NAME}/bin/port-scanner

for cmd in $(${ROOT_NAME}/bin/busybox --list); do
  ln -s busybox ${ROOT_NAME}/bin/${cmd}
done

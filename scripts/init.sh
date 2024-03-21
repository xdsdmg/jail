#!/bin/sh -e

root_name=rootfs-tmp
attachment_dir=attachments

mkdir ${root_name}
mkdir -p ${root_name}/bin ${root_name}/sbin ${root_name}/etc ${root_name}/proc ${root_name}/sys ${root_name}/usr/bin ${root_name}/usr/sbin 

cp ${attachment_dir}/busybox ${root_name}/bin/busybox
cp ${attachment_dir}/bash ${root_name}/bin/bash 
cp ${attachment_dir}/port-scanner ${root_name}/bin/port-scanner

for cmd in $(${root_name}/bin/busybox --list); do
  ln -s busybox ${root_name}/bin/$cmd
done


#!/bin/bash

create_bridge() {
  local bridge=$1
  local bridge_addr=$2

  ip link add $bridge type bridge
  ip addr add dev $bridge $bridge_addr
  ip link set $bridge up
  # Ref: https://opengers.github.io/openstack/openstack-base-virtual-network-devices-bridge-and-vlan/
  iptables -A FORWARD -i $bridge -j ACCEPT
}

create_container_network() {
  local veth=$1
  local addr=$2
  local veth_peer=$3
  local ns=$4
  local bridge=$5
  local default_route=$6

  ip link add $veth type veth peer name $veth_peer
  ip link set dev $veth_peer master $bridge
  ip link set $veth_peer up
  ip netns add $ns
  ip link set $veth netns $ns
  ip -n $ns addr add dev $veth $addr
  ip -n $ns link set $veth up
  ip -n $ns link set lo up
  ip -n $ns route add default via $default_route
}

create_container_ufs() {
  local rootfs=$1
  local rootfs_upper=$2
  local rootfs_work=$3
  local rootfs_merged=$4

  mkdir -p $rootfs_upper $rootfs_work $rootfs_merged
  sudo mount -t overlay overlay \
    -o lowerdir=${rootfs},upperdir=${rootfs_upper},workdir=${rootfs_work} ${rootfs_merged}
}

remove_container_ufs() {
  local rootfs_upper=$1
  local rootfs_work=$2
  local rootfs_merged=$3

  umount proc
  umount ${rootfs_merged}
  rm -rf ${rootfs_upper} ${rootfs_work} ${rootfs_merged}
}

if [[ $1 == "create_bridge" ]]; then
  create_bridge $2 $3
elif [[ $1 == "create_container_network" ]]; then
  create_container_network $2 $3 $4 $5 $6 $7
elif [[ $1 == "create_container_ufs" ]]; then
  create_container_ufs $2 $3 $4 $5
elif [[ $1 == "remove_container_ufs" ]]; then
  remove_container_ufs $2 $3 $4 $5
else
  echo "invalid arg"
fi

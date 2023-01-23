#

##

### docker network

```shell
sudo lsns -t net
ip addr

docker run -it --rm --network=none ubuntu:net /bin/bash

docker container inspect --format='{{ json .NetworkSettings }}' xxxx | jq

docker run -it --rm --network=host ubuntu:net /bin/bash

docker container inspect --format='{{ json .NetworkSettings }}' xxxx | jq


docker run -it --rm --network=bridge ubuntu:net /bin/bash

docker container inspect --format='{{ json .NetworkSettings }}' 123 | jq


docker run --rm -it --network=container:123 ubuntu:net /bin/bash

docker container inspect --format='{{ json .NetworkSettings }}' xxxx | jq

```

### veth test


```shell

## add veth
ip link add <p1-name> type veth peer name <p2-name>


# create namespace
sudo ip netns add ns1
sudo ip netns add ns2
sudo ip netns list

# creat veth pair
sudo ip link add vethaaa type veth peer name vethbbb

# add to namespeac
sudo ip link set vethaaa netns ns1
sudo ip link set vethbbb netns ns2
sudo ip netns list

# host namespace
sudo ip link list


# set ip address
sudo ip netns exec ns1 ip addr add 172.18.0.2/24 dev vethaaa
sudo ip netns exec ns1 ip link set vethaaa up

sudo ip netns exec ns2 ip addr add 172.18.0.3/24 dev vethbbb
sudo ip netns exec ns2 ip link set vethbbb up


sudo ip netns exec ns1 ip addr
sudo ip netns exec ns2 ip addr

(base) lv@lv:network$ sudo ip netns exec ns1 ip addr
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
7959: vethaaa@if7958: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether 16:9b:7e:83:1c:87 brd ff:ff:ff:ff:ff:ff link-netnsid 1
    inet 172.18.0.2/24 scope global vethaaa
       valid_lft forever preferred_lft forever
    inet6 fe80::149b:7eff:fe83:1c87/64 scope link 
       valid_lft forever preferred_lft forever
(base) lv@lv:network$ sudo ip netns exec ns2 ip addr
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
7958: vethbbb@if7959: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether ae:50:f6:d7:41:bb brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.18.0.3/24 scope global vethbbb
       valid_lft forever preferred_lft forever
    inet6 fe80::ac50:f6ff:fed7:41bb/64 scope link 
       valid_lft forever preferred_lft forever

(base) lv@lv:network$ sudo ip netns exec ns1 ping 172.18.0.3 -c 2
PING 172.18.0.3 (172.18.0.3) 56(84) bytes of data.
64 bytes from 172.18.0.3: icmp_seq=1 ttl=64 time=0.030 ms
64 bytes from 172.18.0.3: icmp_seq=2 ttl=64 time=0.037 ms

--- 172.18.0.3 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1021ms
rtt min/avg/max/mdev = 0.030/0.033/0.037/0.006 ms


sudo ip netns delete ns1
sudo ip netns delete ns2
```

### bridge test

```shell
# create namespace
sudo ip netns add ns1
sudo ip netns add ns2
sudo ip netns list

# create bridge
sudo ip link add docker1 type bridge
sudo ip link set docker1 up
ip link list
brctl show

# greate veth
sudo ip link add veth0 type veth peer name veth0-br
sudo ip link add veth1 type veth peer name veth1-br

# veth0-br@veth0 <---->  veth0@veth0-br
# veth1-br@veth1 <--->  veth1@veth1-br

# add them to ns
sudo ip link set veth0 netns ns1
sudo ip link set veth0-br master docker1
sudo ip link set veth0-br up

sudo ip link set veth1 netns ns2
sudo ip link set veth1-br master docker1
sudo ip link set veth1-br up


# set ip address
sudo ip netns exec ns1 ip addr add 172.18.0.2/24 dev veth0
sudo ip netns exec ns1 ip link set veth0 up
sudo ip netns exec ns1 ip addr

sudo ip netns exec ns2 ip addr add 172.18.0.3/24 dev veth1
sudo ip netns exec ns2 ip link set veth1 up
sudo ip netns exec ns2 ip addr

# ping
(base) lv@lv:network$ sudo ip netns exec ns1 ping 172.18.0.3 -c 2
PING 172.18.0.3 (172.18.0.3) 56(84) bytes of data.
64 bytes from 172.18.0.3: icmp_seq=1 ttl=64 time=0.108 ms
64 bytes from 172.18.0.3: icmp_seq=2 ttl=64 time=0.110 ms

--- 172.18.0.3 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 0.108/0.109/0.110/0.001 ms


# tcpdump -i veth0 -nn
# sudo tcpdump -i ens5 icmp -nn

# check connect outside
sudo ip netns exec ns1 bash
(base) lv@lv:network$ sudo ip netns exec ns1 bash
(base) root@lv:network# ping 192.168.2.117
connect: 网络不可达

# check route
(base) lv@lv:network$ sudo ip netns exec ns1 bash
(base) root@lv:network# ip route show
172.18.0.0/24 dev veth0 proto kernel scope link src 172.18.0.2 

(base) root@lv:network# route -n
内核 IP 路由表
目标            网关            子网掩码        标志  跃点   引用  使用 接口
172.18.0.0      0.0.0.0         255.255.255.0   U     0      0        0 veth0

#  routing table
# add ip address to docker1
sudo ip addr add 172.18.0.1/24 dev docker1

(base) lv@lv:network$ ping 172.18.0.2 -c 2
PING 172.18.0.2 (172.18.0.2) 56(84) bytes of data.
64 bytes from 172.18.0.2: icmp_seq=1 ttl=64 time=0.066 ms
64 bytes from 172.18.0.2: icmp_seq=2 ttl=64 time=0.047 ms

--- 172.18.0.2 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1011ms
rtt min/avg/max/mdev = 0.047/0.056/0.066/0.012 ms


(base) lv@lv:network$ sudo ip netns exec ns1 bash
(base) root@lv:network# ip route add default via 172.18.0.1
(base) root@lv:network# route -n
内核 IP 路由表
目标            网关            子网掩码        标志  跃点   引用  使用 接口
0.0.0.0         172.18.0.1      0.0.0.0         UG    0      0        0 veth0
172.18.0.0      0.0.0.0         255.255.255.0   U     0      0        0 veth0

(base) root@lv:network# ping 192.168.2.117 -c 2
PING 192.168.2.117 (192.168.2.117) 56(84) bytes of data.
64 bytes from 192.168.2.117: icmp_seq=1 ttl=64 time=0.129 ms
64 bytes from 192.168.2.117: icmp_seq=2 ttl=64 time=0.039 ms

--- 192.168.2.117 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1032ms
rtt min/avg/max/mdev = 0.039/0.084/0.129/0.045 ms


### ?????

(base) root@lv:network# ping 114.114.114.114 -c 2
PING 114.114.114.114 (114.114.114.114) 56(84) bytes of data.

--- 114.114.114.114 ping statistics ---
2 packets transmitted, 0 received, 100% packet loss, time 1006ms

# sudo tcpdump -i docker1 -nn
# sudo tcpdump -i enp7s0 icmp  -nn
# sudo iptables -nvL --line-numbers

7        0     0 ACCEPT     all  --  *      docker0  0.0.0.0/0            0.0.0.0/0            ctstate RELATED,ESTABLISHED
8        0     0 DOCKER     all  --  *      docker0  0.0.0.0/0            0.0.0.0/0           
9      116  7068 ACCEPT     all  --  docker0 !docker0  0.0.0.0/0            0.0.0.0/0           
10       0     0 ACCEPT     all  --  docker0 docker0  0.0.0.0/0            0.0.0.0/0 


sudo iptables -t filter -A FORWARD -o docker1 -j ACCEPT
sudo iptables -t filter -A FORWARD -i docker1 ! -o docker1 -j ACCEPT
sudo iptables -t filter -A FORWARD -i docker1  -o docker1 -j ACCEPT

```
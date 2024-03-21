# Jail

## Ref

- https://coolshell.cn/articles/17010.html
- https://bytetech.info/articles/7160875563781980190

运用 Linux namespace 和联合文件系统（overlay）模拟一个简单的容器实例，也可算作进程 jail（关押）工具。

## 使用方法

需要在 Linux 环境下运行。

``` bash
# 运行
make all

# 运行后需要清理编译产物
make clean
```

## 容器网络需要满足哪些要求？

``` bash
sudo unshare --fork --pid --mount-proc bash
```


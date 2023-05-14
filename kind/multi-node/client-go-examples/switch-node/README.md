#### 运行
```shell
export NAMESPACE="default"
export TARGET_SVC_NAME="nginx"
export TARGET_STS_MASTER_NAME="nginx-master"
export TARGET_STS_BACKUP_NAME="nginx-backup"
export TARGET_STS_CLIENT_NAME="worker"
export MASTER_NODE_LABEL='{"app-master":"true"}'
export BACKUP_NODE_LABEL='{"app-backup":"true"}'
export MASTER_SVC_SELECTOR='{"app": "nginx-master"}'
export BACKUP_SVC_SELECTOR='{"app":"nginx-backup"}'

export MASTER_NODE_PVC="master-data"
export BACKUP_NODE_PVC="backup-data"

export LEASE_LOCK_NAME="test"
# 没有权限的情况下，设置为"false"
export LEASE_MODE="false"
export MASTER_NODE_PVC="master-data"
export BACKUP_NODE_PVC="backup-data"

go run app.go
```


#### 使用到的命令
```shell
kubectl taint nodes node1 key1=value1:NoSchedule
kubectl taint nodes node1 key1=value1:NoSchedule-

kubectl label nodes dev4-worker app-master=true
kubectl label nodes dev4-worker2 app-backup=true
kubectl get nodes --show-labels

kubectl create serviceaccount test

kubectl api-resources --no-headers --sort-by name -o wide | sed 's/.*\[//g' | tr -d "]" | tr " " "\n" | sort | uniq
create
delete
deletecollection
get
list
patch
update
watch


kubectl drain --ignore-daemonsets <node name>

kubectl uncordon <node name>

```


#### 构建镜像

```shell
GOOS=linux go build  -o ./app app.go

docker build -t test-in-cluster:v1 .

kind load docker-image test-in-cluster:v1 --name=dev4

k delete  -f monitor.yaml

k apply  -f monitor.yaml
```



GOOS=linux go build  -o ./app app.go && docker build -t test-in-cluster:v1 . && kind load docker-image test-in-cluster:v1 --name=dev4


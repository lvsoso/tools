```shell
GOOS=linux go build -o ./app .

docker build -t in-cluster:v1 .

kind load docker-image in-cluster:v1 --name=dev4

kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

kubectl run --rm -i in-cluster --image=in-cluster:v1

kubectl delete deployment in-cluster
```
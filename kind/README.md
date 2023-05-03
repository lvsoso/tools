#### install k8s
```shell
kind create cluster --config config.yaml
kubectl cluster-info --context kind-kind
```
#### test ha
```shell
# -- 和 replica 之间有空格
kubectl create deployment hello-world-flask --image=lyzhang1999/hello-world-flask:latest --replica=2
kubectl create ingress hello-world-flask --rule="/=hello-world-flask:5000"
kubectl create -f https://ghproxy.com/https://raw.githubusercontent.com/lyzhang1999/resource/main/ingress-nginx/ingress-nginx.yaml

while true; do sleep 1; curl http://127.0.0.1; echo -e '\n'$(date);done
```

#### test autoscale
```shell
kubectl apply -f https://ghproxy.com/https://raw.githubusercontent.com/lyzhang1999/resource/main/metrics/metrics.yaml

kubectl wait deployment -n kube-system metrics-server --for condition=Available=True --timeout=90s

# –cpu-percent 表示 CPU 使用率阈值，当 CPU 超过 50% 时将进行自动扩容，–min 代表最小的 Pod 副本数，–max 代表最大扩容的副本数
kubectl autoscale deployment hello-world-flask --cpu-percent=50 --min=2 --max=10

kubectl patch deployment hello-world-flask --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/resources", "value": {"requests": {"memory": "100Mi", "cpu": "100m"}}}]'
deployment.apps/hello-world-flask patched

kubectl get pod --field-selector=status.phase==Running

# -c 代表 50 个并发数，-n 代表一共请求 10000 次
ab -c 50 -n 10000 http://127.0.0.1:5000/

kubectl get pods --watch
```
#### update

```shell
kubectl set image deployment/hello-world-flask hello-world-flask=lyzhang1999/hello-world-flask:v1

kubectl apply -f new-hello-worlad-flask.yaml

kubectl edit deployment hello-world-flask
```

#### gitops

```shell
# fluxcd
kubectl apply -f https://ghproxy.com/https://raw.githubusercontent.com/lyzhang1999/resource/main/fluxcd/fluxcd.yaml

kubectl wait --for=condition=available --timeout=300s --all deployments -n flux-system

mkdir fluxcd-demo && cd fluxcd-demo

$ ls
deployment.yaml

git init
git add -A && git commit -m "Add deployment"
git branch -M main
git remote add origin git@github.com:lvsoso/fluxcd-demo.git
git push -u origin main


kubectl apply -f fluxcd-repo.yaml
kubectl get gitrepository
kubectl apply -f fluxcd-kustomize.yaml
kubectl get kustomization
```

更新部署

```shell
git add -A && git commit -m "Update image tag to v1"
git push origin main
kubectl describe kustomization hello-world-flask
```

```shell
git log
git reset --hard 75f39dc58101b2406d4aaacf276e4d7b2d429fc9
git push origin main -f
kubectl describe kustomization hello-world-flask
```
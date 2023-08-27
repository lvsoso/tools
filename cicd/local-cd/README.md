# local-cd

###  运行
```shell
python3 main.py 9999
```

### 请求
```shell
#echo "root:" | base64
# cm9vdDo=

curl -X POST  http://127.0.0.1:9999  \
-H 'Authorization: Basic cm9vdDo='  \
-d '{"image":"nginx", "tag":"1.17"}'

Recreating local-cd_nginx_1 ... done
HTTP/1.0 200 OK
Server: BaseHTTP/0.6 Python/3.7.4
Date: Sat, 26 Aug 2023 18:39:20 GMT
```


### systemd
编辑好 local-cd.service 文件后，执行
```shell
sudo cp local-cd.service /lib/systemd/system/l
sudo chmod 644 /lib/systemd/system/local-cd.service
sudo systemctl daemon-reload
sudo systemctl enable local-cd.service
```

其他

```shell
# start a service
sudo systemctl start local-cd.service
sudo systemctl status local-cd.service

# stop a service
sudo systemctl stop local-cd.service

# restart a service
sudo systemctl restart local-cd.service

# reload a service
sudo systemctl reload local-cd.service

# enable a service
sudo systemctl enable local-cd.service

# disable a service
sudo systemctl disable local-cd.service

# get the status log of a service
systemctl status local-cd.service
```

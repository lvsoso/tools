# 修改配置
修改 '/etc/hosts'
```shell
127.0.0.1 hi.test
127.0.0.1 hello.test
```

运行

```shell
docker-compose up -d

# test
curl http://hello.test:9999
curl http://hi.test:9999
```


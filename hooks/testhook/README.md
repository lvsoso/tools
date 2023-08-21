```shell
helm install demohook ./demohook/

(base) lv@lv:testhook$ kubectl describe pod/hook-preinstall | grep -E 'Anno|Started:|Finished:'
Annotations:  helm.sh/hook: pre-install
      Started:      Tue, 22 Aug 2023 00:17:27 +0800
      Finished:     Tue, 22 Aug 2023 00:17:42 +0800
(base) lv@lv:testhook$ kubectl describe pod/hook-postinstall | grep -E 'Anno|Started:|Finished:'
Annotations:  helm.sh/hook: post-install
      Started:      Tue, 22 Aug 2023 00:17:53 +0800
      Finished:     Tue, 22 Aug 2023 00:18:03 +0800
(base) lv@lv:testhook$ kubectl describe pod/demohook-fd9d4f498-rpp75 | grep -E 'Anno|Started:|Finished:'
Annotations:  <none>
      Started:      Tue, 22 Aug 2023 00:18:01 +0800
```

```
#  检测端口是否开放
echo -e '\x1dclose\x0d' | telnet baidu.com 80

# 超时退出
echo QUIT > quit.txt && timeout --signal=9 3 telnet baidu.com 444 < quit.txt

# 域名是否解析正常
ping -c3 baidu.com

# 检查是否可访问
curl https://baidu.com/no-exist
```
### TODO

- job's crud；[x]
- state query；[x]
- error handle；[x]
- stopped immediately；[x]
- conflict handle；[x]
- retry handle;[x]
- metrics;[x]

[Task-Retry](https://github.com/hibiken/asynq/wiki/Task-Retry)
[Life-of-a-Task](https://github.com/hibiken/asynq/wiki/Life-of-a-Task)
[Monitoring-and-Alerting](https://github.com/hibiken/asynq/wiki/Monitoring-and-Alerting)



```shell
GIT_LFS_SKIP_SMUDGE=1 GIT_TRACE=1 GIT_TRACE_PERFORMANCE=1 GIT_CURL_VERBOSE=1 git clone https://gitea.com/lvoooo/init.git

(base) lv@lv:init$ time sha256sum boot5.img 
b04c3bfdae341e01d42d8edc26bd76fc72a5b260eb8720749417d7baa806b28a  boot5.img

real	0m40.265s
user	0m39.360s
sys	0m0.900s
```
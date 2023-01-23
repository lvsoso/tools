docker run nginx:latest

export image_id=605c77e624dd

docker image inspect --format='{{json .RootFS}}' $image_id | jq
{
  "Type": "layers",
  "Layers": [
    "sha256:2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f",
    "sha256:e379e8aedd4d72bb4c529a4ca07a4e4d230b5a1d3f7a61bc80179e8f02421ad8",
    "sha256:b8d6e692a25e11b0d32c5c3dd544b71b1085ddc1fddad08e68cbd7fda7f70221",
    "sha256:f1db227348d0a5e0b99b15a096d930d1a69db7474a1847acbc31f05e4ef8df8c",
    "sha256:32ce5f6a5106cc637d09a98289782edf47c32cb082dc475dd47cbf19a4f866da",
    "sha256:d874fd2bc83bb3322b566df739681fbd2248c58d3369cb25908d68e7ed6040a6"
  ]
}


docker image inspect --format='{{json .GraphDriver}}' $image_id | jq
{
  "Data": {
    "LowerDir": "/var/lib/docker/overlay2/346fb089d49f771be0739d0fc22e30095e1af34a6fe10bb380ff8361bada2484/diff:/var/lib/docker/overlay2/f5d8294def4df17e21fdb0075c18f8231b424b0e7afb92b1521938ac0950718d/diff:/var/lib/docker/overlay2/caea4d8c54204e9328eeba6d0fcac8801c5937ca70d080aa91d36a7964498740/diff:/var/lib/docker/overlay2/afaad20467159a6b6ce868983531038d27593cfc1d2b36af342d65494b084137/diff:/var/lib/docker/overlay2/e47b02256cca90f5bb9443bcd9b8cf2b43e42139c885893903e5146bd3d79c67/diff",
    "MergedDir": "/var/lib/docker/overlay2/f82b5bfc31faf1d146783db4a19e4b3d97f1366c781562da243b2ca8e323f3aa/merged",
    "UpperDir": "/var/lib/docker/overlay2/f82b5bfc31faf1d146783db4a19e4b3d97f1366c781562da243b2ca8e323f3aa/diff",
    "WorkDir": "/var/lib/docker/overlay2/f82b5bfc31faf1d146783db4a19e4b3d97f1366c781562da243b2ca8e323f3aa/work"
  },
  "Name": "overlay2"
}

"LowerDir": 
"/var/lib/docker/overlay2/346fb089d49f771be0739d0fc22e30095e1af34a6fe10bb380ff8361bada2484/diff
:/var/lib/docker/overlay2/f5d8294def4df17e21fdb0075c18f8231b424b0e7afb92b1521938ac0950718d/diff
:/var/lib/docker/overlay2/caea4d8c54204e9328eeba6d0fcac8801c5937ca70d080aa91d36a7964498740/diff
:/var/lib/docker/overlay2/afaad20467159a6b6ce868983531038d27593cfc1d2b36af342d65494b084137/diff
:/var/lib/docker/overlay2/e47b02256cca90f5bb9443bcd9b8cf2b43e42139c885893903e5146bd3d79c67/diff",


(base) lv@lv:image$  docker run -it --name nginx-a nginx:latest /bin/bash
root@5931f21f9fd3:/# touch a.txt
root@5931f21f9fd3:/# echo "a" > a.txt

(base) lv@lv:image$ docker commit nginx-a mynginx:a
sha256:cbe43db4a038ef4f2bfca137bb2902432a81ebb3b79920f95cd1904fca5a35f9
(base) lv@lv:image$ docker image inspect --format='{{ json .RootFS}}' mynginx:a | jq
{
  "Type": "layers",
  "Layers": [
    "sha256:2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f",
    "sha256:e379e8aedd4d72bb4c529a4ca07a4e4d230b5a1d3f7a61bc80179e8f02421ad8",
    "sha256:b8d6e692a25e11b0d32c5c3dd544b71b1085ddc1fddad08e68cbd7fda7f70221",
    "sha256:f1db227348d0a5e0b99b15a096d930d1a69db7474a1847acbc31f05e4ef8df8c",
    "sha256:32ce5f6a5106cc637d09a98289782edf47c32cb082dc475dd47cbf19a4f866da",
    "sha256:d874fd2bc83bb3322b566df739681fbd2248c58d3369cb25908d68e7ed6040a6",
    "sha256:fec6f71667ce8933894ffcb781cff06a4e865b3cba5b106d66fbdee2d2b7a93d"
  ]
}

# 最后一层是多出来的

docker image inspect --format='{{json .GraphDriver}}' mynginx:a | jq
{
  "Data": {
    "LowerDir": "/var/lib/docker/overlay2/f82b5bfc31faf1d146783db4a19e4b3d97f1366c781562da243b2ca8e323f3aa/diff:/var/lib/docker/overlay2/346fb089d49f771be0739d0fc22e30095e1af34a6fe10bb380ff8361bada2484/diff:/var/lib/docker/overlay2/f5d8294def4df17e21fdb0075c18f8231b424b0e7afb92b1521938ac0950718d/diff:/var/lib/docker/overlay2/caea4d8c54204e9328eeba6d0fcac8801c5937ca70d080aa91d36a7964498740/diff:/var/lib/docker/overlay2/afaad20467159a6b6ce868983531038d27593cfc1d2b36af342d65494b084137/diff:/var/lib/docker/overlay2/e47b02256cca90f5bb9443bcd9b8cf2b43e42139c885893903e5146bd3d79c67/diff",
    "MergedDir": "/var/lib/docker/overlay2/30cc5ce8b117c73b963bd44525914fe963a75dc8f526284fc1f70031d2ee2a4b/merged",
    "UpperDir": "/var/lib/docker/overlay2/30cc5ce8b117c73b963bd44525914fe963a75dc8f526284fc1f70031d2ee2a4b/diff",
    "WorkDir": "/var/lib/docker/overlay2/30cc5ce8b117c73b963bd44525914fe963a75dc8f526284fc1f70031d2ee2a4b/work"
  },
  "Name": "overlay2"
}


"LowerDir": "/var/lib/docker/overlay2/f82b5bfc31faf1d146783db4a19e4b3d97f1366c781562da243b2ca8e323f3aa/diff
:/var/lib/docker/overlay2/346fb089d49f771be0739d0fc22e30095e1af34a6fe10bb380ff8361bada2484/diff
:/var/lib/docker/overlay2/f5d8294def4df17e21fdb0075c18f8231b424b0e7afb92b1521938ac0950718d/diff
:/var/lib/docker/overlay2/caea4d8c54204e9328eeba6d0fcac8801c5937ca70d080aa91d36a7964498740/diff
:/var/lib/docker/overlay2/afaad20467159a6b6ce868983531038d27593cfc1d2b36af342d65494b084137/diff
:/var/lib/docker/overlay2/e47b02256cca90f5bb9443bcd9b8cf2b43e42139c885893903e5146bd3d79c67/diff",

# 最上一层是多出来的


(base) lv@lv:image$ sudo ls /var/lib/docker/overlay2/30cc5ce8b117c73b963bd44525914fe963a75dc8f526284fc1f70031d2ee2a4b/diff
a.txt
(base) lv@lv:image$ sudo cat /var/lib/docker/overlay2/30cc5ce8b117c73b963bd44525914fe963a75dc8f526284fc1f70031d2ee2a4b/diff/a.txt
a


(base) lv@lv:overlay-test$ tree .
.
├── learn.sh
├── lower1
│   ├── a.txt
│   └── x.txt
├── lower2
│   ├── b.txt
│   └── x.txt
├── merged
├── upper
│   └── c.txt
└── work


sudo mount -t overlay overlay -o lowerdir=lower1:lower2,upperdir=upper,workdir=work merged


(base) lv@lv:overlay-test$ tree .
.
├── learn.sh
├── lower1
│   ├── a.txt
│   └── x.txt
├── lower2
│   ├── b.txt
│   └── x.txt
├── merged
│   ├── a.txt
│   ├── b.txt
│   ├── c.txt
│   └── x.txt
├── upper
│   └── c.txt
└── work
    └── work [error opening dir]


# 编辑和删除 merged 文件夹中文件，会相应的体现在upper文件夹中。
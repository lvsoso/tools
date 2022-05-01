

https://github.com/nshen/learn-neovim-lua

https://www.nerdfonts.com/font-downloads
https://github.com/wbthomason/packer.nvim
https://github.com/junegunn/vim-plug

```shell
# 安装 packer nvim 包管理器
git clone --depth 1 https://github.com/wbthomason/packer.nvim\
 ~/.local/share/nvim/site/pack/packer/start/packer.nvim
```

## 包管理插件
- :PackerCompile： 每次改变插件配置时，必须运行此命令或 PackerSync, 重新生成编译的加载文件
- :PackerClean ： 清除所有不用的插件
- :PackerInstall ： 清除，然后安装缺失的插件
- :PackerUpdate ： 清除，然后更新并安装插件
- :PackerSync : 执行 PackerUpdate 后，再执行 PackerCompile
- :PackerLoad : 立刻加载 opt 插件


```shell
# 查看标准数据目录
# Packer 会将插件默认安装在 标准数据目录/site/pack/packer/start
:h base-directories
~/.local/share/nvim/

```


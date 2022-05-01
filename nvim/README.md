
https://github.com/nshen/learn-neovim-lua

https://www.nerdfonts.com/font-downloads
https://github.com/wbthomason/packer.nvim
https://github.com/junegunn/vim-plug

```shell
# 安装 packer nvim 包管理器
git clone --depth 1 https://github.com/wbthomason/packer.nvim\
 ~/.local/share/nvim/site/pack/packer/start/packer.nvim

# 全局搜索 Telescope 依赖
sudo add-apt-repository ppa:x4121/ripgrep
sudo apt-get update
sudo apt install ripgrep

# fd-find
https://github.com/sharkdp/fd/releases

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

# show runtime dir
:echo $VIMRUNTIME
/usr/share/nvim/runtime/

```


## 主题
- nvim-treesitter 语法高亮 https://github.com/nvim-treesitter/nvim-treesitter
- lua插件 https://github.com/folke/tokyonight.nvim#plugin-support


## 插件扩展
https://github.com/nvim-telescope/telescope.nvim/wiki/Extensions
先在 plugins 里进行添加，再再对应插件配置里引入；


## 语法高亮
https://github.com/nvim-treesitter/nvim-treesitter#supported-languages

## LSP
```shell
:h lsp
```
- 安装 nvim-lspconfig
- 安装对应 language server
- 配置对应语言 require('lspconfig').xx.setup{…}
- :lua print(vim.inspect(vim.lsp.buf_get_clients())) 查看 LSP 连接状态

### 相关插件
- https://github.com/neovim/nvim-lspconfig
- https://github.com/williamboman/nvim-lsp-installer

### 相关操作
```shell
:LspInstallInfo
```
- 大写的 X 是卸载该 server
- u 是更新 server
- 大写 U 更新所有 servers
- c 检查 server 新版本
- 大写 C 检查所有 servers 的新版本
- ESC 关闭窗口
- ? 显示其他帮助信息

### 相关语言配置项
https://github.com/neovim/nvim-lspconfig/blob/master/doc/server_configurations.md

## 自动补全
- 补全引擎：提供代码补全核心插件
- 补全源：补全引擎数据来源
- snippet引擎： 自定义代码段引擎

nvim-cmp 是使用 Lua 编写的 补全引擎 插件。

可以配置多种外部的补全源，

支持 vsnip、luasnip、snippy、 ultisnips 4 种 snippet 引擎 。



## 代码格式化

### mhartington/formatter.nvim

- lua https://github.com/JohnnyMorganz/StyLua

```shell
cargo install stylua

# 配置完成后
:Format
```
### nul-ls.vim
- https://github.com/jose-elias-alvarez/null-ls.nvim/tree/main/lua/null-ls/builtins/formatting

# 前端代码高亮
```shell
:TSInstall css scss json html vue javascript typescript
```

# Rust 语法高亮与代码提示

### rust 相关安装
```shell
curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh
rustup component add rustfmt

```


```shell
:TSInstall rust
```

# Debug

想要在 Neovim 中打断点来调试应用程序，需要基于 DAP (Debug Adapter Protocol) 协议，这是一个跨平台的调试协议，可以在不同的平台上调试不同的程序。

在 Neovim 中基于 DAP 的插件有两个选择，一个是 vimspector ，另一个是 nvim-dap。

```shell
# 安装完成，安装 vimspector 支持的 debug 端适配器
cd ~/.local/share/nvim/site/pack/packer/start/vimspector/
./install_gadget.py --all
```

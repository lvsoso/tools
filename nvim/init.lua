-- 基础设置
require("basic")
require("keybindings")
-- Packer 插件管理
require("plugins")
-- 主题设置 （新增）
require("colorscheme")
-- 插件配置
-- 左边文件树
require("plugin-config.nvim-tree")
-- 窗口buffer
require("plugin-config.bufferline")
-- 底下栏
require("plugin-config.lualine")

-- 搜索插件
require("plugin-config.telescope")

-- 启动页面
require("plugin-config.dashboard")
require("plugin-config.project")
require("plugin-config.nvim-treesitter")

-- lsp server
require("lsp.setup")
require("lsp.cmp")
require("lsp.ui")
-- 对齐线
require("plugin-config.indent-blankline")
-- 格式化
-- require("lsp.formatter")
require("lsp.null-ls")

-- require("dap.vimspector") -- lua/dap/vimspector/init.lua
require("dap.nvim-dap")

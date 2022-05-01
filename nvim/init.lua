-- 基础设置
require('basic')
require('keybindings')
-- Packer 插件管理
require('plugins')
-- 主题设置 （新增）
require("colorscheme")
-- 插件配置
-- 左边文件树
require("plugin-config.nvim-tree")
-- 窗口buffer
require("plugin-config.bufferline")
-- 底下栏
require("plugin-config.lualine")

require("plugin-config.telescope")

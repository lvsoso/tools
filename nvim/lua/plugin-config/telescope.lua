local status, telescope = pcall(require, "telescope")
if not status then
  vim.notify("没有找到 telescope")
  return
end

telescope.setup({
  defaults = {
    -- 打开弹窗后进入的初始模式，默认为 insert，也可以是 normal
    initial_mode = "insert",
    -- 窗口内快捷键
    mappings = require("keybindings").telescopeList,
  },
  pickers = {
    -- 内置 pickers 配置
    find_files = {
      -- 查找文件换皮肤，支持的参数有： dropdown, cursor, ivy
      -- dropdown 垂直展示
      -- cursor 光标处展示
      -- ivy 全屏展示
      -- theme = "dropdown", 
    }
  },
  extensions = {
     -- 扩展插件配置
  },
})

-- 查看系统环境变量
pcall(telescope.load_extension, "env")

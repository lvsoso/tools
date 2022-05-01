local status, treesitter = pcall(require, "nvim-treesitter.configs")
if not status then
    vim.notify("没有找到 nvim-treesitter")
    return
end

treesitter.setup({
  -- 安装 language parser
  -- :TSInstallInfo 命令查看支持的语言
  -- all/maintained 一次下载所有的parsers
  -- ensure_installed 确认安装
  ensure_installed = { "markdown","python","go","json", "html", "css", "vim", "lua", "javascript", "typescript", "tsx" },
  -- 启用代码高亮模块
  highlight = {
    enable = true,
    -- 关闭vim 的正则语法高亮
    additional_vim_regex_highlighting = false, 
  },
})


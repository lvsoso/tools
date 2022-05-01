local colorscheme = "tokyonight"
-- pcall 在 Lua 中用于捕获错误
local status_ok, _ = pcall(vim.cmd, "colorscheme " .. colorscheme)
if not status_ok then
  vim.notify("colorscheme " .. colorscheme .. " 没有找到！")
  return
end



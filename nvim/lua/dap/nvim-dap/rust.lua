local extension_path =
"/home/lv/.local/share/nvim/site/pack/packer/start/vimspector/gadgets/linux/download/CodeLLDB/v1.6.10/root/extension/"
local codelldb_path = extension_path .. "adapter/codelldb"
local liblldb_path = extension_path .. "lldb/lib/liblldb.so"

return {
  adapter = require("rust-tools.dap").get_codelldb_adapter(codelldb_path, liblldb_path),
}

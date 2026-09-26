local lsp = "gopls"
local opts = { capabilities = require('cmp_nvim_lsp').default_capabilities() }

opts.settings = { gopls = { buildFlags = { "-tags=bare,netlify" } } }
vim.lsp.config(lsp, opts)

for _, client in ipairs(vim.lsp.get_clients({ name = "gopls" })) do
  client:stop({ force = true })
end

vim.lsp.enable(lsp)

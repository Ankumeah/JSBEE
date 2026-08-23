local lsp = "gopls"
local capabilities = require('cmp_nvim_lsp').default_capabilities()
local opts = { capabilities = capabilities }
opts.settings = { gopls = { buildFlags = { "-tags=sqlite" } } }
vim.lsp.config(lsp, opts)
for _, client in ipairs(vim.lsp.get_clients({ name = "gopls" })) do
  client:stop({ force = true })
end

local registry = require("mason-registry")
local function ensure_installed(name)
  local ok, pkg = pcall(registry.get_package, name)
  if not ok then
    vim.notify("Mason package not found: " .. name, vim.log.levels.WARN)
    return
  end
  if not pkg:is_installed() then
    pkg:install():once("closed", function()
      if pkg:is_installed() then
        vim.notify(name .. " installed successfully")
      else
        vim.notify(name .. " installation failed", vim.log.levels.ERROR)
      end
    end)
  end
end
ensure_installed("templ")

require('nvim-treesitter').install { "templ" }

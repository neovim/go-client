local chan

local function ensure_job()
  if chan then
    return chan
  end
  chan = vim.fn.jobstart({ 'helloremote' }, { rpc = true })
  return chan
end

vim.api.nvim_create_user_command('Hello', function(args)
  vim.fn.rpcrequest(ensure_job(), 'hello', args.fargs)
end, { nargs = '*' })


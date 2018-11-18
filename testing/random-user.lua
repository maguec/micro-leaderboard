-- wrk -c 1000 -d 10s -s random-user.lua http://localhost:8080


math.randomseed(os.time())


randomPath = function()
  local r = math.random(1,1000000)
  return string.format("/inc/testset/user-%s", r)
end

request = function()
  path = randomPath()
  return wrk.format(nil, path)
end

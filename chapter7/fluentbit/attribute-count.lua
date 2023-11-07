local function elementCounter(record)
  -- this function can be used recursively so we can count nested elements
  counter = 0
  for key, value in pairs(record) do
    if (type(value) == "table") then
      counter = counter + elementCounter(value)
    else
      counter = counter + 1
      print(string.format("-->[%d] %s --> %s", counter, key, value))
    end
  end
  return counter
end

function cb_addElementCount(tag, timestamp, record)
  -- we need to indicate back to Fluent Bit that the record will have changed, but not the timestamp
  local code = 2;

  if (record['remote_user'] ~= nil) then
    -- if the remote_user attribute exists let's change it to be Lua
    record['remote_user'] = "Lua"
  end

  -- add a new element with the count of elements in the structure passed
  record["element_count"] = (elementCounter(record) + 1)

  return code, timestamp, record
end

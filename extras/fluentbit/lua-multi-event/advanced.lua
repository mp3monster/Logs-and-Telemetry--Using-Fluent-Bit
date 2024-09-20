local function printDetails(record, indent)
  -- this function can be used recursively so we can count nested elements
  local counter = 0
  for key, value in pairs(record) do
    local elementType = type(value)
    if (elementType == "table") then
      print(string.format("table %s { %s = ", indent, key))
      printDetails(value, indent .. " ")
      print("}")
    else
      print(string.format("%s %s = %s --> %s", indent, key, tostring(value), elementType))
    end
  end
end

-- perform a deep copy if the received parameter is a table
function copy(obj)
  if type(obj) ~= 'table' then return obj end
  local res = {}
  for k, v in pairs(obj) do res[copy(k)] = copy(v) end
  return res
end

function cb_advanced(tag, timestamp, record)
  local code = 1
  print("Lua script - ", tag, "  ", timestamp, " record is a", type(record))
  -- printDetails(record, "")
  record1 = copy(record)
  record2 = copy(record)
  record2["BLAHH"] = "piiiinnnnnnnggggggggg"
  record1["remoteuser"] = "another user"
  newRecord = { record1, record2 }
  return code, timestamp, newRecord
end

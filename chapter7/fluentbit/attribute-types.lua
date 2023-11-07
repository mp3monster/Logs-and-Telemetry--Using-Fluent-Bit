local function printDetails(record, indent)
  -- this function can be used recursively so we can count nested elements
  local counter = 0
  for key, value in pairs(record) do
    local elementType = type(value)
    if (elementType == "table") then
      print(string.format("%s { %s = ", indent, key))
      printDetails(value, indent .. " ")
      print("}")
    else
      print(string.format("%s %s = %s --> %s", indent, key, tostring(value), elementType))
    end
  end
end


function cb_displayDataAndTypes(tag, timestamp, record)
  local code = 0
  if (type(timestamp) == "table") then
    print(tag, ":", timestamp['sec'], " . ", timestamp['nsec'])
  else
    print(tag, "  ", timestamp)
  end
  printDetails(record, "")
  return code, timestamp, record
end

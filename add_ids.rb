#!/usr/bin/env ruby

require 'json'

data = JSON.parse(File.read('BikeParking.geojson'))
data["features"].each_with_index do |f,i|
  f["properties"]["id"] = i
  data["features"][i] = f
end

puts JSON.dump(data)

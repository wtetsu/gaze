
fname = "test.rb.log"

puts "Append a line to #{fname}"

File.open(fname, "a") {|file|
  file.puts(Time.now)
}

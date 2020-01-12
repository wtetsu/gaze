import datetime

fname = "test.py.log"

print(f"Append a line to {fname}")

with open(fname, "a") as file:
    file.write(str(datetime.datetime.now()))
    file.write("\n")

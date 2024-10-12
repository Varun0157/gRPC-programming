# write N random numbers (floats) to a file

import random

N = int(1e6)
MAX_MAGNITUDE = 1000
filename = "data.txt"

with open(filename, "w") as f:

    def get_float(num_floating: int = 8):
        return round(random.uniform(0, MAX_MAGNITUDE) * 2 - MAX_MAGNITUDE, num_floating)

    for i in range(N):
        f.write(str(get_float()) + "\n")

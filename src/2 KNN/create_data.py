# write N random numbers (floats) to a file

import random

N = 1000
MAX_MAGNITUDE = 100
filename = "data.txt"

with open(filename, "w") as f:

    def get_float(num_floating: int = 4):
        return round(random.uniform(0, MAX_MAGNITUDE) * 2 - MAX_MAGNITUDE, num_floating)

    for i in range(N):
        f.write(str(get_float()) + "\n")

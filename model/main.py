import numpy as np

count = 5
nums = np.random.rand(10, 128)

for num in nums:
  print(num.tobytes(), sep=';')

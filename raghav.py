a = [1,4,1,2,3,5,2,3]
n = 7
k = 2
x = 28
cost = 0

import sys

i = 0
while i != n:
    potential_moves = []
    same = False
    for k in range(k,0,-1):
        if i+k == n:
            cost += x
            print(cost)
            sys.exit()
        if a[k] == a[i]:
            i = k
            same = True
            break

    if same == False:
        i = i + k
        cost += x

print(cost)
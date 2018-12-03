#!/usr/bin/python3

sum = 0
with open('input') as f:
    for line in f:
        sum += int(line)
        print(line + " " + str(int(line)))
        
print("result: ", sum)


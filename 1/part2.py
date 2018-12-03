#!/usr/bin/python3

with open('input') as f:
    changes = [int(change) for change in f]

seen = set()
i = 0
current = 0
while True:
    current += changes[i]
    if current in seen:
        print(current)
        exit()

    seen.add(current)
    
    i = i + 1
    if i >= len(changes):
        print('starting over')
        i = 0
       
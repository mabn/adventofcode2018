#!/usr/bin/python3


def has_repeated_letter(word, num_repeats):
    counts = dict()
    for letter in word:
        counts[letter] = counts.get(letter, 0) + 1
    
    print ("word: ", word, "counts: ", counts)
    return num_repeats in counts.values()
    
twos = 0
threes = 0
with open('input') as f:
    for line in f:
        if has_repeated_letter(line, 2):
            twos += 1
        if has_repeated_letter(line, 3):
            threes += 1

print("checksum:", twos*threes)
import sys
import os

length = int(raw_input())
height = int(raw_input())
text = raw_input()
chars = list("ABCDEFGHIJKLMNOPQRSTUVWXYZ?")
print >> sys.stderr, chars
answer = ""

for i in xrange(height):
    row = raw_input()
    for c in text:
        try:
            index = chars.index(c.upper())
        except ValueError:
            index = chars.index("?")
        answer += row[index * length:index * length + length]
    answer += os.linesep

print answer

import string
import sys
import math


# Don't let the machines win. You are humanity's last hope...

def display_node(char, x, y):
    # type: (string, int, int) -> string

    if char == ".":
        return "-1 -1"
    else:
        return str(x) + " " + str(y)


def get_next_node(line, start):
    # type: (string, int) -> int

    return line.index("0", start)


width = int(raw_input())  # the number of cells on the X axis
height = int(raw_input())  # the number of cells on the Y axis
m1 = []
m2 = []
for x2 in range(width):
    m2.append("")

for y in xrange(height):
    current_line = raw_input()
    m1.append(current_line)
    for x in xrange(width):
        m2[x] += current_line[x]

print >> sys.stderr, m1, m2

current = bottom = right = ""

for j in xrange(height):
    for i in xrange(width):

        if m1[j][i] == ".":
            continue

        current = "{} {}".format(i, j)
        try:
            # print >> sys.stderr, i, j, m1[j]
            right = "{} {}".format(get_next_node(m1[j], i + 1), j)
        except ValueError:
            right = "{} {}".format(-1, -1)

        try:
            print >> sys.stderr, i, j
            print >> sys.stderr, m2[i]
            bottom = "{} {}".format(i, get_next_node(m2[i], j + 1))
        except ValueError:
            bottom = "{} {}".format(-1, -1)

        print "{} {} {}".format(current, right, bottom)  # current_line = list(raw_input())

id:7
employee1
id:8
employee2
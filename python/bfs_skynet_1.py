import sys


def bfs_paths(graph, start, goal):
    queue = [(start, [start])]
    while queue:
        (vertex, path) = queue.pop(0)
        for possible in graph[vertex] - set(path):
            if possible == goal:
                yield path + [possible]
            else:
                queue.append((possible, path + [possible]))


def shortest_path(graph, start, goal):
    try:
        return next(bfs_paths(graph, start, goal))
    except StopIteration:
        return None


# n: the total number of nodes in the level, including the gateways
# l: the number of links
# e: the number of exit gateways
n, l, e = [int(i) for i in raw_input().split()]
graph = dict()


def add_link(graph, start, to):
    if start in graph:
        graph[start].add(to)
    else:
        graph[start] = {to}


def shortest_list(lists):
    return min(lists, key=lambda x: len(x))


for i in xrange(l):
    # n1: N1 and N2 defines a link between these nodes
    n1, n2 = [int(j) for j in raw_input().split()]
    add_link(graph, n1, n2)
    add_link(graph, n2, n1)
ei = []
for i in xrange(e):
    ei.append(int(raw_input()))  # the index of a gateway node

# game loop
while True:
    si = int(raw_input())  # The index of the node on which the Skynet agent is positioned this turn

    # Write an action using print
    # To debug: print >> sys.stderr, "Debug messages..."
    paths = []

    for gateway in ei:
        paths.append(shortest_path(graph, si, gateway))

    print >> sys.stderr, paths
    short_path = shortest_list(paths)

    print >> sys.stderr, short_path
    print " ".join(map(str, [short_path[0], short_path[1]]))

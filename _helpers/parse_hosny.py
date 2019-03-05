import re
import os


def distance(p1, p2):
    return round(((p2[0] - p1[0])**2 + (p2[0] - p1[0])**2)**0.5)


def parse(instance, instance_psa):
    with open(instance) as fd:
        i = 0
        points = {}
        demands = []
        tws = []
        cons = []
        node = 0
        nodes = []

        lines = fd.readlines()

        for i, _ in enumerate(lines):
            if i == 0:
                capacity, tasks = lines[i].split()
                tasks = int(tasks)
            elif i == 1:
                id1, x1, y1, demand, tw_from1, tw_to1 = re.findall(
                    r'\d+', lines[i])
                points[int(id1)] = (int(x1), int(y1))
                demands.append(demand)
                nodes.append(points[int(id1)])
                node += 1
            elif i < tasks+2:
                # 0	40	179	0	0	100000
                id1, x1, y1, demand, tw_from1, tw_to1 = re.findall(
                    r'\d+', lines[i])

                points[int(id1)] = (int(x1), int(y1))
                demands.append(demand)
                nodes.append(points[int(id1)])

                id2, x2, y2, demand2, tw_from2, tw_to2 = re.findall(
                    r'\d+', lines[i+tasks])

                points[int(id2)] = (int(x2), int(y2))

                nodes.append(points[int(id2)])

                cons.append("{} {} {} {} {} {} {}".format(
                    node, node+1, demand, tw_from1, tw_to1,
                    tw_from2, tw_to2))
                node += 2

    print(points)
    print(tws)

    numNodes = (int(tasks)*2)+1

    matrix = [None]*(numNodes)
    for i in range(numNodes):
        matrix[i] = []
        for j in range(numNodes):
            matrix[i].append(distance(nodes[i], nodes[j]))

    with open(instance_psa, "w") as fd:
        fd.write("{} {} {}\n".format(
            str(numNodes), str(capacity), "0"))
        for line in matrix:
            print(line)
            fd.write(' '.join(map(str, line)))
            fd.write('\n')
        fd.write('\n'.join(cons))


INSTANCES_PATH = "./test"

if __name__ == "__main__":
    for instance_name in os.listdir(INSTANCES_PATH):
        parse(os.path.join(INSTANCES_PATH, instance_name), os.path.join(
            INSTANCES_PATH, "{}.psa".format(instance_name)))

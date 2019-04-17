#!/usr/bin/env python3
# -*- coding: utf-8 -*-

##################################################
# This is helper script that transforms instances
# of PDPTW into PSA format
##################################################
# Author: Samuel Mitas
##################################################

import re
import os

INSTANCES_PATH = "./_instances/wan-rong-jih"


def distance(p1, p2):
    return round(((p2[0] - p1[0])**2 + (p2[0] - p1[0])**2)**0.5)


def parse(instance, instance_psa):
    with open(instance) as fd:
        i = 0
        points = {}
        cons = []
        node = 0
        nodes = []

        for line in fd:
            if line[0] != '#':
                i += 1

                if i == 1:
                    tasks, capacity = line.split()
                elif i > 2 and i < int(numNodes) + 2:
                    id, x, y = re.findall(r'\d+', line)
                    points[int(id)] = (int(x), int(y))
                elif i == int(numNodes) + 2:
                    nodes.append(points[int(line)])
                    start_node = 0
                    node += 1
                else:
                    _, pickup, delivery, tw_pickup_from, tw_pickup_to, tw_delivery_from, tw_delivery_to, demand = re.findall(
                        r'\d+', line)
                    nodes.append(points[int(pickup)])
                    nodes.append(points[int(delivery)])
                    cons.append("{} {} {} {} {} {} {}".format(
                        node, node+1, demand, tw_pickup_from, tw_pickup_to,
                        tw_delivery_from, tw_delivery_to))
                    node += 2
            elif "locations" in line:
                _, _, numNodes = line.split()
                numNodes = int(numNodes) + 1
                i += 1

        matrix = [None]*(numNodes-1)
        for i in range(numNodes-1):
            matrix[i] = []
            for j in range(numNodes-1):
                matrix[i].append(distance(nodes[i], nodes[j]))

    with open(instance_psa, "w") as fd:
        fd.write("{} {} {}\n".format(
            str(numNodes-1), str(capacity), start_node))
        for line in matrix:
            print(line)
            fd.write(' '.join(map(str, line)))
            fd.write('\n')
        fd.write('\n'.join(cons))


if __name__ == "__main__":
    for instance_name in os.listdir(INSTANCES_PATH):
        parse(os.path.join(INSTANCES_PATH, instance_name), os.path.join(
            INSTANCES_PATH, "{}.psa".format(instance_name)))

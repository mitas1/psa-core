import tkinter
import re

WIDTH = 600
HEIGHT = 600
PADDING = 50

##################################################
# TOTO: Clean this code and add some docs
##################################################


class Canvas:
    def __init__(self):
        self.canvas = tkinter.Canvas(width=WIDTH, height=HEIGHT)
        self.canvas.pack()
        self.points = []

    def load_instance_file(self):
        with open("_instances/ctsppdtw/test40") as fd:
            i = 0
            numNodes = 0
            for line in fd:
                if line[0] != '#':
                    i += 1

                    if i == 1:
                        tasks, capacity = line.split()
                    elif i == 2:
                        _, _, numNodes = line.split()
                        numNodes = int(numNodes)
                        print(numNodes)
                    elif i < int(numNodes) + 3:
                        id, x, y = re.findall(r'\d+', line)
                        self.points.append((int(x), int(y), id))

    def normalize_points(self):
        self.normalized = []
        max_x = 0
        max_y = 0
        for point in self.points:
            if point[0] > max_x:
                max_x = point[0]
            if point[1] > max_y:
                max_y = point[1]

        for point in self.points:
            self.normalized.append(
                ((point[0]/max_x) * (WIDTH-PADDING), (point[1]/max_x) * (WIDTH-PADDING), point[2]))

    def draw_points(self):
        for point in self.normalized:
            self.canvas.create_text(
                point[0], point[1]+20, text='{}'.format(point[2]))
            self.canvas.create_oval(
                point[0], point[1], point[0]+5, point[1]+5, fill='orange')

    def callback(self, e):
        pass

    def load_solution(self):
        with open("_solutions/test40.psa") as fd:
            self.solution = list(map(int, fd.read().split()))

    def draw_solution(self):
        for i in range(len(self.solution)-1):
            x, y, _ = self.normalized[self.solution[i]]
            x2, y2, _ = self.normalized[self.solution[i+1]]
            self.canvas.create_line(
                x, y, x2, y2
            )


c = Canvas()
c.load_instance_file()
c.normalize_points()
c.draw_points()
c.load_solution()
c.draw_solution()

tkinter.mainloop()

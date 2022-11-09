import numpy as np

Grade = []
Point = []

# 课程数量
ClassNum = int(input("共几门课:"))

for i in range(ClassNum):
    a = int(input("输入成绩:"))
    b = int(input("输入学分:"))
    Grade.append(a)
    Point.append(b)
print("输入成绩为", Grade)
print("输入学分为", Point)
Up = np.dot(Grade, Point)
Down = np.sum(Point)
Result = Up*0.6/Down
print("总得分是", Result)
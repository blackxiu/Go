from selenium import webdriver as wd
import time

edge = wd.Edge(executable_path="D:\Acode\python\mylearn\edgedriver_win32\msedgedriver.exe")
edge.get("http://www.baidu.com")
time.sleep(5)
edge.close()

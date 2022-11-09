import requests as req

url = 'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fimg.2qqtouxiang.com%2Fpic%2FTX9782_05.jpg&refer=http%3A' \
      '%2F%2Fimg.2qqtouxiang.com&app=2002&size=f9999,' \
      '10000&q=a80&n=0&g=0n&fmt=auto?sec=1653123794&t=4ada2b407f41cc65220c99b93a86a6f5 '

res = req.get(url)
print(res.content)

with open('2.png', 'wb') as f:
    f.write(res.content)



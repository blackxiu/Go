import requests as req

headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) ' 'Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44 '
}

params = {
    'wd': '长城',
    'c': 'b'
}

url = 'https://www.baidu.com/s?'

res = req.get(url, headers=headers, params=params)
# print(response.text)  # 得到的文本内容,默认的编码方式
# print(response.content)  # 得到的网页的二进制相应源码
print(res.content.decode('utf-8'))  # 得到的网页的二进制相应源码的解码
print(res.request.url)  # 得到响应的状态码  200是ok的意思
# print(response.request.headers)  # 会被网站识破是爬虫

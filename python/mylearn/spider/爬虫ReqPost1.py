import requests

url = 'https://fanyi.baidu.com/langdetect'
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) '
                  'Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44 '
}
data = {
    'query': '吃饭'
}
proxy = {
    "http": "HTTP://183.151.231.201:8080"
}

print(requests.post(url, headers=headers, data=data, proxies=proxy).text)

import requests as req
from lxml import etree

headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) '
                  'Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.50 '
}

url = 'https://movie.douban.com/'
text = req.get(url, headers=headers).text
dom = etree.HTML(text)
# ret = dom.xpath('//*[@id="screening"]/div[2]/ul/li[16]/ul/li[2]/a/text()')[0].strip()
ret = dom.xpath('//li[@class="ui-slide-item"]/ul/li[2]/a/text()')
for i in ret:
    print(i.strip())

import requests
from bs4 import BeautifulSoup

import requests

ip="220.201.84.233:9999"

proxies = {
  'http': 'http://' + ip,
  'https': 'https://' + ip,
}


def validate(proxies):
    https_url = 'https://ip.cn'
    http_url = 'http://ip111.cn/'
    headers = {'User-Agent': 'curl/7.29.0'}
    https_r = requests.get(https_url, headers=headers, proxies=proxies, timeout=10)
    http_r = requests.get(http_url, headers=headers, proxies=proxies, timeout=10)
    soup = BeautifulSoup(http_r.content, 'html.parser')
    result = soup.find(class_='card-body').get_text().strip().split('''\n''')[0]

    print(f"当前使用代理：{proxies.values()}")
    print(f"访问https网站使用代理：{https_r}")
    print(f"访问http网站使用代理：{result}")

validate(proxies)
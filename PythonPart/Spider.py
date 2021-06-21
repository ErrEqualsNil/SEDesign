import json
import random
import time
import requests
import Conns
import Config


class ProxyPool:
    def __init__(self):
        self.get_proxy_url = "http://dev.qydailiip.com/api/?apikey=abdb1bb36ae142ecefa286040d9f0625ed99840c&num=55&type=text&line=unix&proxy_type=putong&sort=3&model=all&protocol=https&address=%E4%B8%AD%E5%9B%BD&kill_address=&port=&kill_port=&today=false&abroad=1&isp=&anonymity=2"
        self.count = 50
        self.ips = []

    def get_random_proxy(self):
        if self.count % 50 == 0:
            resp = requests.get(self.get_proxy_url)
            self.ips = resp.text.split("\n")
            print("Get New Proxy Group: {}".format(self.ips))
            self.count = 1
            time.sleep(5)
        self.count += 1
        return self.ips[self.count]


class Spider:
    def __init__(self):
        self.conns = Conns.Conns()
        self.user_agent_list = [
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) "
            "Chrome/90.0.4430.212 Safari/537.36",
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 "
            "Safari/537.36 "
        ]
        # 参数说明：productId: 商品itemId ; score: 全部0/好评3/差评1/中评2 筛选 ；sortType：推荐排序5/时间排序6 ; page & pagesize(max=10) 位置
        self.CommentUrl = "https://club.jd.com/comment/productPageComments.action?productId={}&score={}&sortType={" \
                          "}&page={}&pageSize=10"

        self.categories = Config.categories
        self.use_proxy = Config.use_proxy
        self.proxyPool = ProxyPool()

    def run(self, task_id, item_id):
        requests.packages.urllib3.disable_warnings()
        comment_cnt = 0
        good_rate = 0
        for i in range(len(self.categories)):
            cnt, good_rate = self.get_comment(i, task_id, item_id)
            comment_cnt += cnt
        if not self.use_proxy:
            time.sleep(60)
        return comment_cnt, good_rate

    def get_comment(self, category_rank: int, task_id, item_id):
        comment_cnt = 0
        good_rate = 0
        category = self.categories[category_rank]
        if self.use_proxy:
            proxy = self.proxyPool.get_random_proxy()

        for page in range(category.page_offset, category.page_offset + category.count // 10):
            if self.use_proxy:
                if page % 10 == 0:
                    proxy = self.proxyPool.get_random_proxy()
                    time.sleep(0.5)
            else:
                if page % 2 == 0:
                    time.sleep(2)
            current_url = self.CommentUrl.format(item_id, category.score, category.sortType, page, 10)
            headers = {
                "User-Agent": random.choice(self.user_agent_list)
            }

            try:
                if self.use_proxy:
                    proxies = {
                        "http": "http://" + proxy,
                        "https:": "https://" + proxy
                    }
                    print("page: {} ; catching url: {}; proxy: {}".format(page, current_url, proxy))
                    resp = requests.get(current_url, headers=headers, proxies=proxies, verify=False, timeout=10)
                else:
                    print("catching url: {}".format(current_url))
                    resp = requests.get(current_url, headers=headers)
            except Exception as e:
                print("request err: {}".format(e))
                proxy = self.proxyPool.get_random_proxy()
                page -= 1
                continue

            try:
                data = json.loads(resp.text)
            except json.JSONDecodeError:
                print("json decode err, content: {}".format(resp.text))
                if self.use_proxy:
                    proxy = self.proxyPool.get_random_proxy()
                else:
                    print("Sleep 300s")
                    time.sleep(300)
                continue

            comments = data["comments"]
            for comment in comments:
                if comment["content"] == "此用户未填写评价内容":
                    continue
                self.conns.write_comment(comment["content"], comment["score"], task_id, comment["usefulVoteCount"])
                comment_cnt += 1

            good_rate = data["productCommentSummary"]["goodRateShow"]

        return comment_cnt, good_rate

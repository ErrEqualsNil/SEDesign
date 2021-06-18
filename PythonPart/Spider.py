import json
import random
import time
import requests
from PythonPart import Conns


class Category:
    def __init__(self, score, sortType, count, page_offset):
        self.score = score
        self.sortType = sortType
        self.count = count
        self.page_offset = page_offset


class Spider:
    def __init__(self):
        self.conns = Conns.Conns()
        self.headers = {
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) "
                          "Chrome/90.0.4430.212 Safari/537.36 "
        }
        # 参数说明：productId: 商品itemId ; score: 全部0/好评3/差评1/中评2 筛选 ；sortType：推荐排序5/时间排序6 ; page & pagesize(max=10) 位置
        self.CommentUrl = "https://club.jd.com/comment/productPageComments.action?productId={}&score={}&sortType={}&page={}&pageSize={}"
        self.config = {
            "pagesize": 10
        }
        self.categories = [
            Category(0, 5, 100, 0),
            Category(0, 5, 100, 10),
            Category(0, 5, 100, 20),
            Category(0, 5, 100, 30)
        ]

    def run(self, task_id, item_id):
        for i in range(len(self.categories)):
            self.get_comment(i, task_id, item_id)

    def get_comment(self, category_rank: int, task_id, item_id):
        category = self.categories[category_rank]
        for page in range(category.page_offset, category.page_offset + category.count // 10):
            if page % 2 == 0:
                time.sleep(random.randint(0, 2))
            print("catching: page:{}".format(page))
            try:
                currentUrl = self.CommentUrl.format(item_id, category.score, category.sortType, page, 10)
                resp = requests.get(currentUrl, headers=self.headers)
                data = json.loads(resp.text)
                comments = data["comments"]
                for comment in comments:
                    if comment["content"] == "此用户未填写评价内容":
                        continue
                    self.conns.write_comment(comment["content"], comment["score"], task_id, comment["usefulVoteCount"])
            except Exception as e:
                print("err when page {}; err: {}".format(page, e))
                continue

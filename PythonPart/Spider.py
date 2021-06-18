from Model import Task, Comment
import requests
from enum import Enum
import json
import time
import random
from pathos.pools import ProcessPool as Pool


class Category:
    def __init__(self, score, sortType, count, page_offset):
        self.score = score
        self.sortType = sortType
        self.count = count
        self.page_offset = page_offset


class Spider:
    def __init__(self):
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
        self.task_id = 0
        self.task_item_id=0


    def get_comment(self, category:Category):
        for page in range(category.page_offset, category.page_offset + category.count // 10):
            if page % 2 == 0:
                time.sleep(random.randint(0, 2))
            print("catching: category_score:{} ; page:{}".format(category.score, page))
            try:
                currentUrl = self.CommentUrl.format(self.task_item_id, category.score, category.sortType, page, 10)
                resp = requests.get(currentUrl, headers=self.headers)
                data = json.loads(resp.text)
                comments = data["comments"]
                for comment in comments:
                    if comment["content"] == "此用户未填写评价内容":
                        continue
                    _ = Comment(comment=comment["content"], score=comment["score"],
                                taskId=self.task_id, usefulVoteCount=comment["usefulVoteCount"])
            except Exception as e:
                print("err when category {}; page {}; err: {}".format(category, page, e))
                continue

    def run_spider(self, taskId, itemId):
        self.task_id = taskId
        self.task_item_id = itemId
        self.pool = Pool()
        self.pool.map(self.get_comment, self.categories)

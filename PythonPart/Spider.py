from Model import Task, Comment
import requests
from enum import Enum
import json
import time
import random


class Category:
    def __init__(self, score, sortType, count):
        self.score = score
        self.sortType = sortType
        self.count = count


class Spider:
    def __init__(self):
        self.headers = {
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) "
                          "Chrome/90.0.4430.212 Safari/537.36 "
        }
        # 参数说明：productId: 商品itemId ; score: 好评0/差评1/中评2 筛选 ；sortType：推荐排序5/时间排序6 ; page & pagesize(max=10) 位置
        self.CommentUrl = "https://club.jd.com/comment/productPageComments.action?productId={}&score={}&sortType={}&page={}&pageSize={}"
        self.config = {
            "pagesize": 10
        }

        self.categories = [
            Category(0, 5, 300),
            Category(1, 5, 20),
            Category(2, 5, 20)
        ]

    def get_comment(self, task: Task):
        for category in self.categories:
            for page in range(category.count // 10):
                if page % 10 == 0:
                    time.sleep(random.randint(0, 2))
                try:
                    currentUrl = self.CommentUrl.format(task.itemId, category.score, category.sortType, page, 10)
                    resp = requests.get(currentUrl, headers=self.headers)
                    data = json.loads(resp.text)
                    comments = data["comments"]
                    for comment in comments:
                        if comment["content"] == "此用户未填写评价内容":
                            continue
                        _ = Comment(comment=comment["content"], score=comment["score"],
                                    taskId=task.id, usefulVoteCount=comment["usefulVoteCount"])
                except:
                    print("err when category {}; page {}".format(category, page))
                    continue

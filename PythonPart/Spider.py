from Model import Task, Comment
import requests
from enum import Enum
import json

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
            Category(0, 5, 100),
            Category(0, 6, 100),
            Category(1, 5, 30),
            Category(2, 5, 30)
        ]

    def get_comment(self, task: Task):
        for category in self.categories:
            for page in range(category.count // 10):
                try:
                    currentUrl = self.CommentUrl.format(task.itemId, category.score, category.sortType, page, 10)
                    resp = requests.get(currentUrl, headers=self.headers)
                    data = json.loads(resp.text)
                    comments = data["comments"]
                    for comment in comments:
                        tmp = Comment(comment=comment["content"], score=comment["score"],
                                      taskId=task.id, usefulVoteCount=comment["usefulVoteCount"])
                except:
                    print("err when category {}; page {}".format(category, page))
                    continue








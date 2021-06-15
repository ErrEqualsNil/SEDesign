from Model import Task, Comment
import requests


class Spider:
    def __init__(self):
        self.headers = {
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) "
                          "Chrome/90.0.4430.212 Safari/537.36 "
        }

    def getUrl(self, task: Task):
        if len(task.url) != 0:
            return

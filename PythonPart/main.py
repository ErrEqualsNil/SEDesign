import json

import redis
import sqlobject
from Model import Task, Comment
import requests
import Spider


class Settings:
    def __init__(self):
        self.Ip = "121.36.31.113"
        self.RedisPort = 6379
        self.MySQLPort = 3306
        self.User = "pengpeng"
        self.Password = "hdeilm1718"
        self.DbName = "sedesign"


class Listener:
    def __init__(self):
        self.settings = Settings()
        self.RedisConn = redis.StrictRedis(host=self.settings.Ip, port=self.settings.RedisPort)
        mysqlUrl = "mysql://{}:{}@{}:{}/{}".format(self.settings.User, self.settings.Password, self.settings.Ip, self.settings.MySQLPort, self.settings.DbName)
        print(mysqlUrl)
        self.MySQLConn = sqlobject.connectionForURI(mysqlUrl)
        sqlobject.sqlhub.processConnection = self.MySQLConn
        self.spider = Spider.Spider()

    def get_message(self):
        print("Listener Listening!")
        while 1:
            message = self.RedisConn.brpop("Tasks", timeout=None)
            print("Get Message {}".format(message))
            task_id = int(message[1])
            task = Task.get(id=task_id)
            # TODO: Call services
            task.status = 3
            self.spider.run_spider(task.id, task.itemId)
            # todo: call analysis services
            task.status = 4



if __name__ == '__main__':
    l = Listener()
    l.get_message()
import redis
import sqlobject
from PythonPart import Analyzer
from Model import Task, Comment


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
        self.Analysiser = Analyzer

    def get_message(self):
        print("Listener Listening!")
        while 1:
            message = self.RedisConn.brpop("Tasks", timeout=None)
            print("Get Message {}".format(message))
            task_id = int(message[0])
            task = Task.get(id=task_id)
            # TODO: Call services


if __name__ == '__main__':
    listener = Listener()
    listener.get_message()

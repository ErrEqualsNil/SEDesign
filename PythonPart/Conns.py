import sqlobject
import redis
import Model


class Settings:
    def __init__(self):
        self.Ip = "121.36.31.113"
        self.RedisPort = 6379
        self.MySQLPort = 3306
        self.User = "pengpeng"
        self.Password = "hdeilm1718"
        self.DbName = "sedesign"


class Conns:
    def __init__(self):
        self.settings = Settings()
        mysql_url = "mysql://{}:{}@{}:{}/{}".format(self.settings.User, self.settings.Password, self.settings.Ip,
                                                   self.settings.MySQLPort, self.settings.DbName)
        # set mysql conn
        self.MySQLConn = sqlobject.connectionForURI(mysql_url)
        sqlobject.sqlhub.processConnection = self.MySQLConn

        # set redis conn
        self.RedisConn = redis.StrictRedis(host=self.settings.Ip, port=self.settings.RedisPort)

    def get_task(self) -> Model.Task:
        print("Waiting for Message")
        message = self.RedisConn.brpop("Tasks", timeout=None)
        print("Get Message {}".format(message))
        return Model.Task.get(id=int(message[1]))

    def write_comment(self, comment, score, task_id, usefulVoteCount):
        _ = Model.Comment(comment=comment, score=score, taskId=task_id,
                          usefulVoteCount=usefulVoteCount)
    def get_all_comment(self, task_id):
        return Model.Comment.filter(taskId=task_id).comment
    def write_hotword(self, task_id, hotword) -> Model.Task:
        Model.Task.filter(id=task_id).update(hotWords=hotword)
    def write_report(self, task_id, rep) -> Model.Task:
        Model.Task.filter(id=task_id).update(report=rep)

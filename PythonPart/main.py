import Conns
import Spider


class Main:
    def __init__(self):
        self.conn = Conns.Conns()
        self.spider = Spider.Spider()

    def Run(self):
        print("Start Spider")
        while 1:
            task = self.conn.get_task()

            # Spider Step
            task.status = 3  # update status
            comment_cnt, good_rate = self.spider.run(task.id, task.itemId)
            task.commentCount = comment_cnt
            task.goodRate = good_rate
            task.status = 4

            # todo: call analysis services



if __name__ == '__main__':
    l = Main()
    l.Run()
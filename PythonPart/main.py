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
            task.status = 3  # update status
            self.spider.run(task.id, task.itemId)
            # todo: call analysis services
            task.status = 4


if __name__ == '__main__':
    l = Main()
    l.Run()
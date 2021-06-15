import sqlobject


class Task(sqlobject.SQLObject):
    itemName = sqlobject.StringCol()
    status = sqlobject.IntCol()
    url = sqlobject.StringCol()
    report = sqlobject.StringCol()
    wordCloudUrl = sqlobject.StringCol()
    hotWords = sqlobject.StringCol()


class Comment(sqlobject.SQLObject):
    itemName = sqlobject.StringCol()
    comment = sqlobject.StringCol()
    taskId = sqlobject.IntCol()
    score = sqlobject.IntCol()

import sqlobject


class Task(sqlobject.SQLObject):
    itemName = sqlobject.StringCol()
    itemId = sqlobject.IntCol()
    goodRate = sqlobject.IntCol()
    commentCount = sqlobject.IntCol()
    status = sqlobject.IntCol()
    report = sqlobject.StringCol()
    wordCloudUrl = sqlobject.StringCol()
    hotWords = sqlobject.StringCol()


class Comment(sqlobject.SQLObject):
    comment = sqlobject.StringCol()
    score = sqlobject.IntCol()
    usefulVoteCount = sqlobject.IntCol()
    taskId = sqlobject.IntCol()


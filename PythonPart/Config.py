class Category:
    def __init__(self, score, sortType, count, page_offset):
        self.score = score
        self.sortType = sortType
        self.count = count
        self.page_offset = page_offset

use_proxy = False
categories = [
    Category(3, 5, 200, 0),
    Category(1, 5, 30, 0)
]

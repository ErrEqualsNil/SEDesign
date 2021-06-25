from Model import Task, Comment
import pandas as pd
from snownlp import SnowNLP
import jieba.posseg as pseg
import jieba
import jieba.analyse
from collections import Counter
import matplotlib.pyplot as plt
from wordcloud import WordCloud
import Model
import json
import oss2

class Analyzer:
    def __init__(self):
        self.auth = oss2.Auth('LTAI5t6DvAgcgXmiHf8tLAbx', 'SonzRIEFb5K3slYBhsrGC7uJlWlA3V')
        self.bucket = oss2.Bucket(self.auth, 'oss-cn-beijing.aliyuncs.com', 'sedesign')

    # 数据去重
    def quchong(self, comments):
        list_1 = set(comments)
        res = []
        for i in list_1:
            res.append(i)
        return res


    # 创建停用词list
    def stopwordslist(self, filepath):
        stopwords = [line.strip() for line in open(filepath, 'r', encoding='utf-8').readlines()]
        return stopwords

    # 对评论内容进行分词
    def seg_sentence(self, sentence):
        sentence_seged = pseg.cut(sentence.strip())
        stopwords = self.stopwordslist('stopwords.txt')  # 这里加载停用词的路径
        outstr = ''
        for word in sentence_seged:
            if word.word not in stopwords:
                if word.word != '\t' and word.flag == "a":
                    outstr += word.word
            outstr += " "
        return outstr
    # 获取评论
    def get_all_comment(self, task_id):
        comments = []
        selected_comments = Model.Comment.select(Model.Comment.q.taskId==task_id)
        for comment in selected_comments:
            comments.append(comment.comment)
        return comments

    # 返回json
    def jsonfile(self, data):
        n = 12
        L = sorted(data.items(), key=lambda item: item[1], reverse=True)
        L = L[:n]
        dictdata = {}
        for l in L:
            if l[0] == " " or l[0] == "\n":
                continue
            dictdata[l[0]] = l[1]
        return json.dumps(dictdata, ensure_ascii=False)

    def report(self, comments, flagname):
        outstr = ''
        for line in comments:
            sentence_seged = pseg.cut(line.strip())
            stopwords = self.stopwordslist('stopwords.txt')
            for word in sentence_seged:
                if word.word not in stopwords:
                    if word.word != '\t' and word.flag == flagname:
                        outstr += word.word
        Key = jieba.analyse.extract_tags(outstr, topK=10)
        dot = ','
        for i in Key:
            report0 = i + dot
        return report0

    def process(self, taskid):
        comments = self.get_all_comment(taskid)
        comments = self.quchong(comments)
        # 保存情感极性值小于等于0.1的结果为负面情感结果
        comments_neg = []
        for j in comments:
            s = SnowNLP(j)
            if s.sentiments <= 0.1:
                comments_neg.append(j + '\n')
        # 保存情感极性值大于0.1的结果为正面情感结果
        comments_pos = []
        for j in comments:
            s = SnowNLP(j)
            if s.sentiments > 0.1:
                comments_pos.append(j + '\n')
        # 正面评价分词
        comments_pos_split = ""
        for line in comments_pos:
            line_seg = self.seg_sentence(line)
            comments_pos_split += (line_seg + '\n')
        # 负面评价分词
        comments_neg_split = ""
        for line in comments_neg:
            line_seg = self.seg_sentence(line)
            comments_neg_split += (line_seg + '\n')
        # 词频统计
        # 正面评价
        data = jieba.cut(comments_pos_split)
        data = dict(Counter(data))
        data1 = self.jsonfile(data)

        comments_pos_count = ""
        for k, v in data.items():
            if k == " " or k == "\n":
                continue
            comments_pos_count += ('%s, %d\n' % (k, v))

        # 负面评价
        data = jieba.cut(comments_neg_split)
        data = dict(Counter(data))
        data2 = self.jsonfile(data)

        comments_neg_count = ""
        for k, v in data.items():
            if k == " " or k == "\n":
                continue
            comments_neg_count += ('%s, %d\n' % (k, v))
        # 生成词云

        wc = WordCloud(
            background_color="white",
            max_words=2000,
            height=800,
            font_path='simfang.ttf',
            width=1000,
            max_font_size=100,
            random_state=30, )
        wc.generate(comments_pos_count)
        wc.to_file('pos.png')
        self.bucket.put_object_from_file("{}_pos.png".format(taskid), 'pos.png')

        wc.generate(comments_neg_count)
        wc.to_file('neg.png')
        self.bucket.put_object_from_file("{}_neg.png".format(taskid), 'neg.png')

        task = Model.Task.get(id=taskid)
        task.wordCloudUrl = "[{}, {}]".format("https://sedesign.oss-cn-beijing.aliyuncs.com/{}_pos.png".format(taskid),
                                             "https://sedesign.oss-cn-beijing.aliyuncs.com/{}_neg.png".format(taskid))
        
        report = 'Advantage:' + self.report(comments_pos, "a") + 'Disadvantage:' + self.report(
            comments_neg, "a") \
                 + 'Focus:' + self.report(comments, "n")

        task.hotWords = "{},{}".format(data1, data2)
        task.report = report

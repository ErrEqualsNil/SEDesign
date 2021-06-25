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

class Analyzer:
    def __init__(self):
        pass

    # ����ȥ��
    def quchong(self, infile, outfile):
        infopen = open(infile, 'r', encoding='utf-8')
        outopen = open(outfile, 'w', encoding='utf-8')
        lines = infopen.readlines()
        list_1 = []
        for line in lines:
            if line not in list_1:
                list_1.append(line)
                outopen.write(line)
        infopen.close()
        outopen.close()

    # ����ͣ�ô�list
    def stopwordslist(self, filepath):
        stopwords = [line.strip() for line in open(filepath, 'r', encoding='utf-8').readlines()]
        return stopwords

    # ���������ݽ��зִ�
    def seg_sentence(self, sentence):
        sentence_seged = pseg.cut(sentence.strip())
        stopwords = Analyzer.stopwordslist('stopwords.txt')  # �������ͣ�ôʵ�·��
        outstr = ''
        for word in sentence_seged:
            if word.word not in stopwords:
                if word.word != '\t' and word.flag == "a":
                    outstr += word.word
            outstr += " "
        return outstr
    # ��ȡ����
    def get_all_comment(self, task_id):
        f = open('contents.txt', 'w', encoding='utf-8')
        f.write(Model.Comment.filter(taskId=task_id)[1]+'\n')
        f.close()

    # ����json
    def jsonfile(data):
        n = 10
        L = sorted(data.items(), key=lambda item: item[1], reverse=True)
        L = L[:n]
        dictdata = {}
        for l in L:
            dictdata[l[0]] = l[1]
        return json.dumps(dictdata, ensure_ascii=False)

    def process(self, taskid):
        # TODO: do process on item_name
        Analyzer.get_all_comment(taskid)
        Analyzer.quchong('contents.txt', 'contentsquchong.txt')
        data = open('contentsquchong.txt', 'r', encoding='utf-8')
        # ������м���ֵС�ڵ���0.1�Ľ��Ϊ������н��
        f = open('comments_neg.txt', 'w', encoding='utf-8')
        for j in data:
            s = SnowNLP(j)
            if s.sentiments <= 0.1:
                f.write(j + '\n')
        f.close()
        data.close()
        # ������м���ֵ����0.1�Ľ��Ϊ������н��
        data = open('contentsquchong.txt', 'r', encoding='utf-8')
        f = open('comments_pos.txt', 'w', encoding='utf-8')
        for j in data:
            s = SnowNLP(j)
            if s.sentiments > 0.1:
                f.write(j + '\n')
        f.close()
        data.close()
        # �������۷ִ�
        inputs = open('comments_pos.txt', 'r', encoding='utf-8')
        outputs = open('contentsfencipos.txt', 'w', encoding='utf-8')
        for line in inputs:
            line_seg = Analyzer.seg_sentence(line)
            outputs.write(line_seg + '\n')
        outputs.close()
        inputs.close()
        # �������۷ִ�
        inputs = open('comments_neg.txt', 'r', encoding='utf-8')
        outputs = open('contentsfencineg.txt', 'w', encoding='utf-8')
        for line in inputs:
            line_seg = Analyzer.seg_sentence(line)
            outputs.write(line_seg + '\n')
        outputs.close()
        inputs.close()
        # ��Ƶͳ��
        # ��������
        with open('contentsfencipos.txt', 'r', encoding='utf-8') as fr:
            data = jieba.cut(fr.read())
            data = dict(Counter(data))
            print(Analyzer.jsonfile(data))


        with open('contentscipinpos.txt', 'w', encoding='utf-8') as fw:
            for k, v in data.items():
                fw.write('%s, %d\n' % (k, v))
        # ��������
        with open('contentsfencineg.txt', 'r', encoding='utf-8') as fr:
            data = jieba.cut(fr.read())
            data = dict(Counter(data))
            print(Analyzer.jsonfile(data))

        with open('contentscipinneg.txt', 'w', encoding='utf-8') as fw:
            for k, v in data.items():
                fw.write('%s, %d\n' % (k, v))
        # ���ɴ���
        with open('contentsfencipos.txt', encoding='utf-8') as f:
            data = f.read()
            wc = WordCloud(
                background_color="black",
                max_words=2000,
                font_path='C:/Windows/Fonts/simfang.ttf',
                height=800,
                width=1000,
                max_font_size=100,
                random_state=30, )
            myword = wc.generate(data)
        plt.imshow(myword)
        plt.axis("off")
        plt.show()
        wc.to_file('pos.png')
        with open('contentsfencineg.txt', encoding='utf-8') as f:
            data = f.read()
            wc = WordCloud(
                background_color="black",
                max_words=2000,
                font_path='C:/Windows/Fonts/simfang.ttf',
                height=800,
                width=1000,
                max_font_size=100,
                random_state=30, )
            myword = wc.generate(data)
        plt.imshow(myword)
        plt.axis("off")
        plt.show()
        wc.to_file('neg.png')
        # �����û���ע��
        inputs = open('contentsquchong.txt', 'r', encoding='utf-8')
        outstr = ''
        for line in inputs:
            sentence_seged = pseg.cut(line.strip())
            stopwords = Analyzer.stopwordslist('stopwords.txt')
            for word in sentence_seged:
                if word.word not in stopwords:
                    if word.word != '\t' and word.flag == "n":
                        outstr += word.word
        Key = jieba.analyse.extract_tags(outstr, topK=10)
        print(Key)

    def report(self, pathname, flagname):
        inputs = open(pathname, 'r', encoding='utf-8')
        outstr = ''
        for line in inputs:
            sentence_seged = pseg.cut(line.strip())
            stopwords = Analyzer.stopwordslist('stopwords.txt')
            for word in sentence_seged:
                if word.word not in stopwords:
                    if word.word != '\t' and word.flag == flagname:
                        outstr += word.word
        Key = jieba.analyse.extract_tags(outstr, topK=10)
        dot = ','
        for i in Key:
            report0 = i + dot
        return report0

    def totalreport(self):
        report = '�û���Ϊ�����ƣ�' + Analyzer.report("comments_pos.txt", "a") + '�û���Ϊ�����ƣ�' + Analyzer.report(
            "comments_neg.txt", "a") \
                 + '�û���Ϊ��õķ��棺' + Analyzer.report("contentsquchong.txt", "n")
        return report

        inputs.close()

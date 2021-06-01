from nltk.corpus import words
word_list = words.words()
lenWords = len(word_list)

def getPercentages(i):
    dic = {}
    for word in word_list:
        x = word[i].lower()
        if x not in dic:
            dic[x] = 1
        else:
            dic[x] += 1
    for x in dic:
        dic[x] = round((dic[x] / lenWords)* 10000)
    return dic


print(getPercentages(0))
print(getPercentages(-1))

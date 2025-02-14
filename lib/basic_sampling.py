import string
from random import randrange, randint

class NameSampling:

    def __init__(self, filepath, content_length=5):
        self.context_length = content_length

        with open(filepath, "r") as r:
            self.names = ["." + f + "." for f in r.read().split()]

        self.letters = [l for l in string.ascii_lowercase]

        self.itos = {0: "."}
        self.stoi = {".": 0}

        for i, l in enumerate(self.letters):
            offset = i+1
            self.stoi[l] = offset
            self.itos[offset] = l

        self.names_length = len(self.names)

        self.train, self.dev, self.test = [], [], []

        for i in range(self.names_length):
            match randint(0, 10):
                case 0:
                    self.test.append(self.names[i])
                case 1:
                    self.dev.append(self.names[i])
                case _:
                    self.train.append(self.names[i])

    def sample_names(self, size=5, sample_set=None):
        if sample_set is None:
            sample_set = self.train
        sample_length = len(sample_set)
        batch_names = []
        for i in range(size):
            ni = randrange(sample_length-1)
            name = sample_set[ni]
            batch_names.append(name)
        return batch_names

    def word_contexts(self, word):
        samples = []
        max_length = len(word)
        fill = '.' * self.context_length
        for i in range(1,max_length):
            st = max(0, i-self.context_length)
            filled = fill[i:] + word[st:i]
            samples.append(filled[:self.context_length])
        return samples

    def get_xys(self, samples):
        xs, ys = [], []
        for s in samples:
            for ctx in self.word_contexts(s):
                x =  [self.stoi[c] for c in ctx]
                xs.append(x)
            y = [self.stoi[c] for c in s[1:]]
            ys += y

        return xs, ys
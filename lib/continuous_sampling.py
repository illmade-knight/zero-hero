import string
from random import randint, sample
from enum import Enum

import torch

class Split(Enum):
    TRAIN = 0
    DEV = 1
    TEST = 2

class WordSampling:

    def __init__(self):
        self.data = [[], [], []]
        self.letters = [l for l in string.ascii_lowercase]

        self.itos = {0: "."}
        self.stoi = {".": 0}

        for i, l in enumerate(self.letters):
            offset = i+1
            self.stoi[l] = offset
            self.itos[offset] = l

    def from_file(self, filepath):
        with open(filepath, "r") as r:
            self.words = [f for f in r.read().split()]

        self.__create_data()

    def from_words(self, words):
        self.words = words
        self.__create_data()

    def __create_data(self):
        self.words_length = len(self.words)

        self.splits = [[], [], []]

        for i in range(self.words_length):
            match randint(0, 10):
                case 0:
                    self.splits[1].append(self.words[i])
                case 1:
                    self.splits[2].append(self.words[i])
                case _:
                    self.splits[0].append(self.words[i])

        for i in range(3):
            d = []
            for word in self.splits[i]:
                #prepend each word with space (0)
                d += [0]
                for c in word:
                    d.append(self.stoi[c])

            self.data[i] = d

    def sample_batch(self, split, k, context_length):
        xs, ys = [], []

        sd = self.data[split.value]
        sr = len(sd) - context_length - 1

        for offset in sample(range(sr), k):
            xs += sd[offset:offset+context_length]
            ys.append(sd[offset+context_length])

        return torch.tensor(xs).view(k, context_length), torch.tensor(ys)
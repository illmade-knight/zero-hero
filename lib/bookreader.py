import re
import random
import torch

class BookReader:

    
    def __init__(self, lower=True, pattern = r'[^\d+a-zA-Z \n?!]', chunk_size = 25):
        self.pattern = pattern
        self.lower = lower
        self.chunk_size = chunk_size
        self.vocab_size = 0

    def read(self, *books):
        
        lines = []
        for book in books:
            with open(book, 'r', encoding='utf-8') as f:
                lines += f.readlines()

        self.all_lines = []
        self.train = []
        self.dev   = []
        self.test  = []

        self.chunk = []
        for line in lines:
            # we seemed to be getting excessive new lines in some books
            if line == '\n' and has_newline:
                has_newline = False
                continue
            else:
                has_newline = True
            line = re.sub(r'(\n)+', '\n', line)
            
            self.all_lines += line
            self.chunk += line

            if len(self.chunk) % self.chunk_size == 0:
                match random.randint(0, 10):
                    case 0:
                        self.test += self.chunk
                    case 1:
                        self.dev += self.chunk
                    case _:
                        self.train += self.chunk

                self.chunk = []

        self.text = re.sub(self.pattern, ' ', "".join(self.all_lines))
        self.train = re.sub(self.pattern, ' ', "".join(self.train))
        self.dev = re.sub(self.pattern, ' ', "".join(self.dev))
        self.test = re.sub(self.pattern, ' ', "".join(self.test))

        if self.lower is True:
            print("lowercase only")
            self.text = self.text.lower()
            self.train = self.train.lower()
            self.dev = self.dev.lower()
            self.test = self.test.lower()

        # here are all the unique characters that occur in this text
        self.chars = sorted(list(set(self.text)))
        
        self.vocab_size = len(self.chars)
        # create a mapping from characters to integers
        self.stoi = { ch:i for i,ch in enumerate(self.chars) }
        self.itos = { i:ch for i,ch in enumerate(self.chars) }
        self.encode = lambda s: [self.stoi[c] for c in s] # encoder: take a string, output a list of integers
        self.decode = lambda l: ''.join([self.itos[i] for i in l]) # decoder: take a list of integers, output a string

        self.data = [self.encode(self.train), self.encode(self.dev), self.encode(self.test)]

    def get_batch(self, data, context_length=1, batch_size=5):
        # generate a small batch of data of inputs x and targets y
        ix = torch.randint(len(data) - context_length, (batch_size,))
        b = torch.stack([data[i:i+context_length] for i in ix])
        return b
import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim

import collections

class SimpleBook(nn.Module):

    def __init__(self, embedding_size = 6, hidden_size = 100, context_length = 5, vocab_size = 29):
        super().__init__()

        self.embedding_size = embedding_size
        self.hidden_size = hidden_size
        self.context_length = context_length

        layer_dict = collections.OrderedDict([
            ('embed', nn.Embedding(vocab_size, embedding_size)),
            ('flatten', nn.Flatten(1)),
            ('feed_forward', nn.Linear(embedding_size*context_length, hidden_size, bias=True)),
            ('non_linearity', nn.Tanh()),
            ('logits', nn.Linear(hidden_size, vocab_size, bias=True)),
        ])

        nn.init.xavier_uniform_(layer_dict['feed_forward'].weight, gain=5/3)
        nn.init.zeros_(layer_dict['feed_forward'].bias)

        nn.init.zeros_(layer_dict['logits'].bias)

        self.model = nn.Sequential(
            layer_dict
        )
        self.layer_dict = layer_dict
        print(self.embedding_size, self.hidden_size, self.context_length, sum(p.nelement() for p in self.model.parameters())) # number of parameters in total


    def forward(self, idx, targets=None):
        logits = self.model(idx)

        if targets is None:
            loss = None
        else:
            loss = F.cross_entropy(logits, targets)

        return logits, loss

    def generate(self, max_characters, decoder):
        out = []
        ix = [[0 for _ in range(self.context_length)]]

        for _ in range(max_characters):
            logits = self.model(torch.tensor(ix))

            p = F.softmax(logits, dim=1)

            prediction = torch.multinomial(p, num_samples=1).item()

            del ix[0][0]
            ix[0].append(prediction)

            out.append(decoder.itos[prediction])

        return out

class TestRig:

    def __init__(self, embedding_size, hidden_size, context_length, vocab_size, learning_rate=None):
        if learning_rate is None:
            learning_rate = {'initial': .2, 'lambda': lambda epoch: 0.98}

        self.context_length = context_length
        model = SimpleBook(embedding_size, hidden_size, context_length, vocab_size)

        # parameters = model.parameters()
        # print(model.embedding_size, model.hidden_size, model.context_length, sum(p.nelement() for p in parameters)) # number of parameters in total
        print("initial learning rate", learning_rate['initial'])
        self.optimizer = optim.SGD(model.parameters(), lr=learning_rate['initial'], momentum=0.9)
        self.scheduler = optim.lr_scheduler.MultiplicativeLR(self.optimizer, lr_lambda=learning_rate['lambda'])

        self.model = model

        self.track = {'gradients': [], 'loss': [], 'learning_rate': []}

    def train(self, epochs, samples, batch_size, sampler):
        for ep in range(epochs):
            epoch_loss = 0
            for s in range(samples):
                x, y = sampler.sample_batch(batch_size, self.context_length)

                logits, loss = self.model.forward(x, y)

                epoch_loss += loss.detach()
                self.model.zero_grad(set_to_none=True)
                loss.backward()
                self.optimizer.step()

            self.scheduler.step()

            if ep % 10 == 0:
                self.track['loss'].append(epoch_loss)
                self.track['learning_rate'].append(self.scheduler.get_last_lr())
                print("learning rate", ep, self.scheduler.get_last_lr())
                print("running loss", epoch_loss/samples)

    @torch.no_grad()
    def val_split(self, x, y):

        sample_len = len(x) // 100

        val_loss = 0
        for s in range(sample_len):
            X = torch.tensor(x)
            Y = torch.tensor(y)

            logits, loss = self.model.forward(X, Y)

            val_loss += loss

        return val_loss / sample_len
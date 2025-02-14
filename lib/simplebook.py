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

    def forward(self, idx, targets=None):
        logits = self.model(idx)

        if targets is None:
            loss = None
        else:
            loss = F.cross_entropy(logits, targets)

        return logits, loss

    def generate(self, num_characters):
        print("not implemented")

class TestRig:

    def __init__(self, embedding_size, hidden_size, learning_rate = .2):
        self.model = SimpleBook(embedding_size, hidden_size)

        parameters = self.model.parameters()
        print(self.model.embedding_size, self.model.hidden_size, self.model.context_length, sum(p.nelement() for p in parameters)) # number of parameters in total

        self.optimizer = optim.SGD(self.model.parameters(), lr=learning_rate, momentum=0.9)
        self.scheduler = optim.lr_scheduler.MultiplicativeLR(self.optimizer, lr_lambda=lmbda)

        self.track = {'gradients': [], 'loss': [], 'learning_rate': []}

    def train(self, epochs, batch_size, samples, content_length):
        for ep in range(epochs):
            epoch_loss = 0
            for s in range(samples):
                x, y = names.get_xys(names.sample_names(batch_size, train))
                X = torch.tensor(x)
                Y = torch.tensor(y)

                logits, loss = self.model.forward(X, Y)

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
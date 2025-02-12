## Embeddings

We're going to add more dimensions to the character embedding space.

### 2D embedding

We'll introduce embeddings with a 2D embedding which we can later plot using pyplot

[2D embedding](embedding.ipynb)

### softmax - cross entropy

originally in our bigram model we showed the full process of finding the loss
logits -> counts -> normalized probabilities (counts/counts.sum) -> sum.log.mean

we'll compress that down first using torch's cross entropy and finally softmax

### result
unfortunately after all that we don't improve our loss very much 

to do that we'll have to start looking at predictions from [sequences](../sequences)




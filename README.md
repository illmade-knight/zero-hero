## zero-hero
working through Andrej Karpathy's neural-net/gpt starter 

["nn-zero-to-hero"](https://github.com/karpathy/nn-zero-to-hero)

### Description

The lectures build right from a bare bones neural network architecture towards a gpt like model

The jupyter notebooks are written follow-along style

i.e they are *not* lessons/lectures in themselves

### Towards gpt
the first step towards gpt is introducing a simple character level bigram model 
created from counts of one letter followed by another taken from a large list of names

the model then tries to 'generate' names from these counted occurrences

### Neural networking it

the bigram model is then re-stated as a neural-network 

this I think works really well: the neural-network 'learns' approximately the same weights as the 
original model - so our neural net does at least as well as the basic model

### Torch 

another really nice aspect is creating raw versions of torch neural network classes
nn.Linear etc 

this shows there no real magic going on - it's just weight tensors and dot products underneath


## What's this?

a resource for me chiefly, my own work through as I watch the lectures: this isn't linear.
* at times I'd watch some, go away, open up a jupyter notebook try and recreate what I thought I remembered
* go back, correct
* finish the lecture - go back, see what made sense now etc

I'd used tensorflow mostly before so this is also a good way to get into torch:

## The sections

### Torch Tensors
I've created some torch workout workbooks as well to get a better understanding of what's going on in torch
* when new tensors are created
* how to take slices without new tensor memeory
* when to reshape or resize
[tensors](torchbooks)

### Bigrams
Looking at bigram models for creating names given a training set

Work through the bigram lecture and review after
[bigrams](bigram)

### Embeddings
Previously we've used a character to int embedding 'a'=>1, 'b'->2 etc

now we use a higher dimensional space to represent the character and allow the model during training to adjust
the position of the characters in the space

[embeddings](embeddings)

### Sequences
Embeddings by themselves don't really improve our loss

to do that we start looking at more than one letter preceding the prediction

we also start to look at how the weights and bias initialize and look at what's 
happening inside our neurons as we train.

[sequences](sequences)

### A few things I might change

I'll comment more in the sections but briefly:

#### how zero is it? 
* It's certainly a really good review of material
* There's a lack of intuition, motivationy, insighty type discussion (Andrej has plenty of that elsewhere so it seems a little odd, though he's going through a lot here)

It's not spoon-fed certainly, there are plenty of references, papers you can look up etc

and of course there are the stanford online lectures if you need more depth

so to get the most out of it there's a fair bit of homework

#### uses torch but not the best introduction

the lectures start using torch and introduce new features at various points as they progress.

now there are plenty of introduction to torch tutorials but it might have been worthwhile to go over 
a block of things relevant to each section at once

#### notebooks are more scratch-books than explanations

cleanup the jupyter notebooks a bit, maybe use markdown cells every now and then to make things clearer

#### the names file
the names.txt file runs past its usefulness too early:

generating random names is a nice way to ease in but as the networks get deeper and more complex
their ability to generate 'new' names doesn't progress very far.

I'd introduce tiny-shakespeare earlier, the smaller networks start producing weird shakespeare-ish-ness
earlier than one might expect...



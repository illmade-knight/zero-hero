# zero-hero
working through Andrej Karpathy's neural-net/gpt starter 

["nn-zero-to-hero"](https://github.com/karpathy/nn-zero-to-hero )

### Description

The lectures build right from a bear-bones neural network architecture towards a gpt like model

The jupyter notebooks are written follow-along style

i.e they are *not* lessons/lectures in themselves

### Towards gpt
the first step towards gpt is introducing a simple character level bigram model 
created from counts of one letter followed by another taken from a large list of names

the model then tries to 'generate' names from these counted occurrences

### Neural networking it

the bigram model is then re-stated as a neural-network 

this I think works really well: the neural-network 'learns' the approximately same weights as the 
original model - so our neural net does at least as well as the basic model

### Torch 

another really nice aspect is creating raw versions of torch neural network classes
nn.Linear etc 

this shows there no real magic going on - it's just weight tensors and dot products underneath

### A few things I might change

the lectures start using torch and introduce new features at various points as they progress.

now there are plenty of introduction to torch tutorials but it might have been worthwhile to go over 
a block of things relevant to each section at once

cleanup the jupyter notebooks a bit, maybe use markdown cells every now and then to make things clearer

the names.txt file runs past its usefulness too early:

generating random names is a nice way to ease in but as the networks get deeper and more complex
their ability to generate 'new' names doesn't progress very far.

I'd introduce tiny-shakespeare earlier, the smaller networks start producing weird shakespeare-ish-ness
earlier than one might expect...

### What's this?

a resource for me chiefly, my own work through as I watch the lectures: this isn't linear.
* at times I'd watch some, go away, open up a jupyter notebook try and recreate what I thought I remembered
* go back, correct
* finish the lecture - go back, see what made sense now etc

I'd used tensorflow mostly before so this is also a good way to get into torch:

I've created some torch workout workbooks as well to get a better understanding of what's going on in torch
* when new tensors are created
* how to take slices without new tensor memeory
* when to reshape or resize



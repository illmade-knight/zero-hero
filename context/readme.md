## Structure and context

Moving away from a list of names, 
we've already seen that a quite simple network starts demonstrating some structure.

Now we'll look at different approaches to combining characters to obtain better representation
of source texts.

### Wavenet(ish)

Wavenet uses network layers to successively combine fragments
(the original is designed for audio but we'll use it to combine characters) so our output
is works over a large context length.

Now we've got a deeper network our intuition is that normalization layers will become useful
...

well neither our Wavenet alike or LayerNormalization exactly set the world on fire with 
the tiny names dataset.

we have learned how to build a deep network though and how to give ourselves a large context 
without too many parameters (for my cpu only laptop)

### Bigger corpus

we tried reading [alice in wonderland](../nn/alice.ipynb) with our simple neural network,
it did something, but not much

lets see if our Wavenet model does a better job down the [rabbit holw](wavenet_alice.ipynb)

### Onwards to attention

actually pretty promising - a good benchmark for our [next section](../attention) where we'll implement the
architecture from the 'Attention is all you need' paper




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

well neither our Wavenet alike or LayerNormalization exactly set the world on fire, even with 
the tiny names dataset.



## diverging

OK we start to diverge a bit from Andrej here

in the lectures he introduces BatchNorm and applies it to the network

I think this is the wrong place - it doesn't help the model, it's slow, better to wait for a failure in our network which a
normalization layer solves

Lets show the cross over to the torch.nn classes, make our code neater,
introduce a different initialization all in our nn [setup](neural_setup.ipynb)

### Names no more

another diversion: the names dataset is starting to run out of steam for us

I think even our simple network is starting to get about as much information as is
reasonable from names - and names are tricksy things anyway

to move past names, first we'll change our approach the [names.txt](../resources/names.txt) file.

we've taken each name individually, instead we'll address it as a [continuous](continuous_text.ipynb) and see 
how that affects our ability to learn.

Then we can move onto [Alice](alice.ipynb) in Wonderland (which is in fact the first digital book I read 
remotely in 1992...) as see if we get anything Lewis Carroll like our the other end
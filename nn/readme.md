## diverging

OK we start to diverge a bit from Andrej here

in the lectures he introduces BatchNorm and applies it to the network

I think this is the wrong place - it doesn't help the model, it's slow

Lets show the cross over to the torch.nn classes, make our code neater,
introduce a different initialization all in our nn [setup](neural_setup.ipynb)

### Names no more

another diversion: the names dataset is starting to run out of steam for us

I think even our simple network is starting to get about as much information as is
reasonable from names - names are tricksy things anyway

So instead lets look at [Alice](alice.ipynb) in Wonderland which is in fact the first digital book I read 
remotely in 1992...
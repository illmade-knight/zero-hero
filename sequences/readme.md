## Predict characters using more than one previous character

for 'emma' before we had

* e -> m
* m -> m
* m -> a
* a -> .

this becomes

* ...e  -> m
* ..em  -> m
* .emm  -> a
* emma  -> .

### look at these basic sequences

we'll start by basically reusing the network from the previous section on [embeddings](../embeddings)

our more powerful network should now be able to take advantage of that power to lower our loss

we'll start with [sequences](sequences.ipynb) in our 2D embedding space

### give ourselves more dimensions

we can start then start playing with our model parameters starting with [more dimensions](more_dimensions.ipynb)

we'll try and visualise what happens in the embedding space and now that our learning is improving we'll also 
start keeping track of how fast we learn and what we can do better 

### next: initialization

we look at improving learning through weight [initialization](../initialization)

### n.b

we also include some rough analysis of different model size in [analysis](sequence-play-analysis.ipynb)
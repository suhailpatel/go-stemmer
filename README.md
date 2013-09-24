Go Stemmer
===========

Implementation of the Porter Stemmer Algorithm in Go as defined by Martin Porter (author of the algorithm paper)

## Usage ##

```go
import "github.com/suhailpatel/stemmer"
```

Import the library into your application and then call the `Stem(string)` function. Stemmed words may not be perfect on context (please see the introduction original paper for information and discussion about this)

The library includes tests containing a large vocabulary provided by Martin Porter. 

## More Information ##

More information about the algorithm can be found on Martin Porter's Algorithm page (http://tartarus.org/martin/PorterStemmer/) and paper.

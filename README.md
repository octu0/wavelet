# `wavelet`

simple wavelet implementation.

### Example

```go
package main

import (
	"math"

	"github.com/octu0/wavelet"
)

func main() {
	high, low := wavelet.Haar([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 8.0})
	println(high) // => [2.12.., 4.94.., 9.19]
	println(low)  // => [-0.707.., -0.707.., -2.12..]

	out := wavelet.InverseHaar(high, low)
	println(out)  // => [0.999, 1.999, 3.000, 3.999, 5.000, 7.999]

	for i, v := range out {
		out[i] = math.Ceil(v)
	}
	println(out)  // => [1, 2, 3, 5, 8]
}
```

### Example RGBA

An example of converting an image to an intermediate format is implemented in [_example](https://github.com/octu0/wavelet/tree/master/_example)

source image 

| original      |                                     |
| :-----------: | :---------------------------------: |
| source        | ![img](_example/src.png)            |
| intermediate  | ![img](_example/intermediate.png)   |
| substract 55% | ![img](_example/substract_0.55.png) |
| substract 75% | ![img](_example/substract_0.75.png) |
| substract 85% | ![img](_example/substract_0.85.png) |
| substract 95% | ![img](_example/substract_0.95.png) |
| inverse       | ![img](_example/inverse.png)        |
| runlength-enc | 104KB                               |
| runlength-dec | ![img](_example/decoded.png)        |

# License

MIT, see LICENSE file for details.

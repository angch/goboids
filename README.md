# goboids

Swee Meng said he wanted to write a boids in PyGame.

And he did it! https://github.com/sweemeng/pygame-boids

## Running

No idea why Windows is giving me corrupted out put for more than 15 boids. Converting it to wasm and running it on
Chrome gives us at least 120fps/60tps.

```bash
go install github.com/hajimehoshi/wasmserve@latest
wasmserve .
open http://localhost:8080/
```

## In action

![](goboids.gif)

See also it running in your browser here: https://angch.github.io/goboids/

## References used

https://ebitengine.org/en/examples/vector.html
https://ebitengine.org/en/documents/webassembly.html
https://vergenet.net/~conrad/boids/pseudocode.html

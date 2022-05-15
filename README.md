# goboids

Swee Meng said he wanted to write a boids in PyGame

## Running

No idea why Windows is giving me corrupted out put for more than 15 boids. Converting it to wasm and running it on
Chrome gives us at least 120fps/60tps.

```bash
go install github.com/hajimehoshi/wasmserve@latest
wasmserve .
open http://localhost:8080/
```

## References used

https://ebiten.org/examples/vector.html
https://ebiten.org/documents/webassembly.html
https://vergenet.net/~conrad/boids/pseudocode.html
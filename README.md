# rgba-channel-merge

A little command line tool to merge specific layers of multiple images into one RGBA image.

## installation

```
go get github.com/blblblu/rgba-channel-merge
```

## usage example

Use the red, green and blue channel from first.png as green blue and red (in that order), and the alpha channel from second.png as alpha channel:

```
> rgba-channel-merge first.png gbrx second.png xxxa output.png
```

TODO: more readme...
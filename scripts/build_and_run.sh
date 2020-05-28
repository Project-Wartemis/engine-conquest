#!/bin/bash

(
    go get github.com/Project-Wartemis/pw-engine/cmd/engine &&
    $GOPATH/bin/engine
)

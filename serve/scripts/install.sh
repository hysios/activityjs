#!/usr/bin/env bash

export GOPATH=`realpath $1`
export GO111MODULE=off

echo $GOPATH
cd $1
go get github.com/gopherjs/gopherjs/js
go get github.com/gopherjs/jsbuiltin
go get github.com/kr/pretty
go get github.com/Landoop/tableprinter
go get github.com/olekukonko/tablewriter
go get github.com/juju/ansiterm/tabwriter
# go install activityjs.io/serve
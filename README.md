# Applications Packager

## Description

Utilite that can convert import json bundle from/to files of sim, ptl, csv, json.

### struct.dot

Is created in the process of packing or unpacking. Shows the structure of an application. Can be opened using [graphviz](http://graphviz.org/download/) or [webgraphviz](http://webgraphviz.com/)

## build

Windows tip. If you do not plan to work in the console, add "-ldflags -H=windowsgui"

>go build  

Binary files can be found in the current directory 

### on linux for windows

 >env GOARCH=amd64 GOOS=windows CGO_ENABLED=1 CC=/usr/bin/x86_64-w64-mingw32-gcc CXX=/usr/bin/x86_64-w64-mingw32-g++  go build


## Use the examples 

### 1. Unpack file from **basic.json** file to **basic/** folder

**basic.json** - The structure inside the file should look like *./app-tool/basic.json*, It is generated automatically in the Export module of the platform ecosystem.

Execute the command

```shell
$ ./app-tool basic.json
```

The **basic** folder is generated in the current directory

### 2. Pack files from **basic/** folder to **basic.json** file

**basic/** folder - The folder structure inside the file should look like *./app-tool/basic/*.

and the **config.json** structure in the directory look like *./app-tool/basic/config.json*.

Execute the command
```shell
$ ./app-tool basic/
```

The **basic.json** folder is generated in the current directory

### 3. Split json file by number

Execute the command, where **40** - number of files.
```shell
$ ./app-tool -s -n 40 basic.json

output:
   basic1.json
   basic2.json 
```

## Command line options

if you invoke the command without line arguments, the help will be displayed.

```shell
$ ./app-tool -h
Usage of ./app-tool:
  -d    debug
  -g    make graphical structure in dot-file
  -i string
        input (default ".")
  -input string
        -i, path for input files, filename for pack and dirname/ (slashed) for unpack (default ".")
  -n uint
        split json file by number
  -o string
        -output (default "output")
  -output string
        -o, output filename for JSON if input file name not pointed (default "output")
  -s    split json file by type
  -u    -unpack
  -unpack
        -u, unpacking mode
  -v    -version
```

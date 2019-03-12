# JDK verison switcher

This utility switches to a different version of the JDK, prepares the process
environment and chains into a provided command within that evironment.

Since the environment is changed for the duration of the provided command, this
utilty never changes the state of the machine.

## Install

    go get github.com/francescomari/jdk

## Usage

    jdk [version] [command...]

The version is either `8` or `11`. The command is required. `jdk` redirects its
stdin, stdout and stderr to the invoked command. Every environment variable
other than `JAVA_HOME` is passed unchanged to the invoked command. On a
successful invocation, the exist status of `jdk` is the same as the invoked
command's.

## License

This software is released under the MIT license.
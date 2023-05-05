# CentralisedVersionControl

Version control client-server app.
States of the client workplace are stored on the server and can be pushed or pulled by the client.
Supported operations are:
- push commit with message and creator
- list commits on server
- pull commit with specified id

The server stores the file tree and the compressed data of the files of every pushed commit.
When commits are pushed or pulled only files that are not already present for the other side are sent.

The client workplace is the current directory where the client app is ran.
When running the client for the first time the user will be asked to provide its username and remote address of the server.After that in the root of the directory will be created .cvc file that contains this information as well as filenames of ignored files. By default the first 2 ignored files are the .cvc file and the client app binary(if its present).

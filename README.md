# FS: A small file-server inspired by HTTP Server

## Introduction
**FS** is a simple _HTTP File Server_ which will serve the contents of the current directory. If you do not specify a
directory FS will serve your current `home` directory.

```bash
$ ./fs ~      // This will serve /Users/{username}/home/

$ ./fs /home/Users/will/Downloads/    // This will serve the specified folder
```
FS is primitive in its construction. It does not support multiple clients, multiple downloads nor does it support
pausing or restoration of a partial downloads the client. 

At the time of writing my main 'todo' is to increase the amount of unit testing.


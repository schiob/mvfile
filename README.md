# mvfile

mvfile is a simple program that uses github.com/fsnotify/fsnotify to watch a set of directories and move all the new files to a corresponding directory.
It waits until the files are done uploading before move them. This is specially useful for ftp servers where you have to wait for big files without having a cron job running every x time.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine.

### Prerequisites

You'll need a go environment, follow instructions in [go website](https://golang.org/doc/install).

### Installing

Download and install package:

```
go get github.com/schiob/mvfile
```

### Set up configuration files

Set the .json file with the directories to watch, the directories to send the files and the user who will be the owner of the file:

```
{
  "paths":
  [
    {
      "in_path": "/tmp/in",
      "out_path": "/tmp/exit",
      "user": "user"
    },
    {
      "in_path": "/tmp/in2",
      "out_path": "/tmp/exit2",
      "user": "user:group"
    }
  ]
}
```

Set the .toml file with the path of the .json file, the path of the log file and the time between checks for a file to end uploading.

```
# Path to json file
jsonpaths="dirPaths.json"

# Path of logfile
logfile="mvfile.log"

# Time in sec for wait between file size checks
wait=5
```

### Run it
```
mvfile -conf <path of .toml file>
```

## Contributing

Fell free to send pull requests or issues.

## Versioning

I use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/schiob/mvfile/tags).

## License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

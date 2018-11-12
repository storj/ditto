# Ditto

Ditto is a mirroring tool. You can connect two s3 compatible backends and replicate all work to both of them.
You can use Ditto as standalone S3 server with two S3 compatible backends connected.

## Installation

Make sure your PATH includes the $GOPATH/bin directory so your commands can be easily used:

```bash
export PATH=$PATH:$GOPATH/bin
```

### Go package manager

```bash
$ go get github.com/storj/ditto
```

### Build from sources

1. Create folder 'storj' at your $GOPATH/src folder:
```bash
$ mkdir $GOPATH/src/storj
```

2. Navigate to newly created folder and clone ditto repository:
```bash
$ cd $GOPATH/src/storj
$ git clone https://github.com/storj/ditto.git
$ cd ditto
```

3. Install ditto:
```bash
go install 
```

## Usage

First of all you have to set up your S3 Credentials for both backends.
You can do this using your preferred text edittor or by Ditto's config command.

Config file located at $HOME/.ditto/config.json

All you have to setup is Server1 and Server2 sections.
```json
{
  "server1": {
    "accesskey": "AccessKeyForServer1",
    "endpoint": "EndpointURLforServer1",
    "secretkey": "SecretKeyForServer1"
  },
  "server2": {
    "accesskey": "AccessKeyForServer2",
    "endpoint": "EndpointURLforServer2",
    "secretkey": "SecretKeyForServer2"
  }
}
```

The rest will be used from default settings, which you can modify at this config file, or by 'ditto config set'

### CLI usage

#### List(ls)

Prints list of buckets and content of specified bucket.
If no arguments set, prints bucket list.
If bucket name specified(as first argument), prints files list at that bucket

```bash
  -s, --default_source string   Defines source server to display list (default "server2")
  -d, --delimiter string        Char or char sequence that should be used as prefix delimiter (default "/")
  -m, --merge                   Display list from both servers and merge to single list
  -p, --prefix                  Folder simulation path
  -r, --recursive               Shows all nested folder if prefix set
  -t, --throw_immediately       In case of error, throw error immediately, or retry from other server
```
##### Examples

Assume you have two backends: Server1 and Server2 with this file structure on each:

```
├── Server1
│   └── bucket0
│       └── b1
│           ├── bb1
│           │   ├── ff10
│           │   └── ff11
│           ├── f10
│           └── f11
└── Server2
    └── bucket0
        └── b0
            ├── bb0
            │   ├── ff0
            │   └── ff1
            └── f0
```

```bash
$ ditto list 
Buckets:
    bucket0
```

```bash
$ ditto list bucket0
Files at bucket0
	 b1/bb1/ff10
	 b1/bb1/ff11
	 b1/bb1/ff12
	 b1/f10
	 b1/f11
```

```bash
$ditto list bucket0 --default_source="server2"
Files at bucket0
	 b0/bb0/ff0
	 b0/bb0/ff1
	 b0/f0
```

```bash
$ditto list --merge bucket0
Merged files list from bucket0
	 b0/bb0/ff0
	 b0/bb0/ff1
	 b0/f0
	 b1/bb1/ff10
	 b1/bb1/ff11
	 b1/bb1/ff12
	 b1/f10
	 b1/f11
```

```bash
# Note that you should specify delimiter char at the end of the prefix string
$ditto list --prefix bucket0 b0/
Files at bucket0:
	 b0/f0
```

```bash
$ditto list --prefix --recursive bucket0 b0/
Files at bucket0:
	 b0/f0
	 b0/bb0/ff0
	 b0/bb0/ff1
```

#### Put(p)
#### Get(g)
#### Delete(rm)

Gives ability to delete file, all files that match prefix, or delete whole bucket with content.

The only output you can receive from this Commad is an Error message. In case of successfull bucket or file removal there will be no output.

```
  -s, --default_source string   Defines source server to start from (default "server1")
  -d, --delimiter string        Char or char sequence that should be used as prefix delimiter (default "/")
  -f, --force                   if force flag applied - all files without prefixes in bucket will be removed.
  -h, --help                    help for delete
  -p, --prefix                  Folder simulation path
  -r, --recursive               User force flag to delete bucket
  -t, --throw_immediately       in case of error, throw error immediately, or retry from other server (default true)
```
##### Examples
Assume you have two backends: Server1 and Server2 with this file structure on each:

```
├── Server1
│   └── bucket0
│       └── b1
│           ├── bb1
│           │   ├── ff10
│           │   └── ff11
│           ├── f10
│           └── f11
└── Server2
    └── bucket0
        └── b0
            ├── bb0
            │   ├── ff0
            │   └── ff1
            └── f0
```

```bash
$ditto rm bucket0 b0/f0

Will remove only file b0/f0 from bucket0
```

```bash
$ditto rm -rf bucket0 

Will remove everything from bucket0 and delete bucket0.
```

```bash
$ditto rm -p bucket0 b0/

Will remove only files at b0/ prefix. Nested files(prefixes) will not be affected.
```

```
$ditto rm -pr bucket0 b0/

Will remove everythig at b0/ prefix and will remove all nested files(prefixes) Bucket 'bucket0' will not be deleted
```
#### Copy(cp)

Creates copy of an object.
```bash
$ditto cp sourceBucket sourceObject destinationBucket destinationObject(OPTIONAL)
```
##### Note! 
- DestinationBucket must be created before cp operation
- If no destinationObject provided, sourceObject will be used instead

##### Examples

```bash
$ditto cp bucket0 file0 bucket1
Object bucket0/file0 copied
```

```bash
$ditto cp bucket0 file0 bucket1 file1
Object bucket0/file0 copied

This will copy file0 from bucket0 as file1 at bucket1
```
#### Make-Bucket(mb)

Creates bucket with given name. 

##### Examples

```bash
$ditto make-bucket bucket0
Bucket 'bucket0' created
```

#### Config
#### Version

### Server usage

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
See LICENSE file for Info

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

Gives ability to upload file or folders recursively.

```bash
  -s, --default_source string   Defines source server to display list (default "server2")
  -d, --delimiter string        separates object names from prefixes (default "/")
  -f, --force                   truncate object if one exists
  -p, --prefix string           root prefix
  -r, --recursive               recursively upload contents of the specified folder
  -t, --throw_immediately       In case of error, throw error immediately, or retry from other server
```

```bash
$ ditto p bucket0 ~/Downloads/testImage.png

Successfully uploaded 'testImage.png'
```

```bash
$ ditto p bucket0 ~/Downloads

Found 5 folders, add -r flag to recursively upload them
successfully uloaded 'testFile1.png'
successfully uloaded 'testFile2.png'
successfully uloaded 'testFile3.png'
successfully uloaded 'testFile4.png'
successfully uloaded 'testFile5.png'
successfully uloaded 'testFile6.png'
successfully uloaded 'testFile7.png'
successfully uloaded 'testFile8.png'
successfully uloaded 'testFile9.png'
successfully uloaded 'testFile10.png'
Err: Object allready exists 'testImage.png'
```

```bash
$ditto p -r bucket0 ~/Downloads

Recursively uploading folder /Users/user/Downloads/b0
successfully uloaded 'b0/f0'
successfully uloaded 'b0/f1'
successfully uloaded 'b0/f2'
Recursively uploading folder /Users/user/Downloads/b0/bb0
successfully uloaded 'b0/bb0/f0'
successfully uloaded 'b0/bb0/f1'
successfully uloaded 'b0/bb0/f2'
...
```

#### Get(g)

Gives ability to download files or folder recursively

```bash
  -s, --default_source string   Defines source server to display list (default "server2")
  -d, --delimiter string        separates object names from prefixes (default "/")
  -f, --force                   truncate object if one exists
  -h, --help                    help for put
  -p, --prefix string           root prefix
  -r, --recursive               recursively upload contents of the specified folder
  -t, --throw_immediately       In case of error, throw error immediately, or retry from other server
```

Assume you have this bucket structure at 'bucket0'
```bash
bucket0
    ├── bb0
    │   ├── ff0
    │   ├── ff1
    │   └── ff2
    ├── bb1
    │   ├── ff10
    │   ├── ff11
    │   └── ff12
    ├── f0
    ├── f1
    ├── f2
    └── s3-api.pdf 
```

- Download single file
```bash
$ ditto g bucket0 s3-api.pdf

successfully downloaded 's3-api.pdf'
```

- Download all files at bucket0 non-recursively
```bash
$ ditto get bucket0 

successfully downloaded f0
successfully downloaded f1
successfully downloaded f2
Err: s3-api.pdf: file exists
Found new prefix bb0/, missing -r flag to download it recursively
Found new prefix bb1/, missing -r flag to download it recursively
```

- Download all files at bucket0 including nested folders recursively
```bash
$ ditto get -r bucket0

successfully downloaded f0
successfully downloaded f1
successfully downloaded f2
successfully downloaded bb0/ff0
successfully downloaded bb0/ff1
successfully downloaded bb0/ff2
successfully downloaded bb1/ff10
successfully downloaded bb1/ff11
successfully downloaded bb1/ff12
```
#### Delete(rm)

Gives ability to delete file, all files that match prefix, or delete whole bucket with content.

The only output you can receive from this Commad is an Error message. In case of successfull bucket or file removal there will be no output.

```
  -s, --default_source string   Defines source server to start from (default "server1")
  -d, --delimiter string        Char or char sequence that should be used as prefix delimiter (default "/")
  -f, --force                   if force flag applied - all files without prefixes in bucket will be removed.
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
You can modify config file located at ``` $HOME/.ditto/config.json``` using ditto configurator with value check.
You can enter any data to 'Server1' and 'Server2' sections. Rest will accept only predefined set of options.(See examples sections for further instructions).

Available commands:
```
  get         Displays value set for requested key
  list        Displays list of possible to change options
  set         Change value at saves it to config file
```

##### Examples

- Request list of available options
```bash
$ ditto config list

Options, which can be set via `config set`:
	Server1.Endpoint
	Server1.AccessKey
	Server1.SecretKey
	Server2.Endpoint
	Server2.AccessKey
	Server2.SecretKey
	DefaultOptions.DefaultSource
	DefaultOptions.ThrowImmediately
	ListOptions.DefaultOptions.DefaultSource
	ListOptions.DefaultOptions.ThrowImmediately
	ListOptions.Merge
	PutOptions.DefaultOptions.DefaultSource
	PutOptions.DefaultOptions.ThrowImmediately
	PutOptions.CreateBucketIfNotExist
	GetObjectOptions.DefaultOptions.DefaultSource
	GetObjectOptions.DefaultOptions.ThrowImmediately
	CopyOptions.DefaultOptions.DefaultSource
	CopyOptions.DefaultOptions.ThrowImmediately
	DeleteOptions.DefaultOptions.DefaultSource
	DeleteOptions.DefaultOptions.ThrowImmediately
```

- Set Endpoint URL for Server1
```bash
$ ditto config set Server1.Endpoint http://localhost:9007
```
You will see output only in case of an error.

- Set default source for list cmd(Wrong approach)
```bash
$ ditto config set ListOptions.DefaultOptions.DefaultSource anySource

Error: Only these arguments accepted: [server1 server2]
Usage:
  ditto config set [key] [value] [flags]

Flags:
  -h, --help   help for set
```
As you can see, internal validator accepts only predefined set of values. If you make an error, ditto will gently suggest most suitable value.

```bash
$ ditto config set ListOptions.DefaultOptions.DefaultSource server1
```

- Get value assigned to option at config file
```bash
$ ditto config get ListOptions.DefaultOptions.DefaultSource

	server1
```

### Server usage

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
See LICENSE file for Info

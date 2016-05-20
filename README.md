# appmanifest

`appmanifest` takes a pkg and prints an [application manifest](http://help.apple.com/deployment/osx/#/ior5df10f73a)
the current version only creates the `assets` array.

The documentation says the metadata is required, but installs work without the metadata dict.
Adding one only affects what shows up in `Launchpad`

```
appmanifest [options] /path/to/some.pkg
  -url string
    	url of the pkg as it will be on the server
  -version
    	prints the version
```

# Install
see Releases for binary builds
to build from source, use the make Makefile

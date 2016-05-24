# Install
see Releases for binary builds

to build from source, use the make Makefile


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


# certhelper
create and manage push certificate
```
# creates an MDM CSR and private key for vendor cert
# upload the created MDM CSR to enterprise portal to get a push certificate
certhelper vendor -csr -cn=mdm-certtool -password=secret -country=US -email=foo@gmail.com

# create a "provider" or a "customer" csr. This will be signed by the vendor cert and submitted to apple to get a push cert
certhelper provider -csr -cn=mdm-certtool -password=secret -country=US -email=foo@gmail.com

# sign the provider csr with the vendor private key
# assumes `mdm.cer` is in the folder with all the other files. You can specify each path separately as well.
certhelper vendor -sign -password=secret

# Now upload the PushCertificateRequest to https://identity.apple.com/pushcert
```
see the `certhelper` README for more details

# poke
send mdm push notification to APNS  
see `poke` README.md for usage.


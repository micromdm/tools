poke - sends a MDM push notification to a device using the HTTP2 APNS gateway


An utility for testing of MDM push notifications; asking devices to connect to their MDM server.

# How to use
```
Usage of poke:
  -magic string
    	pushmagic
  -push-cert string
    	path to push certificate
  -push-pass string
    	push certificate password
  -token string
    	deviceToken
  -version
    	print version information
```

We can limit the number of flags we set by declaring some of the configuration as environment variables:

```
# ~/push_env

# export from keychain as PKCS#12, not .cer
export MDM_PUSH_CERT=/path/to/pushcert.p12
export MDM_PUSH_PASS=secret
```

```
source push_env
./poke -magic=2AD29D04-2440-4816-B9C2-935F2F0AC1C4 -token=f8b4ccc6da57207807fcff9767a0e15aec204721fee7ff2cd6cd5c16402b1ad5
```

`poke` will send a MDM formatted push notification to the APNS gateway and return back an UUID

# Debugging information
`poke` will print additional verbose APNS logs with `GODEBUG=http2debug=2` environment variable.
```
GODEBUG=http2debug=2 ./poke -magic=6E2881DC-6088-43E5-8E1F-EF267FC8B71D -token=f824ccc6da57207807fcff9767a0e15aee404721fee7ff28d69d5016604a08d5
```

Example debug output:
```
2016/03/03 11:36:06 http2: Transport failed to get client conn for api.push.apple.com:443: http2: no cached connection was available
2016/03/03 11:36:06 http2: Transport creating client conn to 17.172.234.15:443
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: wrote SETTINGS len=18, settings: ENABLE_PUSH=0, INITIAL_WINDOW_SIZE=4194304, MAX_HEADER_LIST_SIZE=10485760
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: wrote WINDOW_UPDATE len=4 (conn) incr=1073741824
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: read SETTINGS len=24, settings: HEADER_TABLE_SIZE=4096, MAX_CONCURRENT_STREAMS=500, MAX_FRAME_SIZE=16384, MAX_HEADER_LIST_SIZE=8000
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: wrote SETTINGS flags=ACK len=0
2016/03/03 11:36:06 Unhandled Setting: [HEADER_TABLE_SIZE = 4096]
2016/03/03 11:36:06 Unhandled Setting: [MAX_HEADER_LIST_SIZE = 8000]
2016/03/03 11:36:06 http2: Transport encoding header ":authority" = "api.push.apple.com"
2016/03/03 11:36:06 http2: Transport encoding header ":method" = "POST"
2016/03/03 11:36:06 http2: Transport encoding header ":path" = "/3/device/xxx"
2016/03/03 11:36:06 http2: Transport encoding header ":scheme" = "https"
2016/03/03 11:36:06 http2: Transport encoding header "content-type" = "application/json"
2016/03/03 11:36:06 http2: Transport encoding header "content-length" = "46"
2016/03/03 11:36:06 http2: Transport encoding header "accept-encoding" = "gzip"
2016/03/03 11:36:06 http2: Transport encoding header "user-agent" = "Go-http-client/2.0"
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: wrote HEADERS flags=END_HEADERS stream=1 len=132
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: wrote DATA stream=1 len=46 data="{\"mdm\":\"\"}"
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: wrote DATA flags=END_STREAM stream=1 len=0 data=""
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: read GOAWAY len=46 LastStreamID=0 ErrCode=NO_ERROR Debug="{\"reason\":\"BadCertificateEnvironment\"}"
2016/03/03 11:36:06 http2: Transport received GOAWAY len=46 LastStreamID=0 ErrCode=NO_ERROR Debug="{\"reason\":\"BadCertificateEnvironment\"}"
2016/03/03 11:36:06 http2: Framer 0xc8204e9800: read GOAWAY len=8 LastStreamID=0 ErrCode=NO_ERROR Debug=""
2016/03/03 11:36:06 http2: Transport received GOAWAY len=8 LastStreamID=0 ErrCode=NO_ERROR Debug=""
2016/03/03 11:36:06 Transport readFrame error: (*errors.errorString) EOF
2016/03/03 11:36:06 RoundTrip failure: unexpected EOF
```

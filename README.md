# go-tiny

A "tiny url" generator and redirector, written in Go, using the Goji framework, and go-cache-lru. 

Developed as a testbed for go-cache-lru, but released since a couple people thought they'd like to play around with it. 
It has a couple cute tricks just because I didn't feel like releasing something as entirely boring as it was. TIMTOWTDI

 To create: `/set/http://lmgtfy.com/?q=I+saw+a+turtle`

 To use: `/5041b150`  (or whatever the above returned)

NOTES: The cache only lasts as long as the server is running. Nothing persists. go-cache-lru does support dumping its 
contents to a file, and priming from a file, but that's your hack to add. The "tiny url" returned is the crc32 of the 
URL, plus random offset. Every time the server is run, there will be a new offset, so the "tiny" URL will be different 
for the same URL.

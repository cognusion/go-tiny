package main

import (
	"io"
	"fmt"
	"github.com/cognusion/go-cache-lru"
	"time"
	"crypto/rand"
	"math/big"
	"math"
	"strings"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/cognusion/tinysum"
)

var VERSION = "go-tiny 1.0.5"
var C *cache.Cache
var GLOBALOFFSET uint32

// Simply return the version
func version(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", VERSION)
}

// Simply return the number of cached urls
func count(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", C.ItemCount())
}

// Grab an item from the cache and redirect to the stored URL
func fetch(c web.C, w http.ResponseWriter, r *http.Request) {
	
	tiny := c.URLParams["tiny"]
	fmt.Printf("Getting %v\n", tiny)
	
	turl, found := C.Get(tiny)
	if ! found {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	
	url := turl.(string)
	if strings.HasPrefix(url,"http:") || strings.HasPrefix(url,"https:") {
		url = strings.Replace(url, ":/", "://", 1)	
	} else {
		url = "http://" + url
	}
	fmt.Printf("Got %v as %v\n", tiny, turl)
	
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Generate the offset crc32, and set it in the cache
func set(c web.C, w http.ResponseWriter, r *http.Request) {
	
	url := r.RequestURI
	url = url[5:len(url)] // get rid of leading "/set"
	
	fmt.Printf("Setting %v\n", url)
	
	// We take the crc32 of the URL add a random offset
	v := tinysum.OffsetStringSum(url, GLOBALOFFSET)
	
	C.Set(v, url, cache.DefaultExpiration)
	
	fmt.Printf("Set %v to %v\n", url,v)
	
	msg := `<html>
	<body>
	<a href="%v">%v</a>
	</body>
	</html>
	`
	msg = fmt.Sprintf(msg, "/" + v, "/" + v)
	
	io.WriteString(w, msg)
}


func main() {

	// So the URLs aren't pure crc32s of the URI
	thirtyTwo := math.Pow(2,32)-1
	max := *big.NewInt(int64(thirtyTwo))
	roff,_ := rand.Int(rand.Reader,&max)
	GLOBALOFFSET = uint32(roff.Uint64())
	fmt.Printf("Offset is %v\n",GLOBALOFFSET)
	
	// Keep items for 24hours, clean every 30s, cap at 50k items.
	C = cache.New(24*time.Hour, 30*time.Second, 50000)

	// Set the URI handlers, and go!
	goji.Get("/version", version)
	goji.Get("/count", count)
	goji.Get("/set/*", set)
	goji.Get("/:tiny", fetch)
	goji.Serve()

}

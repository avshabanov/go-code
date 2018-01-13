package main

import (
	"flag"
	"fmt"
	"net/url"
)

var demo = flag.String("demo", "1", "Demo parameter")

func main() {
	flag.Parse()
	if *demo == "1" {
		urlConstructDemo()
	} else {
		urlParseDemo()
		urlConstructDemo()
	}
}

//
// Demo
//

func urlParseDemo() {
	fmt.Println("*** urlParseDemo ***")

	s := "http://example.com/v1/resource?param1=123&param2=test"
	u, err := url.Parse(s)
	if err != nil {
		fmt.Printf("Can't parse url: %v", err)
		return
	}

	fmt.Printf(
		"Url: rawQuery=%s, fragment=%s, host=%s\n",
		u.RawQuery,
		u.Fragment,
		u.Host,
	)

	for k, v := range u.Query() {
		fmt.Printf("\tparam: %s = %s\n", k, v)
	}
}

func urlConstructDemo() {
	fmt.Println("*** urlConstructDemo ***")

	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	s := "http://example.com/v1/resource?param1=123&param2=test"
	u, err := url.Parse(s)
	if err != nil {
		fmt.Printf("Can't parse url: %v", err)
		return
	}

	v := u.Query()
	v.Set("param3", "somethingelse")
	u.RawQuery = v.Encode()
	u.Fragment = "en/ru/beep"

	fmt.Printf("new url = %s\n", u.String())
}

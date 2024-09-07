package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// We need a way to take an RSS url/ feed url, and kindof parse it into a response body,
// which in this case will be a struct

// RSS is basically structured data, in xml format
// And xml is really just crappy json

// The way we parse xml in go is very similar to how we parse json
//This RSS Feed struct is gonna represent a giant xml document, in the rss feed
// U basically just scan all the tags that rss feed has, and u put them as ID's in structs
// See them here in : https://www.wagslane.dev/index.xml
// Use chrome to see it though, safari doesn't show it properly

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// These RSSItems are basically the posts. Again, see the xml file and ur gonna understand
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// This will take the url to the feed as input, and will return a new type, called RSSFeed, and an error
func urlToFeed(url string) (RSSFeed, error) {
	// First we r gonna need a new httpClient
	// An httpClient sends HTTP requests and receives HTTP responses from a resource identified by a URI
	// Basically we can use to make requests to the web server

	// Timeout(type time.Duration): your client will open a connection to the server via HTTP; the server may take some time to answer your request.
	// This field will define a maximum time to wait for the response. If you do not specify a timeout, there is no timeout by default.
	// This can be dangerous for the user-perceived performance of your service.
	// It’s better, in my opinion, to display an error than to make the user wait indefinitely. Note that this time includes:
	// Connection time (the time needed to connect to the distant server)
	// Time taken by redirects (if any)
	// Time is taken to read the response body.

	// We can also specify parameters for the http.Client like Transport, checkRedirect, Jar, but we don't need them now
	// https://www.practical-go-lessons.com/chap-35-build-an-http-client

	// We r gonna set the timeout to 10s coz fun
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// We r now gonna use that client to make a get request to the url of the feed
	// This will return an HTPP Response, and an err
	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	// If the response is fine, we will defer a close on the response Body
	// resp.Body is a stream of data read by the http client.
	// Do not forget to add this closing instruction; otherwise, the client might not reuse a potential persistent connection to the server
	defer resp.Body.Close()

	// Now, lets read all the data from the response body. So u io.ReadAll() it
	// This returns a slice of bytes and an error
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	// I logged it, and it's just a bunch of numbers, u can see it if u want
	// log.Printf("%x", dat)

	// Now, we will parse the slice of bytes into our RSSFeed type
	// So u xml.Unmarshal instead of json.Unmarshalling
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	// This returns the wholeass feed
	// We now need to write a scraper that will get the feed, and then take the indivdual posts out of it
	// Which are basically the RSSItems
	return rssFeed, nil

}

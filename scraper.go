package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Yendelevium/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

// The scraper is a long running job. This scraper function will run on the bg of our server
// as long as the server is up

// This will take 3 inputs, a connection to our database, how many diff goroutines we wanna do the scraping on,
// And the time delay between each request to scrape a new RSSFeed
// It won't return anything as it will be running forever as long as our server is up

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBewteenRequest time.Duration,
) {
	// Because this scraper is gonna run in the bg of our server, it will be good if we know what is going on while its doing its thing
	// Hence, we will need a lot of good logging, to know what's up
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBewteenRequest)

	// We need to figure out how we wanna make requests in the given time interval
	// This is where a ticker comes to play
	// NewTicker returns a new Ticker containing a channel that will send the current time on the channel after each tick.
	// The period of the ticks is specified by the duration argument. The ticker will adjust the time interval or drop ticks to make up for slow receivers.
	// The duration d must be greater than zero; if not, NewTicker will panic.

	ticker := time.NewTicker(timeBewteenRequest)

	// ticker.C is the channel of the ticker
	// The reason we r passing in ; ;<blah> in the forloop, is so that the first time, the for loop
	// starts immediately, and then it waits for 1 min(assuming that's the time interval) till the channel returns a value(remember, channels can block)
	// And then it goes again
	// Basically, so that we initially don't have to wait for 1 min before we start scraping
	// If we just did : for range ticker.C, it will wait for the time first, then scrape
	for ; ; <-ticker.C {

		// Every interval, we wanna go grab the next batch of feeds to fetch
		// The function takes a context, and the no.of feeds u wanna fetch, which will be the no.of goroutines running at the same time
		// So we pass concurrency here

		// Since we don't have a request body, we can't just do r.Context()
		// So we use context.Background(), which is basically a global context, and it's what we use when we don't
		// have scoped contexts, like our individually http requests(r.context())

		// This will return the RSSFeeds and an error
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		// While handling the error, we aren't terminating the function, but instead, continuing to the next scrape
		// This is because this function, should always be running in the bg as our server operates
		// So we instead, log the error and wait for the next scrape
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		// Now that we have a slice of feeds, we need to fetch these feeds individually, and more importantly
		// It fetches the feeds at the same time. SO we r gonna neeed a synchronization method, so we r gonna use a waitGroup
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			// The way that waitGroup works, is that anytime u wanna make a new goroutine in the context of that wg,
			// U wg.Add(<number>) where number is the no.of goroutines ur making
			// Since in every iteration of this for loop, we are making 1 goroutine to fetch 1 feed, we add 1 to the wg
			wg.Add(1)

			// Now lets spawn a new goroutine to get the feed. Here, we will pass the wg in as one of the params
			// And within the function, we will defer wg.Done(), so it will know that that goroutine is finished
			// The wg will allow us to call various goroutines at the same time, and will block the function, till all of them r done
			// Which is what we wanna do as we don't wanna continue to the next iteration of the loop until we r sure we have scraped all the feeds
			go scrapeFeed(db, wg, feed)
		}
		// Now at the end of the loop, we add a wg.Wait(), which will wait till all the goroutines are done
		// Only then will it proceed
		wg.Wait()
	}
}

// This function will iterate through all the POSTS, (RSSItems), in the feed
// Within the function, we will defer wg.Done(), so wg will know that that goroutine is finished
// This function will need a db connection, and also a specific feed to fetch
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	// The first thing this function should do, is to mark that we r fetching this feed
	// Here also we will pass context.Bg() as again, we don't have access to a scoped context
	// This returns the feed which we marked as fetched, but we don't need it coz we already have it
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	// If we can't mark it as fetched , we'll just log there was an issue and return nothing
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	// Now we have tp scrape the feed, and we hv already written the function for that, so lets use it
	// We log any errors
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	// Now, we can iterate through each post, which will be in the slice: rssFeed.channel.item
	for _, item := range rssFeed.Channel.Item {

		// This is just creating a post
		// We ran into dome problems, like here, descirption was an NullString
		// This nullString is a struct, which has 2 values, the string, and if it's null or not
		// So we can't directly pass the item.Description as its just a string, not a sql.nullstring
		// So we gotta create that nullstring
		description := sql.NullString{}
		// If the desciption is not empty, set valid as true
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// same thing here, PubDate is a String, but for the params we need a time.Time type
		// So we are parsing the string and giving back this specific time format
		// This is the time layout the guy used on his blog, but to make this project more inclusive,
		// U prolly have to take care of all the different time layouts
		// But for now this is fine
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn;t parse date %v with err %v", item.PubDate, err)
		}
		// Now we gonna create the post in the post table finally
		// This returns the post, which we don't need, and an error
		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdateAt:    time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			})

		// Logging the error
		if err != nil {
			// Here, since we keep scraping all the posts and try to add them to the posts
			// We can get an error as of a post is already scraped, it's url won't be unique, so we get a duplicate key error
			// This is expected behaviour, so we don't need to log this error
			// So if out error contains "dupliucate key", we are just gonna skip the logging part
			// err.Error() will return the error STRING, and then we can use strings.Contains() to check if the err has duplicate key in it.
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("failed to create post:", err)
		}
	}

	// Just doing some logging so we know how many posts we collected and from which feed
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}

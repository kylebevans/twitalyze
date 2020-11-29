// Twitalyze searches twitter for a filter and creates a word to frequency mapping.
// It serves the data over an API that a frontend can read from to create a word cloud.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/antihax/optional"
	"github.com/bbalet/stopwords"
	"github.com/gorilla/mux"
	"github.com/kylebevans/twitapi"
)

// Map that holds words with the frequency of occurrence in tweets.
type WordValues struct {
	sync.RWMutex
	WV map[string]int `json:"wordvalues"`
}

// Initialize a WordValues and return a pointer.
func NewWordValues() *WordValues {
	var w WordValues
	w.WV = make(map[string]int)
	return &w
}

// Grab tweet data from the last 7 days to seed the word cloud.
func SeedData(ctx context.Context, s string, apiClient *twitapi.APIClient, w *WordValues) {
	var searchOpts *twitapi.TweetsRecentSearchOpts
	searchOpts = new(twitapi.TweetsRecentSearchOpts)
	var tweets twitapi.TweetSearchResponse
	var nextToken optional.String
	var err error

	// API is paginated, so print pages of recent tweets in the last week
	// that match the query until they are all done
	for ok := true; ok; ok = (tweets.Meta.NextToken != "") {

		tweets, _, err = apiClient.SearchApi.TweetsRecentSearch(ctx, s, searchOpts)

		if err != nil {
			log.Printf("Could not seed data: %v", err)
			return
		}

		for _, v := range tweets.Data {
			ParseTweet(v.Text, w)
		}

		nextToken = optional.NewString(tweets.Meta.NextToken)
		searchOpts.NextToken = nextToken
		time.Sleep(2 * time.Second) // Twitter API is rate limited to 450 requests per 15 min.
	}
}

// PrintTweet prints out a tweet followed by a little red line.
func PrintTweet(tweet twitapi.FilteredStreamingTweet) {
	fmt.Printf("%v\n\u001b[31m--------\n\u001b[0m", tweet.Data.Text)
}

// ParseTweet removes stop words from a tweet and increments the WordValues for each word.
func ParseTweet(tweet string, w *WordValues) {
	re, _ := regexp.Compile(`@.+? `) // remove @username words.
	cleanTweet := re.ReplaceAllString(tweet, "")
	cleanTweet = stopwords.CleanString(cleanTweet, "en", false) // remove words like "the"
	cleanTweet = strings.ToLower(cleanTweet)
	tweetTokens := strings.Split(cleanTweet, " ")

	for _, v := range tweetTokens {
		// Discard words above 16 letters, below 3 letters, and several undesirable words.
		if len(v) > 16 || len(v) < 3 || v == "terraform" || v == "cloud" || v == "" || v == "https" || v == "'terraform" || v == "hashicorp" {
			continue
		}
		w.Lock()
		w.WV[v]++
		w.Unlock()
	}
}

// StreamTweets removes all the current rules from the Twitter stream rules,
// adds a new rule, receives tweets from the stream, and ships them off to
// the processing function provided by the caller.
func StreamTweets(ctx context.Context, s string, apiClient *twitapi.APIClient, f func(string, *WordValues), w *WordValues) error {
	var tweet twitapi.FilteredStreamingTweet
	var err error
	var searchOpts *twitapi.SearchStreamOpts
	searchOpts = new(twitapi.SearchStreamOpts)
	var rules = []twitapi.RuleNoId{
		twitapi.RuleNoId{
			Value: s,
			Tag:   s,
		},
	}
	var ruleReq = twitapi.AddRulesRequest{
		Add: rules,
	}

	// Delete any existing rules.
	var getRulesResp twitapi.InlineResponse2002
	var ruleIds twitapi.DeleteRulesRequest
	var delResp twitapi.AddOrDeleteRulesResponse
	getRulesResp, _, err = apiClient.TweetsApi.GetRules(ctx, nil)
	for _, v := range getRulesResp.Data {
		ruleIds.Delete.Ids = append(ruleIds.Delete.Ids, v.Id)
	}
	delResp, _, err = apiClient.TweetsApi.AddOrDeleteRules(ctx, ruleIds, nil)
	if err != nil {
		log.Println(ruleIds)
		log.Println(delResp)
		panic(err)
	}

	// Add new rule.
	_, _, err = apiClient.TweetsApi.AddOrDeleteRules(ctx, ruleReq, nil)
	if err != nil {
		log.Printf("Unable to add search filter rule: %v", s)
		panic(err)
	}

	//Search stream for tweets.
	tweets := make(chan twitapi.FilteredStreamingTweet)
	errs := make(chan error)
	for {
		go apiClient.TweetsApi.SearchStream(ctx, searchOpts, tweets, errs)
		for {
			err = <-errs
			if err != nil {
				break
			}
			tweet = <-tweets
			f(tweet.Data.Text, w) //Call tweet processing function.
		}
		log.Println(err)
	}
}

// SaveData saves the data to seed.conf every 10 minutes.
func SaveData(w *WordValues) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.RLock()
			outdata, err := json.Marshal(w)
			w.RUnlock()
			if err != nil {
				log.Printf("Unable to save word values to seed file: %v", err)
			} else {
				err = ioutil.WriteFile("seed.conf", outdata, 0644)
				if err != nil {
					log.Printf("Unable to save word values to seed file: %v", err)
				}
			}
		}
	}
}

func (w *WordValues) WordsHandler(wr http.ResponseWriter, r *http.Request) {
	outdata, err := json.Marshal(w)
	if err != nil {
		log.Printf("Unable to convert words to JSON: %v", err)
		wr.WriteHeader(500)
		io.WriteString(wr, "{\"error\" : \"Unable to provide words\"}")
		return
	}

	wr.WriteHeader(200)
	wr.Write(outdata)
}

func main() {

	ctx := context.WithValue(context.Background(), twitapi.ContextAccessToken, os.Getenv("TWITTER_BEARER_TOKEN"))
	cfg := twitapi.NewConfiguration()
	apiClient := twitapi.NewAPIClient(cfg)
	wordNums := NewWordValues()
	searchFilter := "\"terraform cloud\""
	port := ":8080"

	// Read in seed file if it exists, else call SeedData to get data from last 7 days
	if _, err := os.Stat("seed.conf"); err == nil {
		seed, _ := ioutil.ReadFile("seed.conf")
		_ = json.Unmarshal(seed, wordNums)
	} else {
		SeedData(ctx, searchFilter, apiClient, wordNums)
	}

	go StreamTweets(ctx, searchFilter, apiClient, ParseTweet, wordNums)

	go SaveData(wordNums)

	// Serve the API
	r := mux.NewRouter()
	r.HandleFunc("/words", wordNums.WordsHandler)
	log.Printf("Listening on port %v", port)
	http.ListenAndServe(port, r)

}

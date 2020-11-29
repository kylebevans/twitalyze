# twitalyze
backend for a word cloud based on tweets

## Functions

### func [ParseTweet](/main.go#L65)

`func ParseTweet(tweet string, w *WordValues)`

ParseTweet removes stop words from a tweet and increments the WordValues for each word.

### func [PrintTweet](/main.go#L60)

`func PrintTweet(tweet twitapi.FilteredStreamingTweet)`

PrintTweet prints out a tweet followed by a little red line.

### func [SeedData](/main.go#L32)

`func SeedData(ctx context.Context, s string, apiClient *twitapi.APIClient, w *WordValues)`

Grab tweet data from the last 7 days to seed the word cloud.

### func [StreamTweets](/main.go#L83)

`func StreamTweets(ctx context.Context, s string, apiClient *twitapi.APIClient, f func(string, *WordValues), w *WordValues) error`

StreamTweets removes all the current rules from the Twitter stream rules,
adds a new rule, receives tweets from the stream, and ships them off to
the processing function provided by the caller.

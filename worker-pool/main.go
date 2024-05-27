package main

//TODO: when to use worker pool
// limit the number of go routines or tasks that can be run concurrently and this way we manage the memory
// one use could be if we are having a web scraper  and we scraping the data from the multiple websites then
// there can be a case where our go out of memory while spawning the go routines for each website, to avoid
// this issue we can use worker pool to limit the number of go routines that can be run concurrently.

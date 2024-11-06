package search

import (
	"abhi/search/db"
	"fmt"
	"time"
)

func RunEngine(){
	fmt.Println("started search engine crawl")
	defer fmt.Println("search engine crawl has finished")
	//get settings
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		fmt.Println("something went wrong about the settings")
		return
	}
	crawl := &db.CrawledUrl{}
	nextUrls, err := crawl.GetNextCrawledUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("something went wrong in geting next url")
		return
	}

	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		if !result.Success {
			err := next.UpdatedUrl(db.CrawledUrl{
				ID: next.ID,
				Url: next.Url,
				Success: false,
				CrawledDuration: result.CrawlData.CrawlTime,
				ResponseCode: crawl.ResponseCode,
				PageTitle: result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Headings: result.CrawlData.Headings,
				LastTested: &testedTime,
			})
			if err != nil {
				fmt.Println("somethtin gwent wrong updating a failed url")
			}
			continue
		}
		//success
		err := next.UpdatedUrl(db.CrawledUrl{
			ID: next.ID,
				Url: next.Url,
				Success: false,
				CrawledDuration: result.CrawlData.CrawlTime,
				ResponseCode: crawl.ResponseCode,
				PageTitle: result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Headings: result.CrawlData.Headings,
				LastTested: &testedTime,
		})
		if err != nil {
			fmt.Println("something went worng updating a succes url")
			fmt.Println(next.Url)
		}
		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	}
	//end of range
	if !settings.AddNew {
		return
	}
	//insert new urls 
	for _, newUrl := range newUrls {
		err := newUrl.Save()
		if err != nil {
			fmt.Println("somethtin wnet wring adding new url to the database")
		}
	}
	fmt.Printf("\n Added %d new urls to the databse ", len(newUrls))
}

func RunIndex() {
	fmt.Println("started search indexing...")
	defer fmt.Println("search indexing has finished")
	// Get index settings from DB
	crawled := &db.CrawledUrl{}
	// Get all urls that are not indexed
	notIndexed, err := crawled.GetNotIndexed()
	fmt.Println("not indexed urls: ", len(notIndexed))
	if err != nil {
		fmt.Println("something went wrong getting the not indexed urls")
		return
	}
	// Create a new index
	idx := make(Index)
	// Add the not indexed urls to the index
	idx.Add(notIndexed)
	// Save the index to the database
	searchIndex := &db.SearchIndex{}
	err = searchIndex.Save(idx, notIndexed)
	if err != nil {
		fmt.Println(err)
		fmt.Println("something went wrong saving the index")
		return
	}
	// Update the urls to be indexed=true
	err = crawled.SetIndexedTrue(notIndexed)
	if err != nil {
		fmt.Println("something went wrong updating the indexed urls")
		return
	}

}
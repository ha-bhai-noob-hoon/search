package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	ID              string
	Url             string
	Success         bool
	CrawledDuration time.Duration
	ResponseCode 	int
	PageTitle 		string
	PageDescrition 	string 
	Heading 		string 
	LastTested 		*time.Time
	Indexed 		bool
	CreatedAt 		*time.Time
	UpdatedAt 		time.Time
	DeletedAt 		gorm.DeletedAt
}

func(crawled *CrawledUrl) UpdatedUrl(input CrawledUrl) error {
	tx := DBConn.Select("url" , "success", "crawl_duration" , "response_code" , "page_title" , "page_description" , "headings" , "last_tested", "updated_at").Omit("created_at").Save(&input)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	return nil
} 

func (crawled *CrawledUrl) GetNextCrawledUrls(limit int) ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("last_tested IS NULL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

func (crawled *CrawledUrl) Save() error {
	tx := DBConn.Save(&crawled)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (crawled *CrawledUrl) GetNotIndexed() ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("indexed = ? AND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		fmt.Print(tx.Error)
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

func (crawled *CrawledUrl) SetIndexedTrue(urls []CrawledUrl) error {
	for _, url := range urls {
		url.Indexed = true
		tx := DBConn.Save(&url)
		if tx.Error != nil {
			fmt.Print(tx.Error)
			return tx.Error
		}
	}
	return nil
}
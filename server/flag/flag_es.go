package flag

import (
	"bufio"
	"fmt"
	"os"
	"server/model/elasticsearch"
	"server/service"
)

// Elasticsearch ES 索引操作
func Elasticsearch() error {
	esService := service.ServiceGroupApp.EsService

	indexExists, err := esService.IndexExists(elasticsearch.ArticleIndex())
	if err != nil {
		return err
	}
	if indexExists {
		fmt.Println("The index already exists.Do you want to delete the data and recreate the index? (y/n)")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "y":
			fmt.Println("Proceeding to delete the data and recreate the index...")
			if err := esService.IndexDelete(elasticsearch.ArticleIndex()); err != nil {
				return err
			}
		case "n":
			fmt.Println("Exiting the program.")
			os.Exit(0)
		default:
			fmt.Println("Invalid input. Please enter 'y' to delete and recreated the index,or 'n' to exit.")
			return Elasticsearch()
		}
	}
	return esService.IndexCreate(elasticsearch.ArticleIndex(), elasticsearch.ArticleMapping())
}

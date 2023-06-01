package operations

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fstanis/screenresolution"
	"github.com/reujab/wallpaper"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

func GetDesktopResolution() (int, int, error) {
	resolution := screenresolution.GetPrimary()
	if resolution.String() == "" {
		return 0, 0, fmt.Errorf("Could not retrieve resolution")
	}
	return resolution.Width, resolution.Height, nil
}

func GoogleImageSearch(apiKey, cx, query string) ([]*customsearch.Result, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}

	svc, err := customsearch.New(client)
	if err != nil {
		return nil, err
	}

	// Store the final results
	var allResults []*customsearch.Result

	// Set the number of results per page
	numResultsPerPage := 10

	// Total number of pages needed to collect 30 results
	totalPages := 1

	for i := 0; i < totalPages; i++ {
		startIndex := i*numResultsPerPage + 1
		searchCall := svc.Cse.List().Cx(cx).Q(query).SearchType("image").ImgSize("huge").Num(int64(numResultsPerPage)).Start(int64(startIndex))
		response, err := searchCall.Do()

		if err != nil {
			return nil, err
		}

		allResults = append(allResults, response.Items...)
	}

	return allResults, nil
}

func DownloadImage(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	imageData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SetWallpaper(filePath string) error {
	err := wallpaper.SetFromFile(filePath)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type apiResponse struct {
	Response int       `json:"response,omitempty"`
	StoryURL string    `json:"story_url,omitempty"`
	Error    *apiError `json:"error,omitempty"`
}

type apiError struct {
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorText string `json:"error_text,omitempty"`
}

const (
	apiErrorTypeMissing      = "type missing"
	apiErrorIDMissing        = "id missing"
	apiErrorBackdropMissing  = "backdrop_url is missing"
	apiErrorPosterMissing    = "poster_url is missing"
	apiErrorRatingMissing    = "rating is missing"
	apiErrorTitleMissing     = "title is missing"
	apiErrorStoryMissing     = "story_url is missing"
	apiErrorTypeWrong        = "type should be tv or movie"
	apiErrorBackdropWrong    = "backdrop_url should /somestring.jpg"
	apiErrorPosterWrong      = "poster_url should /somestring.jpg"
	apiErrorRatingWrongType  = "rating is not a number"
	apiErrorRatingWrongValue = "rating is less than 1 or more than 5"
	apiErrorIDWrong          = "id is not a number"
)

func handleResponse(response interface{}, isOK bool, writer http.ResponseWriter) {
	jsonResponse, _ := json.Marshal(response)
	writer.Header().Set("Content-Type", "application/json")
	if !isOK {
		writer.WriteHeader(http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
	_, err := writer.Write(jsonResponse)
	if err != nil {
		log.Println("error handling response")
	}
}

func downloadImage(filepath, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadImages(itemType string, id int, backdropURL, posterURL string) {
	folderPath := fmt.Sprintf("images/%s/%d/", itemType, id)
	err := os.MkdirAll(folderPath, 0777)
	if err != nil {
		log.Println(err.Error())
	}

	savePath := fmt.Sprintf("images/%s/%d%s", itemType, id, backdropURL)
	downloadPath := fmt.Sprintf("%s/%s%s", imagesPatch, "original", backdropURL)
	err = downloadImage(savePath, downloadPath)
	if err != nil {
		log.Println(err.Error())
	}
	savePath = fmt.Sprintf("images/%s/%d%s", itemType, id, posterURL)
	downloadPath = fmt.Sprintf("%s/%s%s", imagesPatch, "original", posterURL)
	err = downloadImage(savePath, downloadPath)
	if err != nil {
		log.Println(err.Error())
	}
}

func handleAPIError(errorCode int, errorText string, w http.ResponseWriter) {
	errorResp := apiError{
		ErrorCode: errorCode,
		ErrorText: errorText,
	}

	response := apiResponse{
		Response: 0,
		Error:    &errorResp,
	}

	handleResponse(response, false, w)
}

func handleAPIOk(w http.ResponseWriter) {
	response := apiResponse{
		Response: 1,
	}

	handleResponse(response, true, w)
}

func handleAPIUploaded(w http.ResponseWriter, story uploadResponse) {
	storyURL := fmt.Sprintf("https://vk.com/story%d_%d", story.Response.Story.OwnerID, story.Response.Story.ID)

	response := apiResponse{
		StoryURL: storyURL,
	}

	handleResponse(response, true, w)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request", r.RequestURI)
	defer r.Body.Close()

	query := r.URL.Query()

	typeQuery := query.Get("type")
	idQuery := query.Get("id")
	backdropURLQuery := query.Get("backdrop_url")
	posterURLQuery := query.Get("poster_url")
	ratingQuery := query.Get("rating")
	titleQuery := query.Get("title")
	storyURLQuery := query.Get("story_url")
	line1Query := query.Get("line1")
	line2Query := query.Get("line2")
	line3Query := query.Get("line3")

	reImageURL := regexp.MustCompile(`^\/\w+.jpg$`)

	log.Println(r.RemoteAddr, typeQuery, idQuery, backdropURLQuery, posterURLQuery, ratingQuery, titleQuery, storyURLQuery, line1Query, line2Query, line3Query)

	// http://localhost:3333/?type=movie&id=500&backdrop_url=/mMZRKb3NVo5ZeSPEIaNW9buLWQ0.jpg&poster_url=/adw6Lq9FiC9zjYEpOqfq03ituwp.jpg&rating=5&title=hey&story_url=https://vvv.ru&review=cool movie

	if typeQuery == "" {
		handleAPIError(123, apiErrorTypeMissing, w)
		return
	}

	if typeQuery != "tv" && typeQuery != "movie" {
		handleAPIError(123, apiErrorTypeWrong, w)
		return
	}

	if idQuery == "" {
		handleAPIError(123, apiErrorIDMissing, w)
		return
	}

	id, err := strconv.Atoi(idQuery)

	if err != nil {
		handleAPIError(123, apiErrorIDWrong, w)
		return
	}

	if backdropURLQuery == "" {
		handleAPIError(123, apiErrorBackdropMissing, w)
		return
	}

	if !reImageURL.MatchString(backdropURLQuery) {
		handleAPIError(123, apiErrorBackdropWrong, w)
		return
	}

	if posterURLQuery == "" {
		handleAPIError(123, apiErrorPosterMissing, w)
		return
	}

	if !reImageURL.MatchString(posterURLQuery) {
		handleAPIError(123, apiErrorPosterWrong, w)
		return
	}

	if ratingQuery == "" {
		handleAPIError(123, apiErrorRatingMissing, w)
		return
	}

	rating, err := strconv.Atoi(ratingQuery)

	if err != nil {
		handleAPIError(123, apiErrorRatingWrongType, w)
		return
	}

	if rating > 5 || rating < 1 {
		handleAPIError(123, apiErrorRatingWrongValue, w)
		return
	}

	if titleQuery == "" {
		handleAPIError(123, apiErrorTitleMissing, w)
		return
	}

	if storyURLQuery == "" && appconfig.Upload {
		handleAPIError(123, apiErrorStoryMissing, w)
		return
	}

	downloadImages(typeQuery, id, backdropURLQuery, posterURLQuery)
	posterPath := fmt.Sprintf("images/%s/%d%s", typeQuery, id, posterURLQuery)
	backdropPath := fmt.Sprintf("images/%s/%d%s", typeQuery, id, backdropURLQuery)
	filename := makeStory(titleQuery, posterPath, backdropPath, line1Query, line2Query, line3Query, rating)

	if appconfig.Upload {
		uploadResp := sendStoryToVK(filename, storyURLQuery)
		handleAPIUploaded(w, uploadResp)
	}

	handleAPIOk(w)
}

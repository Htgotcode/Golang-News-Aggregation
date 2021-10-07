package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//https://newsapi.org/docs/endpoints/top-headlines
var apiKey = "&apiKey=6996b0514b4449bfbb365846da95b15f"
var NewsResults1 = NewsResults{}
var r = gin.Default()

//JSON response to golang struct
type NewsResults struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Url         string    `json:"url"`
		UrlToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	}
}

type Search struct {
	Query string
}

func getTopHeadlines(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := &Search{
			Query: c.Request.FormValue("c"),
		}

		fmt.Println(search)

		//Constructing URL
		url := endpoint + search.Query + apiKey

		client := http.Client{
			Timeout: time.Second * 10, // Timeout after 10 seconds
		}

		//Fetch api from url provided
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}

		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		//Execute res.Body.Close() at the end of the function to avoid memory leak
		if res.Body != nil {
			defer res.Body.Close()
		}

		//Read body requested from url
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		//Unmarshel JSON data from fetched url
		NewsResults1 := NewsResults{}
		jsonErr := json.Unmarshal(body, &NewsResults1)
		if jsonErr != nil {
			panic(jsonErr)
		}

		//Load HTML templates from folder
		r.LoadHTMLGlob("templates/*")

		//Build HTML for /topheadlines using mutliple templates

		//Load head & header of HTML
		c.HTML(http.StatusOK, "header_headlines.tmpl.html", gin.H{
			"title": "News Aggregation/" + search.Query,
			"query": search.Query,
		})

		//Load and duplicate article format based on the amount of articles pulled fomr the API
		c.HTML(http.StatusOK, "articles_container.tmpl.html", gin.H{})
		for _, article := range NewsResults1.Articles {
			c.HTML(http.StatusOK, "articles.tmpl.html", gin.H{
				//Send JSON data to HTML
				"articleSource":      article.Source.Name,
				"articlePubDate":     article.PublishedAt.Format("January 2, 2006"),
				"articleTitle":       article.Title,
				"articleDescription": article.Description,
				"articleImage":       article.UrlToImage,
				"articleUrl":         article.Url,
			})
		}

		//Load footer of HTML
		c.HTML(http.StatusOK, "footer.tmpl.html", gin.H{})
	}
}

func getEverything(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {

		search := &Search{
			Query: c.Query("q"),
		}

		//Constructing URL
		url := endpoint + search.Query + apiKey

		client := http.Client{
			Timeout: time.Second * 10, // Timeout after 10 seconds
		}

		//Fetch api from url provided
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}

		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		//Execute res.Body.Close() at the end of the function to avoid memory leak
		if res.Body != nil {
			defer res.Body.Close()
		}

		//Read body requested from url
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		//Unmarshel JSON data from fetched url
		NewsResults1 := NewsResults{}
		jsonErr := json.Unmarshal(body, &NewsResults1)
		if jsonErr != nil {
			panic(jsonErr)
		}

		//Load HTML templates from folder
		r.LoadHTMLGlob("templates/*")

		//Build HTML for /topheadlines using mutliple templates

		//Load head & header of HTML
		c.HTML(http.StatusOK, "header_everything.tmpl.html", gin.H{
			"title": "News Aggregation/" + search.Query,
			//Call function again to refresh results if server isn't restarted
			"refresh": getEverything,
			"query":   search.Query,
		})

		//Load and duplicate article format based on the amount of articles pulled fomr the API
		c.HTML(http.StatusOK, "articles_container.tmpl.html", gin.H{})
		for _, article := range NewsResults1.Articles {
			c.HTML(http.StatusOK, "articles.tmpl.html", gin.H{
				//Send JSON data to HTML
				"articleSource":      article.Source.Name,
				"articlePubDate":     article.PublishedAt.Format("January 2, 2006"),
				"articleTitle":       article.Title,
				"articleDescription": article.Description,
				"articleImage":       article.UrlToImage,
				"articleUrl":         article.Url,
			})
		}
		//Load footer of HTML
		c.HTML(http.StatusOK, "footer.tmpl.html", gin.H{})
	}
}

func main() {

	r.GET("/topheadlines", getTopHeadlines("https://newsapi.org/v2/top-headlines?country="))
	r.GET("/everything", getEverything("https://newsapi.org/v2/everything?q="))
	//handler for static files
	r.Static("css", "../news_aggregation/css")

	r.Run(":8080")
}

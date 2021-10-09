package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

//https://newsapi.org/docs/endpoints/top-headlines
var apiKey string
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
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Search struct {
	Query   string
	Country string
}

//read apikey from text file
func readAPIKey() {
	file, err := os.Open("apikey.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	apiKey = scanner.Text()

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

//Construct /topheadlines
func getTopHeadlines(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := &Search{
			Query:   c.Request.FormValue("c"),
			Country: "",
		}

		//Switch title based on dropbox selection
		switch search.Query {
		case "za":
			search.Country = "South Africa"
		case "ae":
			search.Country = "United Arab Emirates"
		case "ar":
			search.Country = "Argentina"
		case "at":
			search.Country = "Austria"
		case "au":
			search.Country = "Australia"
		case "be":
			search.Country = "Belgium"
		case "bg":
			search.Country = "Bulgaria"
		case "ca":
			search.Country = "Canada"
		case "ch":
			search.Country = "Switzerland"
		case "cn":
			search.Country = "China"
		case "co":
			search.Country = "Colombia"
		case "cu":
			search.Country = "Cuba"
		case "cz":
			search.Country = "Czechia"
		case "de":
			search.Country = "Germany"
		case "eg":
			search.Country = "Egypt"
		case "fr":
			search.Country = "France"
		case "gb":
			search.Country = "United Kingdom"
		case "gr":
			search.Country = "Greece"
		case "hk":
			search.Country = "Hong Kong"
		case "hu":
			search.Country = "Hungary"
		default:
			search.Country = ""
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

		if NewsResults1.Status == "error" {
			c.HTML(http.StatusOK, "error.tmpl.html", gin.H{
				"title":   "News Aggregation | error",
				"status":  NewsResults1.Status,
				"code":    NewsResults1.Code,
				"message": NewsResults1.Message,
			})
		} else {
			//Build HTML for /topheadlines using mutliple templates
			//Load head & header of HTML
			c.HTML(http.StatusOK, "header_headlines.tmpl.html", gin.H{
				"title":   "News Aggregation | " + search.Query,
				"country": search.Country,
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
}

//Construct /everything
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

		if NewsResults1.Status == "error" {
			c.HTML(http.StatusOK, "error.tmpl.html", gin.H{
				"title":   "News Aggregation | error",
				"status":  NewsResults1.Status,
				"code":    NewsResults1.Code,
				"message": NewsResults1.Message,
			})
		} else {
			//Build HTML for /topheadlines using mutliple templates
			//Load head & header of HTML
			c.HTML(http.StatusOK, "header_everything.tmpl.html", gin.H{
				"title": "News Aggregation | " + search.Query,
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
}

func main() {
	readAPIKey()
	r.GET("/topheadlines", getTopHeadlines("https://newsapi.org/v2/top-headlines?country="))
	r.GET("/everything", getEverything("https://newsapi.org/v2/everything?q="))

	//handler for static files
	r.Static("favicon-16x16.ico", "../Golang-News-Aggregation/")
	r.Static("js", "../Golang-News-Aggregation/js")

	r.Run(":8080")
}

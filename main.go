package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

//API recieved from https://newsapi.org/
var apiKey string
var NewsResults1 = NewsResults{}
var r = gin.Default()

//Building newsapi.org JSON article response to struct.
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

//Struct to handle user query input.
type Search struct {
	Query   string
	Country string
}

/*Read apikey from text file. Exists to not upload my personal API to github.com
Snippet credit - https://www.golang101.com/how-tos/how-to-read-text-file-golang/
*/
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

//Constructs /topheadlines by accepting a country selection box form and building HTML through templates.
func getTopHeadlines(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {

		//Assign search to an instance of Search struct from HTML selection box form values.
		search := &Search{
			Query:   c.Request.FormValue("c"),
			Country: "",
		}

		//Load HTML templates from folder
		r.LoadHTMLGlob("templates/*")

		//First time load to avoid empty API call
		if search.Query == "" {
			c.HTML(http.StatusOK, "header_headlines.tmpl.html", gin.H{
				"title":   "News Aggregation",
				"country": search.Country,
				"query":   search.Query,
			})

			//Build footer HTML
			c.HTML(http.StatusOK, "footer.tmpl.html", gin.H{})

		} else {
			//Switch title based on selection box query.
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
			case "id":
				search.Country = "Indonesia"
			case "ru":
				search.Country = "Russian Federation"
			case "us":
				search.Country = "United States of America"
			default:
				search.Country = ""
			}

			/*Construct the URL using the manaully inputted endpoint, user inputted query, and file read API key.
			The net/http package creates a client and fetches the API's news response body in JSON format.
			Then we assign it to our NewsResults struct.
			*/
			url := endpoint + search.Query + apiKey

			client := http.Client{
				Timeout: time.Second * 10,
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				panic(err)
			}

			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			//Execute res.Body.Close() at the end of the function to avoid memory leak.
			if res.Body != nil {
				defer res.Body.Close()
			}

			//Read body requested from url
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			NewsResults1 := NewsResults{}
			jsonErr := json.Unmarshal(body, &NewsResults1)
			if jsonErr != nil {
				panic(jsonErr)
			}

			//Call HTML for /topheadlines using mutliple templates
			//Call head & header of HTML
			c.HTML(http.StatusOK, "header_headlines.tmpl.html", gin.H{
				"title":   "News Aggregation | " + search.Country,
				"country": search.Country,
				"query":   search.Query,
				"status":  NewsResults1.Status,
				"code":    NewsResults1.Code,
				"message": NewsResults1.Message,
			})
			//Call and duplicate article format based on the amount of articles pulled fomr the API
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
			//Call footer HTML
			c.HTML(http.StatusOK, "footer.tmpl.html", gin.H{})

		}
	}
}

/*Constructs /everything by accepting a keyword query and building HTML through templates.
The function is similar to /getTopHeadlines but recieves it query in a different form.
Comments will not be as verbose.
*/
func getEverything(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {

		//Read text box from HTML and assign it to search.Query
		search := &Search{
			Query: c.Query("q"),
		}

		r.LoadHTMLGlob("templates/*")

		if search.Query == "" {
			c.HTML(http.StatusOK, "header_everything.tmpl.html", gin.H{
				"title": "News Aggregation",
				"query": search.Query,
			})
			c.HTML(http.StatusOK, "footer.tmpl.html", gin.H{})

		} else {

			url := endpoint + search.Query + apiKey

			client := http.Client{
				Timeout: time.Second * 10,
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				panic(err)
			}

			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			if res.Body != nil {
				defer res.Body.Close()
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			NewsResults1 := NewsResults{}
			jsonErr := json.Unmarshal(body, &NewsResults1)
			if jsonErr != nil {
				panic(jsonErr)
			}

			c.HTML(http.StatusOK, "header_everything.tmpl.html", gin.H{
				"title":   "News Aggregation | " + search.Query,
				"query":   search.Query,
				"status":  NewsResults1.Status,
				"code":    NewsResults1.Code,
				"message": NewsResults1.Message,
			})
			c.HTML(http.StatusOK, "articles_container.tmpl.html", gin.H{})
			for _, article := range NewsResults1.Articles {
				c.HTML(http.StatusOK, "articles.tmpl.html", gin.H{
					"articleSource":      article.Source.Name,
					"articlePubDate":     article.PublishedAt.Format("January 2, 2006"),
					"articleTitle":       article.Title,
					"articleDescription": article.Description,
					"articleImage":       article.UrlToImage,
					"articleUrl":         article.Url,
				})
			}
			c.HTML(http.StatusOK, "footer.tmpl.html", gin.H{})
		}
	}
}

func main() {
	readAPIKey()
	r.GET("/topheadlines", getTopHeadlines("https://newsapi.org/v2/top-headlines?country="))
	r.GET("/everything", getEverything("https://newsapi.org/v2/everything?pagesize=20&q="))

	//Gin middleware to support favicon.
	r.Use(favicon.New("./favicon.ico"))

	r.Run(":8080")
}

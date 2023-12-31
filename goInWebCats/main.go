package main

import (
	"encoding/json"
	"fmt"
	"log"

	// "io"
	"net/http"
	"text/template"
)

type Context struct{
	w http.ResponseWriter
	req *http.Request
	t *Template
}

type Template struct {
	templates *template.Template
}

type User struct{
	Username string
	IsUser 	 bool
	Age      int
}

func (c *Context) render(name string, data any) error{
	return c.t.templates.ExecuteTemplate(c.w, name, data)
}

func (c *Context)writeJson( code int, v any) error{
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(code)
	return json.NewEncoder(c.w).Encode(v)
}

func main(){
	facts, err := fetchCatFacts()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(facts)

	t := &Template{
		templates: template.Must(template.ParseGlob("www/*.html")),
	}
	http.HandleFunc("/", makeAPIFunc(handleHome, t))
	http.HandleFunc("/api/user", makeAPIFunc(handleUser, t))
	http.ListenAndServe(": 3000", nil)

}

type apiFunc func(c *Context) error

func makeAPIFunc(fn apiFunc, t *Template) http.HandlerFunc {
	return func( w http.ResponseWriter, r *http.Request){
		
		ctx := &Context{
			t: t,
			w: w, 
			req: r,
		}
		if err := fn(ctx); err != nil{
			ctx.writeJson(http.StatusInternalServerError, map[string] string{"error": err.Error()})
		}
	}
}
// func showCatsFacts(w http.ResponseWriter, r *http.Request) error
func showCatsFacts(c *Context) error{
	facts, err := fetchCatFacts()
	if err != nil{
		return err
	}
	return c.render("fact.html", facts)
}

func handleHome(c *Context) error{
	user := User{
		Username: "Dio",
		IsUser: true,
		Age: 137,
	}
	return c.render( "cats.html", user)
}


func handleUser( c *Context) error{
	return c.writeJson( http.StatusOK, map[string] string{"message":"hello there, Banki"})
}

type catFact struct{
	Fact string `json:"fact"`
}

type CFResp struct{
	Data []catFact `json:"data"`
}
func fetchCatFacts()([]catFact, error){
	resp, err := http.Get("http://catfact.ninja/facts")
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
//	var facts []catFact
	var res CFResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err!= nil{
		return nil, err
	}
	return res.Data, nil
}
package restcl

import (
	"net/http"
	"testing"
)

type Auth struct{}

func (a Auth) Modify(req *http.Request) {
}

func Test_RestClient(t *testing.T) {
	r := NewRest()
	r.SetPrefix("http://jsonplaceholder.typicode.com")
	r.Create("/posts").SetMethod("GET").Build("posts")
	r.Create("/posts").SetMethod("POST").Build("createpost")
	r.Create("/posts/{id}").SetMethod("GET").Build("post")
	r.Create("/posts/{id}").SetMethod("DELETE").Build("deletepost")

	ep := r.Get("post")
	if ep.Url != "http://jsonplaceholder.typicode.com/posts/{id}" {
		t.Errorf("url for posts is wrong: %v", ep.Url)
	}

	if ep.Method != "GET" {
		t.Errorf("method should be \"POST\" is \"%v\"", ep.Method)
	}

	ep = r.Get("post").SetParams("id", 1)
	ep.buildUrl()

	if ep.BuildUrl != "http://jsonplaceholder.typicode.com/posts/1" {
		t.Error(ep.BuildUrl)
	}

}

func Test_RestClientGet(t *testing.T) {
	r := NewRest()
	r.SetPrefix("http://jsonplaceholder.typicode.com").Use(Auth{})
	r.Create("/posts").SetMethod("POST").Build("createpost")

	var it map[string]interface{}

	err := r.Get("createpost").SetBody("userId", 1, "title", "foo", "body", "bar").Exec(&it)
	if err != nil {
		t.Error(err.Error())
	}

	if it["id"].(float64) != 101 {
		t.Errorf("%+v", it["id"])
	}
}

func Test_Use(t *testing.T) {
	r := NewRest()
	r.SetPrefix("http://jsonplaceholder.typicode.com").Use(Auth{})
	r.Create("/posts").SetMethod("POST").Use(Auth{}).Build("createpost")

	ep := r.Get("createpost")

	if len(ep.interceptor) != 2 {
		t.Errorf("should be 2 is %v", len(ep.interceptor))
	}
}

func Test_Multi(t *testing.T) {
	r := NewRest()
	r.SetPrefix("http://jsonplaceholder.typicode.com").Use(Auth{})
	r.Create("/posts/{id}").SetMethod("GET").Build("post")

	ep1 := r.Get("post").SetParams("id", 2)

	ep2 := r.Get("post").SetParams("id", 1)

	if ep1 == ep2 {
		t.Error("should be 2 new endpoint objects")
	}

}

package restcl

import (
	"fmt"
	enc "net/url"
	"strings"

	"github.com/Kemonozume/httpcl"
)

//used to query the rest service
type RestEndPoint struct {
	Url         string
	BuildUrl    string
	Method      string
	Params      map[string]string
	Body        []interface{}
	interceptor []Intercept
}

//sets the parameters used to build the url
func (rep *RestEndPoint) SetParams(b ...interface{}) *RestEndPoint {
	rep.Params = rep.iToMap(b)
	return rep
}

//sets the request body
func (rep *RestEndPoint) SetBody(b ...interface{}) *RestEndPoint {
	rep.Body = b
	return rep
}

//executes the request and transforms the response into a json interface
func (rep *RestEndPoint) Exec(a interface{}) error {
	rep.buildUrl()
	cl := &httpcl.ClientBuilder{}
	cl.Method = rep.Method
	cl.Url = rep.BuildUrl
	cl.Redirect = true
	cl.Body = rep.Body
	c := cl.Build()
	if c.Error != nil {
		return c.Error
	}
	if len(rep.interceptor) > 0 {
		for _, inter := range rep.interceptor {
			inter.Modify(c.GetRequest())
		}
	}
	return c.DoTransform(httpcl.TransformToJson, a)
}

//builds the url
func (rep *RestEndPoint) buildUrl() {
	rep.BuildUrl = rep.Url
	for key, value := range rep.Params {
		rep.BuildUrl = strings.Replace(rep.BuildUrl, fmt.Sprintf("{%s}", key), enc.QueryEscape(value), 1)
	}
}

func (rep *RestEndPoint) iToMap(b []interface{}) map[string]string {
	if len(b) == 1 {
		m, ok := b[0].(map[string]string)
		if ok {
			return m
		}
	} else {
		if len(b)%2 == 0 {
			m := make(map[string]string)
			for i := 0; i < len(b)-1; i += 2 {
				m[b[i].(string)] = fmt.Sprintf("%v", b[i+1])
			}
			return m
		}
	}
	m := make(map[string]string)
	return m
}

package libra

import (
	"fmt"
	"github.com/zhuCheer/libra/balancer"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestProxyStart(t *testing.T) {
	proxy := NewHttpProxySrv("127.0.0.1:5001", "roundrobin", nil)
	proxy.Scheme = ""
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("catch error", r)
			}
		}()
		proxy.Start()
	}()

	time.Sleep(2 * time.Second)
}

func TestProxySrvFun(t *testing.T) {
	proxy := NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)
	if _, ok := proxy.balancer.(*balancer.RoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv loadType have an error #1")
	}
	proxy.ResetCustomHeader(map[string]string{"X-LIBRA": "the smart ReverseProxy"})

	header, ok := proxy.customHeader["X-LIBRA"]
	if ok == false || header != "the smart ReverseProxy" {
		t.Error("ResetCustomHeader func have an error #2")
	}

	proxy.ChangeLoadType("wroundrobin")
	if _, ok := proxy.balancer.(*balancer.WRoundRobinLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #3")
	}

	proxy.ChangeLoadType("random")
	if _, ok := proxy.balancer.(*balancer.RandomLoad); ok == false {
		t.Error("NewHttpProxySrv ChangeLoadType have an error #4")
	}

	b := proxy.GetBalancer()
	if b != proxy.balancer {
		t.Error("NewHttpProxySrv GetBalancer have an error #5")
	}
}

func TestReverseProxySrv(t *testing.T) {
	targetHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing ReverseProxySrv")
	}))
	defer targetHttpServer.Close()

	proxy := NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)
	reverseProxy := proxy.dynamicReverseProxy()
	proxy.ResetCustomHeader(map[string]string{"httptest": "01023"})
	ts := httptest.NewServer(proxy.httpMiddleware(reverseProxy))
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL, strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	req.Header.Del("User-Agent")

	res, err := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode != 500 {
		t.Error("ReverseProxySrv have an error #1")
	}
	testHeader := res.Header.Get("httptest")
	if testHeader != "01023" {
		t.Error("ReverseProxySrv have an error #2")
	}

	tsUrl, _ := url.Parse(ts.URL)
	targetHttpUrl, _ := url.Parse(targetHttpServer.URL)
	proxy.balancer.AddAddr(tsUrl.Host, targetHttpUrl.Host, 0)
	res, err = http.Get(ts.URL + "?abc=123")
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("ReverseProxySrv have an error #3")
	}

	if res.Request.URL.RawQuery != "abc=123" {
		t.Error("ReverseProxySrv have an error #4")
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Error(err)
	}
	if string(greeting) != "testing ReverseProxySrv" {
		t.Error("ReverseProxySrv have an error #5")
	}
}

func TestReverseProxySrvUnStart(t *testing.T) {
	targetHttpServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing ReverseProxySrvUnStart")
	}))
	defer targetHttpServer.Close()

	proxy := NewHttpProxySrv("127.0.0.1:5001", "roundrobin", nil)
	reverseProxy := proxy.dynamicReverseProxy()

	ts := httptest.NewServer(proxy.httpMiddleware(reverseProxy))
	defer ts.Close()

	tsUrl, _ := url.Parse(ts.URL)
	targetHttpUrl, _ := url.Parse(targetHttpServer.URL)
	proxy.balancer.AddAddr(tsUrl.Host, targetHttpUrl.Host, 0)
	res, _ := http.Get(ts.URL)

	if res.StatusCode != 502 {
		t.Error("ReverseProxySrv have an error(UnStart) #1")
	}
}

func TestReverseProxySrvNotFound(t *testing.T) {
	targetHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "testing ReverseProxySrv NotFound")
	}))
	defer targetHttpServer.Close()

	proxy := NewHttpProxySrv("127.0.0.1:5002", "roundrobin", nil)
	reverseProxy := proxy.dynamicReverseProxy()

	ts := httptest.NewServer(proxy.httpMiddleware(reverseProxy))
	defer ts.Close()

	tsUrl, _ := url.Parse(ts.URL)
	targetHttpUrl, _ := url.Parse(targetHttpServer.URL)
	proxy.balancer.AddAddr(tsUrl.Host, targetHttpUrl.Host, 0)
	res, _ := http.Get(ts.URL)

	if res.StatusCode != 404 {
		t.Error("ReverseProxySrv have an error(NotFound) #1")
	}
}

func TestGetBalancerByLoadType(t *testing.T) {

	b := getBalancerByLoadType("xxx")
	if _, ok := b.(*balancer.RandomLoad); ok == false {
		t.Error("getBalancerByLoadType func have an error #1")
	}

	b = getBalancerByLoadType("random")
	if _, ok := b.(*balancer.RandomLoad); ok == false {
		t.Error("getBalancerByLoadType func have an error #2")
	}

	b = getBalancerByLoadType("roundrobin")
	if _, ok := b.(*balancer.RoundRobinLoad); ok == false {
		t.Error("getBalancerByLoadType func have an error #2")
	}

	b = getBalancerByLoadType("wroundrobin")
	if _, ok := b.(*balancer.WRoundRobinLoad); ok == false {
		t.Error("getBalancerByLoadType func have an error #2")
	}
}

func TestSingleJoiningSlash(t *testing.T) {
	target, _ := url.Parse("http://192.168.1.100/")
	path := singleJoiningSlash(target.Path, "/")
	if path != "/" {
		t.Error("singleJoiningSlash func have an error #1")
	}

	target, _ = url.Parse("http://192.168.1.100/abc")
	path = singleJoiningSlash(target.Path, "/")
	if path != "/abc" {
		t.Error("singleJoiningSlash func have an error #2")
	}

	target, _ = url.Parse("http://192.168.1.100/abc")
	path = singleJoiningSlash(target.Path, "/efg")
	if path != "/abc/efg" {
		t.Error("singleJoiningSlash func have an error #3")
	}

	target, _ = url.Parse("http://192.168.1.100/abc/")
	path = singleJoiningSlash(target.Path, "/efg")
	if path != "/abc/efg" {
		t.Error("singleJoiningSlash func have an error #4")
	}

	target, _ = url.Parse("http://192.168.1.100/abc")
	path = singleJoiningSlash(target.Path, "efg")
	if path != "/abc/efg" {
		t.Error("singleJoiningSlash func have an error #5")
	}

}

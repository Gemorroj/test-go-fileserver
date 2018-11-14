package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

func main() {
	if err := fasthttp.ListenAndServeTLS(":8181", "ssl-cert-snakeoil.pem", "ssl-cert-snakeoil.key", requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestParser(path string) (encryptedID string, name string) {
	s := strings.Split(path, "/")
	if len(s) != 3 {
		return
	}

	encryptedID, name = s[1], s[2]
	return
}

func doRequest(url string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{} // вынести в конекшен пулл
	client.Do(req, resp)

	return resp
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path()[:])
	encryptedID, name := requestParser(path)

	if encryptedID == "" || name == "" {
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	video := doRequest("https://www.youtube.com/watch?v=A8kSF_TShzQ")

	fmt.Fprintf(ctx, "video: %q\n\n", video)

	fmt.Fprintf(ctx, "Id: %q, name: %q\n", encryptedID, name)

	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	//ctx.Response.Header.Set("X-My-Header", "my-header-value")
	//ctx.Response.SetBodyStream()

	// Set cookies
	//var c fasthttp.Cookie
	//c.SetKey("cookie-name")
	//c.SetValue("cookie-value")
	//ctx.Response.Header.SetCookie(&c)
}

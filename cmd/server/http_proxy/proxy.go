package http_proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/esonhugh/proxyinbrowser/cmd/server/define"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type proxy struct {
	Conn            *websocket.Conn
	ResponseChannel chan define.RelayCommandResp
}

func (p *proxy) corsMiddleware(wr http.ResponseWriter, req *http.Request) bool {
	if req.Method == http.MethodOptions {
		log.Infoln("Triggered by OPTIONS method")

		var methodStrings = strings.Split("GET,POST,PUT,DELETE,PATCH,OPTIONS", ",")
		if str := req.Header.Get("Access-Control-Request-Method"); str != "" {
			methodStrings = req.Header.Values("Access-Control-Request-Method")
		}
		var headerStrings = []string{"Content-Type"}
		if str := req.Header.Get("Access-Control-Request-Headers"); str != "" {
			headerStrings = req.Header.Values("Access-Control-Request-Headers")
		}
		var Origin = "*"
		if str := req.Header.Get("Origin"); str != "" {
			Origin = str
			wr.Header().Add("Access-Control-Allow-Credentials", "true")
		}
		wr.Header().Add("Access-Control-Allow-Origin", Origin)

		for _, method := range methodStrings {
			wr.Header().Add("Access-Control-Allow-Methods", method)
		}
		for _, header := range headerStrings {
			wr.Header().Add("Access-Control-Allow-Headers", header)
		}
		wr.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func (p *proxy) hostBlacklist(wr http.ResponseWriter, req *http.Request) bool {
	if strings.HasSuffix(req.Host, ".googleapis.com") {
		wr.WriteHeader(http.StatusForbidden)
		return true
	}
	return false
}

func drop(target string, filterer []string) bool {
	for _, v := range filterer {
		if strings.HasPrefix(strings.ToLower(target), strings.ToLower(v)) {
			return true
		}
	}
	return false
}

func (p *proxy) filterHeader(req *http.Request) map[string]string {
	var dataString = make(map[string]string)
	droplist := []string{
		"User-Agent",
		"Upgrade-Insecure-Requests",
		"Cache-Control",
		"Pragma",
		"Priority",
		"Cache-Control",
		"Origin",
	}

	for k, v := range req.Header {
		if drop(k, droplist) {
			continue
		}
		dataString[k] = strings.Join(v, ",")
	}
	return dataString
}

func (p *proxy) prepareURL(req *http.Request) string {
	var url string
	if strings.HasPrefix(req.RequestURI, "http") { // if it is a full URL (http:// or https://)
		url = req.RequestURI //  "//"+ strings.Replace(strings.Replace(req.RequestURI, "http://", "", 1), "https://", "", 1)
	} else {
		url = req.URL.String()
	}
	return url
}

func (p *proxy) readReqBody(req *http.Request) (string, bool) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
		return "", false
	}
	return string(data), true
}

func (p *proxy) preJudgeCORS(req *http.Request) string {
	if req.Header.Get("Sec-Fetch-Mode") == "cors" {
		return "cors"
	} else if req.Header.Get("Sec-Fetch-Mode") == "no-cors" {
		return "no-cors"
	}
	if req.Header.Get("Origin") != "" { // has Origin and not empty
		return "cors"
	} else {
		return "no-cors"
	}
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Infoln("Received request: ", req.Method, req.URL.String(), req.RequestURI)

	// fast response CORS
	if p.corsMiddleware(wr, req) {
		return
	}

	if p.hostBlacklist(wr, req) {
		return
	}

	data, ok := p.readReqBody(req)
	if !ok {
		return
	}

	// generate TaskId
	var taskId string = uuid.New().String()

	FetchCommand := define.RelayCommand{
		CommandId: taskId,
		CommandDetail: define.Fetch{
			Url: p.prepareURL(req),
			Option: define.FetchOption{
				Method:  req.Method,
				Body:    data,
				Headers: p.filterHeader(req),
			},
		},
	}

	FetchCommand.CommandDetail.Option.Mode = p.preJudgeCORS(req)

	if FetchCommand.SendTo(p.Conn) != nil {
		log.Error("Error while sending FetchCommand")
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		return
	}

	for resp := range p.ResponseChannel {
		if resp.CommandId == taskId {
			r := resp.CommandResult
			if r.Error != "" {
				http.Error(wr, r.Error, http.StatusInternalServerError)
				break
			}
			var headerkeys []string
			for _, header := range r.Response.Headers {
				if len(header) >= 1 {
					droplist := []string{
						"Strict-Transport-Security",           // no TLS strict
						"Content-Security-Policy",             // no CSP
						"Content-Security-Policy-Report-Only", // no CSP
						"Report-To",                           // no CSP
						"Content-Encoding",                    // gzip will auto ungziped by victim browser
						"Content-Length",                      // need reCalc
						"Access-Control-Allow-Origin",         // no cors
					}
					if drop(header[0], droplist) {
						// ban HSTS header to force HTTPS and no CSP
						continue
					}
					wr.Header().Add(header[0], header[1])
					headerkeys = append(headerkeys, header[0])
				}
			}
			// reCalc the response length
			headerkeys = append(headerkeys, "Content-Length",
				"Access-Control-Allow-Origin", "Access-Control-Allow-Headers",
				"Access-Control-Expose-Headers", "Access-Control-Allow-Credentials")
			// Allow other header we will add.
			wr.Header().Add("Content-Length", fmt.Sprintf("%v", len(r.Response.Text)))
			if p.preJudgeCORS(req) == "cors" {
				wr.Header().Add("Access-Control-Expose-Headers", strings.Join(headerkeys, ","))
				wr.Header().Add("Access-Control-Allow-Headers", strings.Join(headerkeys, ","))
				wr.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
				wr.Header().Add("Access-Control-Allow-Credentials", "true")
			}
			/*
				if r.Response.FinalUrl != url {
					http.Redirect(wr, req, r.Response.FinalUrl, r.Response.Status)
					break
				}
			*/
			url := p.prepareURL(req)
			if r.Response.FinalUrl != url {
				log.Debugf("Request final url: %s, but request url is %s", r.Response.FinalUrl, url)
			}

			wr.WriteHeader(r.Response.Status)
			// pilog.Debugf("Response: %s", r.Response.Text)
			res := strings.NewReader(r.Response.Text)
			io.Copy(wr, res)
			break
		} else {
			p.ResponseChannel <- resp
		}
	}
}

func createHTTPProxy(TargetConn *websocket.Conn, rch chan define.RelayCommandResp) *proxy {
	handler := &proxy{
		Conn:            TargetConn,
		ResponseChannel: rch,
	}
	return handler
}

//
//  server/send_form_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/push"
)

func NewSendFormWS() *restful.WebService {
	ws := new(restful.WebService).Path("/v1/send_form")
	ws.Route(ws.GET("").To(sendForm))
	return ws
}

// Checks if the server is alive. This is useful for monitoring tools, load-balancers and automated scripts.
func sendForm(request *restful.Request, response *restful.Response) {
	htlm := "<!DOCTYPE html>"
	htlm += "<html lang=\"en\">"
	htlm += formHeader()
	htlm += formBody()
	htlm += "</html>"
	response.Write([]byte(htlm))
}

func formHeader() string {
	headerHtlm := "<head>"
	headerHtlm += "<meta charset=\"utf-8\">"
	headerHtlm += "<meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">"
	headerHtlm += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">"
	headerHtlm += fmt.Sprintf("<title>%s - %s</title>", config.ServiceName, config.ServiceVersion)

	// Bootstrap
	headerHtlm += "<link rel=\"stylesheet\" href=\"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css\" "
	headerHtlm += "integrity=\"sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u\" "
	headerHtlm += "crossorigin=\"anonymous\">"

	headerHtlm += "<script src=\"https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js\"></script>"
	headerHtlm += "<script src=\"https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js\"></script>"

	headerHtlm += "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/1.11.0/jquery.min.js\"></script>"

	headerHtlm += "<script src=\"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js\""
	headerHtlm += " integrity=\"sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa\""
	headerHtlm += " crossorigin=\"anonymous\"></script>"

	headerHtlm += "<style type=\"text/css\">"
	headerHtlm += ".container {"
	headerHtlm += "  padding-top: 5px;"
	headerHtlm += "}"

	headerHtlm += ".title {"
	headerHtlm += "  text-align: center;"
	headerHtlm += "}"

	headerHtlm += "</style>"
	headerHtlm += "</head>"

	return headerHtlm
}

func formBody() string {
	bodyHtlm := "<body>"
	bodyHtlm += fmt.Sprintf("<h1 class=\"title\">%s</h1>", config.ServiceName)

	bodyHtlm += "<div class=\"container\">"
	bodyHtlm += "<form>"

	bodyHtlm += "<div class=\"form-group\">"
	bodyHtlm += "<label for=\"sender_id_label\">Sender ID</label>"
	bodyHtlm += "<select class=\"form-control\" id=\"exampleSelect1\">"
	bodyHtlm += fmt.Sprintf("<option>%s</option>", push.ApnsSenderID)
	bodyHtlm += fmt.Sprintf("<option>%s</option>", push.GcmSenderID)
	bodyHtlm += fmt.Sprintf("<option>%s</option>", push.SafariSenderID)
	bodyHtlm += fmt.Sprintf("<option>%s</option>", push.FirefoxSenderID)
	bodyHtlm += fmt.Sprintf("<option>%s</option>", push.ChromeSenderID)
	bodyHtlm += "</select>"
	bodyHtlm += "</div>"

	bodyHtlm += "<div id=\"to_div\" class=\"form-group\">"
	bodyHtlm += "<label for=\"example-text-input\" class=\"col-2 col-form-label\">To</label>"
	bodyHtlm += "<div class=\"col-10\">"
	bodyHtlm += "<input class=\"form-control\" type=\"text\" placeholder=\"Device token\" id=\"example-text-input\">"
	bodyHtlm += "</div>"
	bodyHtlm += "</div>"

	bodyHtlm += "</form>"
	bodyHtlm += "</div>"
	bodyHtlm += "</body>"
	return bodyHtlm
}

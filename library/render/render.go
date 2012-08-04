// This package is called to render the template.
package render

import (
	"appengine"
	"bytes"
	"fmt"
	"net/http"
	"text/template" // WARNING: template/html removes the <!--[if lt IE 9]>, so we're using text/template - but this is just applied over the header and the footer so it should be safe to use.
)

/*
 * Renders a template (passedTemplate) and inserts the header and footer
 */
func Render(w http.ResponseWriter, r *http.Request, passedTemplate *bytes.Buffer, Statuscode ...int) {
	// Check if we are on "apn.statuscode.ch"
	if r.URL.Host == "apn.statuscode.ch" || appengine.IsDevAppServer() == true {
		if len(Statuscode) == 1 {
			w.WriteHeader(Statuscode[0])
		}

		// Header
		template.Must(template.ParseFiles("templates/header.html")).Execute(w, nil)

		// Now add the passedTemplate
		fmt.Fprintf(w, "%s", string(passedTemplate.Bytes())) // %s = the uninterpreted bytes of the string or slice

		// And now we execute the footer
		template.Must(template.ParseFiles("templates/footer.html")).Execute(w, nil)
	} else {
		// No, let's redirect the user to www.statuscode.ch
		http.Redirect(w, r, "https://apn.statuscode.ch/", http.StatusMovedPermanently)
	}
}

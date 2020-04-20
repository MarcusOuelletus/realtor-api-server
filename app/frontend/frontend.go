package frontend

import (
	"net/http"

	"github.com/MarcusOuelletus/rets-server/database"

	"github.com/golang/glog"
)

type frontend struct{}

func Frontend(w http.ResponseWriter, r *http.Request) {
	f := new(frontend)

	// Since my frontend is stored in a Google Cloud Bucket, I have the url stored in the database
	// so I can easily change it without re-deploying. TODO: use custom domain for storage bucket.
	url, err := f.getFrontendURL()

	if err != nil {
		glog.Errorf("frontend url not found: %s\n", err.Error())
		return
	}

	// Relay request to storage bucket
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (f *frontend) getFrontendURL() (url string, err error) {
	var infoRow map[string]interface{}

	err = database.SelectRow(&database.Query{
		Collection: "info",
		Conditions: map[string]interface{}{
			"Name": "FrontendURL",
		},
		Destination: &infoRow,
	})

	if err != nil {
		glog.Errorf("error querying frontend url: %s\n", err.Error())
		return "", err
	}

	return infoRow["Value"].(string), nil
}

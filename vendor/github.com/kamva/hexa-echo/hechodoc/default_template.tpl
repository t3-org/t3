{{.BeginRouteVal}}
// swagger:route {{.Method}} {{.Path}} {{.TagsString}} {{.ParamsId}}
//
// responses:
//   200: {{.SuccessRespId}}

// swagger:parameters {{.ParamsId}}
type {{.Name}}ParamsWrapper struct {
  {{- range .PathParams}}
     // in:path
     {{.ExportedName}} string `json:"{{.Name}}"`
  {{ end}}
	// in:body
	Body struct{
	    // DOCTODO: place your params body here
	}
}

// success response
// swagger:response {{.SuccessRespId}}
type {{.Name}}ResponseWrapper struct {
	// in:body
	Body struct{
	    // response code
    	Code string `json:"code"`
    	// DOCTODO: place your body data here
    }
}

{{.EndRouteVal}}


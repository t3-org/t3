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

	}
}

// success response
// swagger:response {{.SuccessRespId}}
type {{.Name}}ResponseWrapper struct {
	// in:body
	Body struct{
	    replyCode

    }
}
{{.EndRouteVal}}



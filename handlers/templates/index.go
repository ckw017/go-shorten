package templates

import (
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/guregu/kami"
	"golang.org/x/net/context"
)

type IndexParams struct {
	Short string
	Error error
}

func IndexFromContext(ctx context.Context) (IndexParams, bool) {
	p, ok := ctx.Value("IndexParams").(IndexParams)
	return p, ok
}

func IndexWithContext(ctx context.Context, ip IndexParams) context.Context {
	return context.WithValue(ctx, "IndexParams", ip)
}

func Index(srcBox *rice.Box) kami.HandlerFunc {
	t := template.Must(Root.New("index").Parse(
		srcBox.MustString("templates/index.tmpl"),
	))

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		params, ok := IndexFromContext(ctx)
		if !ok {
			params = IndexParams{}
		}
		err := t.Execute(w, params)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
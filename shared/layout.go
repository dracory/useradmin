package shared

import (
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// Layout is the default layout renderer. It builds a complete HTML page
// using hb.NewWebpage() with Bootstrap + Vue CDN, matching the
// blogadmin/shopadmin pattern. If a FuncLayout is provided in
// AdminOptions, it takes precedence over this default.
func Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	_ = r // kept for layout signature compatibility
	_ = w
	return webpageComplete(webpageTitle, webpageHtml, options).ToHTML()
}

// webpageComplete builds the default webpage template
func webpageComplete(title, content string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) *hb.HtmlWebpage {
	webpage := hb.NewWebpage()
	webpage.SetTitle(title)

	webpage.AddStyleURLs([]string{
		cdn.BootstrapCss_5_3_3(),
		"https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css",
	})
	webpage.AddStyleURLs(options.StyleURLs)
	webpage.AddScriptURLs([]string{
		cdn.BootstrapJs_5_3_3(),
		cdn.Jquery_3_7_1(),
		cdn.VueJs_3(),
		cdn.Sweetalert2_11(),
	})
	webpage.AddScriptURLs(options.ScriptURLs)
	// loadVueIfNeeded is a guard so controller JS can safely call
	// Vue.createApp() even if the layout already loaded Vue (or vice versa).
	webpage.AddScripts(append([]string{VueLoaderJS}, options.Scripts...))
	webpage.AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}`)
	webpage.AddStyles(options.Styles)
	webpage.AddChild(hb.NewHTML(content))
	return webpage
}

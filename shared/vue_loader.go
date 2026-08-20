package shared

import "github.com/dracory/hb"

// VueLoaderJS is a JavaScript snippet that defines loadVueIfNeeded().
// It checks whether Vue is already loaded (e.g. by the layout) and,
// if not, loads it from the CDN dynamically before invoking the callback.
//
// Every controller that needs Vue includes this snippet once (via the
// layout or directly). Component JS files then wrap their mount call in:
//
//	loadVueIfNeeded((err) => {
//	  if (err) { console.error('Vue load failed:', err); return; }
//	  const { createApp } = Vue;
//	  const el = document.getElementById('my-app');
//	  if (el) createApp(MyApp).mount('#my-app');
//	});
//
// This avoids double-loading Vue when the layout already provides it.
const VueLoaderJS = `function loadVueIfNeeded(callback) {
  if (typeof Vue !== 'undefined' || (typeof window !== 'undefined' && window.Vue)) {
    return callback(null);
  }
  var script = document.createElement('script');
  script.src = 'https://cdn.jsdelivr.net/npm/vue@3/dist/vue.global.js';
  script.crossOrigin = 'anonymous';
  script.onload = function() { callback(null); };
  script.onerror = function() { callback(new Error('Vue.js failed to load')); };
  document.head.appendChild(script);
}`

// VueLoaderScript returns the loadVueIfNeeded definition as an hb.Tag
// so controllers can inject it into their container divs. This ensures
// the guard is available even when a host project provides a custom
// FuncLayout that does not include VueLoaderJS.
func VueLoaderScript() hb.TagInterface {
	return hb.Script(VueLoaderJS)
}

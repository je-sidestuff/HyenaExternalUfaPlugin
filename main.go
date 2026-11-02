// hyena-ufa-plugin is a hello-world external UFA plug-in. It follows the
// contract in AI-evo1-spiral's docs/PluginSubApps.md: say hello over
// representable, serve /ui-manifest.json and an icon, and serve one page
// using relative URLs only. External plug-ins run sandboxed
// (docs/PluginTrust.md), so the page gets its data server-side.
//
//	hyena-ufa-plugin --port 8091 --lr-port 8082
package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"representable"
)

//go:embed static
var static embed.FS

const (
	reconnectInterval = 5 * time.Second
	pluginVersion     = "1"
	displayName       = "Hyena"
)

const iconSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M5 9L7 1l4 4 2-.5 2 .5 4-4 2 8-1 6c-1 4-4 7-7 7s-6-3-7-7z" fill="#3a2d22"/><circle cx="5.5" cy="5.5" r="4" fill="#b08a55" stroke="#3a2d22" stroke-width="1"/><circle cx="18.5" cy="5.5" r="4" fill="#b08a55" stroke="#3a2d22" stroke-width="1"/><circle cx="5.5" cy="5.5" r="2" fill="#3a2d22"/><circle cx="18.5" cy="5.5" r="2" fill="#3a2d22"/><path d="M12 5c4 0 7 3 7 8 0 4-3 8-7 8s-7-4-7-8c0-5 3-8 7-8z" fill="#c9a66b" stroke="#3a2d22" stroke-width="1"/><path d="M12 10c3 0 5 2 5 5 0 3-2 6-5 6s-5-3-5-6c0-3 2-5 5-5z" fill="#2b2420"/><circle cx="8.6" cy="11.3" r="1.1" fill="#ffd34d"/><circle cx="15.4" cy="11.3" r="1.1" fill="#ffd34d"/><circle cx="8.6" cy="11.3" r=".5" fill="#111"/><circle cx="15.4" cy="11.3" r=".5" fill="#111"/><path d="M10.2 14.5h3.6l-1.8 1.8z" fill="#111" stroke="#111" stroke-width=".4"/><path d="M8.5 17.2q3.5 3 7 0" fill="#fff" stroke="#111" stroke-width=".8" stroke-linejoin="round"/><circle cx="8" cy="7.5" r=".8" fill="#5b4420"/><circle cx="16" cy="7.5" r=".8" fill="#5b4420"/><circle cx="12" cy="7" r=".7" fill="#5b4420"/><circle cx="6.5" cy="13.5" r=".7" fill="#5b4420"/><circle cx="17.5" cy="13.5" r=".7" fill="#5b4420"/></svg>`

// Hello is what the page is served with.
type Hello struct {
	Message string `json:"message"`
	AppID   string `json:"app_id"`
	Time    string `json:"time"`
}

func main() {
	port := flag.String("port", "8091", "HTTP port to listen on")
	name := flag.String("name", "hyena-external-ufa-plugin", "app id, and the name reported to local-representative: they must match")
	lrHost := flag.String("lr-host", "localhost", "local-representative host/IP")
	lrPort := flag.String("lr-port", "8082", "local-representative representable port")
	devMode := flag.Bool("dev-mode", false, "dev mode (SDLC sense); must match local-representative's")
	flag.Parse()

	if !representable.ValidAppID(*name) {
		log.Fatalf("--name %q is not a valid app id (see representable.ValidAppID)", *name)
	}
	go connectLoop(net.JoinHostPort(*lrHost, *lrPort), *name, *port, *devMode)

	addr := ":" + *port
	log.Printf("hyena-ufa-plugin %q listening on http://localhost%s", *name, addr)
	log.Fatal(http.ListenAndServe(addr, routes(*name)))
}

// manifest is served at /ui-manifest.json. URLs are relative to the
// plug-in's own root: it doesn't know the base path it is mounted under.
func manifest(appID string) representable.UIManifest {
	return representable.UIManifest{
		ManifestVersion: representable.UIManifestVersion,
		AppID:           appID,
		DisplayName:     displayName,
		ShortName:       displayName,
		Entry:           "./",
		Icon:            "." + representable.PluginIconPath,
		Order:           110,
		Version:         pluginVersion,
	}
}

func routes(appID string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(representable.UIManifestPath, representable.ServeUIManifest(func() representable.UIManifest { return manifest(appID) }))
	mux.HandleFunc(representable.PluginIconPath, representable.ServeSVG(iconSVG))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		page, _ := static.ReadFile("static/index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write(withHello(page, Hello{
			Message: "Hello from the Hyena external plug-in!",
			AppID:   appID,
			Time:    time.Now().UTC().Format(time.RFC3339),
		}))
	})
	return mux
}

const helloPlaceholder = "/*HELLO*/null"

// withHello fills the page's placeholder. A sandboxed page's fetches carry
// no cookies, so the data is served inside the page rather than fetched.
// json.Marshal escapes <, > and &, so the result is safe inside <script>.
func withHello(page []byte, h Hello) []byte {
	js, err := json.Marshal(h)
	if err != nil {
		return page
	}
	return bytes.Replace(page, []byte(helloPlaceholder), js, 1)
}

// connectLoop keeps a representable connection to local-representative,
// sending the plug-in hello after every (re)connect. The connection is the
// liveness signal: while it is down LR shows this plug-in offline.
func connectLoop(addr, appID, httpPort string, devMode bool) {
	for {
		c, err := representable.Connect(addr, appID, representable.Mode(devMode), 3*time.Second)
		if err != nil {
			time.Sleep(reconnectInterval)
			continue
		}
		log.Printf("connected to local-representative at %s as %q", addr, appID)
		c.SendPluginHello(representable.PluginHello{
			AppID:       appID,
			DisplayName: displayName,
			HTTPPort:    httpPort,
			DevMode:     devMode,
		})
		<-c.DisconnectCh()
		log.Printf("disconnected from local-representative at %s; retrying every %s", addr, reconnectInterval)
		time.Sleep(reconnectInterval)
	}
}

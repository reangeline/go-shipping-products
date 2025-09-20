package ginadapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterDocs publishes openapi.yaml and a Swagger UI page.
// - specPath: local path to the YAML
// - mount: public prefix (e.g., "/docs")
func RegisterDocs(r *gin.Engine, specPath, mount string) {
	if mount == "" {
		mount = "/docs"
	}
	r.StaticFile(mount+"/openapi.yaml", specPath)

	const page = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Shipping Packs API – Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html,body,#swagger-ui {height:100%; margin:0; background:#fff;}
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: window.location.pathname.replace(/\/$/, '') + '/openapi.yaml',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

	r.GET(mount, func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
	})
}

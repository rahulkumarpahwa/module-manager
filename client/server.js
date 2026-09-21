/**
 * Module Manager Client - optional dev server
 *
 * Serves the static client files AND proxies /api/* to the Go backend.
 * Because browser calls to a different origin without CORS headers are
 * blocked, this server makes everything same-origin so the client works
 * out of the box.
 *
 * Usage:
 *   node server.js                 # http://localhost:8080
 *   PORT=9000 API_TARGET=http://localhost:8778 node server.js
 *
 * Needs Node.js (no npm dependencies).
 */

const http = require("http");
const fs = require("fs");
const path = require("path");

const ROOT = __dirname;
const PORT = Number(process.env.PORT) || 8080;
const API_TARGET = process.env.API_TARGET || "http://localhost:8778";

const MIME = {
  ".html": "text/html; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".svg": "image/svg+xml",
  ".png": "image/png",
  ".ico": "image/x-icon",
  ".md": "text/plain; charset=utf-8",
};

const server = http.createServer((req, res) => {
  if (req.url.startsWith("/api/")) {
    proxyApi(req, res);
  } else {
    serveStatic(req, res);
  }
});

function serveStatic(req, res) {
  const url = new URL(req.url, `http://${req.headers.host || "localhost"}`);
  let filePath = path.normalize(path.join(ROOT, decodeURIComponent(url.pathname)));
  if (url.pathname === "/") filePath = path.join(ROOT, "index.html");

  if (!filePath.startsWith(ROOT)) {
    res.writeHead(403, { "Content-Type": "text/plain; charset=utf-8" });
    res.end("Forbidden");
    return;
  }

  fs.stat(filePath, (err, stat) => {
    if (err || !stat.isFile()) {
      res.writeHead(404, { "Content-Type": "text/plain; charset=utf-8" });
      res.end("Not found");
      return;
    }
    res.writeHead(200, { "Content-Type": MIME[path.extname(filePath)] || "application/octet-stream" });
    fs.createReadStream(filePath).pipe(res);
  });
}

function proxyApi(req, res) {
  const targetUrl = new URL(req.url, API_TARGET);
  const headers = { ...req.headers, host: targetUrl.host };

  const preq = http.request(
    targetUrl,
    { method: req.method, headers },
    (pres) => {
      res.writeHead(pres.statusCode || 502, pres.headers);
      pres.pipe(res);
    }
  );

  preq.on("error", (err) => {
    res.writeHead(502, { "Content-Type": "text/plain; charset=utf-8" });
    res.end(`Proxy error: ${err.message}`);
  });

  req.pipe(preq);
}

server.listen(PORT, () => {
  console.log(`Module Manager client served at  http://localhost:${PORT}`);
  console.log(`API /api/* proxied to            ${API_TARGET}`);
  console.log("Press Ctrl+C to stop.");
});
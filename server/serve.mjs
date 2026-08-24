/** Production server: dist + live /api/v1. Docker and `npm start`. */
import http from "node:http";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { handleApi } from "./malo-api.mjs";

const ROOT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const DIST = process.env.DIST_DIR ? path.resolve(process.env.DIST_DIR) : path.join(ROOT, "dist");
const PORT = Number(process.env.PORT || 8080);
const HOST = process.env.HOST || "0.0.0.0";

const MIME = {
  ".html": "text/html; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".svg": "image/svg+xml",
  ".png": "image/png",
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".webp": "image/webp",
  ".woff2": "font/woff2",
  ".woff": "font/woff",
  ".ico": "image/x-icon",
  ".map": "application/json",
  ".txt": "text/plain; charset=utf-8",
  ".bpmn": "application/bpmn20-xml; charset=utf-8",
  ".xml": "application/xml; charset=utf-8",
};

function sendFile(res, file) {
  const ext = path.extname(file);
  const body = fs.readFileSync(file);
  res.statusCode = 200;
  res.setHeader("Content-Type", MIME[ext] || "application/octet-stream");
  if (ext === ".html") res.setHeader("Cache-Control", "no-cache");
  else if (file.includes(`${path.sep}assets${path.sep}`)) {
    res.setHeader("Cache-Control", "public, max-age=31536000, immutable");
  }
  res.end(body);
}

function safeJoin(root, urlPath) {
  const decoded = decodeURIComponent((urlPath || "/").split("?")[0].split("#")[0]);
  const rel = decoded.replace(/^\/+/, "");
  const file = path.normalize(path.join(root, rel || "index.html"));
  if (!file.startsWith(root)) return null;
  return file;
}

function serveStatic(req, res) {
  const file = safeJoin(DIST, req.url || "/");
  if (!file) {
    res.statusCode = 400;
    res.end("bad path");
    return;
  }
  let target = file;
  try {
    if (fs.existsSync(target) && fs.statSync(target).isDirectory()) {
      target = path.join(target, "index.html");
    }
    if (fs.existsSync(target) && fs.statSync(target).isFile()) {
      sendFile(res, target);
      return;
    }
  } catch {
    /* SPA fallback */
  }
  const index = path.join(DIST, "index.html");
  if (fs.existsSync(index)) {
    sendFile(res, index);
    return;
  }
  res.statusCode = 404;
  res.end("not found");
}

const server = http.createServer((req, res) => {
  handleApi(req, res, () => serveStatic(req, res));
});

server.listen(PORT, HOST, () => {
  console.log(`Malo Academy http://${HOST}:${PORT}  (dist ${DIST})`);
});

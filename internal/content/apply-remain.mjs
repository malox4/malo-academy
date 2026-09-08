/**
 * Apply leftover junior/middle teach() bodies.
 * Run: node internal/content/apply-remain.mjs
 */
import { readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { REMAIN } from "./teach-bodies-remain.mjs";

const root = join(dirname(fileURLToPath(import.meta.url)), "../..");

function replaceArticle(src, id, body) {
  const marker = `id: "${id}"`;
  const idAt = src.indexOf(marker);
  if (idAt < 0) throw new Error("id not found: " + id);
  const teaserAt = src.indexOf("teaser:", idAt);
  if (teaserAt < 0 || teaserAt - idAt > 800) throw new Error("teaser not found: " + id);
  const teaserLineEnd = src.indexOf("\n", teaserAt);
  const bodyStart = teaserLineEnd + 1;
  const closer = src.indexOf("\n  }),", bodyStart);
  if (closer < 0) throw new Error("closer not found: " + id);
  return src.slice(0, bodyStart) + body + src.slice(closer);
}

function patch(rel) {
  const path = join(root, rel);
  let src = readFileSync(path, "utf8");
  const applied = [];
  for (const [id, body] of Object.entries(REMAIN)) {
    if (!src.includes(`id: "${id}"`)) continue;
    src = replaceArticle(src, id, body);
    applied.push(id);
  }
  writeFileSync(path, src);
  return applied;
}

const files = [
  "web/src/content/topicsCatalog.js",
  "web/src/content/topicsGap.js",
  "web/src/content/topicsGapMore.js",
];

for (const f of files) {
  const applied = patch(f);
  console.log(f, applied.length, applied.join(", "));
}

const missing = Object.keys(REMAIN).filter((id) => {
  const all = files.map((f) => readFileSync(join(root, f), "utf8")).join("\n");
  return !all.includes(`id: "${id}"`);
});
if (missing.length) console.error("missing ids", missing);

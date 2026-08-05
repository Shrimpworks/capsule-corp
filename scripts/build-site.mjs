import { copyFile, cp, mkdir, rm } from "node:fs/promises";

const projectRoot = new URL("../", import.meta.url);
const siteRoot = new URL("site/", projectRoot);
const outputRoot = new URL(".site-dist/", projectRoot);

await rm(outputRoot, { force: true, recursive: true });
await mkdir(outputRoot, { recursive: true });

for (const file of ["index.html", "styles.css", ".nojekyll"]) {
  await copyFile(new URL(file, siteRoot), new URL(file, outputRoot));
}

await copyFile(new URL("tokens.css", projectRoot), new URL("tokens.css", outputRoot));
await cp(new URL("assets/", siteRoot), new URL("assets/", outputRoot), { recursive: true });

console.log("Built static site in .site-dist/");

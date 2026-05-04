import { copyFile, mkdir, readFile, readdir, stat, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { productDocPages, products, topics } from "./bilingual-metadata.mjs";
import { productDocs } from "./product-doc-content.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const repoRoot = path.resolve(__dirname, "..");
const docsRoot = path.resolve(repoRoot, "docs");
const zhRoot = path.resolve(docsRoot, "zh");

const exists = async (target) => {
  try {
    await stat(target);
    return true;
  } catch {
    return false;
  }
};

const ensureDir = async (target) => {
  await mkdir(path.dirname(target), { recursive: true });
};

const rewriteChineseLinks = (source) =>
  source
    .replace(/href="\/products\//g, 'href="/zh/products/')
    .replace(/href="\/topics\//g, 'href="/zh/topics/')
    .replace(/\]\(\/products\//g, "](/zh/products/")
    .replace(/\]\(\/topics\//g, "](/zh/topics/")
    .replace(
      /\.\.\/\.vitepress\/theme\/components\/ForumPostIndex\.vue/g,
      "../../.vitepress/theme/components/ForumPostIndex.vue"
    )
    .replace(/(^\s*link:\s*)\/topics\/\s*$/gm, "$1/zh/topics/")
    .replace(/(^\s*link:\s*)\/products\/\s*$/gm, "$1/zh/products/")
    .replace(/(^\s*link:\s*)\/\s*$/gm, "$1/zh/")
    .replace(/\]\(\/\)/g, "](/zh/)");

const copyChineseSource = async (sourceDir, targetDir) => {
  const entries = await readdir(sourceDir, { withFileTypes: true });
  await mkdir(targetDir, { recursive: true });

  for (const entry of entries) {
    if (
      entry.name === ".vitepress" ||
      entry.name === "public" ||
      entry.name === "logo" ||
      entry.name === "zh"
    ) {
      continue;
    }

    const sourcePath = path.join(sourceDir, entry.name);
    const targetPath = path.join(targetDir, entry.name);

    if (entry.isDirectory()) {
      await copyChineseSource(sourcePath, targetPath);
      continue;
    }

    if (entry.name === "posts.data.ts") {
      continue;
    }

    if (entry.name.endsWith(".md")) {
      const content = await readFile(sourcePath, "utf8");
      await ensureDir(targetPath);
      await writeFile(targetPath, rewriteChineseLinks(content), "utf8");
      continue;
    }

    await ensureDir(targetPath);
    await copyFile(sourcePath, targetPath);
  }
};

const writeDoc = async (relativePath, content) => {
  const target = path.resolve(docsRoot, relativePath);
  await ensureDir(target);
  await writeFile(target, `${content.trim()}\n`, "utf8");
};

const productLink = (slug, page = "") =>
  page ? `/products/${slug}/${page}` : `/products/${slug}/`;

const zhProductLink = (slug, page = "") =>
  page ? `/zh/products/${slug}/${page}` : `/zh/products/${slug}/`;

const topicLink = (slug) => `/topics/${slug}`;
const zhTopicLink = (slug) => `/zh/topics/${slug}`;
const topicDocRoot = path.resolve(repoRoot, "scripts", "topic-docs");

const renderHomePage = () => {
  const rows = products
    .map(
      (product, index) => `          <tr data-node="m${String(index + 1).padStart(2, "0")}" data-group="${resolveGroup(
        product.slug
      )}" tabindex="0">
            <td><span class="stellar-home-matrix__id">M${String(index + 1).padStart(2, "0")}</span></td>
            <td><strong>${product.domainEn.toUpperCase()}</strong></td>
            <td>${product.originEn}</td>
            <td><strong><a href="${productLink(product.slug)}">${product.name}</a></strong></td>
            <td>${product.rationaleEn}</td>
          </tr>`
    )
    .join("\n");

  return `---
layout: home

hero:
  name: Stell Hub
  text: Platform Infrastructure Documentation
  tagline: English-first reference material for the Stell Hub middleware stack, product architecture, and engineering design notes.
  actions:
    - theme: brand
      text: Explore Products
      link: /products/
    - theme: alt
      text: Explore Topics
      link: /topics/
---

<div class="stellar-home-wide">
  <section class="stellar-home-matrix">
    <h2>Middleware Capability Matrix</h2>
    <div class="stellar-home-capability-legend" aria-label="Capability legend">
      <span class="stellar-home-capability-legend__item is-access">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>ACCESS / SECURITY</span>
      </span>
      <span class="stellar-home-capability-legend__item is-runtime">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>GOVERNANCE / RUNTIME</span>
      </span>
      <span class="stellar-home-capability-legend__item is-foundation">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>FOUNDATION</span>
      </span>
      <span class="stellar-home-capability-legend__item is-observe">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>OBSERVABILITY / ALERT</span>
      </span>
    </div>
    <div class="stellar-home-matrix__table">
      <table>
        <colgroup>
          <col class="stellar-home-matrix__col-id">
          <col class="stellar-home-matrix__col-domain">
          <col class="stellar-home-matrix__col-origin">
          <col class="stellar-home-matrix__col-name">
          <col class="stellar-home-matrix__col-logic">
        </colgroup>
        <thead>
          <tr>
            <th><span class="stellar-home-matrix__tag">ID</span></th>
            <th><span class="stellar-home-matrix__tag">DOMAIN</span></th>
            <th><span class="stellar-home-matrix__tag">ORIGIN</span></th>
            <th><span class="stellar-home-matrix__tag">STELL NAME</span></th>
            <th><span class="stellar-home-matrix__tag">RATIONALE</span></th>
          </tr>
        </thead>
        <tbody>
${rows}
        </tbody>
      </table>
    </div>
  </section>
</div>

<div class="stellar-home-wide">
  <section class="stellar-home-graph">
    <h2>Dependency Blueprint</h2>
    <p>Use this view to read the stack by platform layers, core contracts, and observability feedback paths.</p>
    <object
      id="stellar-home-blueprint"
      class="stellar-home-blueprint"
      data="/core-middleware-matrix-dependency.svg"
      type="image/svg+xml"
      aria-label="Middleware dependency blueprint"
    ></object>
  </section>
</div>

## Contact

- GitHub: [https://github.com/stellhub](https://github.com/stellhub)
- Email: [xiaoyaoyunlian@gmail.com](mailto:xiaoyaoyunlian@gmail.com)
- Backup email: [chenwenlong_java@163.com](mailto:chenwenlong_java@163.com)`;
};

const renderProductsIndex = () => {
  const cards = products
    .map(
      (product) => `  <a class="product-card" href="${productLink(product.slug)}">
    <div class="product-card-logo"><img src="/logo/${product.slug}.png" alt="${product.name} Logo"></div>
    <span class="product-card-domain">${product.domainEn}</span>
    <h3 class="product-card-title">${product.name} · ${product.originEn}</h3>
    <p class="product-card-desc">${product.summaryEn}</p>
    <span class="product-card-link">Read Reference</span>
  </a>`
    )
    .join("\n");

  return `---
title: Core Product Overview
outline: deep
---

# Core Product Overview

This page turns the Stell Hub middleware stack into product-level entry points. Each card leads to the same documentation structure, so architecture, deployment, integration, and observability material stay predictable across the English reference site.

## Product Entry Points

<div class="product-grid">
${cards}
</div>

## Documentation Convention

Each product follows the same structure:

- Overview page: positioning, core capabilities, and typical scenarios
- Sub-pages: design overview, system architecture, deployment model, getting started, configuration guide, API reference, and observability guide
- Chinese reference: every English page points to its original Chinese counterpart under \`/zh\`

This keeps the English reference layer consistent while preserving the full Chinese knowledge base.`;
};

const renderSections = (sections) =>
  sections
    .map((section) => {
      const lines = [`## ${section.heading}`];

      if (section.paragraphs?.length) {
        lines.push("", ...section.paragraphs);
      }

      if (section.bullets?.length) {
        lines.push("", ...section.bullets.map((item) => `- ${item}`));
      }

      if (section.ordered?.length) {
        lines.push("", ...section.ordered.map((item, index) => `${index + 1}. ${item}`));
      }

      return lines.join("\n");
    })
    .join("\n\n");

const renderProductIndex = (product) => {
  const doc = productDocs[product.slug];
  if (!doc) {
    throw new Error(`Missing product doc metadata for ${product.slug}`);
  }

  const docLinks = productDocPages
    .map((page) => `- [${page.titleEn}](${productLink(product.slug, page.slug)})`)
    .join("\n");

  const boundaries = doc.overview.boundaries.map((item) => `- ${item}`).join("\n");
  const capabilities = doc.overview.capabilities.map((item) => `- ${item}`).join("\n");
  const values = doc.overview.values.map((item) => `- ${item}`).join("\n");
  const businessScenarios = doc.overview.businessScenarios.map((item) => `- ${item}`).join("\n");
  const platformScenarios = doc.overview.platformScenarios.map((item) => `- ${item}`).join("\n");

  return `---
title: ${product.name} Design
outline: deep
---

# ${product.name} · ${product.originEn}

<div class="product-logo">
  <img src="/logo/${product.slug}.png" alt="${product.name} Logo">
</div>

> ${doc.overview.tagline}

## Product Scope

### Objective

${doc.overview.goal}

### Boundaries

${boundaries}

## Core Capabilities

### Capabilities

${capabilities}

### Engineering Value

${values}

## Reference Sections

${docLinks}

## Typical Use Cases

### Business Use Cases

${businessScenarios}

### Platform Use Cases

${platformScenarios}

## Chinese Source

- [Read the original Chinese product page](${zhProductLink(product.slug)})`;
};

const renderProductSubpage = (product, page, index) => {
  const doc = productDocs[product.slug];
  if (!doc) {
    throw new Error(`Missing product doc metadata for ${product.slug}`);
  }

  const pageDoc = doc.pages[page.slug];
  if (!pageDoc) {
    throw new Error(`Missing page metadata for ${product.slug}/${page.slug}`);
  }

  const previousPage = productDocPages[index - 1];
  const nextPage = productDocPages[index + 1];
  const readingOrder = [
    previousPage ? `- Previous: [${previousPage.titleEn}](${productLink(product.slug, previousPage.slug)})` : "",
    nextPage ? `- Next: [${nextPage.titleEn}](${productLink(product.slug, nextPage.slug)})` : ""
  ]
    .filter(Boolean)
    .join("\n");

  return `---
title: ${product.name} ${page.titleEn}
outline: deep
---

# ${product.name} · ${page.titleEn}

> ${doc.overview.tagline}

${renderSections(pageDoc.sections)}

## Continue Reading

- Start with the [${product.name} product overview](${productLink(product.slug)})
${readingOrder}

## Chinese Source

- [Read the original Chinese page](${zhProductLink(product.slug, page.slug)})`;
};

const renderTopicsIndex = () => `---
layout: home
hero:
  name: Stell Notes
  text: Engineering Notes and Design Research
  tagline: An English-first index of infrastructure and middleware writing, centered on system boundaries, tradeoffs, and implementation decisions.
  actions:
    - theme: brand
      text: Browse Latest Posts
      link: "#forum-latest"
    - theme: alt
      text: Return Home
      link: /
---

<script setup>
import ForumPostIndex from "../.vitepress/theme/components/ForumPostIndex.vue";
import { data as topicPosts } from "./posts.data.ts";
</script>

<div id="forum-latest">
  <ForumPostIndex :posts="topicPosts" />
</div>

## What You Will Find Here

- Recurring engineering problems that appear during middleware and platform evolution
- Explanations of why specific infrastructure capabilities exist and where their boundaries should stop
- Tradeoff analysis across mainstream implementation approaches
- Decision frameworks for evaluating these choices in real production environments

## Next Topics in Scope

- Service governance and traffic-control strategy
- Boundary design between configuration centers and registry centers
- Migration paths from queues to event-stream platforms
- Collaboration models between infrastructure and platform engineering teams`;

const renderTopicPage = (topic, body) => {
  const tags = topic.tagsEn.map((tag) => `  - ${tag}`).join("\n");

  return `---
title: ${topic.titleEn}
category: ${topic.categoryEn}
summary: ${topic.summaryEn}
tags:
${tags}
readingDirection: ${topic.readingDirectionEn}
outline: deep
---

# ${topic.titleEn}

## Overview

${topic.summaryEn}

${body}

## Chinese Reference

- [Read the original Chinese article](${zhTopicLink(topic.slug)})`;
};

const renderRootTopicData = () => `import { execFileSync } from "node:child_process";
import { statSync } from "node:fs";
import path from "node:path";
import { createContentLoader } from "vitepress";

export interface TopicPost {
  title: string;
  url: string;
  category: string;
  summary: string;
  tags: string[];
  readingDirection: string;
  publishAt: string;
  updatedAt: string;
}

const REPO_ROOT = path.resolve(__dirname, "../..");
const DOCS_ROOT = path.resolve(REPO_ROOT, "docs");

const getGitDate = (filePath: string, args: string[]) => {
  try {
    const value = execFileSync("git", ["-C", REPO_ROOT, ...args, "--", filePath], {
      encoding: "utf-8"
    }).trim();
    return value || "";
  } catch {
    return "";
  }
};

const getFallbackDate = (filePath: string) => statSync(filePath).mtime.toISOString();

const resolveTopicFile = (url: string) =>
  path.resolve(DOCS_ROOT, \`\${url.replace(/^\\//, "").replace(/\\/$/, "")}.md\`);

export default createContentLoader("topics/*.md", {
  globOptions: {
    ignore: ["topics/index.md"]
  },
  transform(rawPosts) {
    return rawPosts
      .filter((post) => post.url !== "/topics/")
      .map<TopicPost>((post) => {
        const filePath = resolveTopicFile(post.url);
        const publishAt =
          getGitDate(filePath, ["log", "--diff-filter=A", "-1", "--format=%aI"]) ||
          getFallbackDate(filePath);
        const updatedAt =
          getGitDate(filePath, ["log", "-1", "--format=%cI"]) ||
          getFallbackDate(filePath);

        return {
          title: post.frontmatter.title ?? "",
          url: post.url,
          category: post.frontmatter.category ?? "Uncategorized",
          summary: post.frontmatter.summary ?? "",
          tags: Array.isArray(post.frontmatter.tags) ? post.frontmatter.tags : [],
          readingDirection: post.frontmatter.readingDirection ?? "",
          publishAt,
          updatedAt
        };
      })
      .sort((left, right) => right.updatedAt.localeCompare(left.updatedAt));
  }
});`;

const renderZhTopicData = () => `import { execFileSync } from "node:child_process";
import { statSync } from "node:fs";
import path from "node:path";
import { createContentLoader } from "vitepress";

export interface TopicPost {
  title: string;
  url: string;
  category: string;
  summary: string;
  tags: string[];
  readingDirection: string;
  publishAt: string;
  updatedAt: string;
}

const REPO_ROOT = path.resolve(__dirname, "../../..");
const DOCS_ROOT = path.resolve(REPO_ROOT, "docs");

const getGitDate = (filePath: string, args: string[]) => {
  try {
    const value = execFileSync("git", ["-C", REPO_ROOT, ...args, "--", filePath], {
      encoding: "utf-8"
    }).trim();
    return value || "";
  } catch {
    return "";
  }
};

const getFallbackDate = (filePath: string) => statSync(filePath).mtime.toISOString();

const resolveTopicFile = (url: string) =>
  path.resolve(DOCS_ROOT, \`\${url.replace(/^\\//, "").replace(/\\/$/, "")}.md\`);

export default createContentLoader("zh/topics/*.md", {
  globOptions: {
    ignore: ["zh/topics/index.md"]
  },
  transform(rawPosts) {
    return rawPosts
      .filter((post) => post.url !== "/zh/topics/")
      .map<TopicPost>((post) => {
        const filePath = resolveTopicFile(post.url);
        const publishAt =
          getGitDate(filePath, ["log", "--diff-filter=A", "-1", "--format=%aI"]) ||
          getFallbackDate(filePath);
        const updatedAt =
          getGitDate(filePath, ["log", "-1", "--format=%cI"]) ||
          getFallbackDate(filePath);

        return {
          title: post.frontmatter.title ?? "",
          url: post.url,
          category: post.frontmatter.category ?? "未分类",
          summary: post.frontmatter.summary ?? "",
          tags: Array.isArray(post.frontmatter.tags) ? post.frontmatter.tags : [],
          readingDirection: post.frontmatter.readingDirection ?? "",
          publishAt,
          updatedAt
        };
      })
      .sort((left, right) => right.updatedAt.localeCompare(left.updatedAt));
  }
});`;

const resolveGroup = (slug) => {
  if (["stellgate", "stellguard"].includes(slug)) {
    return "access";
  }
  if (["stelltrace", "stellspec", "stellcon", "stellvox"].includes(slug)) {
    return "observe";
  }
  if (["stellmap", "stellnula", "stellpoint", "stellkey"].includes(slug)) {
    return "foundation";
  }
  return "runtime";
};

const bootstrapChineseDocs = async () => {
  if (await exists(zhRoot)) {
    return;
  }

  await mkdir(zhRoot, { recursive: true });
  await copyChineseSource(path.resolve(docsRoot), zhRoot);
};

const generateEnglishDocs = async () => {
  await writeDoc("index.md", renderHomePage());
  await writeDoc("products/index.md", renderProductsIndex());
  await writeDoc("topics/index.md", renderTopicsIndex());
  await writeDoc("topics/posts.data.ts", renderRootTopicData());

  for (const product of products) {
    await writeDoc(`products/${product.slug}/index.md`, renderProductIndex(product));

    for (const [index, page] of productDocPages.entries()) {
      await writeDoc(
        `products/${product.slug}/${page.slug}.md`,
        renderProductSubpage(product, page, index)
      );
    }
  }

  for (const topic of topics) {
    const topicBody = await readFile(path.resolve(topicDocRoot, `${topic.slug}.md`), "utf8");
    await writeDoc(`topics/${topic.slug}.md`, renderTopicPage(topic, topicBody.trim()));
  }
};

const generateChineseHelpers = async () => {
  await writeDoc("zh/topics/posts.data.ts", renderZhTopicData());
};

await bootstrapChineseDocs();
await generateEnglishDocs();
await generateChineseHelpers();

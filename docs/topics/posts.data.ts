import { execFileSync } from "node:child_process";
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
  path.resolve(DOCS_ROOT, `${url.replace(/^\//, "").replace(/\/$/, "")}.md`);

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
});

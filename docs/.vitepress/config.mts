import { defineConfig } from "vitepress";
import { products, topics, productDocPages } from "../../scripts/bilingual-metadata.mjs";

const buildPath = (base: string, value: string) => `${base}${value}`;
const buildHomeLink = (base: string) => (base ? `${base}/` : "/");

const buildProductMenu = (base: string, isZh: boolean) => [
  {
    text: isZh ? "产品总览" : "Product Entry Points",
    link: buildPath(base, "/products/")
  },
  ...products.map((product) => ({
    text: isZh
      ? `${product.name} · ${product.chineseName}`
      : `${product.name} · ${product.originEn}`,
    link: buildPath(base, `/products/${product.slug}/`)
  }))
];

const buildTopicMenu = (base: string, isZh: boolean) => [
  {
    text: isZh ? "论坛首页" : "Topics Home",
    link: buildPath(base, "/topics/")
  },
  ...topics.map((topic) => ({
    text: isZh ? topic.titleZh : topic.titleEn,
    link: buildPath(base, `/topics/${topic.slug}`)
  }))
];

const buildProductSidebar = (base: string, isZh: boolean, slug: string) => [
  {
    text: isZh ? "设计文档" : "Reference Sections",
    items: [
      {
        text: isZh ? "产品首页" : "Product Overview",
        link: buildPath(base, `/products/${slug}/`)
      },
      ...productDocPages.map((page) => ({
        text: isZh ? page.titleZh : page.titleEn,
        link: buildPath(base, `/products/${slug}/${page.slug}`)
      }))
    ]
  }
];

const buildRootSidebar = (base: string, isZh: boolean) => [
  {
    text: isZh ? "文档导航" : "Reference Guide",
    items: [
      {
        text: isZh ? "首页" : "Home",
        link: buildHomeLink(base)
      },
      {
        text: isZh ? "星际论坛" : "Topics",
        link: buildPath(base, "/topics/")
      },
      ...topics.map((topic) => ({
        text: isZh ? topic.titleZh : topic.titleEn,
        link: buildPath(base, `/topics/${topic.slug}`)
      })),
      {
        text: isZh ? "产品总览" : "Product Entry Points",
        link: buildPath(base, "/products/")
      },
      {
        text: isZh ? "核心中间件矩阵图" : "Dependency Blueprint",
        link: "/core-middleware-matrix-dependency.svg"
      }
    ]
  }
];

const buildLocaleTheme = (base: string, isZh: boolean) => {
  const homeLink = buildHomeLink(base);

  return {
    siteTitle: isZh ? "星际枢纽" : "Stell Hub",
    logo: "/logo/logo.png",
    nav: [
      {
        text: isZh ? "首页" : "Home",
        link: homeLink
      },
      {
        text: isZh ? "星际论坛" : "Topics",
        link: buildPath(base, "/topics/")
      },
      {
        text: isZh ? "核心产品" : "Products",
        items: buildProductMenu(base, isZh)
      }
    ],
    sidebar: {
      ...Object.fromEntries(
        products.map((product) => [
          buildPath(base, `/products/${product.slug}/`),
          buildProductSidebar(base, isZh, product.slug)
        ])
      ),
      [buildPath(base, "/products/")]: [
        {
          text: isZh ? "核心产品矩阵" : "Product Entry Points",
          items: buildProductMenu(base, isZh)
        }
      ],
      [buildPath(base, "/topics/")]: [
        {
          text: isZh ? "星际论坛" : "Topics",
          items: buildTopicMenu(base, isZh)
        }
      ],
      [homeLink]: buildRootSidebar(base, isZh)
    },
    outline: {
      level: [2, 3],
      label: isZh ? "页面导航" : "On This Page"
    },
    search: {
      provider: "local"
    },
    socialLinks: [{ icon: "github", link: "https://github.com/stellhub" }],
    docFooter: {
      prev: isZh ? "上一页" : "Previous page",
      next: isZh ? "下一页" : "Next page"
    },
    footer: {
      message: "Released under the Apache License 2.0.",
      copyright: "Copyright © Stell Hub"
    },
    lastUpdated: {
      text: isZh ? "最后更新" : "Last updated"
    },
    darkModeSwitchLabel: isZh ? "主题切换" : "Appearance",
    lightModeSwitchTitle: isZh ? "切换到浅色模式" : "Switch to light theme",
    darkModeSwitchTitle: isZh ? "切换到深色模式" : "Switch to dark theme",
    sidebarMenuLabel: isZh ? "菜单" : "Menu",
    returnToTopLabel: isZh ? "返回顶部" : "Return to top"
  };
};

export default defineConfig({
  title: "Stell Hub",
  description:
    "English-first VitePress documentation for Stell Hub infrastructure, middleware products, and platform engineering notes.",
  cleanUrls: true,
  lastUpdated: true,
  head: [["link", { rel: "icon", type: "image/png", href: "/logo/logo.png" }]],
  vite: {
    server: {
      allowedHosts: [".stellhub.top"]
    }
  },
  locales: {
    root: {
      label: "English",
      lang: "en-US",
      title: "Stell Hub",
      description:
        "Infrastructure and middleware documentation covering the Stell Hub product matrix and platform notes.",
      themeConfig: buildLocaleTheme("", false)
    },
    zh: {
      label: "简体中文",
      lang: "zh-CN",
      link: "/zh/",
      title: "星际枢纽",
      description: "Stell Hub（星际枢纽）聚焦基础架构与微服务中间件最佳实践研究。",
      themeConfig: buildLocaleTheme("/zh", true)
    }
  }
});

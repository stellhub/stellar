import { defineConfig } from "vitepress";

const productItems = [
  { text: "产品总览", link: "/products/" },
  { text: "Stellmap · 星图", link: "/products/stellmap" },
  { text: "Stellnula · 星云", link: "/products/stellnula" },
  { text: "Stelltrace · 星迹", link: "/products/stelltrace" },
  { text: "Stellorbit · 星轨", link: "/products/stellorbit" },
  { text: "Stellpulse · 脉冲", link: "/products/stellpulse" },
  { text: "Stellabe · 星盘", link: "/products/stellabe" },
  { text: "Stellpoint · 奇点", link: "/products/stellpoint" },
  { text: "Stellgate · 视界", link: "/products/stellgate" },
  { text: "Stellflow · 彗流", link: "/products/stellflow" },
  { text: "Stellspec · 星谱", link: "/products/stellspec" },
  { text: "Stellcon · 星座", link: "/products/stellcon" },
  { text: "Stellvox · 星讯", link: "/products/stellvox" },
  { text: "Stellguard · 星盾", link: "/products/stellguard" },
  { text: "Stellkey · 星钥", link: "/products/stellkey" }
];

const productSectionItems = (link: string) => [
  { text: "产品定位", link: `${link}#产品定位` },
  { text: "核心能力", link: `${link}#核心能力` },
  { text: "概要设计", link: `${link}#概要设计` },
  { text: "架构组成", link: `${link}#架构组成` },
  { text: "部署形态", link: `${link}#部署形态` },
  { text: "快速入门", link: `${link}#快速入门` },
  { text: "配置建议", link: `${link}#配置建议` },
  { text: "API 与 SDK", link: `${link}#api-与-sdk` },
  { text: "可观测性", link: `${link}#可观测性` },
  { text: "典型场景", link: `${link}#典型场景` }
];

const productDetailSidebar = (link: string) => [
  {
    text: "当前页面",
    items: productSectionItems(link)
  }
];

export default defineConfig({
  lang: "zh-CN",
  title: "星级枢纽",
  description: "Stell Hub（星级枢纽）体系的品牌、命名、仓库布局与运行规范文档站点。",
  cleanUrls: true,
  lastUpdated: true,
  head: [
    ["link", { rel: "icon", type: "image/png", href: "/logo/logo.png" }]
  ],
  vite: {
    server: {
      allowedHosts: [".stellhub.top"]
    }
  },
  themeConfig: {
    siteTitle: "星级枢纽",
    logo: "/logo/logo.png",
    nav: [
      { text: "首页", link: "/" },
      { text: "体系总览", link: "/overview" },
      { text: "核心产品", items: productItems },
      { text: "环境变量规范", link: "/environment-variable-spec" },
      { text: "Header 与指标规范", link: "/request-header-and-metrics-spec" }
    ],
    sidebar: {
      "/products/stellmap": productDetailSidebar("/products/stellmap"),
      "/products/stellnula": productDetailSidebar("/products/stellnula"),
      "/products/stelltrace": productDetailSidebar("/products/stelltrace"),
      "/products/stellorbit": productDetailSidebar("/products/stellorbit"),
      "/products/stellpulse": productDetailSidebar("/products/stellpulse"),
      "/products/stellabe": productDetailSidebar("/products/stellabe"),
      "/products/stellpoint": productDetailSidebar("/products/stellpoint"),
      "/products/stellgate": productDetailSidebar("/products/stellgate"),
      "/products/stellflow": productDetailSidebar("/products/stellflow"),
      "/products/stellspec": productDetailSidebar("/products/stellspec"),
      "/products/stellcon": productDetailSidebar("/products/stellcon"),
      "/products/stellvox": productDetailSidebar("/products/stellvox"),
      "/products/stellguard": productDetailSidebar("/products/stellguard"),
      "/products/stellkey": productDetailSidebar("/products/stellkey"),
      "/products/": [
        {
          text: "核心产品矩阵",
          items: productItems
        }
      ],
      "/": [
        {
          text: "文档导航",
          items: [
            { text: "首页", link: "/" },
            { text: "体系总览", link: "/overview" },
            { text: "产品总览", link: "/products/" },
            { text: "基础应用环境变量规范", link: "/environment-variable-spec" },
            { text: "全局请求 Header 与指标规范", link: "/request-header-and-metrics-spec" },
            { text: "核心中间件矩阵图", link: "/core-middleware-matrix-dependency.svg" }
          ]
        }
      ]
    },
    outline: {
      level: [2, 3]
    },
    search: {
      provider: "local"
    },
    socialLinks: [
      { icon: "github", link: "https://github.com/stellhub" }
    ],
    docFooter: {
      prev: "上一页",
      next: "下一页"
    },
    footer: {
      message: "Released under the Apache License 2.0.",
      copyright: "Copyright © Stell Hub"
    }
  }
});

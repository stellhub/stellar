import { defineConfig } from "vitepress";

const productItems = [
  {
    text: "核心产品矩阵",
    items: [
      { text: "产品总览", link: "/products/" },
      { text: "Stellmap · 星图", link: "/products/stellmap/" },
      { text: "Stellnula · 星云", link: "/products/stellnula/" },
      { text: "Stelltrace · 星迹", link: "/products/stelltrace/" },
      { text: "Stellorbit · 星轨", link: "/products/stellorbit/" },
      { text: "Stellpulse · 脉冲", link: "/products/stellpulse/" },
      { text: "Stellabe · 星盘", link: "/products/stellabe/" },
      { text: "Stellpoint · 奇点", link: "/products/stellpoint/" },
      { text: "Stellgate · 视界", link: "/products/stellgate/" },
      { text: "Stellflow · 彗流", link: "/products/stellflow/" },
      { text: "Stellspec · 星谱", link: "/products/stellspec/" },
      { text: "Stellcon · 星座", link: "/products/stellcon/" },
      { text: "Stellvox · 星讯", link: "/products/stellvox/" },
      { text: "Stellguard · 星盾", link: "/products/stellguard/" },
      { text: "Stellkey · 星钥", link: "/products/stellkey/" }
    ]
  }
];

const topicItems = [
  { text: "论坛首页", link: "/topics/" },
  { text: "最佳 DSL 语言：CUE", link: "/topics/dsl" },
  { text: "面向超大型企业的微服务命名体系研究", link: "/topics/service-naming" },
  { text: "可观测规范", link: "/topics/observability-spec" },
  { text: "大型企业跨语言微服务链路追踪技术调研方案", link: "/topics/traces" },
  { text: "错误码规范", link: "/topics/error-code-spec" },
  { text: "为什么企业要自研中间件", link: "/topics/middleware-evolution" },
  { text: "分布式系统中的一致性挑战及其解决路径", link: "/topics/distributed-consistency" },
  { text: "分布式系统注册中心意义、问题与主流实现", link: "/topics/distributed-system-registry-centers" }
];

const productDetailSidebar = (link: string) => [
  {
    text: "设计文档",
    items: [
      { text: "产品首页", link },
      { text: "概要设计", link: `${link}summary-design` },
      { text: "架构组成", link: `${link}architecture` },
      { text: "部署形态", link: `${link}deployment` },
      { text: "快速入门", link: `${link}quick-start` },
      { text: "配置建议", link: `${link}configuration` },
      { text: "API 与 SDK", link: `${link}api-and-sdk` },
      { text: "可观测性", link: `${link}observability` }
    ]
  }
];

export default defineConfig({
  lang: "zh-CN",
  title: "星际枢纽",
  description: "Stell Hub（星际枢纽）聚焦基础架构与微服务中间件最佳实践研究。",
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
    siteTitle: "星际枢纽",
    logo: "/logo/logo.png",
    nav: [
      { text: "首页", link: "/" },
      { text: "星际论坛", link: "/topics/" },
      { text: "核心产品", items: productItems }
    ],
    sidebar: {
      "/products/stellmap/": productDetailSidebar("/products/stellmap/"),
      "/products/stellnula/": productDetailSidebar("/products/stellnula/"),
      "/products/stelltrace/": productDetailSidebar("/products/stelltrace/"),
      "/products/stellorbit/": productDetailSidebar("/products/stellorbit/"),
      "/products/stellpulse/": productDetailSidebar("/products/stellpulse/"),
      "/products/stellabe/": productDetailSidebar("/products/stellabe/"),
      "/products/stellpoint/": productDetailSidebar("/products/stellpoint/"),
      "/products/stellgate/": productDetailSidebar("/products/stellgate/"),
      "/products/stellflow/": productDetailSidebar("/products/stellflow/"),
      "/products/stellspec/": productDetailSidebar("/products/stellspec/"),
      "/products/stellcon/": productDetailSidebar("/products/stellcon/"),
      "/products/stellvox/": productDetailSidebar("/products/stellvox/"),
      "/products/stellguard/": productDetailSidebar("/products/stellguard/"),
      "/products/stellkey/": productDetailSidebar("/products/stellkey/"),
      "/products/": [
        {
          text: "核心产品矩阵",
          items: productItems[0].items
        }
      ],
      "/topics/": [
        {
          text: "星际论坛",
          items: topicItems
        }
      ],
      "/": [
        {
          text: "文档导航",
          items: [
            { text: "首页", link: "/" },
            { text: "星际论坛", link: "/topics/" },
            { text: "最佳 DSL 语言：CUE", link: "/topics/dsl" },
            { text: "面向超大型企业的微服务命名体系研究", link: "/topics/service-naming" },
            { text: "可观测规范", link: "/topics/observability-spec" },
            { text: "大型企业跨语言微服务链路追踪技术调研方案", link: "/topics/traces" },
            { text: "错误码规范", link: "/topics/error-code-spec" },
            { text: "为什么企业要自研中间件", link: "/topics/middleware-evolution" },
            { text: "分布式系统中的一致性挑战及其解决路径", link: "/topics/distributed-consistency" },
            { text: "分布式系统注册中心意义、问题与主流实现", link: "/topics/distributed-system-registry-centers" },
            { text: "产品总览", link: "/products/" },
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

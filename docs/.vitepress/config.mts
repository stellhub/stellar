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
  },
  {
    text: "体系规范",
    items: [
      { text: "可观测规范", link: "/products/observability-spec" }
    ]
  }
];

const topicItems = [
  { text: "专题总览", link: "/topics/" },
  { text: "为什么企业要自研中间件", link: "/topics/middleware-evolution" },
  { text: "分布式系统中的一致性挑战及其解决路径", link: "/topics/distributed-consistency" }
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

const observabilitySidebar = [
  {
    text: "可观测规范",
    items: [
      { text: "规范定位", link: "/products/observability-spec#规范定位" },
      { text: "基础环境变量规范", link: "/products/observability-spec#第一部分基础环境变量规范" },
      { text: "请求上下文规范", link: "/products/observability-spec#第二部分全局请求上下文规范" },
      { text: "全局指标规范", link: "/products/observability-spec#第三部分全局指标规范" },
      { text: "落地约束", link: "/products/observability-spec#第四部分平台与业务落地约束" }
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
      { text: "专题研究", link: "/topics/" },
      { text: "体系总览", link: "/overview" },
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
      "/products/observability-spec": observabilitySidebar,
      "/products/": [
        {
          text: "核心产品矩阵",
          items: productItems[0].items
        },
        {
          text: "体系规范",
          items: productItems[1].items
        }
      ],
      "/topics/": [
        {
          text: "专题研究",
          items: topicItems
        }
      ],
      "/": [
        {
          text: "文档导航",
          items: [
            { text: "首页", link: "/" },
            { text: "专题研究总览", link: "/topics/" },
            { text: "为什么企业要自研中间件", link: "/topics/middleware-evolution" },
            { text: "分布式系统中的一致性挑战及其解决路径", link: "/topics/distributed-consistency" },
            { text: "体系总览", link: "/overview" },
            { text: "产品总览", link: "/products/" },
            { text: "可观测规范", link: "/products/observability-spec" },
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

---
layout: home
hero:
  name: Stell Notes
  text: 星际论坛
  tagline: 这里更像一页持续更新的博客索引，用来记录基础架构与微服务中间件在演进过程中遇到的问题、判断过程，以及最后采用了什么样的解决路径。
  actions:
    - theme: brand
      text: 浏览全部文章
      link: "#forum-latest"
    - theme: alt
      text: 返回站点首页
      link: /zh/
---

<script setup>
import ForumPostIndex from "../../.vitepress/theme/components/ForumPostIndex.vue";
import { data as topicPosts } from "./posts.data.ts";
</script>

<div id="forum-latest">
  <ForumPostIndex :posts="topicPosts" />
</div>

## 这页记录什么

- 中间件演进过程中遇到的典型问题
- 某一类基础设施能力为什么会出现、边界在哪里
- 主流解决方案分别做了什么取舍
- 在真实工程里应该怎样判断和落地

## 后续会继续补充

- 服务治理与流量控制
- 配置中心与注册中心的边界变化
- 消息系统从队列到事件流平台的迁移
- 基础架构团队与平台工程的协作方式

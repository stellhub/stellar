<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useData } from "vitepress";
import type { TopicPost } from "../../../topics/posts.data";

const props = defineProps<{
  posts: TopicPost[];
}>();

const { lang } = useData();

const isChinese = computed(() => lang.value.startsWith("zh"));
const copy = computed(() => ({
  allTag: isChinese.value ? "全部" : "All",
  filterTitle: isChinese.value ? "按标签筛选" : "Filter by tag",
  latestTitle: isChinese.value ? "最近更新" : "Recently updated",
  publishLabel: isChinese.value ? "发布" : "Published",
  updateLabel: isChinese.value ? "更新" : "Updated",
  readingDirectionLabel: isChinese.value ? "阅读方向：" : "Reading direction:",
  emptyState: isChinese.value ? "当前标签下还没有文章。" : "No posts match the selected tag yet.",
  categoryTitle: isChinese.value ? "按主题分类" : "Browse by category",
  categoryDesc: isChinese.value
    ? "按问题域归拢，而不是按传统专题策展方式组织。"
    : "Grouped by problem domain instead of a traditional editorial sequence.",
  archiveTitle: isChinese.value ? "按年份归档" : "Archive by year",
  writingScopeTitle: isChinese.value ? "写作范围" : "Writing scope",
  unknownYear: isChinese.value ? "未知" : "Unknown"
}));

const activeTag = ref(copy.value.allTag);

watch(
  () => copy.value.allTag,
  (allTag) => {
    if (!availableTags.value.includes(activeTag.value)) {
      activeTag.value = allTag;
    }
  }
);

const formatDate = (value: string) => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat(isChinese.value ? "zh-CN" : "en-US", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit"
  }).format(date);
};

const availableTags = computed(() => {
  const values = new Set<string>();

  props.posts.forEach((post) => {
    post.tags.forEach((tag) => values.add(tag));
  });

  return [copy.value.allTag, ...Array.from(values)];
});

const filteredPosts = computed(() => {
  if (activeTag.value === copy.value.allTag) {
    return props.posts;
  }

  return props.posts.filter((post) => post.tags.includes(activeTag.value));
});

const groupedPosts = computed(() => {
  const groups = new Map<string, TopicPost[]>();

  filteredPosts.value.forEach((post) => {
    const current = groups.get(post.category) ?? [];
    current.push(post);
    groups.set(post.category, current);
  });

  return Array.from(groups.entries()).map(([category, posts]) => ({
    category,
    posts
  }));
});

const archivedPosts = computed(() => {
  const groups = new Map<string, TopicPost[]>();

  filteredPosts.value.forEach((post) => {
    const date = new Date(post.publishAt);
    const year = Number.isNaN(date.getTime()) ? copy.value.unknownYear : String(date.getFullYear());
    const current = groups.get(year) ?? [];
    current.push(post);
    groups.set(year, current);
  });

  return Array.from(groups.entries())
    .sort((left, right) => right[0].localeCompare(left[0]))
    .map(([year, posts]) => ({
      year,
      posts
    }));
});

const setActiveTag = (tag: string) => {
  activeTag.value = tag;
};
</script>

<template>
  <div class="forum-index">
    <section class="forum-section forum-section--compact">
      <div class="forum-section__head">
        <h2>{{ copy.filterTitle }}</h2>
      </div>
      <div class="forum-topic-chips">
        <button
          v-for="tag in availableTags"
          :key="tag"
          type="button"
          class="forum-topic-chip"
          :class="{ 'is-active': activeTag === tag }"
          @click="setActiveTag(tag)"
        >
          {{ tag }}
        </button>
      </div>
    </section>

    <section class="forum-section">
      <div class="forum-section__head">
        <h2>{{ copy.latestTitle }}</h2>
      </div>
      <div class="forum-post-grid">
        <a
          v-for="post in filteredPosts"
          :key="post.url"
          class="forum-post-card"
          :href="post.url"
        >
          <div class="forum-post-card__meta">
            <span class="forum-post-card__category">{{ post.category }}</span>
            <span class="forum-post-card__date">{{ copy.publishLabel }} {{ formatDate(post.publishAt) }}</span>
            <span class="forum-post-card__date">{{ copy.updateLabel }} {{ formatDate(post.updatedAt) }}</span>
          </div>
          <h3>{{ post.title }}</h3>
          <p class="forum-post-card__summary">{{ post.summary }}</p>
          <p class="forum-post-card__direction">
            <strong>{{ copy.readingDirectionLabel }}</strong>{{ post.readingDirection }}
          </p>
          <div class="forum-post-card__tags">
            <span v-for="tag in post.tags" :key="tag">{{ tag }}</span>
          </div>
        </a>
      </div>
      <p v-if="!filteredPosts.length" class="forum-empty-state">{{ copy.emptyState }}</p>
    </section>

    <section class="forum-section">
      <div class="forum-section__head">
        <h2>{{ copy.categoryTitle }}</h2>
        <p>{{ copy.categoryDesc }}</p>
      </div>
      <div class="forum-category-grid">
        <div v-for="group in groupedPosts" :key="group.category" class="forum-category-card">
          <h3>{{ group.category }}</h3>
          <ul>
            <li v-for="post in group.posts" :key="post.url">
              <a :href="post.url">{{ post.title }}</a>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <section class="forum-section">
      <div class="forum-section__head">
        <h2>{{ copy.archiveTitle }}</h2>
      </div>
      <div class="forum-archive-list">
        <div v-for="archive in archivedPosts" :key="archive.year" class="forum-archive-card">
          <h3>{{ archive.year }}</h3>
          <ul>
            <li v-for="post in archive.posts" :key="post.url">
              <a :href="post.url">{{ post.title }}</a>
              <span>{{ formatDate(post.publishAt) }}</span>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <section class="forum-section forum-section--compact">
      <div class="forum-section__head">
        <h2>{{ copy.writingScopeTitle }}</h2>
      </div>
      <div class="forum-topic-chips">
        <span v-for="tag in availableTags.filter((tag) => tag !== copy.allTag)" :key="tag">{{ tag }}</span>
      </div>
    </section>
  </div>
</template>

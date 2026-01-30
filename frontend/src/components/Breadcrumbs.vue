<template>
  <div class="breadcrumbs">
    <component
      :is="element"
      :to="base || ''"
      :aria-label="t('files.home')"
      :title="t('files.home')"
    >
      <i class="material-icons home-icon">home</i>
    </component>

    <span v-for="(link, index) in items" :key="index">
      <span class="chevron"
        ><i class="material-icons">keyboard_arrow_right</i></span
      >
      <component :is="element" :to="link.url">{{ link.name }}</component>
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";

const { t } = useI18n();

const route = useRoute();

const props = defineProps<{
  base: string;
  noLink?: boolean;
}>();

const items = computed(() => {
  const relativePath = route.path.replace(props.base, "");
  const parts = relativePath.split("/");

  if (parts[0] === "") {
    parts.shift();
  }

  if (parts[parts.length - 1] === "") {
    parts.pop();
  }

  const breadcrumbs: BreadCrumb[] = [];

  for (let i = 0; i < parts.length; i++) {
    if (i === 0) {
      breadcrumbs.push({
        name: decodeURIComponent(parts[i]),
        url: props.base + "/" + parts[i] + "/",
      });
    } else {
      breadcrumbs.push({
        name: decodeURIComponent(parts[i]),
        url: breadcrumbs[i - 1].url + parts[i] + "/",
      });
    }
  }

  if (breadcrumbs.length > 3) {
    while (breadcrumbs.length !== 4) {
      breadcrumbs.shift();
    }

    breadcrumbs[0].name = "...";
  }

  return breadcrumbs;
});

const element = computed(() => {
  if (props.noLink) {
    return "span";
  }

  return "router-link";
});
</script>

<style scoped>
.breadcrumbs {
  display: flex;
  align-items: center;
  padding: 0.75em 1em;
  background: var(--surfaceSecondary);
  border-radius: 8px;
  margin-bottom: 1em;
  font-size: 0.9em;
  overflow-x: auto;
  white-space: nowrap;
  scrollbar-width: none;
}

.breadcrumbs::-webkit-scrollbar {
  display: none;
}

.breadcrumbs a, .breadcrumbs span {
  color: var(--textSecondary);
  text-decoration: none;
  display: flex;
  align-items: center;
  transition: all 0.2s;
}

.breadcrumbs a:hover {
  color: var(--blue);
  background: rgba(33, 150, 243, 0.05);
  border-radius: 4px;
}

.breadcrumbs .chevron {
  margin: 0 4px;
  color: var(--divider);
  opacity: 0.5;
}

.breadcrumbs i {
  font-size: 1.25em;
}

.breadcrumbs .home-icon {
  color: var(--textSecondary);
}

.breadcrumbs a.router-link-active {
  color: var(--textPrimary);
  font-weight: 500;
}
</style>

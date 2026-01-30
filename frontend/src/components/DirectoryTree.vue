<template>
  <div class="directory-tree-container">
    <div class="tree-header">
      <i class="material-icons">account_tree</i>
      <span>{{ rootName }}</span>
    </div>
    <div class="tree-content" ref="scrollContainer">
      <TreeItem
        v-for="node in visibleNodes"
        :key="node.path"
        :name="node.name"
        :path="node.path"
        :depth="node.depth"
        :loading="node.loading"
        :expanded="expandedPaths[node.path]"
        :is-active="isActive(node.path)"
        @toggle="toggleNode(node.path)"
        @navigate="navigate(node.path)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useTreeStore } from '@/stores/tree';
import TreeItem from './TreeItem.vue';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const treeStore = useTreeStore();

const scrollContainer = ref<HTMLElement | null>(null);

const isShare = computed(() => route.path.startsWith('/share'));
const currentShareHash = computed(() => {
  if (!isShare.value) return null;
  const parts = route.path.split('/');
  return parts[2];
});

const rootName = computed(() => isShare.value ? (t('files.shareRoot') || 'Share Root') : (t('files.home') || 'My Files'));

const expandedPaths = computed(() => treeStore.expandedPaths);

const getCurrentPath = () => {
  if (isShare.value) {
    const parts = route.path.split('/').filter(p => p);
    if (parts.length > 2) return '/' + parts.slice(2).join('/');
    return '/';
  }
  return route.path.replace(/^\/files/, '') || '/';
};

const isActive = (path: string) => {
  const current = getCurrentPath();
  const normalizedPath = path === "" ? "/" : path;
  return current === normalizedPath;
};

// Calculate flat list of visible nodes
const visibleNodes = computed(() => {
  const list: any[] = [];
  
  function walk(path: string, depth: number) {
    const node = treeStore.nodes[path];
    if (!node) {
      // If root isn't loaded yet, add a placeholder
      if (path === "") {
        list.push({ path: "", name: rootName.value, depth: 0, loading: true });
      }
      return;
    }
    
    list.push({
      path: node.path,
      name: node.path === "" ? rootName.value : node.name,
      depth,
      loading: node.loading
    });

    if (expandedPaths.value[node.path] && node.children) {
      node.children.forEach(childPath => walk(childPath, depth + 1));
    }
  }

  walk("", 0);
  return list;
});

const toggleNode = (path: string) => {
  treeStore.toggleNode(path, isShare.value, currentShareHash.value);
};

const navigate = (path: string) => {
  const targetPath = isShare.value 
    ? `/share/${currentShareHash.value}${path === "" ? "" : path}`
    : `/files${path}`;
    
  if (route.path === targetPath) {
    // Already active: Toggle expansion
    toggleNode(path);
  } else {
    // Navigating to new path: Ensure it's expanded (to show children/trigger lazy load)
    if (!expandedPaths.value[path]) {
       toggleNode(path);
    }
    router.push({ path: targetPath });
  }
};

// Initial load and auto-expand logic
const handleRouteChange = async () => {
  const path = getCurrentPath();
  const normalizedPath = path === "/" ? "" : path;
  
  // Ensure the tree store has the root
  if (!treeStore.nodes[""]) {
    await treeStore.fetchNode("", isShare.value, currentShareHash.value);
  }

  // Auto-expand parents of current path
  if (normalizedPath !== "") {
    await treeStore.ensurePathExpanded(normalizedPath, isShare.value, currentShareHash.value);
  }
};

watch(() => route.path, handleRouteChange, { immediate: true });

onMounted(() => {
  if (!treeStore.nodes[""]) {
    treeStore.fetchNode("", isShare.value, currentShareHash.value);
  }
});
</script>

<style scoped>
.directory-tree-container {
  margin-top: 0.5rem;
  display: flex;
  flex-direction: column;
  min-height: 0;
  flex: 1;
}

.tree-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0.5rem 0.75rem;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--textTertiary);
  font-weight: 700;
  flex-shrink: 0;
}

.tree-header i {
  font-size: 16px;
}

.tree-content {
  overflow-x: hidden;
  overflow-y: auto;
  flex: 1;
  padding-right: 4px;
}

/* Scrollbar for tree */
.tree-content::-webkit-scrollbar {
  width: 3px;
}

.tree-content::-webkit-scrollbar-thumb {
  background: var(--divider);
  border-radius: 10px;
}
</style>

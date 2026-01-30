<template>
  <div class="tree-item" :class="{ 'is-active': isActive }">
    <div 
      class="tree-item-label" 
      :style="{ paddingLeft: depth * 16 + 10 + 'px' }"
      @click="$emit('navigate')"
    >
      <i 
        class="material-icons toggle-icon" 
        :class="{ 'is-expanded': expanded, 'is-invisible-custom': !hasInteraction }"
        @click.stop="$emit('toggle')"
      >
        chevron_right
      </i>
      <i class="material-icons folder-icon">{{ expanded ? 'folder_open' : 'folder' }}</i>
      <span class="name">{{ name }}</span>
      <div v-if="loading" class="spinner-small"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

defineProps<{
  name: string;
  path: string;
  depth: number;
  expanded?: boolean;
  loading?: boolean;
  isActive?: boolean;
}>();

defineEmits(['toggle', 'navigate']);

// We assume it's a folder for now in Tree View
const hasInteraction = computed(() => true);
</script>

<style scoped>
.tree-item {
  user-select: none;
}

.tree-item-label {
  display: flex;
  align-items: center;
  padding: 6px 10px;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--textSecondary);
  gap: 8px;
  margin: 1px 4px;
}

.tree-item-label:hover {
  background: var(--surfaceSecondary);
  color: var(--textPrimary);
}

.tree-item.is-active > .tree-item-label {
  background: linear-gradient(90deg, rgba(33, 150, 243, 0.15) 0%, rgba(33, 150, 243, 0.05) 100%);
  color: var(--blue);
  font-weight: 600;
  box-shadow: inset 3px 0 0 var(--blue);
}

.toggle-icon {
  font-size: 18px;
  transition: transform 0.2s;
  color: var(--textSecondary);
  opacity: 0.6;
}

.toggle-icon:hover {
  opacity: 1;
  background: var(--divider);
  border-radius: 4px;
}

.toggle-icon.is-expanded {
  transform: rotate(90deg);
}

.is-invisible-custom {
  visibility: hidden;
}

.folder-icon {
  font-size: 20px;
  color: #ffca28; /* Material yellow */
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
}

.name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 0.9em;
  flex: 1;
}

.spinner-small {
  width: 12px;
  height: 12px;
  border: 2px solid var(--divider);
  border-top-color: var(--blue);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-left: auto;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>

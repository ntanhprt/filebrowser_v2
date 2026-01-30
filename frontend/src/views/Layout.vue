<template>
  <div class="app-layout">
    <div v-if="uploadStore.totalBytes" class="upload-progress-global">
      <div
        class="progress-bar"
        v-bind:style="{
          width: sentPercent + '%',
        }"
      ></div>
    </div>
    
    <sidebar v-if="showSidebar"></sidebar>
    
    <div class="content-wrapper" :class="{ 
      'sidebar-collapsed': !showSidebar || !layoutStore.showSidebar,
      'no-sidebar': !showSidebar
    }" :style="layoutStore.showPreviewPanel && !isShareRoute ? { marginRight: previewPanelWidth + 'px' } : {}">
      <main>
        <router-view></router-view>
        <shell
          v-if="enableExec && authStore.isLoggedIn && authStore.user?.perm.execute"
        />
      </main>
      
      <div
        v-if="layoutStore.showPreviewPanel && !isShareRoute"
        class="preview-resize-handle"
        :style="{ right: previewPanelWidth + 'px' }"
        @mousedown="startResize"
      ></div>
      
      <preview-panel 
        v-if="layoutStore.showPreviewPanel && !isShareRoute" 
        :style="{ width: previewPanelWidth + 'px' }" 
      />
    </div>
    
    <prompts></prompts>
    <upload-files></upload-files>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { useFileStore } from "@/stores/file";
import { useUploadStore } from "@/stores/upload";
import Sidebar from "@/components/Sidebar.vue";
import Prompts from "@/components/prompts/Prompts.vue";
import Shell from "@/components/Shell.vue";
import UploadFiles from "@/components/prompts/UploadFiles.vue";
import PreviewPanel from "@/components/PreviewPanel.vue";
import { enableExec } from "@/utils/constants";
import { computed, watch, ref, onUnmounted } from "vue";
import { useRoute } from "vue-router";

const layoutStore = useLayoutStore();
const authStore = useAuthStore();
const fileStore = useFileStore();
const uploadStore = useUploadStore();
const route = useRoute();

const previewPanelWidth = ref(parseInt(localStorage.getItem('previewPanelWidth') || '400'));
const isResizing = ref(false);

const sentPercent = computed(() =>
  ((uploadStore.sentBytes / uploadStore.totalBytes) * 100).toFixed(2)
);

const isShareRoute = computed(() => route.path.startsWith('/share'));

// Hide sidebar for share routes when user is not logged in
const showSidebar = computed(() => {
  if (isShareRoute.value && !authStore.isLoggedIn) {
    return false;
  }
  return true;
});

const startResize = (e: MouseEvent) => {
  isResizing.value = true;
  document.addEventListener('mousemove', doResize);
  document.addEventListener('mouseup', stopResize);
  document.body.style.cursor = 'col-resize';
  e.preventDefault();
};

const doResize = (e: MouseEvent) => {
  if (!isResizing.value) return;
  const newWidth = window.innerWidth - e.clientX;
  if (newWidth > 200 && newWidth < window.innerWidth * 0.8) {
    previewPanelWidth.value = newWidth;
    localStorage.setItem('previewPanelWidth', String(newWidth));
  }
};

const stopResize = () => {
  isResizing.value = false;
  document.removeEventListener('mousemove', doResize);
  document.removeEventListener('mouseup', stopResize);
  document.body.style.cursor = '';
};

onUnmounted(() => {
  document.removeEventListener('mousemove', doResize);
  document.removeEventListener('mouseup', stopResize);
});

watch(route, () => {
  fileStore.selected = [];
  fileStore.multiple = false;
  if (layoutStore.currentPromptName !== "success") {
    layoutStore.closeHovers();
  }
});
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.content-wrapper {
  flex: 1;
  margin-left: 18em;
  transition: margin-left 0.3s;
  display: flex;
  min-width: 0;
}

.content-wrapper.sidebar-collapsed {
  margin-left: 0;
}

.content-wrapper.no-sidebar {
  margin-left: 0;
}

main {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.preview-resize-handle {
  width: 4px;
  cursor: col-resize;
  background: var(--divider);
  z-index: 2000;
  position: fixed;
  top: 4em;
  bottom: 0;
}

.preview-resize-handle:hover {
  background: var(--blue);
}

.upload-progress-global {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  z-index: 2000;
}

.progress-bar {
  height: 100%;
  background: var(--blue);
}

@media (max-width: 1024px) {
  .content-wrapper {
    margin-left: 0 !important;
  }
}
</style>

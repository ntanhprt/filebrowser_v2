<template>
  <div class="file-listing-page">
    <header-bar showMenu showLogo @logo-click="layoutStore.toggleSidebar()">
      <action icon="arrow_back" :label="t('buttons.back')" @action="router.back()" :disabled="!canGoBack" />
      <action icon="arrow_forward" :label="t('buttons.forward')" @action="router.forward()" />
      
      <search />

      <template #actions>
        <action v-if="headerButtons.share" icon="share" :label="t('buttons.share')" show="share" />
        <action v-if="headerButtons.rename" icon="mode_edit" :label="t('buttons.rename')" show="rename" />
        <action v-if="headerButtons.copy" icon="content_copy" :label="t('buttons.copyFile')" show="copy" />
        <action v-if="headerButtons.move" icon="forward" :label="t('buttons.moveFile')" show="move" />
        <action v-if="headerButtons.delete" icon="delete" :label="t('buttons.delete')" show="delete" />
        
        <action v-if="headerButtons.shell" icon="terminal" @action="layoutStore.toggleShell" />
        <action :icon="layoutStore.showPreviewPanel ? 'visibility_off' : 'visibility'" @action="layoutStore.togglePreviewPanel" />
        <action :icon="viewIcon" @action="switchView" />
        <action v-if="headerButtons.download" icon="file_download" @action="download" />
        <action v-if="headerButtons.upload" icon="file_upload" @action="uploadFunc" />
      </template>
    </header-bar>

    <div class="listing-container" ref="listingContainer" @scroll="scrollEvent">
      <div v-if="layoutStore.loading" class="spinner-container">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
      </div>

      <div v-else id="listing" class="file-icons" :class="viewModeClass" @contextmenu="handleContextMenu">
        <div v-if="isEmpty" class="empty-state">
          <i class="material-icons">folder_open</i>
          <p>{{ t("files.lonely") }}</p>
        </div>

        <template v-else>
          <div class="item header">
            <div></div>
            <div>
              <p :class="{ active: nameSorted }" class="name" @click="sort('name')">
                <span>{{ t("files.name") }}</span>
                <i class="material-icons">{{ nameIcon }}</i>
              </p>
              <p :class="{ active: sizeSorted }" class="size" @click="sort('size')">
                <span>{{ t("files.size") }}</span>
                <i class="material-icons">{{ sizeIcon }}</i>
              </p>
              <p :class="{ active: modifiedSorted }" class="modified" @click="sort('modified')">
                <span>{{ t("files.lastModified") }}</span>
                <i class="material-icons">{{ modifiedIcon }}</i>
              </p>
            </div>
          </div>

          <item
            v-for="(item, index) in visibleItems"
            :key="item.path"
            v-bind="item"
            :index="item.originalIndex"
          />
        </template>
      </div>
    </div>

    <input style="display: none" type="file" id="upload-input" @change="uploadInput" multiple />
    <input style="display: none" type="file" id="upload-folder-input" @change="uploadInput" webkitdirectory multiple />

    <context-menu :show="isContextMenuVisible" :pos="contextMenuPos" @hide="isContextMenuVisible = false">
      <div v-if="headerButtons.upload" class="action" @click="uploadFunc(); isContextMenuVisible = false">
        <i class="material-icons">file_upload</i>
        <span>{{ t('buttons.upload') }}</span>
      </div>
      <div v-if="headerButtons.upload" class="action" @click="uploadFolderFunc(); isContextMenuVisible = false">
        <i class="material-icons">create_new_folder</i>
        <span>{{ t('buttons.uploadFolder') || 'Upload Folder' }}</span>
      </div>
      <hr v-if="headerButtons.upload" />
      <div v-if="headerButtons.share" class="action" @click="layoutStore.showHover('share'); isContextMenuVisible = false">
        <i class="material-icons">share</i>
        <span>{{ t('buttons.share') }}</span>
      </div>
      <div v-if="headerButtons.rename" class="action" @click="layoutStore.showHover('rename'); isContextMenuVisible = false">
        <i class="material-icons">mode_edit</i>
        <span>{{ t('buttons.rename') }}</span>
      </div>
      <div v-if="headerButtons.copy" class="action" @click="layoutStore.showHover('copy'); isContextMenuVisible = false">
        <i class="material-icons">content_copy</i>
        <span>{{ t('buttons.copyFile') }}</span>
      </div>
      <div v-if="headerButtons.move" class="action" @click="layoutStore.showHover('move'); isContextMenuVisible = false">
        <i class="material-icons">forward</i>
        <span>{{ t('buttons.moveFile') }}</span>
      </div>
      <div v-if="headerButtons.delete" class="action" @click="layoutStore.showHover('delete'); isContextMenuVisible = false">
        <i class="material-icons">delete</i>
        <span>{{ t('buttons.delete') }}</span>
      </div>
      <div v-if="headerButtons.download" class="action" @click="download(); isContextMenuVisible = false">
        <i class="material-icons">file_download</i>
        <span>{{ t('buttons.download') }}</span>
      </div>
      <div v-if="headerButtons.shell" class="action" @click="layoutStore.toggleShell(); isContextMenuVisible = false">
        <i class="material-icons">terminal</i>
        <span>{{ t('buttons.shell') }}</span>
      </div>
    </context-menu>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { files as api } from "@/api";
import * as upload from "@/utils/upload";
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";

import HeaderBar from "@/components/header/HeaderBar.vue";
import Action from "@/components/header/Action.vue";
import Search from "@/components/Search.vue";
import Item from "@/components/files/ListingItem.vue";
import ContextMenu from "@/components/ContextMenu.vue";
import Breadcrumbs from "@/components/Breadcrumbs.vue";

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();
const authStore = useAuthStore();

const showLimit = ref(50);
const listingContainer = ref<HTMLElement | null>(null);
const isContextMenuVisible = ref(false);
const contextMenuPos = ref({ x: 0, y: 0 });

const isEmpty = computed(() => !fileStore.req?.items || fileStore.req.items.length === 0);
const canGoBack = computed(() => route.path !== '/files' && route.path !== '/files/');

const visibleItems = computed(() => {
  if (!fileStore.req?.items) return [];
  return fileStore.req.items
    .map((item, index) => ({ ...item, originalIndex: index }))
    .slice(0, showLimit.value);
});

const nameSorted = computed(() => fileStore.req?.sorting.by === "name");
const sizeSorted = computed(() => fileStore.req?.sorting.by === "size");
const modifiedSorted = computed(() => fileStore.req?.sorting.by === "modified");
const ascOrdered = computed(() => fileStore.req?.sorting.asc);

const nameIcon = computed(() => (nameSorted.value && !ascOrdered.value) ? "arrow_upward" : "arrow_downward");
const sizeIcon = computed(() => (sizeSorted.value && ascOrdered.value) ? "arrow_downward" : "arrow_upward");
const modifiedIcon = computed(() => (modifiedSorted.value && ascOrdered.value) ? "arrow_downward" : "arrow_upward");

const viewIcon = computed(() => {
  const modes: any = { 
    list: "view_list", 
    compact: "view_headline",
    mosaic: "view_module", 
    "mosaic gallery": "grid_view" 
  };
  return modes[authStore.user?.viewMode || 'list'];
});

const viewModeClass = computed(() => {
  const mode = authStore.user?.viewMode || 'list';
  if (mode === 'mosaic gallery') {
    return 'mosaic gallery';
  }
  return mode;
});

const headerButtons = computed(() => ({
  upload: authStore.user?.perm.create,
  download: authStore.user?.perm.download,
  shell: authStore.user?.perm.execute,
  delete: fileStore.selectedCount > 0 && authStore.user?.perm.delete,
  rename: fileStore.selectedCount === 1 && authStore.user?.perm.rename,
  share: fileStore.selectedCount === 1 && authStore.user?.perm.share,
  move: fileStore.selectedCount > 0 && authStore.user?.perm.rename,
  copy: fileStore.selectedCount > 0 && authStore.user?.perm.create,
}));

const switchView = () => {
  if (!authStore.user) return;
  const modes = ["list", "compact", "mosaic", "mosaic gallery"];
  let index = modes.indexOf(authStore.user.viewMode) + 1;
  if (index >= modes.length) index = 0;
  authStore.user.viewMode = modes[index] as any;
};

const sort = (by: string) => {
  let asc = (fileStore.req?.sorting.by === by) ? !fileStore.req.sorting.asc : false;
  router.replace({ query: { ...route.query, sort: by, order: asc ? 'asc' : 'desc' } });
};

const uploadFunc = () => document.getElementById("upload-input")?.click();
const uploadFolderFunc = () => document.getElementById("upload-folder-input")?.click();
const uploadInput = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && fileStore.req) {
    const files = Array.from(target.files).map(f => ({
      file: f,
      name: f.name,
      size: f.size,
      isDir: false,
      fullPath: undefined
    }));
    upload.handleFiles(files as any, fileStore.req.url);
    target.value = ""; // Reset input
  }
};

const download = () => {
  const selectedItems = fileStore.selected.map(i => fileStore.req!.items[i]);
  if (selectedItems.length === 0) return;
  if (selectedItems.length === 1 && !selectedItems[0].isDir) {
    window.open(api.getDownloadURL(selectedItems[0], false));
  } else {
    api.download('zip', ...selectedItems.map(i => i.path));
  }
};

const handleContextMenu = (e: MouseEvent) => {
  e.preventDefault();
  isContextMenuVisible.value = true;
  contextMenuPos.value = { x: e.clientX, y: e.clientY };
};

const scrollEvent = () => {
  if (!listingContainer.value) return;
  const { scrollTop, scrollHeight, clientHeight } = listingContainer.value;
  if (scrollTop + clientHeight >= scrollHeight - 100) {
    if (fileStore.req?.items && showLimit.value < fileStore.req.items.length) {
      showLimit.value += 50;
    }
  }
};

watch(() => fileStore.req, () => {
  showLimit.value = 50;
  if (listingContainer.value) listingContainer.value.scrollTop = 0;
});
</script>

<style scoped>
.file-listing-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding-top: 4em; /* Offset for fixed header */
}
.listing-container {
  flex: 1;
  overflow-y: auto;
}
.spinner-container {
  display: flex;
  justify-content: center;
  padding: 2em;
}
.header-breadcrumbs {
  flex: 1;
  margin: 0 1em;
  padding: 0.5em 1em;
  background: transparent;
  border-radius: 0;
  margin-bottom: 0;
}
.header-breadcrumbs:deep(.breadcrumbs) {
  background: transparent;
  padding: 0;
  margin: 0;
}
</style>

<template>
  <div class="listing-container" ref="listingContainer" @scroll="scrollEvent">
      <div v-if="loading" class="spinner-container">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
      </div>

      <div v-else id="listing" class="file-icons" :class="viewModeClass" @contextmenu="handleContextMenu">
        <div v-if="items.length === 0" class="empty-state">
           <!-- Handled by parent or just show empty -->
           <slot name="empty">
             <div class="no-results">
                <i class="material-icons">folder_open</i>
                <p>{{ t("files.lonely") }}</p>
             </div>
           </slot>
        </div>

        <template v-else>
          <!-- Header for List Mode -->
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
            v-for="(item, index) in items"
            :key="base64(item.name)"
            v-bind="item"
            :index="item.originalIndex || index"
            readOnly
          />
        </template>
      </div>

      <context-menu :show="isContextMenuVisible" :pos="contextMenuPos" @hide="isContextMenuVisible = false">
        <div v-if="canEdit" class="action" @click="emit('upload'); isContextMenuVisible = false">
          <i class="material-icons">file_upload</i>
          <span>{{ t('buttons.upload') }}</span>
        </div>
        <div v-if="canEdit" class="action" @click="emit('create-folder'); isContextMenuVisible = false">
          <i class="material-icons">create_new_folder</i>
          <span>{{ t('buttons.newFolder') || 'New Folder' }}</span>
        </div>
        <div class="action" @click="emit('download'); isContextMenuVisible = false">
          <i class="material-icons">file_download</i>
          <span>{{ t('buttons.download') }}</span>
        </div>
      </context-menu>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useFileStore } from "@/stores/file";
import { useAuthStore } from "@/stores/auth";
import { useI18n } from "vue-i18n";
import Item from "@/components/files/ListingItem.vue";
import ContextMenu from "@/components/ContextMenu.vue";

const props = defineProps<{
  items: any[];
  viewMode: string;
  loading: boolean;
  canEdit?: boolean;
}>();

const emit = defineEmits<{
  (e: 'load-more'): void;
  (e: 'upload'): void;
  (e: 'create-folder'): void;
  (e: 'download'): void;
}>();

const { t } = useI18n();
const fileStore = useFileStore();
const authStore = useAuthStore();
const listingContainer = ref<HTMLElement | null>(null);

const isContextMenuVisible = ref(false);
const contextMenuPos = ref({ x: 0, y: 0 });

const handleContextMenu = (e: MouseEvent) => {
  e.preventDefault();
  isContextMenuVisible.value = true;
  contextMenuPos.value = { x: e.clientX, y: e.clientY };
};

const base64 = (s: string) => window.btoa(unescape(encodeURIComponent(s)));

// Sort state (local or store?)
// Share view usually relies on store.req.sorting
const nameSorted = computed(() => fileStore.req?.sorting.by === "name");
const sizeSorted = computed(() => fileStore.req?.sorting.by === "size");
const modifiedSorted = computed(() => fileStore.req?.sorting.by === "modified");
const ascOrdered = computed(() => fileStore.req?.sorting.asc);

const nameIcon = computed(() => (nameSorted.value && !ascOrdered.value) ? "arrow_upward" : "arrow_downward");
const sizeIcon = computed(() => (sizeSorted.value && ascOrdered.value) ? "arrow_downward" : "arrow_upward");
const modifiedIcon = computed(() => (modifiedSorted.value && ascOrdered.value) ? "arrow_downward" : "arrow_upward");

const viewModeClass = computed(() => {
  const mode = props.viewMode;
  if (mode === 'mosaic gallery') {
    return 'mosaic gallery';
  }
  return mode;
});

// Emulate FileListing sort logic
// Note: In public view, correct sorting requires backend support or local sorting.
// Assuming backend handles it if we push to router, same as Files.vue
import { useRouter, useRoute } from "vue-router";
const router = useRouter();
const route = useRoute();

const sort = (by: string) => {
  let asc = (fileStore.req?.sorting?.by === by) ? !fileStore.req?.sorting?.asc : false;
  router.replace({ query: { ...route.query, sort: by, order: asc ? 'asc' : 'desc' } });
};

const scrollEvent = () => {
  if (!listingContainer.value) return;
  const { scrollTop, scrollHeight, clientHeight } = listingContainer.value;
  if (scrollTop + clientHeight >= scrollHeight - 100) {
     emit('load-more');
  }
};

watch(() => props.viewMode, () => {
    // maybe scroll to top or re-layout
});

</script>

<style scoped>
.listing-container {
  flex: 1;
  overflow-y: auto;
}
.spinner-container {
  display: flex;
  justify-content: center;
  padding: 2em;
}
</style>

<template>
  <div class="public-shares-view">
    <header-bar showMenu showLogo @logo-click="layoutStore.toggleSidebar()">
      <div class="nav-buttons">
        <action icon="arrow_back" :label="t('buttons.back')" @action="router.back()" />
        <action icon="arrow_forward" :label="t('buttons.forward')" @action="router.forward()" />
      </div>
      <title />
      <action
        icon="refresh"
        :label="t('buttons.reload') || 'Refresh'"
        @action="refreshShares"
      />
      <action
        v-if="layoutStore.showPreviewPanel"
        icon="visibility_off"
        :label="t('buttons.hidePreview') || 'Hide Preview'"
        @action="layoutStore.togglePreviewPanel()"
      />
      <action
        v-else
        icon="preview"
        :label="t('buttons.showPreview') || 'Show Preview'"
        @action="layoutStore.togglePreviewPanel()"
      />
    </header-bar>

    <breadcrumbs base="/public-shares" />

    <div v-if="loading" class="loading">
      <h2 class="message delayed">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
        <span>{{ t("files.loading") }}</span>
      </h2>
    </div>

    <div v-else-if="shares.length === 0" class="empty">
      <h2 class="message">
        <i class="material-icons">sentiment_dissatisfied</i>
        <span>{{ t("files.lonely") }}</span>
      </h2>
    </div>

    <div v-else id="listing" :class="listingClass">
      <div class="shares-header">
        <h3>
          <i class="material-icons">public</i>
          {{ t('sidebar.publicShares') || 'Public Shares' }}
          <span class="count">({{ shares.length }})</span>
        </h3>
      </div>

      <div class="shares-grid">
        <div
          v-for="share in shares"
          :key="share.hash"
          class="share-item-card"
          :class="{ selected: selectedShare === share.hash }"
          @click="openShare(share)"
        >
          <div class="item-icon-wrapper">
            <i class="material-icons item-icon" :class="getIconClass(share.path)">
              {{ getIcon(share.path) }}
            </i>
          </div>
          <div class="item-details">
            <span class="item-name" :title="getShareName(share.path)">
              {{ getShareName(share.path) }}
            </span>
            <span class="item-meta">
              <span class="owner">
                <i class="material-icons">person</i>
                {{ share.username || 'User ' + share.userID }}
              </span>
              <span class="permission" :class="'perm-' + share.permission">
                <i class="material-icons">{{ share.permission === 'view' ? 'visibility' : 'edit' }}</i>
                {{ share.permission === 'view' ? 'View' : 'Edit' }}
              </span>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import HeaderBar from '@/components/header/HeaderBar.vue';
import Action from '@/components/header/Action.vue';
import Breadcrumbs from '@/components/Breadcrumbs.vue';
import { useLayoutStore } from '@/stores/layout';
import { useAuthStore } from '@/stores/auth';
import * as api from '@/api/index';
import { inject } from "vue";

const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const layoutStore = useLayoutStore();
const authStore = useAuthStore();

const loading = ref(true);
const shares = ref<any[]>([]);
const selectedShare = ref<string | null>(null);

const listingClass = computed(() => {
  return authStore.user?.viewMode === 'mosaic' ? 'mosaic' : 'list';
});

onMounted(async () => {
  await refreshShares();
});

async function refreshShares() {
  loading.value = true;
  try {
    shares.value = await api.share.listPublic();
  } catch (error) {
    console.error('Failed to load public shares:', error);
  } finally {
    loading.value = false;
  }
}

function selectShare(share: any) {
  selectedShare.value = share.hash;
  // Set preview panel file info
  layoutStore.setPreviewPanelFile({
    name: getShareName(share.path),
    path: share.path,
    isDir: isFolder(share.path),
    type: getFileType(share.path),
    size: 0,
    modified: '',
    // Custom property to indicate this is a share preview
    _isPublicShare: true,
    _shareHash: share.hash,
    _shareOwner: share.username || `User ${share.userID}`,
    _sharePermission: share.permission,
  });
}

function openShare(share: any) {
  router.push(`/share/${share.hash}`);
}

function getShareName(path: string): string {
  if (path === "/" || path === "") {
    return "Root";
  }
  const parts = path.split("/").filter(p => p);
  return parts[parts.length - 1] || "Root";
}

function isFolder(path: string): boolean {
  if (path === "/" || path === "") return true;
  if (path.endsWith("/")) return true;
  // Check if path has no extension (likely a folder)
  const name = getShareName(path);
  return !name.includes('.');
}

function getFileType(path: string): string {
  if (isFolder(path)) return 'folder';
  const ext = path.split('.').pop()?.toLowerCase() || '';
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'];
  const videoExts = ['mp4', 'mkv', 'avi', 'mov', 'webm'];
  const audioExts = ['mp3', 'wav', 'ogg', 'flac', 'm4a'];
  
  if (imageExts.includes(ext)) return 'image';
  if (videoExts.includes(ext)) return 'video';
  if (audioExts.includes(ext)) return 'audio';
  return 'text';
}

function getIcon(path: string): string {
  if (path === "/" || path === "") return "folder";
  if (path.endsWith("/")) return "folder";
  
  const ext = path.split('.').pop()?.toLowerCase();
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'];
  const videoExts = ['mp4', 'mkv', 'avi', 'mov', 'webm'];
  const audioExts = ['mp3', 'wav', 'ogg', 'flac', 'm4a'];
  const docExts = ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'];
  
  if (imageExts.includes(ext || '')) return "image";
  if (videoExts.includes(ext || '')) return "movie";
  if (audioExts.includes(ext || '')) return "audiotrack";
  if (docExts.includes(ext || '')) return "description";
  
  // Check if it's likely a folder (no extension)
  const name = getShareName(path);
  if (!name.includes('.')) return "folder";
  
  return "insert_drive_file";
}

function getIconClass(path: string): string {
  const icon = getIcon(path);
  if (icon === 'folder') return 'icon-folder';
  if (icon === 'image') return 'icon-image';
  if (icon === 'movie') return 'icon-video';
  if (icon === 'audiotrack') return 'icon-audio';
  return 'icon-file';
}
</script>

<style scoped>
.public-shares-view {
  width: 100%;
}

.nav-buttons {
  display: flex;
  gap: 4px;
  margin-right: 12px;
  padding-right: 12px;
  border-right: 1px solid var(--divider);
}

.nav-buttons :deep(.action) {
  padding: 6px;
  border-radius: 6px;
}

.loading,
.empty {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100vh - 12em);
}

.message {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1em;
  color: var(--textSecondary);
}

.message i {
  font-size: 4em;
  opacity: 0.5;
}

#listing {
  padding: 1em;
}

.shares-header {
  margin-bottom: 1.5em;
  padding-bottom: 0.5em;
  border-bottom: 1px solid var(--divider);
}

.shares-header h3 {
  display: flex;
  align-items: center;
  gap: 0.5em;
  color: var(--textPrimary);
  font-weight: 500;
  margin: 0;
}

.shares-header h3 i {
  color: var(--blue);
}

.shares-header .count {
  color: var(--textSecondary);
  font-weight: normal;
  font-size: 0.9em;
}

.shares-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1em;
}

.share-item-card {
  display: flex;
  align-items: center;
  gap: 1em;
  padding: 1em;
  background: var(--surfacePrimary);
  border: 1px solid var(--divider);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.share-item-card:hover {
  border-color: var(--blue);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.share-item-card.selected {
  border-color: var(--blue);
  background: rgba(33, 150, 243, 0.08);
}

.item-icon-wrapper {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--surfaceSecondary);
  border-radius: 8px;
  flex-shrink: 0;
}

.item-icon {
  font-size: 28px;
}

.icon-folder { color: #f9a825; }
.icon-image { color: #43a047; }
.icon-video { color: #e53935; }
.icon-audio { color: #8e24aa; }
.icon-file { color: var(--textSecondary); }

.item-details {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.35em;
}

.item-name {
  font-weight: 500;
  color: var(--textPrimary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 1em;
  font-size: 0.85em;
  color: var(--textSecondary);
}

.item-meta .owner,
.item-meta .permission {
  display: flex;
  align-items: center;
  gap: 0.25em;
}

.item-meta i {
  font-size: 1em;
}

.item-meta .perm-view {
  color: var(--blue);
}

.item-meta .perm-change {
  color: #43a047;
}

/* List view styles */
#listing.list .shares-grid {
  display: flex;
  flex-direction: column;
  gap: 0.5em;
}

#listing.list .share-item-card {
  padding: 0.75em 1em;
}

#listing.list .item-icon-wrapper {
  width: 40px;
  height: 40px;
}

#listing.list .item-icon {
  font-size: 24px;
}

/* Spinner */
.spinner {
  display: flex;
  gap: 0.3em;
}

.spinner > div {
  width: 10px;
  height: 10px;
  background-color: var(--blue);
  border-radius: 100%;
  animation: bounce 1.4s infinite ease-in-out both;
}

.spinner .bounce1 {
  animation-delay: -0.32s;
}

.spinner .bounce2 {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

/* Responsive */
@media (max-width: 600px) {
  .shares-grid {
    grid-template-columns: 1fr;
  }
}
</style>

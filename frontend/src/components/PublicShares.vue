<template>
  <div class="public-shares-container">
    <div class="public-shares-header" @click="toggleExpanded">
      <div class="header-left">
        <i class="material-icons expand-icon">{{ isExpanded ? 'expand_less' : 'expand_more' }}</i>
        <i class="material-icons">public</i>
        <span>{{ $t("sidebar.publicShares") || "Public Shares" }}</span>
        <span v-if="publicShares.length > 0" class="share-count">({{ filteredShares.length }}/{{ publicShares.length }})</span>
      </div>
      <div class="header-actions">
        <button 
          class="action view-mode-btn" 
          @click.stop="toggleViewMode"
          :aria-label="viewMode === 'list' ? 'View as My Files' : 'View as List'"
          :title="viewMode === 'list' ? ($t('sidebar.viewAsMyFiles') || 'View as My Files') : ($t('sidebar.viewAsList') || 'View as List')"
        >
          <i class="material-icons">{{ viewMode === 'list' ? 'folder_open' : 'list' }}</i>
        </button>
        <button 
          class="action refresh-btn" 
          @click.stop="refreshPublicShares"
          :aria-label="'Refresh'"
          :title="'Refresh'"
        >
          <i class="material-icons" :class="{ spinning: loading }">refresh</i>
        </button>
      </div>
    </div>

    <div v-if="isExpanded" class="public-shares-content">
      <!-- Filter input -->
      <div class="filter-container">
        <i class="material-icons filter-icon">search</i>
        <input
          type="text"
          v-model="filterText"
          :placeholder="$t('search.search') || 'Search...'"
          class="filter-input"
          @click.stop
        />
        <button 
          v-if="filterText"
          class="clear-filter"
          @click.stop="filterText = ''"
          :title="$t('buttons.clear') || 'Clear'"
        >
          <i class="material-icons">close</i>
        </button>
      </div>

      <div v-if="loading" class="shares-loading">
        <i class="material-icons spinning">sync</i>
      </div>

      <div v-else-if="publicShares.length === 0" class="shares-empty">
        <span>{{ $t("files.lonely") || "No shares yet" }}</span>
      </div>

      <div v-else-if="filteredShares.length === 0" class="shares-empty">
        <span>No matches found</span>
      </div>

      <div v-else class="shares-list">
        <router-link
          v-for="share in filteredShares"
          :key="share.hash"
          :to="`/share/${share.hash}`"
          class="share-item"
          :title="getTooltip(share)"
          @click="closeHovers"
        >
          <i class="material-icons item-icon" :class="getIconClass(share.path)">{{ getIcon(share.path) }}</i>
          <div class="item-info">
            <span class="item-name">{{ getShareName(share.path) }}</span>
            <span class="item-owner">by {{ share.username || 'User ' + share.userID }}</span>
          </div>
          <i class="material-icons item-perm" :class="'perm-' + share.permission" :title="share.permission === 'view' ? 'View Only' : 'Can Edit'">
            {{ share.permission === 'view' ? 'visibility' : 'edit' }}
          </i>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from "pinia";
import * as api from "@/api/index";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { useRouter } from "vue-router";

export default {
  name: "publicShares",
  data: function () {
    return {
      publicShares: [],
      loading: true,
      isExpanded: false,
      filterText: "",
      viewMode: "list",
    };
  },
  inject: ["$showError", "$showSuccess"],
  computed: {
    ...mapState(useFileStore, ["req", "selected", "selectedCount", "isListing"]),
    filteredShares() {
      if (!this.filterText.trim()) {
        return this.publicShares;
      }
      const search = this.filterText.toLowerCase();
      return this.publicShares.filter(share => {
        const name = this.getShareName(share.path).toLowerCase();
        const path = share.path.toLowerCase();
        const owner = (share.username + `${share.userName}` || `user ${share.userID}`).toLowerCase();
        return name.includes(search) || path.includes(search) || owner.includes(search);
      });
    },
  },
  async beforeMount() {
    await this.refreshPublicShares();
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers"]),
    toggleExpanded() {
      this.isExpanded = !this.isExpanded;
    },
    async refreshPublicShares() {
      this.loading = true;
      try {
        this.publicShares = await api.share.listPublic();
      } catch (e) {
        this.$showError(e);
      } finally {
        this.loading = false;
      }
    },
    getShareName(path) {
      if (path === "/" || path === "") {
        return "Root";
      }
      const parts = path.split("/").filter(p => p);
      return parts[parts.length - 1] || "Root";
    },
    getIcon(path) {
      if (path === "/" || path === "") {
        return "folder";
      }
      if (path.endsWith("/")) return "folder";
      
      // Check extension for file type
      const ext = path.split('.').pop()?.toLowerCase();
      const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'];
      const videoExts = ['mp4', 'mkv', 'avi', 'mov', 'webm'];
      const audioExts = ['mp3', 'wav', 'ogg', 'flac', 'm4a'];
      const docExts = ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'];
      
      if (imageExts.includes(ext)) return "image";
      if (videoExts.includes(ext)) return "movie";
      if (audioExts.includes(ext)) return "audiotrack";
      if (docExts.includes(ext)) return "description";
      
      return "insert_drive_file";
    },
    getIconClass(path) {
      const icon = this.getIcon(path);
      if (icon === 'folder') return 'icon-folder';
      if (icon === 'image') return 'icon-image';
      if (icon === 'movie') return 'icon-video';
      if (icon === 'audiotrack') return 'icon-audio';
      return 'icon-file';
    },
    getTooltip(share) {
      const name = this.getShareName(share.path);
      const perm = share.permission === 'view' ? 'View Only' : 'Can Edit';
      const owner = share.username || `User ${share.userID}`;
      return `${name}\n📁 ${share.path}\n👤 Shared by: ${owner}\n🔐 Permission: ${perm}\n🔗 Hash: ${share.hash}`;
    },
    toggleViewMode() {
      // Navigate to public shares view page
      this.$router.push('/public-shares');
    },
  },
};
</script>

<style scoped>
.public-shares-container {
  border-top: 1px solid var(--divider);
  background: var(--background);
}

.public-shares-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.public-shares-header:hover {
  background: var(--surfacePrimary);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.85rem;
  color: var(--textPrimary);
}

.header-left .material-icons {
  font-size: 1.1rem;
  color: var(--textSecondary);
}

.expand-icon {
  font-size: 1.25rem !important;
}

.share-count {
  color: var(--textSecondary);
  font-size: 0.75rem;
  margin-left: 0.25rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.view-mode-btn,
.refresh-btn {
  width: 24px;
  height: 24px;
  padding: 0;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--textSecondary);
  transition: all 0.2s;
}

.view-mode-btn:hover,
.refresh-btn:hover {
  color: var(--textPrimary);
  background: var(--surfaceSecondary);
}

.view-mode-btn .material-icons,
.refresh-btn .material-icons {
  font-size: 1rem;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Filter */
.filter-container {
  display: flex;
  align-items: center;
  padding: 0.25rem 0.5rem;
  margin: 0.25rem 0.5rem;
  background: var(--surfacePrimary);
  border-radius: 4px;
  border: 1px solid var(--divider);
}

.filter-icon {
  font-size: 1rem;
  color: var(--textSecondary);
  margin-right: 0.25rem;
}

.filter-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 0.8rem;
  color: var(--textPrimary);
  outline: none;
  padding: 0.25rem;
}

.filter-input::placeholder {
  color: var(--textSecondary);
}

.clear-filter {
  width: 20px;
  height: 20px;
  padding: 0;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--textSecondary);
  border-radius: 50%;
}

.clear-filter:hover {
  color: var(--textPrimary);
  background: var(--surfaceSecondary);
}

.clear-filter .material-icons {
  font-size: 0.9rem;
}

/* Content */
.public-shares-content {
  max-height: 300px;
  overflow-y: auto;
}

.shares-loading,
.shares-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  color: var(--textSecondary);
  font-size: 0.8rem;
}

.shares-list {
  display: flex;
  flex-direction: column;
  padding-bottom: 0.25rem;
}

.share-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.35rem 0.5rem 0.35rem 0.75rem;
  text-decoration: none;
  color: var(--textPrimary);
  transition: background 0.15s;
  border-left: 2px solid transparent;
}

.share-item:hover {
  background: var(--surfacePrimary);
  border-left-color: var(--blue);
}

.item-icon {
  font-size: 1.1rem;
  flex-shrink: 0;
}

.icon-folder { color: #f9a825; }
.icon-image { color: #43a047; }
.icon-video { color: #e53935; }
.icon-audio { color: #8e24aa; }
.icon-file { color: var(--textSecondary); }

.item-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.item-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.85rem;
  line-height: 1.2;
}

.item-owner {
  font-size: 0.7rem;
  color: var(--textSecondary);
  line-height: 1.2;
}

.item-perm {
  font-size: 0.9rem;
  flex-shrink: 0;
  padding: 2px;
  border-radius: 2px;
}

.perm-view {
  color: var(--blue);
}

.perm-change {
  color: #43a047;
}

/* Scrollbar */
.public-shares-content::-webkit-scrollbar {
  width: 4px;
}

.public-shares-content::-webkit-scrollbar-track {
  background: transparent;
}

.public-shares-content::-webkit-scrollbar-thumb {
  background: var(--divider);
  border-radius: 2px;
}

.public-shares-content::-webkit-scrollbar-thumb:hover {
  background: var(--textSecondary);
}
</style>

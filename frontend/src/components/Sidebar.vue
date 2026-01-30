<template>
  <div v-show="active" @click="closeHovers" class="sidebar-overlay"></div>
  <aside :class="{ active, collapsed: !layoutStore.showSidebar }">
    <div class="sidebar-header">
      <div class="user-profile" @click="toAccountSettings">
        <div class="avatar">
          <i class="material-icons">account_circle</i>
        </div>
        <div class="user-info">
          <span class="username">{{ user.username }}</span>
          <span class="role">{{ user.perm.admin ? 'Administrator' : 'User' }}</span>
        </div>
      </div>
    </div>

    <div class="sidebar-scrollable">
      <div class="menu-section">
        <button class="nav-item" :class="{ 'is-active': isHome }" @click="toRoot">
          <i class="material-icons">home</i>
          <span>{{ $t("sidebar.myFiles") }}</span>
        </button>
        
        <button v-if="user.perm.create" class="nav-item action-btn" @click="showHover('newDir')">
          <i class="material-icons">add_box</i>
          <span>{{ $t("sidebar.newFolder") }}</span>
        </button>
      </div>

      <div class="menu-divider"></div>

      <div class="tree-section">
        <directory-tree />
      </div>

      <div class="menu-divider" v-if="user.perm.admin || user.perm.share"></div>

      <div class="menu-section">
        <button v-if="user.perm.admin" class="nav-item" @click="toGlobalSettings">
          <i class="material-icons">settings</i>
          <span>{{ $t("sidebar.settings") }}</span>
        </button>

        <div v-if="user.perm.share" class="shares-wrapper">
          <public-shares />
        </div>
      </div>
    </div>

    <div class="sidebar-footer">
      <div v-if="isFiles && !disableUsedPercentage" class="storage-info">
        <div class="storage-header">
          <span>Storage</span>
          <span>{{ usage.usedPercentage }}%</span>
        </div>
        <div class="progress-track">
          <div class="progress-fill" :style="{ width: usage.usedPercentage + '%' }"></div>
        </div>
        <div class="storage-details">
          {{ usage.used }} / {{ usage.total }}
        </div>
      </div>
      
      <button v-if="canLogout" @click="logout" class="logout-btn">
        <i class="material-icons">logout</i>
        <span>{{ $t("sidebar.logout") }}</span>
      </button>
    </div>
  </aside>
</template>

<script>
import { reactive } from "vue";
import { mapActions, mapState } from "pinia";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import * as auth from "@/utils/auth";
import {
  signup,
  hideLoginButton,
  disableExternal,
  disableUsedPercentage,
  noAuth,
  logoutPage,
  loginPage,
} from "@/utils/constants";
import { files as api } from "@/api";
import PublicShares from "@/components/PublicShares.vue";
import DirectoryTree from "@/components/DirectoryTree.vue";
import prettyBytes from "pretty-bytes";

const USAGE_DEFAULT = { used: "0 B", total: "0 B", usedPercentage: 0 };

export default {
  name: "sidebar",
  setup() {
    const usage = reactive(USAGE_DEFAULT);
    const layoutStore = useLayoutStore();
    return { usage, usageAbortController: new AbortController(), layoutStore };
  },
  components: {
    PublicShares,
    DirectoryTree,
  },
  inject: ["$showError"],
  computed: {
    ...mapState(useAuthStore, ["user", "isLoggedIn"]),
    ...mapState(useFileStore, ["isFiles", "reload"]),
    ...mapState(useLayoutStore, ["currentPromptName"]),
    active() {
      return this.currentPromptName === "sidebar";
    },
    isHome() {
      return this.$route.path === '/files' || this.$route.path === '/files/';
    },
    canLogout: () => !noAuth && (loginPage || logoutPage !== "/login"),
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers", "showHover"]),
    abortOngoingFetchUsage() {
      this.usageAbortController.abort();
    },
    async fetchUsage() {
      const path = this.$route.path.endsWith("/")
        ? this.$route.path
        : this.$route.path + "/";
      let usageStats = USAGE_DEFAULT;
      if (this.disableUsedPercentage) {
        return Object.assign(this.usage, usageStats);
      }
      try {
        this.abortOngoingFetchUsage();
        this.usageAbortController = new AbortController();
        const usage = await api.usage(path, this.usageAbortController.signal);
        usageStats = {
          used: prettyBytes(usage.used, { binary: true }),
          total: prettyBytes(usage.total, { binary: true }),
          usedPercentage: Math.round((usage.used / usage.total) * 100),
        };
      } finally {
        return Object.assign(this.usage, usageStats);
      }
    },
    toRoot() {
      this.$router.push({ path: "/files" });
      this.closeHovers();
    },
    toAccountSettings() {
      this.$router.push({ path: "/settings/profile" });
      this.closeHovers();
    },
    toGlobalSettings() {
      this.$router.push({ path: "/settings/global" });
      this.closeHovers();
    },
    logout: auth.logout,
  },
  watch: {
    $route: {
      handler(to) {
        if (to.path.includes("/files")) {
          this.fetchUsage();
        }
      },
      immediate: true,
    },
  },
  unmounted() {
    this.abortOngoingFetchUsage();
  },
};
</script>

<style scoped>
aside {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 18em;
  background: var(--surfacePrimary);
  border-right: 1px solid var(--divider);
  display: flex;
  flex-direction: column;
  z-index: 1001;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

aside.collapsed {
  transform: translateX(-100%);
}

.sidebar-header {
  padding: 1.5rem 1.25rem;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.user-profile:hover {
  background: var(--surfaceSecondary);
}

.avatar i {
  font-size: 40px;
  color: var(--blue);
}

.user-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.username {
  font-weight: 600;
  color: var(--textPrimary);
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.role {
  font-size: 0.75rem;
  color: var(--textTertiary);
}

.sidebar-scrollable {
  flex: 1;
  overflow-y: auto;
  padding: 0 0.75rem;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 12px;
  border-radius: 10px;
  color: var(--textSecondary);
  background: transparent;
  border: none;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s;
  margin-bottom: 4px;
}

.nav-item:hover {
  background: var(--surfaceSecondary);
  color: var(--textPrimary);
}

.nav-item.is-active {
  background: var(--selection);
  color: var(--blue);
  font-weight: 600;
}

.nav-item i {
  font-size: 22px;
}

.action-btn {
  color: var(--blue);
  background: rgba(59, 130, 246, 0.05);
  margin-top: 1rem;
}

.menu-divider {
  height: 1px;
  background: var(--divider);
  margin: 1rem 0.75rem;
}

.sidebar-footer {
  padding: 1.25rem;
  border-top: 1px solid var(--divider);
}

.storage-info {
  margin-bottom: 1.5rem;
}

.storage-header {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  margin-bottom: 8px;
  color: var(--textSecondary);
  font-weight: 500;
}

.progress-track {
  height: 6px;
  background: var(--surfaceSecondary);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--blue);
  border-radius: 3px;
}

.storage-details {
  font-size: 0.75rem;
  color: var(--textTertiary);
  margin-top: 6px;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 12px;
  border-radius: 10px;
  color: var(--red);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: background 0.2s;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.08);
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  backdrop-filter: blur(4px);
  z-index: 1000;
}

@media (max-width: 1024px) {
  aside {
    transform: translateX(-100%);
  }
  aside.active {
    transform: translateX(0);
  }
}
</style>

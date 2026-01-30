<template>
  <div class="dashboard">
    <header-bar showMenu showLogo @logo-click="layoutStore.toggleSidebar()">
      <div class="nav-buttons">
        <action icon="arrow_back" :label="t('buttons.back')" @action="router.back()" />
        <action icon="arrow_forward" :label="t('buttons.forward')" @action="router.forward()" />
      </div>
    </header-bar>

    <div id="nav">
      <div class="wrapper">
        <ul>
          <router-link to="/settings/profile"
            ><li :class="{ active: $route.path === '/settings/profile' }">
              {{ t("settings.profileSettings") }}
            </li></router-link
          >
          <router-link to="/settings/shares" v-if="user?.perm.share"
            ><li :class="{ active: $route.path === '/settings/shares' }">
              {{ t("settings.shareManagement") }}
            </li></router-link
          >
          <router-link to="/settings/global" v-if="user?.perm.admin"
            ><li :class="{ active: $route.path === '/settings/global' }">
              {{ t("settings.globalSettings") }}
            </li></router-link
          >
          <router-link to="/settings/users" v-if="user?.perm.admin"
            ><li
              :class="{
                active:
                  $route.path === '/settings/users' || $route.name === 'User',
              }"
            >
              {{ t("settings.userManagement") }}
            </li></router-link
          >
        </ul>
      </div>
    </div>

    <div v-if="loading">
      <h2 class="message delayed">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
        <span>{{ t("files.loading") }}</span>
      </h2>
    </div>

    <router-view></router-view>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import HeaderBar from "@/components/header/HeaderBar.vue";
import Action from "@/components/header/Action.vue";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";

const authStore = useAuthStore();
const layoutStore = useLayoutStore();
const router = useRouter();

const { t } = useI18n();
const user = computed(() => authStore.user);
const loading = computed(() => layoutStore.loading);
</script>

<style scoped>
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
</style>

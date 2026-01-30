<template>
  <div class="card floating" id="public-share">
    <div class="card-title">
      <h2>Make Share Public</h2>
    </div>

    <div class="card-content">
      <p>Choose who can access and modify this share:</p>
      
      <div class="permission-options">
        <label class="permission-option">
          <input 
            type="radio" 
            v-model="permission" 
            value="view"
          />
          <span class="label-text">
            <i class="material-icons">visibility</i>
            <strong>View Only</strong>
            <small>Everyone can view and download, but cannot modify</small>
          </span>
        </label>

        <label class="permission-option">
          <input 
            type="radio" 
            v-model="permission" 
            value="change"
          />
          <span class="label-text">
            <i class="material-icons">edit</i>
            <strong>Change - Everyone</strong>
            <small>Everyone can view, download, and modify</small>
          </span>
        </label>

        <label v-if="users.length > 0" class="permission-option">
          <input 
            type="radio" 
            v-model="permission" 
            value="change-select"
          />
          <span class="label-text">
            <i class="material-icons">people</i>
            <strong>Change - Select Users</strong>
            <small>Only selected users can modify</small>
          </span>
        </label>
      </div>

      <div v-if="permission === 'change-select'" class="users-list">
        <p>Select users who can modify:</p>
        <div class="user-checkboxes">
          <label v-for="user in users" :key="user.id" class="user-checkbox">
            <input 
              type="checkbox" 
              :value="user.id"
              v-model="selectedUsers"
            />
            {{ user.username }}
          </label>
        </div>
      </div>
    </div>

    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="cancel"
        :aria-label="'Cancel'"
        :title="'Cancel'"
        tabindex="2"
      >
        Cancel
      </button>
      <button
        id="focus-prompt"
        class="button button--flat button--blue"
        @click="confirm"
        :aria-label="'Make Public'"
        :title="'Make Public'"
        tabindex="1"
      >
        Make Public
      </button>
    </div>
  </div>
</template>

<script>
import { mapActions } from "pinia";
import { useLayoutStore } from "@/stores/layout";
import * as api from "@/api/index";

export default {
  name: "publicShare",
  data: function () {
    return {
      permission: "view",
      selectedUsers: [],
      users: [],
    };
  },
  inject: ["$showError", "$showSuccess"],
  async beforeMount() {
    try {
      this.users = await api.users.getAll();
    } catch (e) {
      // Silently fail - user list is only available to admins
      // "Change - Select Users" option will be hidden when users list is empty
      console.log('Cannot fetch users list (not admin)');
    }
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers"]),
    cancel: function () {
      this.closeHovers();
    },
    confirm: function () {
      let allowedUsers = [];
      if (this.permission === "change-select") {
        allowedUsers = this.selectedUsers;
      }
      
      // Map "change-select" back to "change" for the API
      const finalPermission = this.permission === "change-select" ? "change" : this.permission;
      
      // Get the current prompt and call its confirm callback
      const layoutStore = useLayoutStore();
      const currentPrompt = layoutStore.currentPrompt;
      if (currentPrompt && currentPrompt.confirm && typeof currentPrompt.confirm === 'function') {
        currentPrompt.confirm(finalPermission, allowedUsers);
      }
      this.closeHovers();
    },
  },
};
</script>

<style scoped>
.permission-options {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin: 1rem 0;
}

.permission-option {
  display: flex;
  align-items: flex-start;
  padding: 0.75rem;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.permission-option:hover {
  background-color: #f5f5f5;
}

.permission-option input[type="radio"] {
  margin-right: 0.75rem;
  margin-top: 0.25rem;
  cursor: pointer;
}

.label-text {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

.label-text i {
  font-size: 1.25rem;
  color: #666;
}

.label-text strong {
  font-weight: 500;
  color: #333;
}

.label-text small {
  color: #999;
  font-size: 0.85rem;
}

.users-list {
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.user-checkboxes {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 200px;
  overflow-y: auto;
}

.user-checkbox {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  cursor: pointer;
}

.user-checkbox input[type="checkbox"] {
  margin-right: 0.5rem;
  cursor: pointer;
}
</style>

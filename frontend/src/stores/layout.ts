import { defineStore } from "pinia";
// import { useAuthPreferencesStore } from "./auth-preferences";
// import { useAuthEmailStore } from "./auth-email";

export const useLayoutStore = defineStore("layout", {
  // convert to a function
  state: (): {
    loading: boolean;
    prompts: PopupProps[];
    showShell: boolean | null;
    showSidebar: boolean;
    showPreviewPanel: boolean;
    previewPanelFile: any | null;
    markdownViewMode: 'source' | 'rendered';
    expandedNodes: Record<string, boolean>;
  } => ({
    loading: false,
    prompts: [],
    showShell: false,
    showSidebar: localStorage.getItem('showSidebar') !== 'false',
    showPreviewPanel: localStorage.getItem('showPreviewPanel') === 'true',
    previewPanelFile: null,
    markdownViewMode: (localStorage.getItem('markdownViewMode') as 'source' | 'rendered') || 'rendered',
    expandedNodes: JSON.parse(localStorage.getItem('expandedNodes') || '{}'),
  }),
  getters: {
    currentPrompt(state) {
      return state.prompts.length > 0
        ? state.prompts[state.prompts.length - 1]
        : null;
    },
    currentPromptName(): string | null | undefined {
      return this.currentPrompt?.prompt;
    },
    // user and jwt getter removed, no longer needed
  },
  actions: {
    // no context as first argument, use `this` instead
    toggleShell() {
      this.showShell = !this.showShell;
    },
    toggleSidebar() {
      if (window.innerWidth <= 1024) {
        this.showHover('sidebar');
      } else {
        this.showSidebar = !this.showSidebar;
        localStorage.setItem('showSidebar', String(this.showSidebar));
      }
    },
    toggleNode(path: string) {
      this.expandedNodes[path] = !this.expandedNodes[path];
      localStorage.setItem('expandedNodes', JSON.stringify(this.expandedNodes));
    },
    setCloseOnPrompt(closeFunction: () => Promise<string>, onPrompt: string) {
      const prompt = this.prompts.find((prompt) => prompt.prompt === onPrompt);
      if (prompt) {
        prompt.close = closeFunction;
      }
    },
    showHover(value: PopupProps | string) {
      if (typeof value !== "object") {
        this.prompts.push({
          prompt: value,
          confirm: null,
          action: undefined,
          saveAction: undefined,
          props: null,
          close: null,
        });
        return;
      }

      this.prompts.push({
        prompt: value.prompt,
        confirm: value?.confirm,
        action: value?.action,
        saveAction: value?.saveAction,
        props: value?.props,
        close: value?.close,
      });
    },
    showError() {
      this.prompts.push({
        prompt: "error",
        confirm: null,
        action: undefined,
        props: null,
        close: null,
      });
    },
    showSuccess() {
      this.prompts.push({
        prompt: "success",
        confirm: null,
        action: undefined,
        props: null,
        close: null,
      });
    },
    closeHovers() {
      this.prompts.shift()?.close?.();
    },
    // Preview panel actions
    togglePreviewPanel() {
      this.showPreviewPanel = !this.showPreviewPanel;
      localStorage.setItem('showPreviewPanel', String(this.showPreviewPanel));
      if (!this.showPreviewPanel) {
        this.previewPanelFile = null;
      }
    },
    setPreviewPanelFile(file: any) {
      this.previewPanelFile = file;
    },
    clearPreviewPanelFile() {
      this.previewPanelFile = null;
    },
    setMarkdownViewMode(mode: 'source' | 'rendered') {
      this.markdownViewMode = mode;
      localStorage.setItem('markdownViewMode', mode);
    },
    toggleMarkdownViewMode() {
      const newMode = this.markdownViewMode === 'source' ? 'rendered' : 'source';
      this.setMarkdownViewMode(newMode);
    },
    // easily reset state using `$reset`
    clearLayout() {
      this.$reset();
    },
  },
});

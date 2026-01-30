<template>
  <div class="share-page">
    <ShareEditor
      v-if="isEditing && editorData"
      :hash="hash"
      :filePath="editorData.filePath"
      :fileName="editorData.fileName"
      :fileContent="editorData.fileContent"
      :canEdit="canEdit"
      @close="closeEditor"
      @saved="onEditorSaved"
    />

    <MediaViewer
      v-if="showMediaViewer"
      :items="mediaViewerItems"
      :startIndex="mediaViewerStartIndex"
      :getMediaUrl="getMediaUrl"
      @close="closeMediaViewer"
    />

    <template v-if="!isEditing">
      <header-bar showMenu showLogo @logo-click="toggleDownloadPanel">
        <title />
        <action
          icon="file_download"
          :label="t('buttons.download')"
          @action="download"
        />
        <action
          icon="content_copy"
          :label="t('buttons.copyDownloadLink')"
          @action="() => copyToClipboard(link)"
        />
        <action
          v-if="canEdit && req?.isDir"
          icon="file_upload"
          :label="t('buttons.upload')"
          @action="triggerUpload"
        />
        <action
          v-if="canEdit && req?.isDir"
          icon="create_new_folder"
          :label="t('buttons.newFolder') || 'New Folder'"
          @action="createFolder"
        />
        <action
          v-if="canEdit && !req?.isDir"
          icon="mode_edit"
          :label="t('buttons.edit')"
          @action="edit"
        />
        <action
          :icon="viewIcon"
          :label="t('buttons.switchView') || 'Switch view'"
          @action="switchView"
        />
        <action
          icon="info"
          :class="{ active: showDownloadPanel }"
          @action="toggleDownloadPanel"
        />
        <action
          icon="preview"
          :class="{ active: layoutStore.showPreviewPanel }"
          :label="layoutStore.showPreviewPanel ? t('buttons.hidePreview') : t('buttons.showPreview')"
          @action="togglePreviewPanel"
        />
      </header-bar>

      <div class="share-main">
        <!-- Sidebar Info -->
        <aside class="share-sidebar" :class="{ 'collapsed': !showDownloadPanel }">
          <div class="share-sidebar-content" v-if="req">
            <div class="share-info-header">
              <div class="share-info-icon">
                <i class="material-icons">{{ icon }}</i>
              </div>
              <h2 class="share-info-title">{{ req.name }}</h2>
              <div class="share-info-subtitle">{{ humanSize }} • {{ humanTime }}</div>
            </div>

            <!-- Directory Tree in Share -->
            <div v-if="req.isDir" class="share-tree-section">
              <div class="share-tree-header">
                <i class="material-icons">folder_open</i>
                <span>{{ t('sidebar.directoryTree') || 'Folder Tree' }}</span>
              </div>
              <div class="share-tree-container">
                <directory-tree />
              </div>
            </div>

            <div class="share-qrcode-section">
              <div class="qrcode-wrapper">
                <qrcode-vue :value="link" :size="140" level="M" />
              </div>
              <p class="qrcode-hint">Scan to access this share</p>
            </div>

            <div class="share-sidebar-actions">
              <a :href="link" target="_blank" class="button button--primary">
                <i class="material-icons">file_download</i>
                {{ t("buttons.download") }}
              </a>
            </div>
          </div>
        </aside>

        <!-- Main Content -->
        <div class="share-content" :class="{ 'with-preview': layoutStore.showPreviewPanel }" :style="layoutStore.showPreviewPanel ? { marginRight: previewPanelWidth + 'px' } : {}">
          <div class="share-content-header">
            <breadcrumbs :base="'/share/' + hash" />
            
            <div class="share-search-wrapper">
              <div class="search-box">
                <i class="material-icons">search</i>
                <input 
                  type="text" 
                  v-model="searchQuery" 
                  @keyup.enter="performSearch"
                  :placeholder="t('buttons.search') || 'Search files...'"
                  class="search-input"
                />
                <i v-if="searchQuery" class="material-icons clear-btn" @click="searchQuery = ''">close</i>
              </div>
            </div>
          </div>
          
          <div class="share-content-body">
            <!-- Loading -->
            <div v-if="layoutStore.loading" class="loading-container">
              <div class="spinner">
                <div class="bounce1"></div>
                <div class="bounce2"></div>
                <div class="bounce3"></div>
              </div>
            </div>

            <!-- Error handling -->
            <div v-else-if="error" class="share-error-overlay">
              <div v-if="error.status === 401" class="card floating" style="width: 320px">
                <div v-if="attemptedPasswordLogin" class="share__wrong__password">
                  {{ t("login.wrongCredentials") }}
                </div>
                <div class="card-title">
                  <h2>{{ t("login.password") }}</h2>
                </div>
                <div class="card-content">
                  <input v-focus class="input input--block" type="password" :placeholder="t('login.password')" v-model="password" @keyup.enter="fetchData" />
                </div>
                <div class="card-action">
                  <button class="button button--flat" @click="fetchData">{{ t("buttons.submit") }}</button>
                </div>
              </div>
              <errors v-else :errorCode="error.status" />
            </div>

            <!-- Content -->
            <div v-else-if="req" class="share-listing-container">

               <ShareListing 
                  v-if="req.isDir" 
                  :items="filteredItems.slice(0, showLimit)" 
                  :view-mode="viewMode" 
                  :loading="false"
                  :can-edit="canEdit"
                  @load-more="showLimit += 100"
                  @upload="triggerUpload"
                  @create-folder="createFolder"
                  @download="download"
               >
                  <template #empty>
                     <div class="no-results">
                        <i class="material-icons">manage_search</i>
                        <h3>{{ searchQuery ? t('files.noMatches') || 'No files match your search' : t("files.lonely") }}</h3>
                        <p v-if="searchQuery">Try different keywords</p>
                     </div>
                  </template>
               </ShareListing>

                <!-- File Placeholder (for non-previewable files) -->
                <div v-else class="no-results file-placeholder">
                    <i class="material-icons">{{ icon }}</i>
                    <h3>{{ req.name }}</h3>
                    <div class="actions">
                        <button class="button button--primary" @click="download">
                            <i class="material-icons">file_download</i>
                            {{ t('buttons.download') }}
                        </button>
                    </div>
                </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Right Preview Panel -->
      <SharePreviewPanel :hash="hash" :token="token" @resize="onPanelResize" />
      <input type="file" id="share-upload-input" multiple style="display: none" @change="uploadInput" />
    </template>
  </div>
</template>

<script setup lang="ts">
import Action from "@/components/header/Action.vue";
import Breadcrumbs from "@/components/Breadcrumbs.vue";
import Errors from "@/views/Errors.vue";
import QrcodeVue from "qrcode.vue";
import Item from "@/components/files/ListingItem.vue";
import ShareListing from "@/views/files/ShareListing.vue";
import ShareEditor from "@/views/files/ShareEditor.vue";
import MediaViewer from "@/components/MediaViewer.vue";
import SharePreviewPanel from "@/components/SharePreviewPanel.vue";
import DirectoryTree from "@/components/DirectoryTree.vue";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { StatusError } from "@/api/utils";
import { copy as copyText } from "@/utils/clipboard";
import { pub as api } from "@/api";
import { baseURL } from "@/utils/constants";
import { filesize } from "@/utils"; 
import dayjs from "dayjs";

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const error = ref<StatusError | null>(null);
const showLimit = ref<number>(100);
const searchQuery = ref<string>("");
const password = ref<string>("");
const attemptedPasswordLogin = ref<boolean>(false);
const hash = ref<string>("");
const token = ref<string>("");
const req = ref<Resource | null>(null);
const showDownloadPanel = ref<boolean>(localStorage.getItem('showShareSidebar') !== 'false');
const viewMode = ref<string>(localStorage.getItem('shareViewMode') || 'list');
const previewPanelWidth = ref(parseInt(localStorage.getItem("sharePreviewPanelWidth") || "400"));

const $showError = inject<IToastError>("$showError")!;
const $showSuccess = inject<IToastSuccess>("$showSuccess")!;

const link = computed(() => {
  return `${window.location.origin}${baseURL}${route.path}`;
});

const icon = computed(() => {
  if (req.value?.isDir) return "folder";
  if (req.value?.type === "image") return "image";
  if (req.value?.type === "video") return "movie";
  if (req.value?.type === "audio") return "audiotrack";
  return "insert_drive_file";
});

const humanSize = computed(() => {
  if (!req.value) return "0 B";
  return filesize(req.value.size);
});

const humanTime = computed(() => {
  if (!req.value) return "";
  return dayjs(req.value.modified).fromNow();
});

// View mode switching - supports: list, compact, mosaic, mosaic gallery
const viewModes = ['list', 'compact', 'mosaic', 'mosaic gallery'];
// viewModeClass computed removed, handled in ShareListing

const viewIcon = computed(() => {
  const icons: Record<string, string> = {
    'list': 'view_list',
    'compact': 'view_headline',
    'mosaic': 'view_module',
    'mosaic gallery': 'grid_view'
  };
  return icons[viewMode.value] || 'view_list';
});

const switchView = () => {
  let index = viewModes.indexOf(viewMode.value) + 1;
  if (index >= viewModes.length) index = 0;
  viewMode.value = viewModes[index];
  localStorage.setItem('shareViewMode', viewMode.value);
};

// Remove Vietnamese diacritics for search
const removeVietnameseDiacritics = (str: string): string => {
  return str
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/đ/g, 'd')
    .replace(/Đ/g, 'D');
};

const deepSearchResults = ref<any[] | null>(null);

const filteredItems = computed(() => {
  if (deepSearchResults.value !== null) {
    return deepSearchResults.value;
  }

  if (!req.value?.items) return [];
  if (!searchQuery.value.trim()) return req.value.items;
  
  // Split search query into words, normalize each
  const searchWords = searchQuery.value
    .trim()
    .toLowerCase()
    .split(/\s+/)
    .map(word => removeVietnameseDiacritics(word));
  
  return req.value.items.filter(item => {
    const normalizedName = removeVietnameseDiacritics(item.name.toLowerCase());
    // All search words must be found in the filename (AND logic)
    return searchWords.every(word => normalizedName.includes(word));
  });
});

const performSearch = async () => {
    if (!searchQuery.value.trim()) return;
    
    // Check if we are already seeing deep results for this query?
    // User might want to refresh.
    
    layoutStore.loading = true;
    deepSearchResults.value = [];
    
    // Normalize query
    const query = removeVietnameseDiacritics(searchQuery.value.trim().toLowerCase());
    const searchWords = query.split(/\s+/);
    
    try {
        const results: any[] = [];
        // Start crawling from current view
        await crawl(route.path, searchWords, results);
        deepSearchResults.value = results;
    } catch (e) {
        console.error("Search failed", e);
        $showError(e as Error);
    } finally {
        layoutStore.loading = false;
    }
};

const crawl = async (path: string, searchWords: string[], accumulated: any[]) => {
    try {
        const res = await api.fetch(path, password.value);
        if (!res.items) return;
        
        for (const item of res.items) {
             const normalizedName = removeVietnameseDiacritics(item.name.toLowerCase());
             const match = searchWords.every(word => normalizedName.includes(word));
             
             if (match) {
                 // Clone item to modify name for display context if needed
                 // Showing full relative path can be helpful
                 // item.path e.g. "/folder/file.txt"
                 const displayItem = { ...item };
                 // Display path primarily if it's deep search, so user knows where it is
                 if (displayItem.path.startsWith('/')) {
                     displayItem.name = displayItem.path.substring(1);
                 } else {
                     displayItem.name = displayItem.path;
                 }
                 accumulated.push(displayItem);
             }
             
             if (item.isDir) {
                 await crawl(item.url, searchWords, accumulated);
             }
        }
    } catch (e) {
        console.error("Failed to crawl " + path, e);
    }
};

watch(searchQuery, (val) => {
    if (!val.trim()) {
        deepSearchResults.value = null;
    }
});

const fetchData = async () => {
  error.value = null;
  layoutStore.loading = true;
  
  try {
    const res = await api.fetch(route.path, password.value);
    req.value = res;
    fileStore.req = res;

    // Handle File View (if navigated directly to a file)
    if (!res.isDir) {
       // Determine type - fallback to extension if backend doesn't provide type
       let type = res.type as string;
       if (!type) {
          const ext = res.name.split('.').pop()?.toLowerCase();
          if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'].includes(ext || '')) type = 'image';
          else if (['mp4', 'mkv', 'avi', 'mov', 'webm'].includes(ext || '')) type = 'video';
          else if (['mp3', 'wav', 'ogg', 'flac', 'm4a'].includes(ext || '')) type = 'audio';
          else if (['txt', 'md', 'json', 'xml', 'js', 'ts', 'css', 'html', 'py', 'go', 'java', 'c', 'cpp'].includes(ext || '')) type = 'text';
       }

       // Text / Code
       if (['text', 'textImmutable'].includes(type) || type === 'application/json' || type === 'application/javascript') {
           // Fetch content
           try {
             // For public shares, content is not in metadata. Fetch it from dl url.
             const fileUrl = `${baseURL}/api/public/dl/${hash.value}${res.path}`;
             const contentRes = await fetch(fileUrl);
             if (contentRes.ok) {
                const text = await contentRes.text();
                editorData.value = {
                  filePath: res.path, 
                  fileName: res.name,
                  fileContent: text
                };
                isEditing.value = true;
                showMediaViewer.value = false;
             }
           } catch (e) {
             console.error("Failed to fetch file content", e);
           }
       } 
       // Media
       else if (['image', 'video', 'audio', 'blob'].includes(type)) { // 'blob' sometimes used for images
           mediaViewerItems.value = [res];
           mediaViewerStartIndex.value = 0;
           showMediaViewer.value = true;
           isEditing.value = false;
       }
       else {
           // Unsupported file type for inline viewing
           // Show a message or just let the sidebar handle download
           isEditing.value = false;
           showMediaViewer.value = false;
       }
    } else {
       // Reset views if we are in a directory
       isEditing.value = false;
       showMediaViewer.value = false;
    }

  } catch (e: any) {
    error.value = e;
    if (e.status === 401) {
      attemptedPasswordLogin.value = true;
    }
  } finally {
    layoutStore.loading = false;
  }
};

onMounted(() => {
  const parts = route.path.split("/");
  hash.value = parts[2];
  fetchData();
});

watch(() => route.path, () => {
  fetchData();
});

const download = () => {
  const url = `${baseURL}/api/public/dl/${hash.value}${route.path.replace(`/share/${hash.value}`, "")}`;
  window.open(url);
};

const toggleDownloadPanel = () => {
  showDownloadPanel.value = !showDownloadPanel.value;
  localStorage.setItem('showShareSidebar', String(showDownloadPanel.value));
};

const copyToClipboard = async (text: string) => {
  try {
    await copyText({ text });
    $showSuccess(t("success.linkCopied"));
  } catch (e) {
    $showError(e as Error);
  }
};

const base64 = (s: string) => window.btoa(unescape(encodeURIComponent(s)));

const onPanelResize = (width: number) => {
  previewPanelWidth.value = width;
};

// Editor & Media Viewer skipped for brevity of reconstruct
const isEditing = ref(false);
const editorData = ref<any>(null);
// ... other imports
import { computed, inject, onMounted, ref, watch } from "vue"; // Ensure imports are correct

// ... existing code

// Re-add canEdit computed property
const canEdit = computed(() => {
    // If it's a directory, check if share is writable
    if (req.value?.isDir) {
        return (req.value as any).canEdit === true;
    }
    // If it's a file, check if we are in a writable share
    return (req.value as any)?.canEdit === true;
}); 



const uploadInput = async (event: Event) => {
    const target = event.target as HTMLInputElement;
    if (!target.files || target.files.length === 0) return;
    
    layoutStore.loading = true;
    try {
        const files = Array.from(target.files);
        // Base path relative to share root
        let basePath = route.path.replace(`/share/${hash.value}`, "");
        if (!basePath.endsWith('/')) basePath += '/';
        
        for (const file of files) {
             const filePath = basePath + file.name;
             await api.upload(hash.value, filePath, file, false);
        }
        $showSuccess(t('success.upload'));
        fetchData();
    } catch (e) {
        $showError(e as Error);
    } finally {
        layoutStore.loading = false;
        target.value = ''; // Reset input
    }
};

const createFolder = async () => {
    const name = prompt(t('prompts.folderName') || 'Folder name');
    if (!name) return;
    
    layoutStore.loading = true;
    try {
        let basePath = route.path.replace(`/share/${hash.value}`, "");
        if (!basePath.endsWith('/')) basePath += '/';
        const folderPath = basePath + name;
        
        await api.createFolder(hash.value, folderPath);
        $showSuccess(t('success.folderCreated') || 'Folder created');
        fetchData();
    } catch (e) {
        $showError(e as Error);
    } finally {
        layoutStore.loading = false;
    }
};

const triggerUpload = () => {
    document.getElementById('share-upload-input')?.click();
};

const edit = async () => {
// ... existing edit logic

  try {
    const fileUrl = `${baseURL}/api/public/dl/${hash.value}${route.path.replace(`/share/${hash.value}`, "")}`;
    const response = await fetch(fileUrl, { credentials: 'include' });
    const content = await response.text();

    editorData.value = {
      filePath: route.path.replace(`/share/${hash.value}`, ""),
      fileName: req.value?.name,
      fileContent: content
    };
    isEditing.value = true;
  } catch (e) {
    $showError(e as Error);
  }
};

const closeEditor = () => {
  isEditing.value = false;
  editorData.value = null;
  if (req.value && !req.value.isDir) {
     const parts = route.path.split('/');
     parts.pop();
     router.push(parts.join('/'));
  }
};

const onEditorSaved = () => fetchData();

const togglePreviewPanel = () => layoutStore.togglePreviewPanel();

const showMediaViewer = ref(false);
const mediaViewerItems = ref<any[]>([]);
const mediaViewerStartIndex = ref(0);
const getMediaUrl = (item: any) => `${baseURL}/api/public/dl/${hash.value}${item.path}?inline=true`;
const closeMediaViewer = () => { showMediaViewer.value = false; };


</script>

<style scoped>
.share-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  /* padding-top: 4em; Account for fixed header */
}

/* Override header left position for Share page (no sidebar) */
.share-page :deep(header) {
  left: 0;
}

.share-main {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.share-sidebar {
  width: 320px;
  background: var(--surfacePrimary);
  border-right: 1px solid var(--divider);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow-y: auto;
  z-index: 10;
}

.share-sidebar.collapsed {
  width: 0;
  border-right-width: 0;
}

.share-sidebar-content {
  padding: 2em;
  min-width: 320px;
}

.share-info-header {
  text-align: center;
  margin-bottom: 2em;
}

.share-info-icon i {
  font-size: 64px;
  color: var(--blue);
  margin-bottom: 0.2em;
}

.share-info-title {
  font-size: 1.25em;
  font-weight: 600;
  word-break: break-all;
}

.share-info-subtitle {
  font-size: 0.9em;
  color: var(--textSecondary);
}

.share-content {
  flex: 1;
  overflow-y: auto;
  padding: 2em;
  background: var(--background);
  transition: margin-right 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.share-content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2em;
  gap: 2em;
}

.share-search-wrapper {
  flex: 1;
  max-width: 400px;
}

.search-box {
  display: flex;
  align-items: center;
  background: var(--surfacePrimary);
  border: 1px solid var(--divider);
  border-radius: 8px;
  padding: 0 12px;
  transition: all 0.2s;
}

.search-box:focus-within {
  border-color: var(--blue);
  box-shadow: 0 0 0 3px rgba(33, 150, 243, 0.1);
}

.search-input {
  width: 100%;
  padding: 10px 8px;
  border: none;
  background: transparent;
  outline: none;
  color: var(--textPrimary);
}

.clear-btn {
  font-size: 18px;
  cursor: pointer;
  color: var(--textSecondary);
}

.share-tree-section {
  margin-top: 2em;
  border-top: 1px solid var(--divider);
  padding-top: 1.5em;
}

.share-tree-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 0.8em;
  text-transform: uppercase;
  font-weight: 700;
  color: var(--textSecondary);
}

.share-qrcode-section {
  margin: 2em 0;
  text-align: center;
}

.qrcode-wrapper {
  background: white;
  padding: 12px;
  border-radius: 12px;
  display: inline-block;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.qrcode-hint {
  font-size: 0.8em;
  color: var(--textSecondary);
  margin-top: 0.8em;
}

.no-results {
  text-align: center;
  padding: 4em 0;
  color: var(--textSecondary);
}

.no-results i {
  font-size: 64px;
  opacity: 0.2;
}

.file-placeholder .actions {
    margin-top: 1em;
}

@media (max-width: 736px) {
  .share-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 280px;
    transform: translateX(-100%);
  }
  
  .share-sidebar:not(.collapsed) {
    transform: translateX(0);
    box-shadow: 4px 0 24px rgba(0,0,0,0.2);
  }
  
  .share-content-header {
    flex-direction: column;
    align-items: stretch;
  }
}

/* Scrollbar styling */
.share-sidebar::-webkit-scrollbar,
.share-content::-webkit-scrollbar {
  width: 6px;
}

.share-sidebar::-webkit-scrollbar-thumb,
.share-content::-webkit-scrollbar-thumb {
  background: var(--divider);
  border-radius: 3px;
}

.share-sidebar::-webkit-scrollbar-track,
.share-content::-webkit-scrollbar-track {
  background: transparent;
}
</style>

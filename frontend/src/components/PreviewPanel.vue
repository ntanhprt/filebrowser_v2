<template>
  <div class="preview-panel" v-if="layoutStore.showPreviewPanel">
    <div class="preview-panel-header">
      <div class="preview-panel-title">
        <i class="material-icons">preview</i>
        <span v-if="file">{{ file.name }}</span>
        <span v-else>{{ t('sidebar.preview') }}</span>
      </div>
      <div class="preview-panel-actions">
        <button
          v-if="isMarkdown && file"
          class="action"
          @click="toggleMarkdownView"
          :title="layoutStore.markdownViewMode === 'rendered' ? t('buttons.viewSource') : t('buttons.viewRendered')"
        >
          <i class="material-icons">{{ layoutStore.markdownViewMode === 'rendered' ? 'code' : 'preview' }}</i>
        </button>
        <button class="action" @click="openInNewTab" v-if="file" :title="t('buttons.openFile')">
          <i class="material-icons">open_in_new</i>
        </button>
        <button class="action" @click="closePanel" :title="t('buttons.close')">
          <i class="material-icons">close</i>
        </button>
      </div>
    </div>
    <div class="preview-panel-content">
      <div v-if="loading" class="loading-container">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
      </div>
      <div v-else-if="!file" class="no-file-selected">
        <i class="material-icons">touch_app</i>
        <p>{{ t('files.selectFileToPreview') }}</p>
      </div>
      <div v-else-if="file.isDir" class="folder-preview">
        <i class="material-icons">folder</i>
        <p>{{ file.name }}</p>
        <p class="folder-info">{{ t('files.folders') }}</p>
      </div>
      <template v-else>
        <!-- Image preview -->
        <div v-if="file.type === 'image'" class="image-preview">
          <img :src="previewUrl" :alt="file.name" />
        </div>
        <!-- Video preview -->
        <div v-else-if="file.type === 'video'" class="video-preview">
          <video controls :src="previewUrl"></video>
        </div>
        <!-- Audio preview -->
        <div v-else-if="file.type === 'audio'" class="audio-preview">
          <i class="material-icons audio-icon">audiotrack</i>
          <audio controls :src="previewUrl"></audio>
        </div>
        <!-- PDF preview -->
        <div v-else-if="isPdf" class="pdf-preview">
           <div v-if="loading && pdfPages.length === 0" class="loading-container">
               <div class="spinner">
                   <div class="bounce1"></div>
                   <div class="bounce2"></div>
                   <div class="bounce3"></div>
               </div>
           </div>
           <div v-else-if="pdfPages.length > 0" class="pdf-pages">
               <img v-for="(page, idx) in pdfPages" :key="idx" :src="page" class="pdf-page-img" />
           </div>
            <div v-else class="no-preview">
               <p>{{ t('files.noPreview') }}</p>
            </div>
        </div>
        <!-- Markdown preview -->
        <div v-else-if="isMarkdown" class="markdown-preview">
          <div v-if="layoutStore.markdownViewMode === 'rendered'" class="md_preview" v-html="renderedMarkdown"></div>
          <pre v-else class="source-code"><code>{{ fileContent }}</code></pre>
        </div>
        <!-- Text preview -->
        <div v-else-if="isText" class="text-preview">
          <pre class="source-code"><code>{{ fileContent }}</code></pre>
        </div>
        <!-- No preview available -->
        <div v-else class="no-preview">
          <i class="material-icons">visibility_off</i>
          <p>{{ t('files.noPreview') }}</p>
          <a :href="downloadUrl" class="button button--flat" target="_blank">
            <i class="material-icons">file_download</i>
            {{ t('buttons.download') }}
          </a>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useLayoutStore } from '@/stores/layout';
import { useI18n } from 'vue-i18n';
import { files as api } from '@/api';
import { baseURL } from '@/utils/constants';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import * as pdfjsLib from 'pdfjs-dist';

// Set the worker source for PDF.js
pdfjsLib.GlobalWorkerOptions.workerSrc = new URL(
  'pdfjs-dist/build/pdf.worker.min.mjs',
  import.meta.url
).toString();

const layoutStore = useLayoutStore();
const { t } = useI18n();

const loading = ref(false);
const fileContent = ref('');
const renderedMarkdown = ref('');
const pdfPages = ref<string[]>([]);

const file = computed(() => layoutStore.previewPanelFile);

const isMarkdown = computed(() => {
  if (!file.value) return false;
  const ext = file.value.name.toLowerCase();
  return ext.endsWith('.md') || ext.endsWith('.markdown');
});

const isPdf = computed(() => {
  if (!file.value) return false;
  return file.value.name.toLowerCase().endsWith('.pdf');
});

const isText = computed(() => {
  if (!file.value) return false;
  return file.value.type === 'text' || file.value.type === 'textImmutable';
});

const previewUrl = computed(() => {
  if (!file.value) return '';
  return `${baseURL}/api/raw${file.value.path}`;
});

const downloadUrl = computed(() => {
  if (!file.value) return '';
  return api.getDownloadURL(file.value, false);
});

const closePanel = () => {
  layoutStore.togglePreviewPanel();
};

const toggleMarkdownView = () => {
  layoutStore.toggleMarkdownViewMode();
};

const openInNewTab = () => {
  if (file.value) {
    window.open(previewUrl.value, '_blank');
  }
};

// Fetch and render PDF file
const fetchPdf = async () => {
  loading.value = true;
  pdfPages.value = [];
  try {
    const src = { url: previewUrl.value, withCredentials: true };
    const loadingTask = pdfjsLib.getDocument(src);
    const doc = await loadingTask.promise;
    
    // Limit pages to avoid memory issues in preview
    const numPages = Math.min(doc.numPages, 50);
    
    for (let i = 1; i <= numPages; i++) {
        const page = await doc.getPage(i);
        const viewport = page.getViewport({ scale: 1.5 });
        const canvas = document.createElement('canvas');
        const context = canvas.getContext('2d');
        canvas.height = viewport.height;
        canvas.width = viewport.width;
        
        if (context) {
            await page.render({
                canvasContext: context,
                viewport: viewport,
                canvas: context.canvas
            }).promise;
            pdfPages.value.push(canvas.toDataURL());
        }
    }
  } catch (error) {
    console.error('Failed to parse PDF file:', error);
  } finally {
    loading.value = false;
  }
};

// Fetch file content for text/markdown/pdf files
const fetchFileContent = async () => {
  if (!file.value || file.value.isDir) {
    fileContent.value = '';
    renderedMarkdown.value = '';
    pdfPages.value = [];
    return;
  }

  if (isPdf.value) {
    fileContent.value = '';
    renderedMarkdown.value = '';
    await fetchPdf();
  } else if (isText.value || isMarkdown.value) {
    pdfPages.value = [];
    loading.value = true;
    try {
      const response = await fetch(previewUrl.value, {
        credentials: 'include'
      });
      if (response.ok) {
        fileContent.value = await response.text();
        if (isMarkdown.value) {
          renderedMarkdown.value = DOMPurify.sanitize(await marked(fileContent.value));
        }
      } else {
        fileContent.value = '';
        renderedMarkdown.value = '';
      }
    } catch (error) {
      console.error('Failed to fetch file content:', error);
      fileContent.value = '';
      renderedMarkdown.value = '';
    } finally {
      loading.value = false;
    }
  } else {
    // Clear data
    fileContent.value = '';
    renderedMarkdown.value = '';
    pdfPages.value = [];
  }
};

// Watch for file changes
watch(file, () => {
  fetchFileContent();
}, { immediate: true });

// Watch for markdown view mode changes to re-render
watch(() => layoutStore.markdownViewMode, async () => {
  if (isMarkdown.value && fileContent.value && layoutStore.markdownViewMode === 'rendered') {
    renderedMarkdown.value = DOMPurify.sanitize(await marked(fileContent.value));
  }
});
</script>

<style scoped>
.preview-panel {
  position: fixed;
  top: 4em;
  right: 0;
  width: 400px;
  height: calc(100vh - 4em);
  background: var(--surfacePrimary);
  border-left: 1px solid var(--divider);
  z-index: 1001;
  display: flex;
  flex-direction: column;
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.preview-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75em 1em;
  border-bottom: 1px solid var(--divider);
  background: var(--surfaceSecondary);
  min-height: 3.5em;
  flex-shrink: 0;
  z-index: 10;
}

.preview-panel-title {
  display: flex;
  align-items: center;
  gap: 0.5em;
  font-weight: 500;
  color: var(--textPrimary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.preview-panel-title i {
  font-size: 1.2em;
  color: var(--blue);
  flex-shrink: 0;
}

.preview-panel-title span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-panel-actions {
  display: flex;
  gap: 0.5em;
  flex-shrink: 0;
  margin-left: 0.5em;
}

.preview-panel-actions .action {
  background: var(--surfacePrimary);
  border: 1px solid var(--divider);
  cursor: pointer;
  padding: 0.5em;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--textPrimary);
  transition: all 0.2s;
  min-width: 36px;
  min-height: 36px;
}

.preview-panel-actions .action:hover {
  background: var(--blue);
  color: white;
  border-color: var(--blue);
}

.preview-panel-actions .action i {
  font-size: 1.2em;
}

.preview-panel-content {
  flex: 1;
  overflow: auto;
  padding: 1em;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.no-file-selected,
.folder-preview,
.no-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--textSecondary);
  text-align: center;
}

.no-file-selected i,
.folder-preview i,
.no-preview i {
  font-size: 4em;
  margin-bottom: 0.5em;
  opacity: 0.5;
}

.folder-preview .folder-info {
  font-size: 0.9em;
  opacity: 0.7;
}

.image-preview {
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.image-preview img {
  max-width: 100%;
  max-height: calc(100vh - 150px);
  object-fit: contain;
  border-radius: 4px;
}

.video-preview video {
  width: 100%;
  max-height: calc(100vh - 150px);
  border-radius: 4px;
}

.audio-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2em;
}

.audio-preview .audio-icon {
  font-size: 4em;
  color: var(--iconPrimary);
  margin-bottom: 1em;
}

.audio-preview audio {
  width: 100%;
}

.pdf-preview {
  height: calc(100vh - 150px);
  overflow: hidden;
}

.pdf-preview object {
  width: 100%;
  height: 100%;
  border: none;
}

.pdf-pages {
    height: 100%;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 1em;
    gap: 1em;
    background: var(--surfaceSecondary);
}

.pdf-page-img {
    max-width: 100%;
    box-shadow: 0 2px 8px rgba(0,0,0,0.15);
    border-radius: 4px;
}

.markdown-preview,
.text-preview {
  height: 100%;
}

.source-code {
  background: var(--background);
  padding: 1em;
  border-radius: 4px;
  overflow: auto;
  max-height: calc(100vh - 150px);
  font-size: 0.9em;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--textPrimary);
  margin: 0;
}

.source-code code {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.no-preview .button {
  margin-top: 1em;
  display: flex;
  align-items: center;
  gap: 0.5em;
  padding: 0.75em 1.5em;
}

/* Spinner styles */
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
@media (max-width: 768px) {
  .preview-panel {
    width: 100%;
    border-left: none;
  }
}
</style>

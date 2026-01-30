<template>
  <!-- Resize handle -->
  <div
    v-if="layoutStore.showPreviewPanel"
    class="resize-handle"
    :style="{ right: panelWidth + 'px' }"
    @mousedown="startResize"
  >
    <div class="resize-handle-line"></div>
  </div>
  
  <div class="preview-panel" v-if="layoutStore.showPreviewPanel" :style="{ width: panelWidth + 'px' }">
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
        <button class="action" @click="openInNewTab" v-if="file && !file.isDir" :title="t('buttons.openFile')">
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
        <div v-if="isImage" class="image-preview">
          <img :src="previewUrl" :alt="file.name" />
        </div>
        <!-- Video preview -->
        <div v-else-if="isVideo" class="video-preview">
          <video controls :src="previewUrl"></video>
        </div>
        <!-- Audio preview -->
        <div v-else-if="isAudio" class="audio-preview">
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
        <!-- CSV preview -->
        <div v-else-if="isCsv" class="csv-preview">
          <div class="csv-table-container">
            <table v-if="csvData.length > 0" class="csv-table">
              <thead>
                <tr>
                  <th v-for="(header, index) in csvData[0]" :key="index">{{ header }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, rowIndex) in csvData.slice(1)" :key="rowIndex">
                  <td v-for="(cell, cellIndex) in row" :key="cellIndex">{{ cell }}</td>
                </tr>
              </tbody>
            </table>
            <p v-else class="no-data">{{ t('files.noPreview') || 'No data' }}</p>
          </div>
        </div>
        <!-- Excel preview (local rendering with xlsx) -->
        <div v-else-if="isExcel" class="excel-preview">
          <div v-if="excelData.length > 0" class="csv-table-container">
            <table class="csv-table">
              <thead>
                <tr>
                  <th v-for="(header, index) in excelData[0]" :key="index">{{ header }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, rowIndex) in excelData.slice(1)" :key="rowIndex">
                  <td v-for="(cell, cellIndex) in row" :key="cellIndex">{{ cell }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="office-preview">
            <i class="material-icons">table_chart</i>
            <p>{{ file.name }}</p>
            <p class="file-type">Excel File</p>
            <a :href="previewUrl" class="button button--flat" target="_blank">
              <i class="material-icons">file_download</i>
              {{ t('buttons.download') }}
            </a>
          </div>
        </div>
        <!-- Word preview (local rendering with mammoth) -->
        <div v-else-if="isWord" class="word-preview">
          <div v-if="wordHtml" class="word-content" v-html="wordHtml"></div>
          <div v-else class="office-preview">
            <i class="material-icons">description</i>
            <p>{{ file.name }}</p>
            <p class="file-type">Word Document</p>
            <a :href="previewUrl" class="button button--flat" target="_blank">
              <i class="material-icons">file_download</i>
              {{ t('buttons.download') }}
            </a>
          </div>
        </div>
        <!-- PowerPoint preview (no local library, show download) -->
        <div v-else-if="isPowerPoint" class="office-preview">
          <i class="material-icons">slideshow</i>
          <p>{{ file.name }}</p>
          <p class="file-type">PowerPoint Presentation</p>
          <p class="preview-note">Preview not available for PowerPoint files</p>
          <a :href="previewUrl" class="button button--flat" target="_blank">
            <i class="material-icons">file_download</i>
            {{ t('buttons.download') }}
          </a>
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
          <a :href="previewUrl" class="button button--flat" target="_blank">
            <i class="material-icons">file_download</i>
            {{ t('buttons.download') }}
          </a>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, onUnmounted, nextTick } from 'vue';
import { useLayoutStore } from '@/stores/layout';
import { useI18n } from 'vue-i18n';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import * as XLSX from 'xlsx';
import mammoth from 'mammoth';
import * as pdfjsLib from 'pdfjs-dist';

// Set the worker source for PDF.js
pdfjsLib.GlobalWorkerOptions.workerSrc = new URL(
  'pdfjs-dist/build/pdf.worker.min.mjs',
  import.meta.url
).toString();

const props = defineProps<{
  hash: string;
  token?: string;
}>();

const emit = defineEmits<{
  (e: 'resize', width: number): void;
}>();

const layoutStore = useLayoutStore();
const { t } = useI18n();

const loading = ref(false);
const fileContent = ref('');
const renderedMarkdown = ref('');
const csvData = ref<string[][]>([]);
const excelData = ref<string[][]>([]);
const wordHtml = ref('');
const pdfPages = ref<string[]>([]);

// Panel resize state
const panelWidth = ref(parseInt(localStorage.getItem('sharePreviewPanelWidth') || '400'));
const isResizing = ref(false);

// Resize handlers
const startResize = (e: MouseEvent) => {
  isResizing.value = true;
  document.addEventListener('mousemove', doResize);
  document.addEventListener('mouseup', stopResize);
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';
  e.preventDefault();
};

const doResize = (e: MouseEvent) => {
  if (!isResizing.value) return;
  
  const containerRight = window.innerWidth;
  const newWidth = containerRight - e.clientX;
  
  // Min width 200px, max width 70% of window
  const minWidth = 200;
  const maxWidth = window.innerWidth * 0.7;
  
  if (newWidth >= minWidth && newWidth <= maxWidth) {
    panelWidth.value = newWidth;
    localStorage.setItem('sharePreviewPanelWidth', String(newWidth));
    emit('resize', newWidth);
  }
};

const stopResize = () => {
  isResizing.value = false;
  document.removeEventListener('mousemove', doResize);
  document.removeEventListener('mouseup', stopResize);
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
};

onUnmounted(() => {
  document.removeEventListener('mousemove', doResize);
  document.removeEventListener('mouseup', stopResize);
});

const file = computed(() => layoutStore.previewPanelFile);

const isMarkdown = computed(() => {
  if (!file.value) return false;
  const name = file.value.name?.toLowerCase() || '';
  return name.endsWith('.md') || name.endsWith('.markdown');
});

const isPdf = computed(() => {
  if (!file.value) return false;
  return file.value.name?.toLowerCase().endsWith('.pdf');
});

const isImage = computed(() => {
  if (!file.value) return false;
  if (file.value.type === 'image') return true;
  const name = file.value.name?.toLowerCase() || '';
  return /\.(jpg|jpeg|png|gif|bmp|svg|webp)$/.test(name);
});

const isVideo = computed(() => {
  if (!file.value) return false;
  if (file.value.type === 'video') return true;
  const name = file.value.name?.toLowerCase() || '';
  return /\.(mp4|mkv|avi|mov|webm)$/.test(name);
});

const isAudio = computed(() => {
  if (!file.value) return false;
  if (file.value.type === 'audio') return true;
  const name = file.value.name?.toLowerCase() || '';
  return /\.(mp3|wav|ogg|flac|m4a)$/.test(name);
});

const isText = computed(() => {
  if (!file.value) return false;
  if (file.value.type === 'text' || file.value.type === 'textImmutable') return true;
  const name = file.value.name?.toLowerCase() || '';
  return /\.(txt|md|json|xml|js|ts|css|html|py|go|java|c|cpp|h|hpp|vue|sql|log|conf|yml|yaml|sh|bat)$/.test(name);
});

const isCsv = computed(() => {
  if (!file.value) return false;
  const name = file.value.name?.toLowerCase() || '';
  return name.endsWith('.csv');
});

const isExcel = computed(() => {
  if (!file.value) return false;
  const name = file.value.name?.toLowerCase() || '';
  return name.endsWith('.xlsx') || name.endsWith('.xls');
});

const isWord = computed(() => {
  if (!file.value) return false;
  const name = file.value.name?.toLowerCase() || '';
  return name.endsWith('.docx') || name.endsWith('.doc');
});

const isPowerPoint = computed(() => {
  if (!file.value) return false;
  const name = file.value.name?.toLowerCase() || '';
  return name.endsWith('.pptx') || name.endsWith('.ppt');
});

const previewUrl = computed(() => {
  if (!file.value) return '';
  
  // Use the share hash for public shares
  const shareHash = props.hash;
  if (!shareHash) {
    console.warn('SharePreviewPanel: No hash provided');
    return '';
  }
  
  // Get the file path - directly use the path property from the item
  // This matches how the download panel builds its URL
  let filePath = file.value.path;
  
  // Fallback: if path is empty, try to extract from url or use name
  if (!filePath || filePath === '/' || filePath === '') {
    if (file.value.url) {
      // url format: "/share/HASH/path/to/file"
      const urlParts = file.value.url.split('/');
      const shareIndex = urlParts.indexOf('share');
      if (shareIndex !== -1 && urlParts.length > shareIndex + 2) {
        filePath = '/' + urlParts.slice(shareIndex + 2).join('/');
      } else {
        filePath = `/${file.value.name || ''}`;
      }
    } else {
      filePath = `/${file.value.name || ''}`;
    }
  }
  
  // Ensure path starts with /
  if (filePath && !filePath.startsWith('/')) {
    filePath = '/' + filePath;
  }
  
  // Build URL - same format as download panel's raw computed
  let url = `${window.location.origin}/api/public/dl/${shareHash}${filePath}`;
  
  // Add token if available, and inline for PDF
  const params = new URLSearchParams();
  if (props.token) {
    params.append('token', props.token);
  }
  // Add inline parameter for PDF to display instead of download
  if (isPdf.value) {
    params.append('inline', 'true');
  }
  
  const queryString = params.toString();
  if (queryString) {
    url += `?${queryString}`;
  }
  
  console.log('SharePreviewPanel previewUrl:', { shareHash, filePath, url, file: file.value });
  return url;
});

// URL for Office Online Viewer (Microsoft)
const officeViewerUrl = computed(() => {
  if (!previewUrl.value) return '';
  // Microsoft Office Online Viewer - works with publicly accessible URLs
  return `https://view.officeapps.live.com/op/embed.aspx?src=${encodeURIComponent(previewUrl.value)}`;
});

// URL for Google Docs Viewer (alternative)
const googleViewerUrl = computed(() => {
  if (!previewUrl.value) return '';
  return `https://docs.google.com/gview?embedded=true&url=${encodeURIComponent(previewUrl.value)}`;
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

// Parse CSV content
const parseCsv = (content: string): string[][] => {
  const lines = content.split(/\r?\n/).filter(line => line.trim());
  return lines.map(line => {
    const result: string[] = [];
    let current = '';
    let inQuotes = false;
    
    for (let i = 0; i < line.length; i++) {
      const char = line[i];
      if (char === '"') {
        inQuotes = !inQuotes;
      } else if (char === ',' && !inQuotes) {
        result.push(current.trim());
        current = '';
      } else {
        current += char;
      }
    }
    result.push(current.trim());
    return result;
  });
};

// Clear all preview data
const clearPreviewData = () => {
  fileContent.value = '';
  renderedMarkdown.value = '';
  csvData.value = [];
  excelData.value = [];
  wordHtml.value = '';
  pdfPages.value = [];
};

// Fetch and render Excel file
const fetchExcelFile = async () => {
  loading.value = true;
  try {
    const response = await fetch(previewUrl.value, {
      credentials: 'include'
    });
    if (response.ok) {
      const arrayBuffer = await response.arrayBuffer();
      const workbook = XLSX.read(arrayBuffer, { type: 'array' });
      
      // Get first sheet
      const firstSheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[firstSheetName];
      
      // Convert to array of arrays
      const data = XLSX.utils.sheet_to_json(worksheet, { header: 1 }) as string[][];
      excelData.value = data;
    } else {
      excelData.value = [];
    }
  } catch (error) {
    console.error('Failed to parse Excel file:', error);
    excelData.value = [];
  } finally {
    loading.value = false;
  }
};

// Fetch and render Word file
const fetchWordFile = async () => {
  loading.value = true;
  try {
    const response = await fetch(previewUrl.value, {
      credentials: 'include'
    });
    if (response.ok) {
      const arrayBuffer = await response.arrayBuffer();
      const result = await mammoth.convertToHtml({ arrayBuffer });
      wordHtml.value = DOMPurify.sanitize(result.value);
    } else {
      wordHtml.value = '';
    }
  } catch (error) {
    console.error('Failed to parse Word file:', error);
    wordHtml.value = '';
  } finally {
    loading.value = false;
  }
};

// Fetch and render PDF file
const fetchPdf = async () => {
  loading.value = true;
  pdfPages.value = [];
  try {
    // Determine source
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

// Fetch file content for text/markdown/csv/pdf files
const fetchFileContent = async () => {
  if (!file.value || file.value.isDir) {
    clearPreviewData();
    return;
  }

  // Handle different file types
  if (isExcel.value) {
    clearPreviewData();
    await fetchExcelFile();
  } else if (isWord.value) {
    clearPreviewData();
    await fetchWordFile();
  } else if (isPdf.value) {
    clearPreviewData();
    await fetchPdf();
  } else if (isText.value || isMarkdown.value || isCsv.value) {
    clearPreviewData();
    loading.value = true;
    try {
      const response = await fetch(previewUrl.value, {
        credentials: 'include'
      });
      if (response.ok) {
        const content = await response.text();
        fileContent.value = content;
        
        if (isMarkdown.value) {
          renderedMarkdown.value = DOMPurify.sanitize(await marked(content));
        } else if (isCsv.value) {
          csvData.value = parseCsv(content);
        }
      }
    } catch (error) {
      console.error('Failed to fetch file content:', error);
    } finally {
      loading.value = false;
    }
  } else {
    // Clear data for other file types (image, video, audio)
    clearPreviewData();
  }
};

// Watch for file changes
watch(file, () => {
  fetchFileContent();
}, { immediate: true });

// Watch for hash changes
watch(() => props.hash, () => {
  if (file.value) {
    fetchFileContent();
  }
});

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

.pdf-preview iframe {
  width: 100%;
  height: 100%;
  border: none;
  border-radius: 4px;
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

/* CSV Preview */
.csv-preview {
  height: 100%;
  overflow: auto;
}

.csv-table-container {
  overflow: auto;
  max-height: calc(100vh - 150px);
}

.csv-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85em;
}

.csv-table th,
.csv-table td {
  border: 1px solid var(--divider);
  padding: 0.5em 0.75em;
  text-align: left;
  white-space: nowrap;
}

.csv-table th {
  background: var(--surfaceSecondary);
  font-weight: 600;
  position: sticky;
  top: 0;
  z-index: 1;
}

.csv-table tr:nth-child(even) {
  background: var(--background);
}

.csv-table tr:hover {
  background: rgba(33, 150, 243, 0.1);
}

.csv-preview .no-data {
  text-align: center;
  color: var(--textSecondary);
  padding: 2em;
}

/* Office Preview Container (with iframe viewer) */
.office-preview-container {
  position: relative;
  height: calc(100vh - 150px);
}

.office-viewer-iframe {
  width: 100%;
  height: 100%;
  border: none;
  border-radius: 4px;
  background: var(--surfaceSecondary);
}

.office-preview-fallback {
  display: none;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--textSecondary);
  text-align: center;
  gap: 0.5em;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--surfacePrimary);
}

.office-preview-fallback i {
  font-size: 4em;
  margin-bottom: 0.5em;
  color: var(--iconPrimary);
}

.office-preview-fallback .file-type {
  font-size: 0.9em;
  opacity: 0.7;
  margin-bottom: 1em;
}

.office-preview-fallback .button {
  display: flex;
  align-items: center;
  gap: 0.5em;
  padding: 0.75em 1.5em;
}

/* Original Office Preview (fallback without iframe) */
.office-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--textSecondary);
  text-align: center;
  gap: 0.5em;
}

.office-preview i {
  font-size: 4em;
  margin-bottom: 0.5em;
  color: var(--iconPrimary);
}

.office-preview .file-type {
  font-size: 0.9em;
  opacity: 0.7;
  margin-bottom: 1em;
}

.office-preview .button {
  display: flex;
  align-items: center;
  gap: 0.5em;
  padding: 0.75em 1.5em;
}

.office-preview .preview-note {
  font-size: 0.85em;
  font-style: italic;
  opacity: 0.6;
  margin-bottom: 1em;
}

/* Excel Preview */
.excel-preview {
  height: 100%;
}

/* Word Preview */
.word-preview {
  height: 100%;
  overflow: auto;
}

.word-content {
  padding: 1em;
  background: var(--surfacePrimary);
  border-radius: 4px;
  line-height: 1.6;
  max-height: calc(100vh - 150px);
  overflow: auto;
}

.word-content :deep(p) {
  margin: 0.5em 0;
}

.word-content :deep(h1),
.word-content :deep(h2),
.word-content :deep(h3),
.word-content :deep(h4) {
  margin: 1em 0 0.5em 0;
  color: var(--textPrimary);
}

.word-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 1em 0;
}

.word-content :deep(table td),
.word-content :deep(table th) {
  border: 1px solid var(--divider);
  padding: 0.5em;
}

.word-content :deep(img) {
  max-width: 100%;
  height: auto;
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

/* Resize handle */
.resize-handle {
  width: 8px;
  cursor: col-resize;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.2s;
  position: fixed;
  top: 4em;
  height: calc(100vh - 4em);
  z-index: 1002;
}

.resize-handle:hover {
  background: var(--divider);
}

.resize-handle-line {
  width: 2px;
  height: 40px;
  background: var(--divider);
  border-radius: 2px;
}

.resize-handle:hover .resize-handle-line {
  background: var(--blue);
}

/* Responsive */
@media (max-width: 768px) {
  .preview-panel {
    width: 100%;
    border-left: none;
  }
  
  .resize-handle {
    display: none;
  }
}
</style>

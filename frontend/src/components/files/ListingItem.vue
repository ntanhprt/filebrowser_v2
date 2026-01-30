<template>
  <div
    class="item"
    role="button"
    tabindex="0"
    :draggable="isDraggable"
    @dragstart="dragStart"
    @dragover="dragOver"
    @drop="drop"
    @click="itemClick"
    @mousedown="handleMouseDown"
    @mouseup="handleMouseUp"
    @mouseleave="handleMouseLeave"
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
    @touchcancel="handleTouchCancel"
    @touchmove="handleTouchMove"
    :data-dir="isDir"
    :data-type="type"
    :aria-label="name"
    :aria-selected="isSelected"
    :data-ext="getExtension(name).toLowerCase()"
    @contextmenu="contextMenu"
  >
    <!-- FIRST DIV: Icon/Thumbnail -->
    <div>
      <img
        v-if="!readOnly && type === 'image' && isThumbsEnabled"
        v-lazy="thumbnailUrl"
      />
      <i v-else class="material-icons"></i>
    </div>

    <!-- LAST DIV: Details -->
    <div>
      <p class="name">{{ name }}</p>

      <p v-if="isDir" class="size" data-order="-1">&mdash;</p>
      <p v-else class="size" :data-order="humanSize()">{{ humanSize() }}</p>

      <p class="modified">
        <time :datetime="modified">{{ humanTime() }}</time>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import { enableThumbs } from "@/utils/constants";
import { filesize } from "@/utils/index"; // Fixed path
import dayjs from "dayjs";
import { files as api } from "@/api";
import * as upload from "@/utils/upload";
import { computed, inject, ref } from "vue";
import { useRouter } from "vue-router";

const props = defineProps<{
  name: string;
  isDir: boolean;
  url: string;
  type: string;
  size: number;
  modified: string;
  index: number;
  readOnly?: boolean;
  path?: string;
}>();

const touches = ref<number>(0);
const longPressTimer = ref<number | null>(null);
const longPressTriggered = ref<boolean>(false);
const longPressDelay = ref<number>(500);
const startPosition = ref<{ x: number; y: number } | null>(null);
const moveThreshold = ref<number>(10);

const $showError = inject<any>("$showError")!;
const router = useRouter();

const authStore = useAuthStore();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const singleClick = computed(() => !props.readOnly && authStore.user?.singleClick);
const isSelected = computed(() => fileStore.selected.indexOf(props.index) !== -1);
const isDraggable = computed(() => !props.readOnly && authStore.user?.perm.rename);

const canDrop = computed(() => {
  if (!props.isDir || props.readOnly) return false;
  for (const i of fileStore.selected) {
    if (fileStore.req?.items[i].url === props.url) return false;
  }
  return true;
});

const thumbnailUrl = computed(() => {
  const file = { path: props.path, modified: props.modified };
  return api.getPreviewURL(file as Resource, "thumb");
});

const isThumbsEnabled = computed(() => enableThumbs);

const humanSize = () => props.type == "invalid_link" ? "invalid link" : filesize(props.size);

const humanTime = () => {
  if (!props.readOnly && authStore.user?.dateFormat) {
    return dayjs(props.modified).format("L LT");
  }
  return dayjs(props.modified).fromNow();
};

const dragStart = () => {
  if (fileStore.selectedCount === 0 || !isSelected.value) {
    fileStore.selected = [props.index];
  }
};

const dragOver = (event: Event) => {
  if (!canDrop.value) return;
  event.preventDefault();
};

const drop = async (event: Event) => {
  if (!canDrop.value) return;
  event.preventDefault();
  if (fileStore.selectedCount === 0) return;

  const items: any[] = [];
  for (const i of fileStore.selected) {
    if (fileStore.req) {
      items.push({
        from: fileStore.req.items[i].url,
        to: props.url + encodeURIComponent(fileStore.req.items[i].name),
        name: fileStore.req.items[i].name,
      });
    }
  }

  const action = (overwrite: boolean, rename: boolean) => {
    api.move(items, overwrite, rename).then(() => {
      fileStore.reload = true;
    }).catch($showError);
  };

  const currentDirItems = (await api.fetch(props.url)).items;
  const conflict = upload.checkConflict(items, currentDirItems);

  if (conflict) {
    layoutStore.showHover({
      prompt: "replace-rename",
      confirm: (event: Event, option: any) => {
        event.preventDefault();
        layoutStore.closeHovers();
        action(option === "overwrite", option === "rename");
      },
    });
    return;
  }
  action(false, false);
};

const itemClick = (event: Event | KeyboardEvent) => {
  if (longPressTriggered.value) {
    longPressTriggered.value = false;
    return;
  }
  
  // In readOnly mode (public shares)
  if (props.readOnly) {
    // For folders: navigate directly on single click
    if (props.isDir) {
      open();
      return;
    }
    // For files: select + show in preview panel if open
    click(event);
    // Also update preview panel if it's visible
    if (layoutStore.showPreviewPanel) {
      layoutStore.setPreviewPanelFile({ ...props });
    }
    return;
  }
  
  if (singleClick.value && !(event as KeyboardEvent).ctrlKey && !(event as KeyboardEvent).metaKey && !(event as KeyboardEvent).shiftKey && !fileStore.multiple) {
    open();
  } else {
    click(event);
  }
};

const contextMenu = (event: MouseEvent) => {
  event.preventDefault();
  if (fileStore.selected.indexOf(props.index) === -1) {
    click(event);
  }
};

const click = (event: Event | KeyboardEvent) => {
  if (!singleClick.value && fileStore.selectedCount !== 0) event.preventDefault();
  setTimeout(() => { touches.value = 0; }, 300);
  touches.value++;
  // Double-click: open(), but skip for files in readOnly (Share.vue handles via @dblclick)
  if (touches.value > 1) {
    if (props.readOnly && !props.isDir) {
      open();
      return;
    }
    open();
    return;
  }

  if (fileStore.selected.indexOf(props.index) !== -1) {
    if ((event as KeyboardEvent).ctrlKey || (event as KeyboardEvent).metaKey || fileStore.multiple) {
      fileStore.removeSelected(props.index);
    } else {
      fileStore.selected = [props.index];
    }
    return;
  }

  if ((event as KeyboardEvent).shiftKey && fileStore.selected.length > 0) {
    let fi = Math.min(props.index, fileStore.selected[0]);
    let la = Math.max(props.index, fileStore.selected[0]);
    for (let i = fi; i <= la; i++) {
      if (fileStore.selected.indexOf(i) === -1) fileStore.selected.push(i);
    }
    return;
  }

  if (!(event as KeyboardEvent).ctrlKey && !(event as KeyboardEvent).metaKey && !fileStore.multiple) {
    fileStore.selected = [];
  }
  fileStore.selected.push(props.index);
};

const open = () => {
  // For files (not folders)
  if (!props.isDir) {
    // In readOnly mode (public shares), we still want to navigate to the file URL
    // so Share.vue can handle displaying it (Editor/MediaViewer)
    if (props.readOnly) {
       router.push({ path: props.url });
       return;
    }
    // Normal mode with preview panel
    if (layoutStore.showPreviewPanel) {
      layoutStore.setPreviewPanelFile({ ...props });
      return;
    }
  }
  // For folders, navigate to the URL
  router.push({ path: props.url });
};

const getExtension = (name: string) => {
  const i = name.lastIndexOf(".");
  return i === -1 ? "" : name.substring(i);
};

// Touch/Long press handlers (simplified)
const startLongPress = (x: number, y: number) => {
  startPosition.value = { x, y };
  longPressTimer.value = window.setTimeout(() => {
    longPressTriggered.value = true;
    click(new Event("longpress"));
    cancelLongPress();
  }, longPressDelay.value);
};
const cancelLongPress = () => {
  if (longPressTimer.value) clearTimeout(longPressTimer.value);
  longPressTimer.value = null;
};
const handleMouseDown = (e: MouseEvent) => e.button === 0 && startLongPress(e.clientX, e.clientY);
const handleMouseUp = () => cancelLongPress();
const handleMouseLeave = () => cancelLongPress();
const handleTouchStart = (e: TouchEvent) => e.touches.length === 1 && startLongPress(e.touches[0].clientX, e.touches[0].clientY);
const handleTouchEnd = () => cancelLongPress();
const handleTouchCancel = () => cancelLongPress();
const handleTouchMove = (e: TouchEvent) => {
  if (e.touches.length === 1 && startPosition.value) {
    const dX = Math.abs(e.touches[0].clientX - startPosition.value.x);
    const dY = Math.abs(e.touches[0].clientY - startPosition.value.y);
    if (dX > moveThreshold.value || dY > moveThreshold.value) cancelLongPress();
  }
};
</script>

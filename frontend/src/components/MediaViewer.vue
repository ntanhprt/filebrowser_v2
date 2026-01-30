<template>
  <div v-if="isOpen" class="media-viewer-overlay" @click.self="close">
    <div class="media-viewer-container" :class="{ fullscreen: isFullscreen }">
      <!-- Header -->
      <div class="media-viewer-header">
        <span class="media-title">{{ currentItem?.name || 'Media' }}</span>
        <span class="media-counter">{{ currentIndex + 1 }} / {{ items.length }}</span>
        <div class="media-actions">
          <!-- Slideshow controls (images only) -->
          <div v-if="isImage && items.length > 1" class="slideshow-controls">
            <button @click="toggleSlideshow" :title="slideshowActive ? 'Stop Slideshow' : 'Start Slideshow'">
              <i class="material-icons">{{ slideshowActive ? 'pause_circle' : 'play_circle' }}</i>
            </button>
            <select v-model="slideshowSpeed" title="Slideshow Speed">
              <option value="2000">2s</option>
              <option value="3000">3s</option>
              <option value="5000">5s</option>
              <option value="8000">8s</option>
              <option value="10000">10s</option>
            </select>
          </div>
          <!-- Video speed controls -->
          <div v-if="isVideo" class="video-speed-controls">
            <select v-model="videoSpeed" @change="setVideoSpeed" title="Playback Speed">
              <option value="0.5">0.5x</option>
              <option value="1">1x</option>
              <option value="1.25">1.25x</option>
              <option value="1.5">1.5x</option>
              <option value="1.75">1.75x</option>
              <option value="2">2x</option>
              <option value="2.5">2.5x</option>
              <option value="3">3x</option>
            </select>
          </div>
          <button @click="toggleFullscreen" title="Fullscreen">
            <i class="material-icons">{{ isFullscreen ? 'fullscreen_exit' : 'fullscreen' }}</i>
          </button>
          <button @click="close" title="Close">
            <i class="material-icons">close</i>
          </button>
        </div>
      </div>

      <!-- Content -->
      <div class="media-viewer-content" 
           @touchstart="handleTouchStart" 
           @touchmove="handleTouchMove" 
           @touchend="handleTouchEnd">
        <!-- Previous button -->
        <button v-if="items.length > 1" class="nav-btn nav-prev" @click="prev">
          <i class="material-icons">chevron_left</i>
        </button>

        <!-- Image -->
        <div v-if="isImage" class="media-image-container">
          <img 
            :src="mediaUrl" 
            :alt="currentItem?.name" 
            @load="onImageLoad"
            :style="{ transform: `scale(${zoom}) rotate(${rotation}deg)` }"
          />
          <div class="image-controls">
            <button @click="zoomIn" title="Zoom In">
              <i class="material-icons">zoom_in</i>
            </button>
            <button @click="zoomOut" title="Zoom Out">
              <i class="material-icons">zoom_out</i>
            </button>
            <button @click="resetZoom" title="Reset">
              <i class="material-icons">refresh</i>
            </button>
            <button @click="rotateRight" title="Rotate">
              <i class="material-icons">rotate_right</i>
            </button>
          </div>
        </div>

        <!-- Video -->
        <div v-else-if="isVideo" class="media-video-container">
          <video 
            ref="videoRef"
            :src="mediaUrl" 
            controls 
            autoplay
            @loadedmetadata="onVideoLoad"
          >
            Your browser does not support video playback.
          </video>
        </div>

        <!-- Audio -->
        <div v-else-if="isAudio" class="media-audio-container">
          <i class="material-icons audio-icon">music_note</i>
          <p class="audio-title">{{ currentItem?.name }}</p>
          <audio ref="audioRef" :src="mediaUrl" controls autoplay>
            Your browser does not support audio playback.
          </audio>
        </div>

        <!-- Next button -->
        <button v-if="items.length > 1" class="nav-btn nav-next" @click="next">
          <i class="material-icons">chevron_right</i>
        </button>
      </div>

      <!-- Thumbnails -->
      <div v-if="items.length > 1 && showThumbnails" class="media-thumbnails">
        <div 
          v-for="(item, index) in items" 
          :key="index"
          class="thumbnail"
          :class="{ active: index === currentIndex }"
          @click="goTo(index)"
        >
          <img v-if="item.type === 'image'" :src="getThumbnailUrl(item)" :alt="item.name" />
          <i v-else-if="item.type === 'video'" class="material-icons">video_file</i>
          <i v-else-if="item.type === 'audio'" class="material-icons">audio_file</i>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps<{
  items: Array<{
    name: string;
    path: string;
    type: string;
    url?: string;
  }>;
  startIndex?: number;
  getMediaUrl: (item: any) => string;
  getThumbnailUrl?: (item: any) => string;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

// State
const isOpen = ref(true);
const currentIndex = ref(props.startIndex || 0);
const isFullscreen = ref(false);
const showThumbnails = ref(true);

// Image controls
const zoom = ref(1);
const rotation = ref(0);

// Video controls
const videoRef = ref<HTMLVideoElement | null>(null);
const videoSpeed = ref('1');

// Audio
const audioRef = ref<HTMLAudioElement | null>(null);

// Slideshow
const slideshowActive = ref(false);
const slideshowSpeed = ref('5000');
let slideshowTimer: number | null = null;

// Touch handling
let touchStartX = 0;
let touchStartY = 0;

// Computed
const currentItem = computed(() => props.items[currentIndex.value]);
const isImage = computed(() => currentItem.value?.type === 'image');
const isVideo = computed(() => currentItem.value?.type === 'video');
const isAudio = computed(() => currentItem.value?.type === 'audio');

const mediaUrl = computed(() => {
  if (!currentItem.value) return '';
  return props.getMediaUrl(currentItem.value);
});

// Methods
const close = () => {
  stopSlideshow();
  isOpen.value = false;
  emit('close');
};

const prev = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--;
  } else {
    currentIndex.value = props.items.length - 1;
  }
  resetImageControls();
};

const next = () => {
  if (currentIndex.value < props.items.length - 1) {
    currentIndex.value++;
  } else {
    currentIndex.value = 0;
  }
  resetImageControls();
};

const goTo = (index: number) => {
  currentIndex.value = index;
  resetImageControls();
};

const resetImageControls = () => {
  zoom.value = 1;
  rotation.value = 0;
};

// Zoom controls
const zoomIn = () => {
  zoom.value = Math.min(zoom.value + 0.25, 5);
};

const zoomOut = () => {
  zoom.value = Math.max(zoom.value - 0.25, 0.25);
};

const resetZoom = () => {
  zoom.value = 1;
  rotation.value = 0;
};

const rotateRight = () => {
  rotation.value = (rotation.value + 90) % 360;
};

// Video controls
const setVideoSpeed = () => {
  if (videoRef.value) {
    videoRef.value.playbackRate = parseFloat(videoSpeed.value);
  }
};

const onVideoLoad = () => {
  setVideoSpeed();
};

const onImageLoad = () => {
  // Image loaded
};

// Slideshow
const toggleSlideshow = () => {
  if (slideshowActive.value) {
    stopSlideshow();
  } else {
    startSlideshow();
  }
};

const startSlideshow = () => {
  slideshowActive.value = true;
  slideshowTimer = window.setInterval(() => {
    next();
  }, parseInt(slideshowSpeed.value));
};

const stopSlideshow = () => {
  slideshowActive.value = false;
  if (slideshowTimer) {
    clearInterval(slideshowTimer);
    slideshowTimer = null;
  }
};

watch(slideshowSpeed, () => {
  if (slideshowActive.value) {
    stopSlideshow();
    startSlideshow();
  }
});

// Fullscreen
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen();
    isFullscreen.value = true;
  } else {
    document.exitFullscreen();
    isFullscreen.value = false;
  }
};

// Touch handling for swipe
const handleTouchStart = (e: TouchEvent) => {
  touchStartX = e.touches[0].clientX;
  touchStartY = e.touches[0].clientY;
};

const handleTouchMove = (e: TouchEvent) => {
  // Prevent scroll while swiping
  if (Math.abs(e.touches[0].clientX - touchStartX) > 10) {
    e.preventDefault();
  }
};

const handleTouchEnd = (e: TouchEvent) => {
  const touchEndX = e.changedTouches[0].clientX;
  const diffX = touchEndX - touchStartX;
  
  if (Math.abs(diffX) > 50) {
    if (diffX > 0) {
      prev();
    } else {
      next();
    }
  }
};

// Keyboard handling
const handleKeydown = (e: KeyboardEvent) => {
  switch (e.key) {
    case 'ArrowLeft':
      prev();
      break;
    case 'ArrowRight':
      next();
      break;
    case 'Escape':
      close();
      break;
    case ' ':
      if (isImage.value) {
        toggleSlideshow();
      } else if (isVideo.value && videoRef.value) {
        if (videoRef.value.paused) {
          videoRef.value.play();
        } else {
          videoRef.value.pause();
        }
      }
      e.preventDefault();
      break;
    case '+':
    case '=':
      if (isImage.value) zoomIn();
      break;
    case '-':
      if (isImage.value) zoomOut();
      break;
    case 'f':
      toggleFullscreen();
      break;
  }
};

// Default thumbnail function
const getThumbnailUrl = (item: any) => {
  if (props.getThumbnailUrl) {
    return props.getThumbnailUrl(item);
  }
  return props.getMediaUrl(item);
};

onMounted(() => {
  window.addEventListener('keydown', handleKeydown);
  document.addEventListener('fullscreenchange', () => {
    isFullscreen.value = !!document.fullscreenElement;
  });
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown);
  stopSlideshow();
});
</script>

<style scoped>
.media-viewer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.95);
  z-index: 99999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.media-viewer-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.media-viewer-container.fullscreen {
  background: black;
}

.media-viewer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  background: rgba(0, 0, 0, 0.5);
  color: white;
}

.media-title {
  font-size: 16px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 40%;
}

.media-counter {
  font-size: 14px;
  color: #aaa;
}

.media-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.media-actions button,
.slideshow-controls button,
.video-speed-controls button {
  background: transparent;
  border: none;
  color: white;
  cursor: pointer;
  padding: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background 0.2s;
}

.media-actions button:hover,
.slideshow-controls button:hover {
  background: rgba(255, 255, 255, 0.1);
}

.slideshow-controls,
.video-speed-controls {
  display: flex;
  align-items: center;
  gap: 5px;
}

.slideshow-controls select,
.video-speed-controls select {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}

.media-viewer-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0, 0, 0, 0.5);
  border: none;
  color: white;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  transition: background 0.2s, transform 0.2s;
}

.nav-btn:hover {
  background: rgba(0, 0, 0, 0.8);
  transform: translateY(-50%) scale(1.1);
}

.nav-btn .material-icons {
  font-size: 36px;
}

.nav-prev {
  left: 20px;
}

.nav-next {
  right: 20px;
}

.media-image-container {
  max-width: 100%;
  max-height: calc(100vh - 150px);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.media-image-container img {
  max-width: 90vw;
  max-height: calc(100vh - 180px);
  object-fit: contain;
  transition: transform 0.3s ease;
}

.image-controls {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 10px;
  background: rgba(0, 0, 0, 0.7);
  padding: 8px 15px;
  border-radius: 25px;
}

.image-controls button {
  background: transparent;
  border: none;
  color: white;
  cursor: pointer;
  padding: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.2s;
}

.image-controls button:hover {
  background: rgba(255, 255, 255, 0.2);
}

.media-video-container {
  width: 100%;
  max-height: calc(100vh - 150px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.media-video-container video {
  max-width: 90vw;
  max-height: calc(100vh - 180px);
}

.media-audio-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  color: white;
}

.audio-icon {
  font-size: 100px;
  color: #666;
}

.audio-title {
  font-size: 20px;
}

.media-audio-container audio {
  width: 400px;
  max-width: 90vw;
}

.media-thumbnails {
  display: flex;
  gap: 5px;
  padding: 10px;
  background: rgba(0, 0, 0, 0.5);
  overflow-x: auto;
  justify-content: center;
}

.thumbnail {
  width: 60px;
  height: 60px;
  border: 2px solid transparent;
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
  transition: border-color 0.2s;
}

.thumbnail.active {
  border-color: #2196f3;
}

.thumbnail:hover {
  border-color: rgba(255, 255, 255, 0.5);
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail .material-icons {
  color: white;
  font-size: 24px;
}

/* Mobile responsive */
@media (max-width: 768px) {
  .nav-btn {
    width: 40px;
    height: 40px;
  }

  .nav-btn .material-icons {
    font-size: 24px;
  }

  .nav-prev {
    left: 10px;
  }

  .nav-next {
    right: 10px;
  }

  .media-title {
    max-width: 30%;
    font-size: 14px;
  }

  .image-controls {
    padding: 5px 10px;
  }

  .thumbnail {
    width: 50px;
    height: 50px;
  }
}
</style>

# TA File Browser

**Enhanced File Browser with Public Sharing & Media Viewer**

📧 **Author**: nt.anh.prt@gmail.com
📅 **Last Updated**: December 2025

---

## 🚀 Features Overview

### 1. Public Sharing System (Google Drive-like)

Share files and folders publicly with customizable permissions:

| Feature                       | Description                                         |
| ----------------------------- | --------------------------------------------------- |
| **View Only**           | Recipients can view and download, but cannot modify |
| **Edit - Everyone**     | Everyone with the link can view and modify          |
| **Edit - Select Users** | Only selected users can modify (Admin only)         |

#### How to Use:

1. Select a file/folder in **My Files**
2. Click **Share** button in the header
3. Click **Make Public** to create a public share
4. Choose permission level
5. Copy the share link or scan QR code

### 2. Public Shares Sidebar

View all public shares in a compact sidebar:

- 📂 **Collapsible panel** - Click to expand/collapse
- 🔍 **Search/Filter** - Find shares by name, path, or owner
- 📊 **Counter** - Shows filtered/total shares (5/10)
- 🎨 **Color-coded icons** - Different colors for file types
- 💬 **Rich tooltips** - Shows path, owner, permission, hash

### 3. Share Edit Functions

When viewing a shared folder with edit permission:

| Action                  | Description                          |
| ----------------------- | ------------------------------------ |
| **📁 New Folder** | Create new folder                    |
| **📄 New File**   | Create new empty file                |
| **⬆️ Upload**   | Upload files (button or drag & drop) |
| **🗑️ Delete**   | Delete selected items                |
| **✏️ Rename**   | Rename selected item                 |
| **📝 Edit**       | Open text files in ACE editor        |
| **ℹ️ Info**     | View file information                |
| **⬇️ Download** | Download selected items              |

### 4. Drag & Drop Upload

Upload files by dragging them into the share view:

- **Files**: Drag any file to upload
- **Folders**: Drag folders to upload entire directory structure
- **Visual feedback**: Blue overlay when dragging

### 5. Media Viewer

Full-featured media viewer for images, videos, and audio:

#### Image Features:

| Feature    | Shortcut | Description                           |
| ---------- | -------- | ------------------------------------- |
| Slideshow  | Space    | Auto-play with speed options (2s-10s) |
| Zoom       | +/-      | Zoom in/out                           |
| Rotate     | R        | Rotate image                          |
| Fullscreen | F        | Toggle fullscreen mode                |
| Prev/Next  | ←/→    | Navigate between images               |
| Thumbnails | -        | Thumbnail strip at bottom             |

#### Video Features:

| Feature        | Description                                |
| -------------- | ------------------------------------------ |
| Playback Speed | 0.5x, 1x, 1.25x, 1.5x, 1.75x, 2x, 2.5x, 3x |
| Controls       | Play/Pause, Seek, Volume                   |

#### Audio Features:

| Feature | Description                         |
| ------- | ----------------------------------- |
| Player  | Standard audio player with controls |

### 6. ACE Code Editor

Edit text files with syntax highlighting:

- **Syntax highlighting** - For various languages
- **Multiple themes** - Light and dark themes
- **Read-only mode** - For view-only shares
- **Save changes** - Save directly to shared folder

---

## 🛠️ Technical Details

### Frontend Components

| Component            | Path                                                | Description                       |
| -------------------- | --------------------------------------------------- | --------------------------------- |
| `Share.vue`        | `frontend/src/views/Share.vue`                    | Main share view with all features |
| `PublicShares.vue` | `frontend/src/components/PublicShares.vue`        | Sidebar for public shares         |
| `MediaViewer.vue`  | `frontend/src/components/MediaViewer.vue`         | Full-featured media viewer        |
| `ShareEditor.vue`  | `frontend/src/views/files/ShareEditor.vue`        | ACE editor for shares             |
| `PublicShare.vue`  | `frontend/src/components/prompts/PublicShare.vue` | Make Public prompt                |

### Backend Endpoints

| Endpoint                                         | Method | Description            |
| ------------------------------------------------ | ------ | ---------------------- |
| `/api/public/share/{hash}`                     | GET    | Get share info         |
| `/api/public/share/{hash}{path}`               | POST   | Upload file            |
| `/api/public/share/{hash}{path}/`              | POST   | Create folder          |
| `/api/public/share/{hash}{path}`               | DELETE | Delete file/folder     |
| `/api/public/share/{hash}{path}?action=rename` | PATCH  | Rename file            |
| `/api/public/dl/{hash}{path}`                  | GET    | Download file          |
| `/api/shares/public`                           | GET    | List all public shares |

### API Functions (pub.ts)

```typescript
// Fetch share info
fetch(url: string, password?: string): Promise<Resource>

// Upload file
upload(hash: string, filePath: string, content: Blob | string): Promise<Response>

// Create folder
createFolder(hash: string, folderPath: string): Promise<Response>

// Delete file/folder
remove(hash: string, filePath: string): Promise<Response>

// Rename file
rename(hash: string, srcPath: string, dstPath: string): Promise<Response>

// Save file content
save(hash: string, filePath: string, content: string): Promise<Response>

// Get download URL
getDownloadURL(res: Resource, inline?: boolean): string

// Download file(s)
download(format: DownloadFormat, hash: string, token: string, ...files: string[]): void
```

---

## 📦 Build Instructions

### Prerequisites

- Go 1.21+
- Node.js 22+
- pnpm

### Build Steps

```bash
# Clone the repository
git clone https://github.com/your-repo/filebrowser.git
cd filebrowser

# Build frontend
cd frontend
pnpm install
pnpm build
cd ..

# Build backend
go build -o filebrowser.exe .

# Or use the build script (Windows)
.\build.bat
```

### Run

```bash
# Start the server
./filebrowser -a 127.0.0.1 -p 8080 -r /path/to/files
```

---

## 📖 Usage Guide

### Creating a Public Share

1. Navigate to **My Files**
2. Select a file or folder
3. Click **Share** in the header bar
4. Click **Make Public** button
5. Choose permission:
   - **View Only**: Read-only access
   - **Change - Everyone**: Full edit access for all
   - **Change - Select Users**: Edit access for selected users (Admin)
6. Copy the share link

### Accessing a Public Share

1. Open the share link in a browser
2. If password protected, enter the password
3. View/download files
4. If edit permission: create, upload, edit, delete files

### Using Media Viewer

1. Double-click on an image, video, or audio file
2. Use keyboard shortcuts or buttons to navigate
3. Press **Esc** to close

---

## 🔐 Permissions

| Permission | View | Download | Upload | Edit | Delete | Create |
| ---------- | ---- | -------- | ------ | ---- | ------ | ------ |
| View Only  | ✅   | ✅       | ❌     | ❌   | ❌     | ❌     |
| Change     | ✅   | ✅       | ✅     | ✅   | ✅     | ✅     |

---

## 📝 Changelog

### December 2024 (nt.anh.prt@gmail.com)

#### New Features:

- ✨ **Public Sharing System** - Share files/folders with customizable permissions
- ✨ **Public Shares Sidebar** - Compact list with search/filter
- ✨ **Share Edit Functions** - Create, upload, edit, delete in shared folders
- ✨ **Drag & Drop Upload** - Drag files/folders to upload
- ✨ **Media Viewer** - Full-featured viewer for images/videos/audio
- ✨ **ACE Code Editor** - Edit text files with syntax highlighting
- ✨ **Slideshow Mode** - Auto-play images with speed control
- ✨ **Video Speed Control** - Playback speed from 0.5x to 3x

#### Bug Fixes:

- 🐛 Fixed 403 error when accessing Public Shares
- 🐛 Fixed 404 error when loading images in shared folders
- 🐛 Fixed hash collision between different shares
- 🐛 Fixed single file share not loading correctly

#### Improvements:

- 🎨 Compact design for Public Shares sidebar
- 🎨 Color-coded file type icons
- 🎨 Rich tooltips with detailed info
- 🎨 CSS consistent with system theme

---

## 📚 Base Code Reference

<details>
<summary>Click to expand original README</summary>

<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/filebrowser/master/branding/banner.png" width="550"/>
</p>

[![Build](https://github.com/filebrowser/filebrowser/actions/workflows/ci.yaml/badge.svg)](https://github.com/filebrowser/filebrowser/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/filebrowser/filebrowser/v2)](https://goreportcard.com/report/github.com/filebrowser/filebrowser/v2)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg)](https://github.com/filebrowser/filebrowser/releases/latest)

File Browser provides a file managing interface within a specified directory and it can be used to upload, delete, preview and edit your files. It is a **create-your-own-cloud**-kind of software where you can just install it on your server, direct it to a path and access your files through a nice web interface.

### Documentation

Documentation on how to install, configure, and contribute to this project is hosted at [filebrowser.org](https://filebrowser.org).

### Project Status

This project is a finished product which fulfills its goal: be a single binary web File Browser which can be run by anyone anywhere. That means that File Browser is currently on **maintenance-only** mode.

### Contributing

Contributions are always welcome. To start contributing to this project, read our [guidelines](CONTRIBUTING.md) first.

### License

[Apache License 2.0](LICENSE) © File Browser Contributors

</details>

---

## 📄 License

[Apache License 2.0](LICENSE) © File Browser Contributors

Enhanced by **nt.anh.prt@gmail.com**

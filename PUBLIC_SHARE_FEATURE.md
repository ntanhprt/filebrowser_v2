# Share Public Feature Implementation

## Overview
This document describes the implementation of the "Share Public" feature that allows users to share files/folders with all other users in a shared "Public Shares" section, with configurable permissions.

## Features
1. **Share Public Dialog**: Users can make existing shares public with permission options
2. **Permission Levels**:
   - **View Only**: Everyone can view and download, but cannot modify
   - **Change - Everyone**: Everyone can view, download, and modify
   - **Change - Select Users**: Only selected users can modify
3. **Public Shares List**: A dedicated sidebar section showing all public shares with quick access
4. **Permissions**: Respects existing share permissions

## Files Modified

### Backend (Go)

#### 1. `share/share.go`
- Added `PublicShareBody` struct for API request payload
- Extended `Link` struct with new fields:
  - `IsPublic bool` - marks as public share
  - `Permission string` - "view" or "change"
  - `AllowedUsers []uint` - userIDs allowed to edit

#### 2. `storage/bolt/share.go`
- Added `GetPublic()` - retrieves all public shares
- Added `GetPublicByUserID(id uint)` - retrieves public shares by user

#### 3. `share/storage.go`
- Updated `StorageBackend` interface with new methods:
  - `GetPublic()` - interface method
  - `GetPublicByUserID(id)` - interface method
- Implemented wrapper methods in `Storage` struct with expiration handling

#### 4. `http/share.go`
- Added `publicSharesListHandler` - GET /api/publicshares - retrieves all public shares
- Added `sharePublicHandler` - PUT /api/share/{hash}/public - converts a share to public
- Added `sharePrivateHandler` - PUT /api/share/{hash}/private - removes public sharing

#### 5. `http/http.go`
- Added route: `GET /api/publicshares` → `publicSharesListHandler`
- Added route: `PUT /api/share/{hash}/public` → `sharePublicHandler`
- Added route: `PUT /api/share/{hash}/private` → `sharePrivateHandler`

### Frontend (Vue.js/TypeScript)

#### 1. `frontend/src/api/share.ts`
- Added `listPublic()` - fetches all public shares
- Added `makePublic(hash, permission, allowedUsers)` - converts share to public
- Added `makePrivate(hash)` - removes public sharing

#### 2. `frontend/src/components/prompts/Share.vue`
- Added "Make Public" button (public icon) in share list table
- Added "Remove Public" button (lock icon) for public shares
- Added `openPublicShareDialog(link)` method
- Added `makePrivate(link)` method

#### 3. `frontend/src/components/prompts/PublicShare.vue` (NEW)
- New component for public share permission dialog
- Shows three permission options with descriptions
- Conditional user selection dropdown for "Change - Select Users" option
- Integrated with layout store for modal interactions

#### 4. `frontend/src/components/PublicShares.vue` (NEW)
- New sidebar component showing all public shares
- Displays share name, permission level, and owner
- Features:
  - Loading state with spinner
  - "No public shares" message
  - Copy link button
  - Open in new tab button
  - Automatic refresh capability

#### 5. `frontend/src/components/prompts/Prompts.vue`
- Imported `PublicShare` component
- Registered `publicShare` in components map

#### 6. `frontend/src/components/Sidebar.vue`
- Imported `PublicShares` component
- Registered component in sidebar
- Added public shares section (visible if user has share permission)

## Data Flow

### Making a Share Public
1. User clicks "Make Public" button on an existing share
2. `PublicShare` dialog opens with permission options
3. User selects permission level and optional users
4. Dialog calls `api.share.makePublic(hash, permission, allowedUsers)`
5. Backend updates the share with `IsPublic=true` and permission details
6. Frontend updates the share list

### Viewing Public Shares
1. `PublicShares` component loads on sidebar mount
2. Calls `api.share.listPublic()` to fetch all public shares
3. Displays shares with owner info and permission badges
4. User can click link to view/access or open in new tab
5. Can copy share link to clipboard

## API Endpoints

### GET /api/publicshares
Returns all public shares (permissions checked by `withPermShare`)
```json
[
  {
    "hash": "aB3cDe",
    "path": "/documents/report.pdf",
    "userID": 2,
    "expire": 0,
    "is_public": true,
    "permission": "view",
    "allowed_users": []
  }
]
```

### PUT /api/share/{hash}/public
Makes a share public
Request:
```json
{
  "permission": "view|change",
  "allowed_users": [1, 2, 3]  // optional, for "change" permission
}
```
Response: Updated share object

### PUT /api/share/{hash}/private
Removes public sharing from a share

## Permission Model

- **Admin users**: Can see all public shares, make any share public/private
- **Regular users**: Can only make their own shares public/private, can see all public shares
- **Public share access**: Depends on permission level:
  - **view**: Can download and view
  - **change**: Can modify based on `AllowedUsers` array

## Testing Instructions

1. **Build the application**: Run `build.bat` to compile Go backend and Vue frontend
2. **Run filebrowser.exe**: Start the server
3. **Create a test share**: 
   - Go to Files view
   - Select a file/folder
   - Click Share button
   - Create a share link (with or without password/expiration)
4. **Make it public**:
   - Click the "public" icon on the share
   - Select permission level
   - Optionally select users for edit access
   - Click "Make Public"
5. **View in Public Shares**:
   - Check sidebar for "Public Shares" section
   - See the shared item listed with permission level
   - Click to open or copy link
6. **Test different permission levels**:
   - View Only: Try to modify (should be denied)
   - Change-Everyone: Modify freely
   - Change-Select Users: Modify with selected users only

## Future Enhancements

1. Add search/filter in public shares list
2. Add categories or grouping in public shares sidebar
3. Add statistics on public share access
4. Add bulk operations for public shares
5. Add audit logging for public share access
6. Add expiration dates for public shares

## Notes

- All public shares respect existing expiration times (automatically deleted when expired)
- Changes to share permissions trigger database updates
- Public shares are visible to all users with share permission
- The `AllowedUsers` array is stored in the database and used for access control
- Share links remain accessible as long as they're not deleted or expire

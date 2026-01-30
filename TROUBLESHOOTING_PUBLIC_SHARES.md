# Troubleshooting: User Không Truy Cập Được Public Shares

## Vấn Đề

```
Log lỗi:
/auth/: 404 ... GetFileAttributesEx D:\...\users\tata\auth: The system cannot find the file specified.
/branding/: 404 ... GetFileAttributesEx D:\...\users\tata\branding: The system cannot find the file specified.
```

**User `tata` không thể truy cập các folders đã được public share.**

## Nguyên Nhân: User Scope Configuration

### File Browser Scope System Hoạt Động Như Thế Nào?

Mỗi user trong File Browser có một **Scope** (phạm vi truy cập):
```
User: tata
Scope: /users/tata
↓
User chỉ có thể truy cập: D:\...\filebrowser\users\tata\
```

Khi user truy cập `/auth/`:
```
1. User request: /files/auth/
2. Backend áp dụng scope: Scope + Path = /users/tata + /auth = /users/tata/auth
3. Backend tìm file: D:\...\filebrowser\users\tata\auth
4. Folder không tồn tại → 404 Error
```

## Giải Pháp: 3 Cách Fix

### ✅ Giải Pháp 1: Share Folders TRONG User Scope (KHUYẾN NGHỊ)

**Nguyên tắc:**
- Chỉ share folders/files nằm TRONG scope mà users có thể truy cập
- Tạo shared folders trong common area

**Các bước:**

#### A. Tạo Shared Folder Structure
```
project_root/
├── users/
│   ├── user1/         ← User 1 scope
│   ├── user2/         ← User 2 scope
│   └── tata/          ← Tata scope
└── shared/            ← ⭐ COMMON AREA
    ├── documents/
    ├── projects/
    └── resources/
```

#### B. Cấu Hình User Scope
```
1. Login as Admin
2. Settings → Users → Edit User (tata)
3. Scope: Đổi từ "/users/tata" → "/" hoặc "/shared"
4. Save
```

**Ví dụ cụ thể:**
```
Admin:
1. Tạo folder: D:\...\filebrowser\shared\team-docs
2. Upload files vào folder này
3. Tạo share link cho /shared/team-docs
4. Make Public → Change-Everyone

User tata (scope = "/shared" hoặc "/"):
1. Sidebar → Public Shares → Click "team-docs"
2. Navigate đến /files/shared/team-docs
3. ✅ Access granted! (vì scope bao gồm /shared)
```

### ⚠️ Giải Pháp 2: Thay Đổi User Scope Thành Root

**Cách làm:**
```
Admin → Settings → Users → tata
Scope: "/" (root)
```

**Ưu điểm:**
- ✅ User có thể truy cập mọi folder được share
- ✅ Đơn giản, không cần tạo structure mới

**Nhược điểm:**
- ❌ Mất tính isolation giữa users
- ❌ User có thể thấy folders của users khác
- ❌ Security risk

**⚠️ CHỈ DÙNG cho:**
- Trusted users
- Small teams
- Internal environments

### 🔧 Giải Pháp 3: Sử Dụng Rules Để Restrict Access

**Advanced solution - Kết hợp root scope + rules:**

```
1. Set user scope = "/"
2. Add Rules để restrict:
   - ALLOW: /shared/*
   - DENY: /users/other_user/*
   - ALLOW: /users/tata/*
```

**Cách cấu hình:**
```
Admin → Settings → Users → tata
Scope: "/"
Rules:
  - regex: ^/users/(?!tata).*$
    allow: false
  - regex: ^/shared/.*$
    allow: true
```

## So Sánh Các Giải Pháp

| Giải pháp | Khó | An toàn | Linh hoạt | Use Case |
|-----------|-----|---------|-----------|----------|
| 1. Common folder | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | Team collaboration |
| 2. Root scope | ⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ | Small trusted team |
| 3. Scope + Rules | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Enterprise |

## Quick Fix (Cho Testing)

**Nhanh nhất để test Public Shares:**

1. **Login as Admin**
2. **Vào Settings → Users → tata**
3. **Thay đổi:**
   ```
   Scope: "/" (thay vì "/users/tata")
   ```
4. **Save**
5. **Logout & Login lại**
6. **Test Public Shares** → Giờ sẽ work!

## Workflow Đúng (Production)

### Setup lần đầu:

#### Admin:
```bash
# 1. Tạo shared folder structure
mkdir D:\...\filebrowser\shared
mkdir D:\...\filebrowser\shared\documents
mkdir D:\...\filebrowser\shared\projects

# 2. Copy files vào shared
# Files should be in D:\...\filebrowser\shared\...
```

#### File Browser Settings:
```
1. Admin login
2. Settings → Users
3. For each user:
   - Scope: "/shared" (hoặc "/" nếu muốn flexible)
   - Permissions: Bật "Create", "Modify", "Delete" nếu cần
```

#### Tạo Public Shares:
```
1. Admin tạo share cho /shared/documents
2. Make Public → Change-Everyone
3. Tất cả users giờ thấy trong Public Shares
4. Click vào → Navigate đến /files/shared/documents
5. ✅ Access granted!
```

## Debug Checklist

Khi user không truy cập được public share:

- [ ] **Check User Scope**
  ```
  Settings → Users → {username} → Scope = ?
  ```
  
- [ ] **Check File Path**
  ```
  Share path: /shared/docs
  User scope: /shared
  Result: ✅ /shared + /docs = /shared/docs (OK)
  
  Share path: /shared/docs
  User scope: /users/tata
  Result: ❌ /users/tata + /shared/docs = không match
  ```

- [ ] **Check Physical Files**
  ```
  Share path: /shared/docs
  Physical location: D:\...\filebrowser\shared\docs
  Exists? Yes ✅ / No ❌
  ```

- [ ] **Check Permissions**
  ```
  User permissions: Create, Modify, Delete enabled?
  ```

## Lưu Ý Quan Trọng

### ⚠️ Share Path PHẢI Nằm Trong User Scope

```
❌ WRONG:
User scope: /users/tata
Share path: /auth/  → 404 Error (ngoài scope)

✅ CORRECT:
User scope: /
Share path: /auth/  → OK

✅ CORRECT:
User scope: /shared
Share path: /shared/docs  → OK

✅ CORRECT:
User scope: /users/tata
Share path: /users/tata/files  → OK
```

### 📁 Recommended Folder Structure

```
project_root/
├── shared/              ← For public shares
│   ├── company/
│   ├── projects/
│   └── resources/
├── users/               ← Private user folders
│   ├── user1/
│   ├── user2/
│   └── tata/
└── admin/               ← Admin only
    └── backups/

User Scopes:
- Admin: "/"
- Regular users: "/" or "/shared,/users/{username}"
- Restricted users: "/users/{username}"
```

## Kết Luận

**Vấn đề KHÔNG phải bug của Public Shares feature.**

**Nguyên nhân:** User scope configuration sai.

**Fix nhanh:** Đổi user scope thành "/" (root)

**Fix đúng:** Tạo /shared folder structure và config scopes properly.

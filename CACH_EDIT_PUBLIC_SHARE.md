# Cách Edit File/Folder Trong Public Share

## Hiện Trạng

Tính năng "Public Shares" hiện tại **CHỈ** cung cấp:
- ✅ Hiển thị danh sách shares công khai trong sidebar
- ✅ Copy link và mở share link
- ✅ Lưu metadata permission (view/change)

**NHƯNG:**
- ❌ Chưa implement permission checking khi user truy cập share link
- ❌ Share link hiện tại là READ-ONLY (giống như share link thông thường)

## Cách Hoạt Động Của Share Link (Hiện Tại)

### 1. User Truy Cập Share Link
```
User click vào Public Share → Mở /share/{hash}
→ Xem được files/folders
→ Download được
→ NHƯNG KHÔNG thể upload/edit/delete
```

### 2. Tại Sao Không Edit Được?

Share link trong File Browser được thiết kế là **READ-ONLY** view:
- **File `http/public.go`** xử lý `/share/{hash}` 
- Chỉ cho phép:
  - Browse (xem danh sách files)
  - Download files
  - View file content
- **KHÔNG** cho phép:
  - Upload files
  - Delete files
  - Rename files
  - Create folders

## Giải Pháp: 2 Cách Để User Edit Được

### ⭐ Giải Pháp 1: Truy Cập Trực Tiếp Qua Account (KHUYẾN NGHỊ)

**Cách hoạt động:**
1. User phải login vào File Browser
2. Vào phần "My Files" 
3. Navigate đến folder đã được public share
4. Edit trực tiếp (vì user đã đăng nhập và có quyền)

**Ưu điểm:**
- ✅ Không cần thay đổi code
- ✅ An toàn hơn (có authentication)
- ✅ Đầy đủ tính năng
- ✅ Audit trail (biết ai đã edit)

**Nhược điểm:**
- ❌ Mỗi user phải có account
- ❌ Phải biết đường dẫn đến folder

**Workflow:**
```
1. Admin tạo public share: /shared/team-documents → Permission: Change-Everyone
2. User A login → My Files → Navigate đến /shared/team-documents
3. User A có thể upload/edit/delete (vì có quyền trên folder gốc)
4. Public share chỉ để "discover" - cho mọi người biết folder này tồn tại
```

### ⭐⭐ Giải Pháp 2: Implement Edit Permission Trong Share View (CẦN CODE)

**Cần implement thêm:**

#### A. Backend Changes

1. **Modify `http/public.go`** - Add permission checking:
```go
// Khi user truy cập /share/{hash}
// Check xem share.Permission == "change"
// Nếu có, enable upload/edit actions
```

2. **Add new handlers** cho edit operations trong share view:
```go
// POST /api/public/share/{hash}/upload
// PUT /api/public/share/{hash}/rename
// DELETE /api/public/share/{hash}/delete
```

3. **Permission validation**:
```go
func canEdit(share *share.Link, userID uint) bool {
    if !share.IsPublic || share.Permission != "change" {
        return false
    }
    if len(share.AllowedUsers) == 0 {
        return true // Everyone can edit
    }
    return contains(share.AllowedUsers, userID)
}
```

#### B. Frontend Changes

1. **Modify `views/Share.vue`**:
- Hiển thị upload button nếu có quyền edit
- Enable context menu với delete/rename
- Show "Drop files here" area

2. **Add permission indicator**:
```vue
<div v-if="canEdit" class="edit-enabled">
  ✏️ You can edit this share
</div>
```

## Workflow Đề Xuất (Giải Pháp 1 - Không Cần Code)

### Scenario: Team Collaboration Folder

#### Setup (Admin):
```
1. Tạo folder: /shared/team-alpha
2. Cấp quyền cho users: user1, user2, user3 (via user permissions)
3. Tạo share link cho folder
4. Make public với permission: View Only
```

#### User Workflow:
```
User 1:
1. Login vào account
2. Sidebar → Public Shares → Click "Team Alpha"
3. Copy link hoặc ghi nhớ path: /shared/team-alpha
4. Vào My Files → Navigate đến /shared/team-alpha
5. Upload/Edit files (vì user1 có quyền trên folder gốc)

User 2, User 3: Làm tương tự
```

#### Lợi ích:
- ✅ "Public Shares" = Discovery tool (tìm folder shared)
- ✅ Actual editing = Qua My Files với proper authentication
- ✅ Security: Chỉ users có quyền mới edit được
- ✅ Không cần viết code mới

## So Sánh 2 Giải Pháp

| Tiêu chí | Giải Pháp 1 (Via Account) | Giải Pháp 2 (Edit In Share) |
|----------|---------------------------|------------------------------|
| Code mới | Không cần | Cần nhiều code |
| Security | ⭐⭐⭐⭐⭐ Tốt | ⭐⭐⭐ Trung bình |
| User experience | ⭐⭐⭐ Cần login | ⭐⭐⭐⭐⭐ Tiện lợi |
| Authentication | Bắt buộc login | Có thể anonymous |
| Audit trail | ⭐⭐⭐⭐⭐ Đầy đủ | ⭐⭐ Giới hạn |
| Implementation time | ✅ Ngay lập tức | ⏱️ 2-3 ngày |

## Khuyến Nghị

### Cho Use Case Hiện Tại: **Dùng Giải Pháp 1**

**Tại sao?**
1. ✅ Không cần viết code mới
2. ✅ An toàn hơn
3. ✅ Đã có đầy đủ tính năng
4. ✅ Proper audit trail

**Cách sử dụng "Public Shares" đúng:**
- **Public Shares** = "Bulletin Board" - nơi announcement các folders/files quan trọng
- **My Files** = Nơi thực sự làm việc với files

### Nếu Cần Giải Pháp 2:

**Khi nào cần:**
- Users không muốn login
- Cần anonymous collaboration
- Guest users cần edit
- External partners không có account

**Công việc cần làm:**
1. Modify `http/public.go` - 200+ lines code
2. Add permission checking middleware
3. Update frontend Share.vue - 150+ lines
4. Add upload handler cho share view
5. Security testing
6. Update documentation

**Thời gian ước tính:** 2-3 ngày làm việc

## Kết Luận

**Câu trả lời ngắn gọn:**
- ✅ Public Share **hiện tại** chỉ để "discover" (tìm và copy link)
- ✅ Để **edit**, user cần login và truy cập qua **My Files**
- ✅ Permission "change" trong public share = metadata cho tương lai
- ⚠️ Nếu cần edit trực tiếp trong share view, cần implement thêm code

**Next steps:**
1. Xác định: Có cần edit trực tiếp trong share view không?
2. Nếu KHÔNG → Dùng workflow hiện tại (login → My Files)
3. Nếu CÓ → Implement Giải pháp 2 (tốn 2-3 ngày)

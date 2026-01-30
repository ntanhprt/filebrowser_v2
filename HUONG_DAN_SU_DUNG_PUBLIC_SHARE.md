# Hướng Dẫn Sử Dụng Tính Năng "Public Shares"

## Cách Đưa File/Folder Vào "Public Shares" (Sidebar)

### Bước 1: Tạo Share Link Thông Thường
1. Vào phần **Files** (My Files)
2. Chọn 1 file hoặc folder bạn muốn chia sẻ
3. Click nút **Share** (biểu tượng chia sẻ) trên thanh công cụ
4. Trong dialog "Share", click nút **New** để tạo share link mới
5. Cấu hình:
   - **Share Duration**: Thời gian hết hạn (có thể để 0 = vĩnh viễn)
   - **Optional Password**: Mật khẩu bảo vệ (tùy chọn)
6. Click **Share** để tạo

### Bước 2: Biến Share Link Thành "Public Share"
1. Sau khi tạo share link, bạn sẽ thấy danh sách các share links
2. Ở mỗi share link, có các nút:
   - 📋 Copy link
   - 📥 Copy download link
   - **🌐 Public icon** (nút này để "Make Public")
   - 🗑️ Delete
3. Click vào nút **🌐 (Public icon)** 
4. Một dialog mới sẽ hiện lên: "Make Share Public"

### Bước 3: Chọn Quyền Truy Cập
Trong dialog "Make Share Public", chọn một trong ba tùy chọn:

#### Option 1: View Only
- ✅ Tất cả users có thể xem và tải xuống
- ❌ Không ai có thể chỉnh sửa
- **Dùng cho**: Tài liệu công khai chỉ đọc

#### Option 2: Change - Everyone
- ✅ Tất cả users có thể xem, tải xuống
- ✅ Tất cả users có thể chỉnh sửa
- **Dùng cho**: Thư mục làm việc chung

#### Option 3: Change - Select Users
- ✅ Tất cả users có thể xem, tải xuống
- ✅ Chỉ những users được chọn có thể chỉnh sửa
- Sau khi chọn option này, danh sách users sẽ hiện ra
- Check vào các users bạn muốn cho phép chỉnh sửa
- **Dùng cho**: Dự án với nhóm người cụ thể

### Bước 4: Xác Nhận
1. Click nút **Make Public**
2. Thông báo "Shared publicly" sẽ hiện ra
3. Icon 🌐 sẽ chuyển thành 🔒 (cho biết đã public, click để remove)

### Bước 5: Xem Trong "Public Shares"
1. Mở **Sidebar** (click vào menu hamburger ☰)
2. Cuộn xuống phần **"Public Shares"** (có icon 🌐)
3. Bạn sẽ thấy:
   - Tên file/folder vừa share
   - Permission badge (View hoặc Edit)
   - Tên người tạo share (by username)
   - Nút copy link và open

## Lưu Ý Quan Trọng

### Ai Có Thể Thấy "Public Shares"?
- ✅ **Tất cả users** trong hệ thống đều thấy phần "Public Shares" trong sidebar
- ✅ Mọi user đều thấy tất cả các items được đánh dấu "public"
- ✅ Chỉ cần user có quyền **Share** permission

### Khác Biệt Giữa Share Link Thông Thường và Public Share

| Đặc điểm | Share Link Thông Thường | Public Share |
|----------|------------------------|--------------|
| Hiển thị trong sidebar | ❌ Không | ✅ Có (cho tất cả users) |
| Cần biết link | ✅ Phải có link | ❌ Không cần, thấy trong sidebar |
| Quản lý tập trung | ❌ Mỗi user tự quản lý | ✅ Hiển thị chung cho mọi người |
| Quyền chỉnh sửa | ❌ Không có | ✅ Có thể cấu hình |

## Ví Dụ Thực Tế

### Ví Dụ 1: Chia Sẻ Tài Liệu Công Ty
```
1. Tạo folder: /documents/company-policies
2. Share folder với thời hạn permanent
3. Click 🌐 Make Public
4. Chọn: View Only
5. → Tất cả nhân viên thấy trong "Public Shares"
6. → Mọi người có thể đọc nhưng không sửa
```

### Ví Dụ 2: Thư Mục Làm Việc Nhóm
```
1. Tạo folder: /projects/team-alpha
2. Share folder với thời hạn permanent
3. Click 🌐 Make Public
4. Chọn: Change - Select Users
5. Chọn các thành viên team: user1, user2, user3
6. → Tất cả users thấy trong "Public Shares"
7. → Chỉ user1, user2, user3 có thể chỉnh sửa
```

### Ví Dụ 3: Thư Mục Upload Chung
```
1. Tạo folder: /shared/uploads
2. Share folder với thời hạn permanent
3. Click 🌐 Make Public
4. Chọn: Change - Everyone
5. → Tất cả users thấy và có thể upload/edit
```

## Quản Lý Public Shares

### Xóa Khỏi Public Shares
1. Vào phần **Share** của file/folder
2. Click vào icon **🔒** (lock icon) bên cạnh share đã public
3. Share sẽ biến thành share link thông thường
4. Biến mất khỏi "Public Shares" sidebar

### Xóa Hoàn Toàn
1. Vào phần **Share** của file/folder
2. Click nút **🗑️ Delete** để xóa share link hoàn toàn

## Troubleshooting

### Không Thấy Nút "Public" Icon
- ✅ Kiểm tra user có quyền **Share** permission
- ✅ Đảm bảo đã tạo share link trước

### Không Thấy "Public Shares" Trong Sidebar
- ✅ Kiểm tra user có quyền **Share** permission
- ✅ Click vào menu hamburger ☰ để mở sidebar
- ✅ Cuộn xuống, phần "Public Shares" ở giữa Settings và Logout

### Public Share Không Hiển Thị Trong Sidebar
- ✅ Refresh lại trang (F5)
- ✅ Click nút **Refresh** (↻) trong phần "Public Shares"
- ✅ Kiểm tra backend có chạy đúng không

### Lỗi "(void 0) is not a function"
- ✅ Đảm bảo đã build lại frontend sau khi sửa code
- ✅ Chạy lại `build.bat`
- ✅ Restart filebrowser.exe

## Workflow Tổng Quát

```
User A                          System                          User B, C, D
  |                               |                                  |
  | 1. Create Share Link          |                                  |
  |------------------------------>|                                  |
  |                               |                                  |
  | 2. Click "Make Public"        |                                  |
  |------------------------------>|                                  |
  |                               |                                  |
  | 3. Select Permission          |                                  |
  |------------------------------>|                                  |
  |                               |                                  |
  |                               | 4. Update Database               |
  |                               |    (IsPublic = true)             |
  |                               |                                  |
  |                               | 5. Notify All Users              |
  |                               |----------------------------------| 
  |                               |                                  |
  |                               |      6. Show in "Public Shares"  |
  |                               |<---------------------------------|
  |                               |                                  |
  | 7. Confirm Success            |                                  |
  |<------------------------------|                                  |
  |                               |                                  |
```

## Kết Luận

"Public Shares" là một khu vực tập trung trong sidebar nơi tất cả users có thể:
- 👀 Xem tất cả files/folders được chia sẻ công khai
- 🔗 Truy cập nhanh qua link
- 📋 Copy link để chia sẻ
- ✏️ Chỉnh sửa (nếu có quyền)

Đây là giải pháp tuyệt vời cho:
- Tài liệu công ty
- Dự án nhóm
- Thư mục upload chung
- Knowledge base
- Resource sharing

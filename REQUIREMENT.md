# Yêu cầu Thực Hành — Khóa Gin + RealWorld API

## 1. Cấu Trúc Dự Án

- Tạo project Go sử dụng module (`go mod init ...`)
- Sử dụng Gin làm HTTP framework
- Thiết kế cấu trúc thư mục rõ ràng (ví dụ: `controllers/`, `services/`, `repositories/`, `models/`, `middlewares/`, `dto/`, ...)
- Dùng cơ sở dữ liệu (tùy chọn: Postgres, MySQL, SQLite), sử dụng ORM (ví dụ: GORM) hoặc SQL thuần, tùy học viên/giảng viên quyết định
- Cấu hình biến môi trường để quản lý các cấu hình như kết nối DB, JWT secret

## 2. Authentication (Xác Thực Người Dùng)

Theo spec RealWorld:

### Đăng Ký (Registration)

- **POST** `/api/users`
- **Body JSON:** `{ "user": { "username", "email", "password" } }`
- Sau khi đăng ký thành công, trả về đối tượng User (theo format RealWorld)

### Đăng Nhập (Login)

- **POST** `/api/users/login`
- **Body:** `email`, `password`
- Nếu đăng nhập thành công, trả về User + JWT token (token dùng để gọi các endpoint khác)

### Lấy Thông Tin Người Dùng Hiện Tại

- **GET** `/api/user`

## 3. Articles (Bài Viết)

### Feed Bài Viết

- **GET** `/api/articles/feed`
- Chỉ dành cho người dùng đã xác thực (authentication required)
- Trả feed các bài viết từ người dùng mà người dùng hiện tại theo dõi

### Lấy Chi Tiết Một Bài Viết

- **GET** `/api/articles/:slug`
- Không cần authentication
- Trả về toàn bộ chi tiết bài viết (title, body, description, tagList, author, ngày tạo, cập nhật, ...)

### Tạo Bài Viết Mới

- **POST** `/api/articles`
- **Body:** article object theo spec: `title`, `description`, `body`, tùy chọn `tagList` (mảng string)
- Cần xác thực (token)
- Trả về bài viết vừa tạo (article)

### Cập Nhật Bài Viết

- **PUT** `/api/articles/:slug`
- **Body:** article object với các trường optional: `title`, `description`, `body`
- Cần xác thực
- Nếu đổi title, slug cũng cần cập nhật theo spec
- Trả về bài viết được cập nhật

### Xoá Bài Viết

- **DELETE** `/api/articles/:slug`
- Cần xác thực
- Xoá bài viết

## 4. Comments (Bình Luận)

### Thêm Bình Luận Vào Bài Viết

- **POST** `/api/articles/:slug/comments`
- **Body:** `{ "comment": { "body": "..." } }`
- Cần xác thực
- Trả về đối tượng Comment tạo thành công

### Lấy Bình Luận Từ Bài Viết

- **GET** `/api/articles/:slug/comments`
- Authentication optional
- Trả về danh sách các comment (theo spec)

### Xoá Bình Luận

- **DELETE** `/api/articles/:slug/comments/:id`
- Cần xác thực
- Xoá comment theo id

## 5. Favorites (Yêu Thích / Thích Bài Viết)

### Thích Bài Viết

- **POST** `/api/articles/:slug/favorite`
- Cần xác thực
- Trả về bài viết (có cập nhật số lượng `favoritesCount` và trạng thái `favorited`)

### Bỏ Thích Bài Viết

- **DELETE** `/api/articles/:slug/favorite`
- Cần xác thực
- Trả về bài viết (cập nhật lại)

## 6. Tags

### Lấy Danh Sách Tags

- **GET** `/api/tags`
- Không cần xác thực
- Trả về danh sách tag (mảng string)

## 7. Yêu Cầu Kỹ Thuật & Bổ Sung

- **JWT:** Sử dụng JWT để tạo token khi đăng nhập, và middleware Gin để validate token trong các endpoint cần authentication
- **Validation:** Validate dữ liệu request (ví dụ: email, password, title, body) — nếu data thiếu hoặc sai định dạng, trả lỗi phù hợp (HTTP status code + message)
- **Error Handling:** Thiết kế lỗi rõ ràng theo JSON (không chỉ HTTP 500) khi có lỗi business (ví dụ: user tồn tại, không có permission, slug không tìm thấy...)
- **Database Schema:** Thiết kế bảng/collection để lưu Users, Profiles, Articles, Comments, Favorites (quan hệ giữa user và bài viết) và Tags
- **Pagination:** Cho endpoint list articles (`/api/articles`) xử lý limit và offset như spec yêu cầu
- **Slug Generation:** Tự động tạo slug cho bài viết từ title khi tạo bài mới, và cập nhật slug khi đổi title
- **Unit / Integration Test:** Viết test cho các layer: service + repository + controller (handler) + middleware. Kiểm thử các endpoint chính: đăng ký, login, CRUD bài viết, bình luận, favorite, profile
- **Postman / Swagger:** Tạo Postman collection để test các API hoặc dùng Swagger / OpenAPI (nếu muốn) để document API của bạn
- **Docker (tùy chọn):** Có thể tạo Dockerfile và docker-compose để dễ deploy DB + ứng dụng backend

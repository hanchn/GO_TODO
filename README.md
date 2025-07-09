# 学生管理系统

基于Go语言开发的现代化学生信息管理系统，采用MVC架构模式，使用Gin框架、GORM和SQLite数据库。

## 技术栈

- **后端框架**: Gin (Go语言最流行的Web框架)
- **数据库**: SQLite (轻量级本地数据库)
- **ORM**: GORM (Go语言最流行的ORM库)
- **前端**: Bootstrap 5 + 原生JavaScript
- **架构模式**: MVC (Model-View-Controller)

## 功能特性

- ✅ 学生信息的增删改查 (CRUD)
- ✅ 多条件搜索 (姓名、专业、年级)
- ✅ 响应式Web界面
- ✅ RESTful API设计
- ✅ 数据验证和错误处理
- ✅ 软删除功能
- ✅ 自动数据库迁移

## 项目结构

```
student-management-system/
├── main.go                 # 程序入口
├── go.mod                  # Go模块依赖
├── config/
│   └── database.go         # 数据库配置
├── models/
│   └── student.go          # 学生数据模型
├── controllers/
│   └── student_controller.go # 学生控制器
├── routes/
│   └── routes.go           # 路由配置
└── views/
    └── templates/
        ├── index.html      # 首页模板
        └── students.html   # 学生管理页面
```

## 安装和运行

### 前置要求

- Go 1.21 或更高版本
- Git

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd student-management-system
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **运行项目**
   ```bash
   go run main.go
   ```

4. **访问应用**
   - Web界面: http://localhost:8080
   - 学生管理: http://localhost:8080/students
   - API文档见下方

## API文档

### 基础URL
```
http://localhost:8080/api/v1
```

### 学生管理接口

| 方法 | 路径 | 描述 | 参数 |
|------|------|------|------|
| GET | `/students` | 获取所有学生 | - |
| GET | `/students/:id` | 根据ID获取学生 | id (路径参数) |
| GET | `/students/search` | 搜索学生 | name, major, grade (查询参数) |
| POST | `/students` | 创建新学生 | JSON请求体 |
| PUT | `/students/:id` | 更新学生信息 | id (路径参数) + JSON请求体 |
| DELETE | `/students/:id` | 删除学生 | id (路径参数) |

### 学生数据模型

```json
{
  "id": 1,
  "name": "张三",
  "age": 20,
  "gender": "男",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "major": "计算机科学与技术",
  "grade": "2023级",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### API使用示例

#### 1. 创建学生
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "age": 20,
    "gender": "男",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "major": "计算机科学与技术",
    "grade": "2023级"
  }'
```

#### 2. 获取所有学生
```bash
curl http://localhost:8080/api/v1/students
```

#### 3. 搜索学生
```bash
curl "http://localhost:8080/api/v1/students/search?name=张&major=计算机"
```

#### 4. 更新学生
```bash
curl -X PUT http://localhost:8080/api/v1/students/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三丰",
    "age": 21
  }'
```

#### 5. 删除学生
```bash
curl -X DELETE http://localhost:8080/api/v1/students/1
```

## 数据验证规则

- **姓名**: 必填，最大长度100字符
- **年龄**: 必填，范围1-150
- **性别**: 必填，只能是"男"或"女"
- **邮箱**: 必填，必须是有效的邮箱格式，唯一
- **电话**: 必填，最大长度20字符
- **专业**: 必填，最大长度100字符
- **年级**: 必填，最大长度50字符

## 开发说明

### MVC架构说明

- **Model (models/)**: 定义数据结构和数据库操作
- **View (views/)**: 前端页面模板
- **Controller (controllers/)**: 处理业务逻辑和HTTP请求

### 数据库

项目使用SQLite作为本地数据库，数据库文件为`students.db`，会在首次运行时自动创建。

### 扩展功能

如需添加新功能，建议按照以下步骤：

1. 在`models/`中定义新的数据模型
2. 在`controllers/`中实现业务逻辑
3. 在`routes/`中添加新的路由
4. 在`views/`中创建对应的前端页面

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。
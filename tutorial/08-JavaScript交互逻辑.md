# 第8章：JavaScript交互逻辑

本章将详细介绍如何使用现代JavaScript开发前端交互功能，包括API调用、数据处理、用户界面交互和错误处理。

## 8.1 JavaScript技术栈

### 技术选择

我们使用的JavaScript技术包括：

- 🚀 **ES6+**：现代JavaScript语法和特性
- 🌐 **Fetch API**：现代HTTP请求处理
- 📦 **模块化**：ES6模块系统
- 🎯 **事件驱动**：DOM事件处理
- 📊 **数据绑定**：双向数据绑定
- 🔄 **异步处理**：Promise和async/await
- 🛡️ **错误处理**：统一错误处理机制
- 📱 **响应式**：移动端交互优化

### 项目结构

```
static/js/
├── app.js              # 主应用文件
├── config.js           # 配置文件
├── api/
│   ├── client.js       # API客户端
│   ├── student.js      # 学生API
│   └── auth.js         # 认证API
├── components/
│   ├── student.js      # 学生组件
│   ├── dashboard.js    # 仪表板组件
│   ├── modal.js        # 模态框组件
│   └── table.js        # 表格组件
├── utils/
│   ├── helpers.js      # 工具函数
│   ├── validation.js   # 验证工具
│   ├── storage.js      # 存储工具
│   └── dom.js          # DOM操作工具
└── services/
    ├── notification.js # 通知服务
    ├── cache.js        # 缓存服务
    └── analytics.js    # 分析服务
```

## 8.2 核心应用框架

### 主应用文件

**文件路径：** `static/js/app.js`

```javascript
/**
 * 学生管理系统主应用
 * @author Your Name
 * @version 1.0.0
 */

// 应用配置
const App = {
    config: {
        apiBaseUrl: '/api/v1',
        version: '1.0.0',
        debug: true,
        pageSize: 10,
        maxFileSize: 2 * 1024 * 1024, // 2MB
        supportedImageTypes: ['image/jpeg', 'image/png', 'image/gif']
    },
    
    // 应用状态
    state: {
        currentUser: null,
        currentPage: 1,
        totalPages: 0,
        selectedStudents: new Set(),
        filters: {},
        sortBy: 'created_at',
        sortOrder: 'desc'
    },
    
    // 缓存
    cache: new Map(),
    
    // 事件总线
    events: new EventTarget(),
    
    // 初始化应用
    async init() {
        try {
            this.log('应用初始化开始...');
            
            // 初始化组件
            await this.initComponents();
            
            // 绑定全局事件
            this.bindGlobalEvents();
            
            // 初始化路由
            this.initRouter();
            
            // 加载用户信息
            await this.loadUserInfo();
            
            // 初始化通知系统
            NotificationService.init();
            
            this.log('应用初始化完成');
            
            // 触发应用就绪事件
            this.events.dispatchEvent(new CustomEvent('app:ready'));
            
        } catch (error) {
            this.error('应用初始化失败:', error);
            this.showError('应用初始化失败，请刷新页面重试');
        }
    },
    
    // 初始化组件
    async initComponents() {
        const components = [
            'StudentManager',
            'Dashboard',
            'ModalManager',
            'TableManager'
        ];
        
        for (const componentName of components) {
            if (window[componentName] && typeof window[componentName].init === 'function') {
                await window[componentName].init();
                this.log(`组件 ${componentName} 初始化完成`);
            }
        }
    },
    
    // 绑定全局事件
    bindGlobalEvents() {
        // 页面加载完成
        document.addEventListener('DOMContentLoaded', () => {
            this.hideLoadingScreen();
        });
        
        // 全局错误处理
        window.addEventListener('error', (event) => {
            this.error('全局错误:', event.error);
        });
        
        // 未处理的Promise拒绝
        window.addEventListener('unhandledrejection', (event) => {
            this.error('未处理的Promise拒绝:', event.reason);
            event.preventDefault();
        });
        
        // 网络状态变化
        window.addEventListener('online', () => {
            this.showSuccess('网络连接已恢复');
        });
        
        window.addEventListener('offline', () => {
            this.showWarning('网络连接已断开');
        });
        
        // 键盘快捷键
        document.addEventListener('keydown', (event) => {
            this.handleKeyboardShortcuts(event);
        });
        
        // 页面可见性变化
        document.addEventListener('visibilitychange', () => {
            if (document.visibilityState === 'visible') {
                this.events.dispatchEvent(new CustomEvent('app:focus'));
            } else {
                this.events.dispatchEvent(new CustomEvent('app:blur'));
            }
        });
    },
    
    // 处理键盘快捷键
    handleKeyboardShortcuts(event) {
        // Ctrl/Cmd + K: 全局搜索
        if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
            event.preventDefault();
            const searchInput = document.getElementById('global-search');
            if (searchInput) {
                searchInput.focus();
            }
        }
        
        // Ctrl/Cmd + N: 新建学生
        if ((event.ctrlKey || event.metaKey) && event.key === 'n') {
            event.preventDefault();
            if (window.StudentManager) {
                StudentManager.showAddModal();
            }
        }
        
        // ESC: 关闭模态框
        if (event.key === 'Escape') {
            const openModal = document.querySelector('.modal.show');
            if (openModal) {
                const modal = bootstrap.Modal.getInstance(openModal);
                if (modal) {
                    modal.hide();
                }
            }
        }
    },
    
    // 初始化路由
    initRouter() {
        // 简单的路由系统
        const path = window.location.pathname;
        const routes = {
            '/': 'dashboard',
            '/students': 'students',
            '/dashboard': 'dashboard'
        };
        
        const currentRoute = routes[path] || 'dashboard';
        this.state.currentRoute = currentRoute;
        
        // 更新导航状态
        this.updateNavigation(currentRoute);
    },
    
    // 更新导航状态
    updateNavigation(activeRoute) {
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
        });
        
        const activeLink = document.querySelector(`[href="/${activeRoute}"]`);
        if (activeLink) {
            activeLink.classList.add('active');
        }
    },
    
    // 加载用户信息
    async loadUserInfo() {
        try {
            // 这里可以从API加载用户信息
            this.state.currentUser = {
                id: 1,
                name: '管理员',
                email: 'admin@example.com',
                avatar: '/static/images/avatar.png'
            };
        } catch (error) {
            this.error('加载用户信息失败:', error);
        }
    },
    
    // 隐藏加载屏幕
    hideLoadingScreen() {
        const loadingScreen = document.getElementById('loading-screen');
        if (loadingScreen) {
            loadingScreen.style.opacity = '0';
            setTimeout(() => {
                loadingScreen.style.display = 'none';
            }, 300);
        }
    },
    
    // 显示成功消息
    showSuccess(message, options = {}) {
        NotificationService.show(message, 'success', options);
    },
    
    // 显示错误消息
    showError(message, options = {}) {
        NotificationService.show(message, 'error', options);
    },
    
    // 显示警告消息
    showWarning(message, options = {}) {
        NotificationService.show(message, 'warning', options);
    },
    
    // 显示信息消息
    showInfo(message, options = {}) {
        NotificationService.show(message, 'info', options);
    },
    
    // 日志记录
    log(...args) {
        if (this.config.debug) {
            console.log('[App]', ...args);
        }
    },
    
    // 错误记录
    error(...args) {
        console.error('[App Error]', ...args);
        
        // 发送错误到分析服务
        if (window.AnalyticsService) {
            AnalyticsService.trackError(args[1] || args[0]);
        }
    },
    
    // 格式化日期
    formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
        if (!date) return '-';
        
        const d = new Date(date);
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        const hours = String(d.getHours()).padStart(2, '0');
        const minutes = String(d.getMinutes()).padStart(2, '0');
        const seconds = String(d.getSeconds()).padStart(2, '0');
        
        return format
            .replace('YYYY', year)
            .replace('MM', month)
            .replace('DD', day)
            .replace('HH', hours)
            .replace('mm', minutes)
            .replace('ss', seconds);
    },
    
    // 防抖函数
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },
    
    // 节流函数
    throttle(func, limit) {
        let inThrottle;
        return function() {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }
};

// 全局暴露App对象
window.App = App;

// 页面加载完成后初始化应用
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => App.init());
} else {
    App.init();
}
```

### API客户端

**文件路径：** `static/js/api/client.js`

```javascript
/**
 * API客户端
 * 统一处理HTTP请求和响应
 */

class APIClient {
    constructor(baseURL = '/api/v1') {
        this.baseURL = baseURL;
        this.defaultHeaders = {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        };
        this.interceptors = {
            request: [],
            response: []
        };
        
        // 添加默认拦截器
        this.addDefaultInterceptors();
    }
    
    // 添加默认拦截器
    addDefaultInterceptors() {
        // 请求拦截器：添加认证头
        this.interceptors.request.push((config) => {
            const token = localStorage.getItem('auth_token');
            if (token) {
                config.headers['Authorization'] = `Bearer ${token}`;
            }
            return config;
        });
        
        // 请求拦截器：添加请求ID
        this.interceptors.request.push((config) => {
            config.headers['X-Request-ID'] = this.generateRequestId();
            return config;
        });
        
        // 响应拦截器：处理认证错误
        this.interceptors.response.push((response) => {
            if (response.status === 401) {
                this.handleAuthError();
            }
            return response;
        });
        
        // 响应拦截器：处理网络错误
        this.interceptors.response.push((response) => {
            if (!response.ok && response.status >= 500) {
                App.showError('服务器错误，请稍后重试');
            }
            return response;
        });
    }
    
    // 生成请求ID
    generateRequestId() {
        return 'req_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    }
    
    // 处理认证错误
    handleAuthError() {
        localStorage.removeItem('auth_token');
        App.showError('登录已过期，请重新登录');
        // 重定向到登录页面
        setTimeout(() => {
            window.location.href = '/login';
        }, 2000);
    }
    
    // 添加请求拦截器
    addRequestInterceptor(interceptor) {
        this.interceptors.request.push(interceptor);
    }
    
    // 添加响应拦截器
    addResponseInterceptor(interceptor) {
        this.interceptors.response.push(interceptor);
    }
    
    // 应用请求拦截器
    applyRequestInterceptors(config) {
        return this.interceptors.request.reduce((config, interceptor) => {
            return interceptor(config) || config;
        }, config);
    }
    
    // 应用响应拦截器
    applyResponseInterceptors(response) {
        return this.interceptors.response.reduce((response, interceptor) => {
            return interceptor(response) || response;
        }, response);
    }
    
    // 构建完整URL
    buildURL(endpoint) {
        if (endpoint.startsWith('http')) {
            return endpoint;
        }
        return `${this.baseURL}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;
    }
    
    // 通用请求方法
    async request(endpoint, options = {}) {
        const url = this.buildURL(endpoint);
        
        // 准备请求配置
        const config = {
            method: 'GET',
            headers: { ...this.defaultHeaders },
            ...options
        };
        
        // 合并自定义头部
        if (options.headers) {
            config.headers = { ...config.headers, ...options.headers };
        }
        
        // 应用请求拦截器
        const finalConfig = this.applyRequestInterceptors(config);
        
        try {
            App.log(`API请求: ${finalConfig.method} ${url}`);
            
            // 发送请求
            const response = await fetch(url, finalConfig);
            
            // 应用响应拦截器
            const finalResponse = this.applyResponseInterceptors(response);
            
            // 解析响应
            const data = await this.parseResponse(finalResponse);
            
            App.log(`API响应: ${finalResponse.status}`, data);
            
            return {
                data,
                status: finalResponse.status,
                headers: finalResponse.headers,
                ok: finalResponse.ok
            };
            
        } catch (error) {
            App.error('API请求失败:', error);
            
            // 网络错误处理
            if (error.name === 'TypeError' && error.message.includes('fetch')) {
                throw new Error('网络连接失败，请检查网络设置');
            }
            
            throw error;
        }
    }
    
    // 解析响应
    async parseResponse(response) {
        const contentType = response.headers.get('content-type');
        
        if (contentType && contentType.includes('application/json')) {
            return await response.json();
        }
        
        if (contentType && contentType.includes('text/')) {
            return await response.text();
        }
        
        return await response.blob();
    }
    
    // GET请求
    async get(endpoint, params = {}) {
        const url = new URL(this.buildURL(endpoint));
        
        // 添加查询参数
        Object.keys(params).forEach(key => {
            if (params[key] !== undefined && params[key] !== null) {
                url.searchParams.append(key, params[key]);
            }
        });
        
        return this.request(url.toString());
    }
    
    // POST请求
    async post(endpoint, data = {}) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }
    
    // PUT请求
    async put(endpoint, data = {}) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }
    
    // DELETE请求
    async delete(endpoint) {
        return this.request(endpoint, {
            method: 'DELETE'
        });
    }
    
    // PATCH请求
    async patch(endpoint, data = {}) {
        return this.request(endpoint, {
            method: 'PATCH',
            body: JSON.stringify(data)
        });
    }
    
    // 文件上传
    async upload(endpoint, file, onProgress = null) {
        const formData = new FormData();
        formData.append('file', file);
        
        const config = {
            method: 'POST',
            body: formData,
            headers: {} // 不设置Content-Type，让浏览器自动设置
        };
        
        // 如果需要进度回调
        if (onProgress && typeof onProgress === 'function') {
            return new Promise((resolve, reject) => {
                const xhr = new XMLHttpRequest();
                
                xhr.upload.addEventListener('progress', (event) => {
                    if (event.lengthComputable) {
                        const percentComplete = (event.loaded / event.total) * 100;
                        onProgress(percentComplete);
                    }
                });
                
                xhr.addEventListener('load', () => {
                    if (xhr.status >= 200 && xhr.status < 300) {
                        resolve(JSON.parse(xhr.responseText));
                    } else {
                        reject(new Error(`上传失败: ${xhr.statusText}`));
                    }
                });
                
                xhr.addEventListener('error', () => {
                    reject(new Error('上传失败'));
                });
                
                xhr.open('POST', this.buildURL(endpoint));
                
                // 添加认证头
                const token = localStorage.getItem('auth_token');
                if (token) {
                    xhr.setRequestHeader('Authorization', `Bearer ${token}`);
                }
                
                xhr.send(formData);
            });
        }
        
        return this.request(endpoint, config);
    }
    
    // 下载文件
    async download(endpoint, filename = null) {
        try {
            const response = await this.request(endpoint, {
                headers: {
                    'Accept': 'application/octet-stream'
                }
            });
            
            if (!response.ok) {
                throw new Error('下载失败');
            }
            
            // 创建下载链接
            const blob = response.data;
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = filename || 'download';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
            
        } catch (error) {
            App.error('下载失败:', error);
            throw error;
        }
    }
}

// 创建全局API客户端实例
window.apiClient = new APIClient();
```

### 学生API服务

**文件路径：** `static/js/api/student.js`

```javascript
/**
 * 学生API服务
 * 处理所有学生相关的API调用
 */

class StudentAPI {
    constructor(client) {
        this.client = client;
        this.endpoint = '/students';
    }
    
    // 获取所有学生
    async getAll(params = {}) {
        try {
            const response = await this.client.get(this.endpoint, params);
            
            if (!response.ok) {
                throw new Error(response.data.message || '获取学生列表失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('获取学生列表失败:', error);
            throw error;
        }
    }
    
    // 根据ID获取学生
    async getById(id) {
        try {
            const response = await this.client.get(`${this.endpoint}/${id}`);
            
            if (!response.ok) {
                throw new Error(response.data.message || '获取学生信息失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('获取学生信息失败:', error);
            throw error;
        }
    }
    
    // 创建学生
    async create(studentData) {
        try {
            // 数据验证
            this.validateStudentData(studentData);
            
            const response = await this.client.post(this.endpoint, studentData);
            
            if (!response.ok) {
                throw new Error(response.data.message || '创建学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('创建学生失败:', error);
            throw error;
        }
    }
    
    // 更新学生
    async update(id, studentData) {
        try {
            // 数据验证
            this.validateStudentData(studentData, false);
            
            const response = await this.client.put(`${this.endpoint}/${id}`, studentData);
            
            if (!response.ok) {
                throw new Error(response.data.message || '更新学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('更新学生失败:', error);
            throw error;
        }
    }
    
    // 删除学生
    async delete(id) {
        try {
            const response = await this.client.delete(`${this.endpoint}/${id}`);
            
            if (!response.ok) {
                throw new Error(response.data.message || '删除学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('删除学生失败:', error);
            throw error;
        }
    }
    
    // 搜索学生
    async search(query, filters = {}) {
        try {
            const params = {
                ...filters
            };
            
            // 添加搜索关键词
            if (query && query.trim()) {
                params.name = query.trim();
            }
            
            const response = await this.client.get(`${this.endpoint}/search`, params);
            
            if (!response.ok) {
                throw new Error(response.data.message || '搜索学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('搜索学生失败:', error);
            throw error;
        }
    }
    
    // 批量创建学生
    async batchCreate(studentsData) {
        try {
            // 验证批量数据
            if (!Array.isArray(studentsData) || studentsData.length === 0) {
                throw new Error('学生数据不能为空');
            }
            
            if (studentsData.length > 100) {
                throw new Error('批量创建数量不能超过100个');
            }
            
            // 验证每个学生数据
            studentsData.forEach((student, index) => {
                try {
                    this.validateStudentData(student);
                } catch (error) {
                    throw new Error(`第${index + 1}个学生数据验证失败: ${error.message}`);
                }
            });
            
            const response = await this.client.post(`${this.endpoint}/batch`, studentsData);
            
            if (!response.ok) {
                throw new Error(response.data.message || '批量创建学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('批量创建学生失败:', error);
            throw error;
        }
    }
    
    // 批量删除学生
    async batchDelete(ids) {
        try {
            if (!Array.isArray(ids) || ids.length === 0) {
                throw new Error('请选择要删除的学生');
            }
            
            if (ids.length > 100) {
                throw new Error('批量删除数量不能超过100个');
            }
            
            const response = await this.client.request(`${this.endpoint}/batch`, {
                method: 'DELETE',
                body: JSON.stringify(ids)
            });
            
            if (!response.ok) {
                throw new Error(response.data.message || '批量删除学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('批量删除学生失败:', error);
            throw error;
        }
    }
    
    // 激活学生
    async activate(id) {
        try {
            const response = await this.client.put(`${this.endpoint}/${id}/activate`);
            
            if (!response.ok) {
                throw new Error(response.data.message || '激活学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('激活学生失败:', error);
            throw error;
        }
    }
    
    // 禁用学生
    async deactivate(id) {
        try {
            const response = await this.client.put(`${this.endpoint}/${id}/deactivate`);
            
            if (!response.ok) {
                throw new Error(response.data.message || '禁用学生失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('禁用学生失败:', error);
            throw error;
        }
    }
    
    // 获取统计信息
    async getStatistics() {
        try {
            const response = await this.client.get(`${this.endpoint}/statistics`);
            
            if (!response.ok) {
                throw new Error(response.data.message || '获取统计信息失败');
            }
            
            return response.data;
        } catch (error) {
            App.error('获取统计信息失败:', error);
            throw error;
        }
    }
    
    // 导入学生数据
    async import(file, onProgress = null) {
        try {
            // 验证文件
            this.validateImportFile(file);
            
            const response = await this.client.upload(`${this.endpoint}/import`, file, onProgress);
            
            if (!response.success) {
                throw new Error(response.message || '导入失败');
            }
            
            return response;
        } catch (error) {
            App.error('导入学生数据失败:', error);
            throw error;
        }
    }
    
    // 导出学生数据
    async export(format = 'xlsx', filters = {}) {
        try {
            const params = {
                format,
                ...filters
            };
            
            const filename = `students_${new Date().toISOString().split('T')[0]}.${format}`;
            
            await this.client.download(`${this.endpoint}/export`, filename);
            
        } catch (error) {
            App.error('导出学生数据失败:', error);
            throw error;
        }
    }
    
    // 验证学生数据
    validateStudentData(data, isCreate = true) {
        const errors = [];
        
        // 必填字段验证（仅创建时）
        if (isCreate) {
            if (!data.name || !data.name.trim()) {
                errors.push('姓名不能为空');
            }
            
            if (!data.email || !data.email.trim()) {
                errors.push('邮箱不能为空');
            }
            
            if (!data.major || !data.major.trim()) {
                errors.push('专业不能为空');
            }
            
            if (!data.grade || !data.grade.trim()) {
                errors.push('年级不能为空');
            }
        }
        
        // 数据格式验证
        if (data.name && data.name.length > 50) {
            errors.push('姓名长度不能超过50个字符');
        }
        
        if (data.email && !this.isValidEmail(data.email)) {
            errors.push('邮箱格式不正确');
        }
        
        if (data.age && (data.age < 16 || data.age > 30)) {
            errors.push('年龄必须在16-30之间');
        }
        
        if (data.phone && !this.isValidPhone(data.phone)) {
            errors.push('手机号码格式不正确');
        }
        
        if (errors.length > 0) {
            throw new Error(errors.join(', '));
        }
    }
    
    // 验证导入文件
    validateImportFile(file) {
        if (!file) {
            throw new Error('请选择要导入的文件');
        }
        
        const allowedTypes = [
            'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
            'application/vnd.ms-excel',
            'text/csv'
        ];
        
        if (!allowedTypes.includes(file.type)) {
            throw new Error('文件格式不支持，请选择Excel或CSV文件');
        }
        
        const maxSize = 10 * 1024 * 1024; // 10MB
        if (file.size > maxSize) {
            throw new Error('文件大小不能超过10MB');
        }
    }
    
    // 验证邮箱格式
    isValidEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }
    
    // 验证手机号码格式
    isValidPhone(phone) {
        const phoneRegex = /^1[3-9]\d{9}$/;
        return phoneRegex.test(phone);
    }
}

// 创建全局学生API实例
window.studentAPI = new StudentAPI(window.apiClient);
```

## 8.3 学生管理组件

### 学生管理器

**文件路径：** `static/js/components/student.js`

```javascript
/**
 * 学生管理组件
 * 处理学生列表、添加、编辑、删除等功能
 */

class StudentManager {
    constructor() {
        this.students = [];
        this.filteredStudents = [];
        this.currentPage = 1;
        this.pageSize = 10;
        this.totalCount = 0;
        this.sortBy = 'created_at';
        this.sortOrder = 'desc';
        this.filters = {};
        this.selectedStudents = new Set();
        
        // DOM元素
        this.elements = {};
        
        // 事件处理器
        this.handlers = {};
        
        // 初始化
        this.init();
    }
    
    // 初始化
    async init() {
        try {
            App.log('StudentManager 初始化开始...');
            
            // 缓存DOM元素
            this.cacheElements();
            
            // 绑定事件
            this.bindEvents();
            
            // 加载初始数据
            await this.loadStudents();
            
            // 加载统计信息
            await this.loadStatistics();
            
            App.log('StudentManager 初始化完成');
            
        } catch (error) {
            App.error('StudentManager 初始化失败:', error);
        }
    }
    
    // 缓存DOM元素
    cacheElements() {
        this.elements = {
            // 搜索和过滤
            searchInput: document.getElementById('searchInput'),
            searchBtn: document.getElementById('searchBtn'),
            clearSearchBtn: document.getElementById('clearSearchBtn'),
            gradeFilter: document.getElementById('gradeFilter'),
            majorFilter: document.getElementById('majorFilter'),
            statusFilter: document.getElementById('statusFilter'),
            resetFilters: document.getElementById('resetFilters'),
            
            // 表格和列表
            studentTableBody: document.getElementById('studentTableBody'),
            studentCardContainer: document.getElementById('studentCardContainer'),
            loadingIndicator: document.getElementById('loadingIndicator'),
            emptyState: document.getElementById('emptyState'),
            studentCount: document.getElementById('studentCount'),
            
            // 分页
            pagination: document.getElementById('pagination'),
            pageStart: document.getElementById('pageStart'),
            pageEnd: document.getElementById('pageEnd'),
            totalCount: document.getElementById('totalCount'),
            
            // 视图切换
            tableView: document.getElementById('tableView'),
            cardView: document.getElementById('cardView'),
            tableViewContainer: document.getElementById('tableViewContainer'),
            cardViewContainer: document.getElementById('cardViewContainer'),
            
            // 批量操作
            selectAll: document.getElementById('selectAll'),
            batchToolbar: document.getElementById('batchToolbar'),
            selectedCount: document.getElementById('selectedCount'),
            batchActivateBtn: document.getElementById('batchActivateBtn'),
            batchDeactivateBtn: document.getElementById('batchDeactivateBtn'),
            batchDeleteBtn: document.getElementById('batchDeleteBtn'),
            cancelBatchBtn: document.getElementById('cancelBatchBtn'),
            
            // 模态框
            addStudentModal: document.getElementById('addStudentModal'),
            studentForm: document.getElementById('studentForm'),
            saveStudentBtn: document.getElementById('saveStudentBtn'),
            
            // 其他按钮
            refreshBtn: document.getElementById('refreshBtn'),
            importBtn: document.getElementById('importBtn'),
            exportBtn: document.getElementById('exportBtn')
        };
    }
    
    // 绑定事件
    bindEvents() {
        // 搜索事件
        if (this.elements.searchInput) {
            this.handlers.search = App.debounce(this.handleSearch.bind(this), 300);
            this.elements.searchInput.addEventListener('input', this.handlers.search);
            this.elements.searchInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    this.handleSearch();
                }
            });
        }
        
        if (this.elements.searchBtn) {
            this.elements.searchBtn.addEventListener('click', this.handleSearch.bind(this));
        }
        
        if (this.elements.clearSearchBtn) {
            this.elements.clearSearchBtn.addEventListener('click', this.clearSearch.bind(this));
        }
        
        // 过滤器事件
        [this.elements.gradeFilter, this.elements.majorFilter, this.elements.statusFilter].forEach(filter => {
            if (filter) {
                filter.addEventListener('change', this.handleFilterChange.bind(this));
            }
        });
        
        if (this.elements.resetFilters) {
            this.elements.resetFilters.addEventListener('click', this.resetFilters.bind(this));
        }
        
        // 视图切换
        if (this.elements.tableView) {
            this.elements.tableView.addEventListener('change', () => {
                this.switchView('table');
            });
        }
        
        if (this.elements.cardView) {
            this.elements.cardView.addEventListener('change', () => {
                this.switchView('card');
            });
        }
        
        // 全选
        if (this.elements.selectAll) {
            this.elements.selectAll.addEventListener('change', this.handleSelectAll.bind(this));
        }
        
        // 批量操作
        if (this.elements.batchActivateBtn) {
            this.elements.batchActivateBtn.addEventListener('click', () => {
                this.handleBatchOperation('activate');
            });
        }
        
        if (this.elements.batchDeactivateBtn) {
            this.elements.batchDeactivateBtn.addEventListener('click', () => {
                this.handleBatchOperation('deactivate');
            });
        }
        
        if (this.elements.batchDeleteBtn) {
            this.elements.batchDeleteBtn.addEventListener('click', () => {
                this.handleBatchOperation('delete');
            });
        }
        
        if (this.elements.cancelBatchBtn) {
            this.elements.cancelBatchBtn.addEventListener('click', this.cancelBatchSelection.bind(this));
        }
        
        // 表单提交
        if (this.elements.saveStudentBtn) {
            this.elements.saveStudentBtn.addEventListener('click', this.handleSaveStudent.bind(this));
        }
        
        // 其他按钮
        if (this.elements.refreshBtn) {
            this.elements.refreshBtn.addEventListener('click', this.refresh.bind(this));
        }
        
        if (this.elements.importBtn) {
            this.elements.importBtn.addEventListener('click', this.showImportModal.bind(this));
        }
        
        if (this.elements.exportBtn) {
            this.elements.exportBtn.addEventListener('click', this.handleExport.bind(this));
        }
        
        // 排序事件
        document.addEventListener('click', (e) => {
            if (e.target.closest('.sort-link')) {
                e.preventDefault();
                const sortBy = e.target.closest('.sort-link').dataset.sort;
                this.handleSort(sortBy);
            }
        });
        
        // 模态框事件
        if (this.elements.addStudentModal) {
            this.elements.addStudentModal.addEventListener('hidden.bs.modal', () => {
                this.resetForm();
            });
        }
    }
    
    // 加载学生列表
    async loadStudents(showLoading = true) {
        try {
            if (showLoading) {
                this.showLoading();
            }
            
            const params = {
                page: this.currentPage,
                page_size: this.pageSize,
                sort_by: this.sortBy,
                sort_order: this.sortOrder,
                ...this.filters
            };
            
            const response = await studentAPI.getAll(params);
            
            if (response.success) {
                this.students = response.data.students || [];
                this.totalCount = response.data.total || 0;
                
                this.renderStudents();
                this.renderPagination();
                this.updateStatistics();
            } else {
                throw new Error(response.message || '加载学生列表失败');
            }
            
        } catch (error) {
            App.error('加载学生列表失败:', error);
            App.showError(error.message);
            this.showEmptyState();
        } finally {
            this.hideLoading();
        }
    }
    
    // 渲染学生列表
    renderStudents() {
        if (this.students.length === 0) {
            this.showEmptyState();
            return;
        }
        
        this.hideEmptyState();
        
        // 根据当前视图模式渲染
        if (this.elements.tableView && this.elements.tableView.checked) {
            this.renderTableView();
        } else {
            this.renderCardView();
        }
        
        // 更新计数
        if (this.elements.studentCount) {
            this.elements.studentCount.textContent = this.totalCount;
        }
    }
    
    // 渲染表格视图
    renderTableView() {
        if (!this.elements.studentTableBody) return;
        
        const tbody = this.elements.studentTableBody;
        tbody.innerHTML = '';
        
        this.students.forEach(student => {
            const row = this.createTableRow(student);
            tbody.appendChild(row);
        });
    }
    
    // 创建表格行
    createTableRow(student) {
        const row = document.createElement('tr');
        row.dataset.studentId = student.id;
        
        if (this.selectedStudents.has(student.id)) {
            row.classList.add('selected');
        }
        
        row.innerHTML = `
            <td>
                <input type="checkbox" class="form-check-input student-checkbox" 
                       value="${student.id}" ${this.selectedStudents.has(student.id) ? 'checked' : ''}>
            </td>
            <td>
                <img src="${student.avatar || '/static/images/default-avatar.png'}" 
                     alt="${student.name}" class="avatar rounded-circle">
            </td>
            <td>
                <div class="fw-bold">${this.escapeHtml(student.name)}</div>
                <small class="text-muted">${student.student_id || ''}</small>
            </td>
            <td>${student.age || '-'}</td>
            <td>
                <span class="badge bg-${student.gender === '男' ? 'primary' : 'pink'}">
                    ${student.gender || '-'}
                </span>
            </td>
            <td>
                <a href="mailto:${student.email}" class="text-decoration-none">
                    ${this.escapeHtml(student.email)}
                </a>
            </td>
            <td>${this.escapeHtml(student.major || '-')}</td>
            <td>
                <span class="badge bg-info">${student.grade || '-'}</span>
            </td>
            <td>
                <span class="status-indicator ${student.status === 1 ? 'active' : 'inactive'}"></span>
                <span class="badge status-badge ${student.status === 1 ? 'status-active' : 'status-inactive'}">
                    ${student.status === 1 ? '正常' : '禁用'}
                </span>
            </td>
            <td>
                <small class="text-muted">
                    ${App.formatDate(student.created_at, 'YYYY-MM-DD')}
                </small>
            </td>
            <td>
                <div class="action-buttons">
                    <button class="btn btn-sm btn-outline-primary" 
                            onclick="StudentManager.showStudentDetail(${student.id})" 
                            title="查看详情">
                        <i class="fas fa-eye"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-secondary" 
                            onclick="StudentManager.showEditModal(${student.id})" 
                            title="编辑">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-${student.status === 1 ? 'warning' : 'success'}" 
                            onclick="StudentManager.toggleStatus(${student.id})" 
                            title="${student.status === 1 ? '禁用' : '激活'}">
                        <i class="fas fa-${student.status === 1 ? 'ban' : 'check'}"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger" 
                            onclick="StudentManager.deleteStudent(${student.id})" 
                            title="删除">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </td>
        `;
        
        // 绑定复选框事件
        const checkbox = row.querySelector('.student-checkbox');
        checkbox.addEventListener('change', (e) => {
            this.handleStudentSelect(student.id, e.target.checked);
        });
        
        return row;
    }
    
    // 渲染卡片视图
    renderCardView() {
        if (!this.elements.studentCardContainer) return;
        
        const container = this.elements.studentCardContainer;
        container.innerHTML = '';
        
        this.students.forEach(student => {
            const card = this.createStudentCard(student);
            container.appendChild(card);
        });
    }
    
    // 创建学生卡片
    createStudentCard(student) {
        const col = document.createElement('div');
        col.className = 'col-md-6 col-lg-4';
        
        col.innerHTML = `
            <div class="card student-card h-100" data-student-id="${student.id}">
                <div class="card-body text-center">
                    <div class="position-relative d-inline-block mb-3">
                        <img src="${student.avatar || '/static/images/default-avatar.png'}" 
                             alt="${student.name}" class="student-avatar rounded-circle">
                        <span class="position-absolute bottom-0 end-0 status-indicator ${student.status === 1 ? 'active' : 'inactive'}"></span>
                    </div>
                    
                    <h6 class="card-title mb-2">${this.escapeHtml(student.name)}</h6>
                    
                    <div class="student-info text-start">
                        <p><i class="fas fa-envelope me-2"></i>${this.escapeHtml(student.email)}</p>
                        <p><i class="fas fa-graduation-cap me-2"></i>${this.escapeHtml(student.major || '-')}</p>
                        <p><i class="fas fa-calendar me-2"></i>${student.grade || '-'}级</p>
                        <p><i class="fas fa-birthday-cake me-2"></i>${student.age || '-'}岁</p>
                    </div>
                    
                    <div class="d-flex justify-content-between align-items-center mt-3">
                        <div class="form-check">
                            <input class="form-check-input student-checkbox" type="checkbox" 
                                   value="${student.id}" ${this.selectedStudents.has(student.id) ? 'checked' : ''}>
                            <label class="form-check-label">选择</label>
                        </div>
                        
                        <div class="btn-group btn-group-sm">
                            <button class="btn btn-outline-primary" 
                                    onclick="StudentManager.showStudentDetail(${student.id})" 
                                    title="查看详情">
                                <i class="fas fa-eye"></i>
                            </button>
                            <button class="btn btn-outline-secondary" 
                                    onclick="StudentManager.showEditModal(${student.id})" 
                                    title="编辑">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-outline-danger" 
                                    onclick="StudentManager.deleteStudent(${student.id})" 
                                    title="删除">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        // 绑定复选框事件
        const checkbox = col.querySelector('.student-checkbox');
        checkbox.addEventListener('change', (e) => {
            this.handleStudentSelect(student.id, e.target.checked);
        });
        
        // 绑定卡片点击事件
        const card = col.querySelector('.student-card');
        card.addEventListener('click', (e) => {
            if (!e.target.closest('button') && !e.target.closest('.form-check')) {
                this.showStudentDetail(student.id);
            }
        });
        
        return col;
    }
    
    // 处理搜索
    async handleSearch() {
        const query = this.elements.searchInput?.value?.trim() || '';
        
        // 重置到第一页
        this.currentPage = 1;
        
        // 如果有搜索词，使用搜索API
        if (query) {
            try {
                this.showLoading();
                
                const response = await studentAPI.search(query, this.filters);
                
                if (response.success) {
                    this.students = response.data || [];
                    this.totalCount = this.students.length;
                    
                    this.renderStudents();
                    this.renderPagination();
                } else {
                    throw new Error(response.message || '搜索失败');
                }
                
            } catch (error) {
                App.error('搜索失败:', error);
                App.showError(error.message);
            } finally {
                this.hideLoading();
            }
        } else {
            // 没有搜索词，加载所有数据
            await this.loadStudents();
        }
    }
    
    // 清除搜索
    clearSearch() {
        if (this.elements.searchInput) {
            this.elements.searchInput.value = '';
        }
        this.handleSearch();
    }
    
    // 处理过滤器变化
    handleFilterChange() {
        this.filters = {};
        
        // 收集过滤条件
        if (this.elements.gradeFilter?.value) {
            this.filters.grade = this.elements.gradeFilter.value;
        }
        
        if (this.elements.majorFilter?.value) {
            this.filters.major = this.elements.majorFilter.value;
        }
        
        if (this.elements.statusFilter?.value !== '') {
            this.filters.status = this.elements.statusFilter.value;
        }
        
        // 重置到第一页
        this.currentPage = 1;
        
        // 重新加载数据
        this.loadStudents();
    }
    
    // 重置过滤器
    resetFilters() {
        this.filters = {};
        
        // 重置表单
        if (this.elements.gradeFilter) this.elements.gradeFilter.value = '';
        if (this.elements.majorFilter) this.elements.majorFilter.value = '';
        if (this.elements.statusFilter) this.elements.statusFilter.value = '';
        
        // 重置到第一页
        this.currentPage = 1;
        
        // 重新加载数据
        this.loadStudents();
    }
    
    // 处理排序
    handleSort(sortBy) {
        if (this.sortBy === sortBy) {
            // 切换排序顺序
            this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
        } else {
            // 新的排序字段
            this.sortBy = sortBy;
            this.sortOrder = 'asc';
        }
        
        // 更新排序图标
        this.updateSortIcons();
        
        // 重新加载数据
        this.loadStudents();
    }
    
    // 更新排序图标
    updateSortIcons() {
        document.querySelectorAll('.sort-link').forEach(link => {
            const icon = link.querySelector('i');
            const sortBy = link.dataset.sort;
            
            link.classList.remove('active');
            icon.className = 'fas fa-sort';
            
            if (sortBy === this.sortBy) {
                link.classList.add('active');
                icon.className = `fas fa-sort-${this.sortOrder === 'asc' ? 'up' : 'down'}`;
            }
        });
    }
    
    // 切换视图模式
    switchView(mode) {
        if (mode === 'table') {
            this.elements.tableViewContainer?.classList.remove('d-none');
            this.elements.cardViewContainer?.classList.add('d-none');
        } else {
            this.elements.tableViewContainer?.classList.add('d-none');
            this.elements.cardViewContainer?.classList.remove('d-none');
        }
        
        // 重新渲染
        this.renderStudents();
    }
    
    // 显示加载状态
    showLoading() {
        if (this.elements.loadingIndicator) {
            this.elements.loadingIndicator.classList.remove('d-none');
        }
        
        if (this.elements.studentTableBody) {
            this.elements.studentTableBody.style.opacity = '0.5';
        }
        
        if (this.elements.studentCardContainer) {
            this.elements.studentCardContainer.style.opacity = '0.5';
        }
    }
    
    // 隐藏加载状态
    hideLoading() {
        if (this.elements.loadingIndicator) {
            this.elements.loadingIndicator.classList.add('d-none');
        }
        
        if (this.elements.studentTableBody) {
            this.elements.studentTableBody.style.opacity = '1';
        }
        
        if (this.elements.studentCardContainer) {
            this.elements.studentCardContainer.style.opacity = '1';
        }
    }
    
    // 显示空状态
    showEmptyState() {
        if (this.elements.emptyState) {
            this.elements.emptyState.classList.remove('d-none');
        }
        
        if (this.elements.tableViewContainer) {
            this.elements.tableViewContainer.classList.add('d-none');
        }
        
        if (this.elements.cardViewContainer) {
            this.elements.cardViewContainer.classList.add('d-none');
        }
    }
    
    // 隐藏空状态
    hideEmptyState() {
        if (this.elements.emptyState) {
            this.elements.emptyState.classList.add('d-none');
        }
        
        if (this.elements.tableViewContainer) {
            this.elements.tableViewContainer.classList.remove('d-none');
        }
        
        if (this.elements.cardViewContainer) {
            this.elements.cardViewContainer.classList.remove('d-none');
        }
    }
    
    // HTML转义
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    
    // 刷新数据
    async refresh() {
        this.currentPage = 1;
        await this.loadStudents();
        App.showSuccess('数据已刷新');
    }
}

// 创建全局实例
window.StudentManager = new StudentManager();
```

## 8.4 通知服务

### 通知系统

**文件路径：** `static/js/services/notification.js`

```javascript
/**
 * 通知服务
 * 统一处理应用内的通知消息
 */

class NotificationService {
    constructor() {
        this.container = null;
        this.notifications = new Map();
        this.defaultOptions = {
            duration: 5000,
            position: 'top-right',
            showProgress: true,
            allowClose: true,
            pauseOnHover: true
        };
        
        this.init();
    }
    
    // 初始化通知系统
    init() {
        this.createContainer();
        this.bindEvents();
    }
    
    // 创建通知容器
    createContainer() {
        this.container = document.createElement('div');
        this.container.id = 'notification-container';
        this.container.className = 'notification-container';
        document.body.appendChild(this.container);
    }
    
    // 绑定事件
    bindEvents() {
        // 监听页面可见性变化
        document.addEventListener('visibilitychange', () => {
            if (document.visibilityState === 'hidden') {
                this.pauseAll();
            } else {
                this.resumeAll();
            }
        });
    }
    
    // 显示通知
    show(message, type = 'info', options = {}) {
        const config = { ...this.defaultOptions, ...options };
        const id = this.generateId();
        
        const notification = this.createNotification(id, message, type, config);
        this.container.appendChild(notification);
        
        // 添加到管理列表
        this.notifications.set(id, {
            element: notification,
            config,
            timer: null,
            paused: false
        });
        
        // 设置自动关闭
        if (config.duration > 0) {
            this.setAutoClose(id, config.duration);
        }
        
        // 触发显示动画
        requestAnimationFrame(() => {
            notification.classList.add('show');
        });
        
        return id;
    }
    
    // 创建通知元素
    createNotification(id, message, type, config) {
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.dataset.id = id;
        
        const icon = this.getTypeIcon(type);
        const progressBar = config.showProgress ? '<div class="notification-progress"></div>' : '';
        const closeButton = config.allowClose ? '<button class="notification-close" aria-label="关闭"><i class="fas fa-times"></i></button>' : '';
        
        notification.innerHTML = `
            <div class="notification-content">
                <div class="notification-icon">
                    <i class="${icon}"></i>
                </div>
                <div class="notification-message">${message}</div>
                ${closeButton}
            </div>
            ${progressBar}
        `;
        
        // 绑定事件
        this.bindNotificationEvents(notification, id, config);
        
        return notification;
    }
    
    // 获取类型图标
    getTypeIcon(type) {
        const icons = {
            success: 'fas fa-check-circle',
            error: 'fas fa-exclamation-circle',
            warning: 'fas fa-exclamation-triangle',
            info: 'fas fa-info-circle'
        };
        return icons[type] || icons.info;
    }
    
    // 绑定通知事件
    bindNotificationEvents(notification, id, config) {
        // 关闭按钮
        const closeBtn = notification.querySelector('.notification-close');
        if (closeBtn) {
            closeBtn.addEventListener('click', () => {
                this.close(id);
            });
        }
        
        // 鼠标悬停暂停
        if (config.pauseOnHover) {
            notification.addEventListener('mouseenter', () => {
                this.pause(id);
            });
            
            notification.addEventListener('mouseleave', () => {
                this.resume(id);
            });
        }
        
        // 点击关闭
        notification.addEventListener('click', (e) => {
            if (!e.target.closest('.notification-close')) {
                this.close(id);
            }
        });
    }
    
    // 设置自动关闭
    setAutoClose(id, duration) {
        const notificationData = this.notifications.get(id);
        if (!notificationData) return;
        
        const progressBar = notificationData.element.querySelector('.notification-progress');
        
        if (progressBar) {
            progressBar.style.animationDuration = `${duration}ms`;
            progressBar.classList.add('running');
        }
        
        notificationData.timer = setTimeout(() => {
            this.close(id);
        }, duration);
    }
    
    // 暂停通知
    pause(id) {
        const notificationData = this.notifications.get(id);
        if (!notificationData || notificationData.paused) return;
        
        notificationData.paused = true;
        
        if (notificationData.timer) {
            clearTimeout(notificationData.timer);
        }
        
        const progressBar = notificationData.element.querySelector('.notification-progress');
        if (progressBar) {
            progressBar.style.animationPlayState = 'paused';
        }
    }
    
    // 恢复通知
    resume(id) {
        const notificationData = this.notifications.get(id);
        if (!notificationData || !notificationData.paused) return;
        
        notificationData.paused = false;
        
        const progressBar = notificationData.element.querySelector('.notification-progress');
        if (progressBar) {
            progressBar.style.animationPlayState = 'running';
            
            // 计算剩余时间
            const computedStyle = window.getComputedStyle(progressBar);
            const animationDuration = parseFloat(computedStyle.animationDuration) * 1000;
            const animationDelay = parseFloat(computedStyle.animationDelay) * 1000;
            const elapsedTime = Date.now() - (notificationData.startTime || Date.now());
            const remainingTime = Math.max(0, animationDuration - elapsedTime);
            
            if (remainingTime > 0) {
                notificationData.timer = setTimeout(() => {
                    this.close(id);
                }, remainingTime);
            }
        }
    }
    
    // 关闭通知
    close(id) {
        const notificationData = this.notifications.get(id);
        if (!notificationData) return;
        
        const { element, timer } = notificationData;
        
        // 清除定时器
        if (timer) {
            clearTimeout(timer);
        }
        
        // 添加关闭动画
        element.classList.add('closing');
        
        // 动画结束后移除元素
        setTimeout(() => {
            if (element.parentNode) {
                element.parentNode.removeChild(element);
            }
            this.notifications.delete(id);
        }, 300);
    }
    
    // 暂停所有通知
    pauseAll() {
        this.notifications.forEach((_, id) => {
            this.pause(id);
        });
    }
    
    // 恢复所有通知
    resumeAll() {
        this.notifications.forEach((_, id) => {
            this.resume(id);
        });
    }
    
    // 清除所有通知
    clearAll() {
        this.notifications.forEach((_, id) => {
            this.close(id);
        });
    }
    
    // 生成唯一ID
    generateId() {
        return 'notification_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    }
    
    // 便捷方法
    success(message, options = {}) {
        return this.show(message, 'success', options);
    }
    
    error(message, options = {}) {
        return this.show(message, 'error', { duration: 8000, ...options });
    }
    
    warning(message, options = {}) {
        return this.show(message, 'warning', { duration: 6000, ...options });
    }
    
    info(message, options = {}) {
        return this.show(message, 'info', options);
    }
}

// 创建全局实例
window.NotificationService = new NotificationService();
```

## 8.5 工具函数

### 验证工具

**文件路径：** `static/js/utils/validation.js`

```javascript
/**
 * 验证工具函数
 * 提供各种数据验证功能
 */

class ValidationUtils {
    // 验证邮箱
    static isValidEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }
    
    // 验证手机号
    static isValidPhone(phone) {
        const phoneRegex = /^1[3-9]\d{9}$/;
        return phoneRegex.test(phone);
    }
    
    // 验证身份证号
    static isValidIdCard(idCard) {
        const idCardRegex = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/;
        return idCardRegex.test(idCard);
    }
    
    // 验证学号
    static isValidStudentId(studentId) {
        const studentIdRegex = /^\d{8,12}$/;
        return studentIdRegex.test(studentId);
    }
    
    // 验证密码强度
    static validatePassword(password) {
        const result = {
            isValid: false,
            score: 0,
            feedback: []
        };
        
        if (!password) {
            result.feedback.push('密码不能为空');
            return result;
        }
        
        if (password.length < 8) {
            result.feedback.push('密码长度至少8位');
        } else {
            result.score += 1;
        }
        
        if (!/[a-z]/.test(password)) {
            result.feedback.push('密码需包含小写字母');
        } else {
            result.score += 1;
        }
        
        if (!/[A-Z]/.test(password)) {
            result.feedback.push('密码需包含大写字母');
        } else {
            result.score += 1;
        }
        
        if (!/\d/.test(password)) {
            result.feedback.push('密码需包含数字');
        } else {
            result.score += 1;
        }
        
        if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
            result.feedback.push('密码需包含特殊字符');
        } else {
            result.score += 1;
        }
        
        result.isValid = result.score >= 3 && result.feedback.length === 0;
        
        return result;
    }
    
    // 验证年龄
    static isValidAge(age) {
        return Number.isInteger(age) && age >= 16 && age <= 30;
    }
    
    // 验证必填字段
    static isRequired(value) {
        if (typeof value === 'string') {
            return value.trim().length > 0;
        }
        return value !== null && value !== undefined;
    }
    
    // 验证字符串长度
    static isValidLength(str, min = 0, max = Infinity) {
        if (typeof str !== 'string') return false;
        const length = str.trim().length;
        return length >= min && length <= max;
    }
    
    // 验证数字范围
    static isInRange(num, min = -Infinity, max = Infinity) {
        return typeof num === 'number' && num >= min && num <= max;
    }
    
    // 验证文件类型
    static isValidFileType(file, allowedTypes) {
        if (!file || !allowedTypes) return false;
        return allowedTypes.includes(file.type);
    }
    
    // 验证文件大小
    static isValidFileSize(file, maxSize) {
        if (!file) return false;
        return file.size <= maxSize;
    }
    
    // 表单验证
    static validateForm(formData, rules) {
        const errors = {};
        
        Object.keys(rules).forEach(field => {
            const value = formData[field];
            const fieldRules = rules[field];
            const fieldErrors = [];
            
            fieldRules.forEach(rule => {
                if (rule.required && !this.isRequired(value)) {
                    fieldErrors.push(rule.message || `${field}是必填项`);
                    return;
                }
                
                if (value && rule.type) {
                    switch (rule.type) {
                        case 'email':
                            if (!this.isValidEmail(value)) {
                                fieldErrors.push(rule.message || '邮箱格式不正确');
                            }
                            break;
                        case 'phone':
                            if (!this.isValidPhone(value)) {
                                fieldErrors.push(rule.message || '手机号格式不正确');
                            }
                            break;
                        case 'age':
                            if (!this.isValidAge(parseInt(value))) {
                                fieldErrors.push(rule.message || '年龄必须在16-30之间');
                            }
                            break;
                    }
                }
                
                if (value && rule.minLength && !this.isValidLength(value, rule.minLength)) {
                    fieldErrors.push(rule.message || `${field}长度不能少于${rule.minLength}个字符`);
                }
                
                if (value && rule.maxLength && !this.isValidLength(value, 0, rule.maxLength)) {
                    fieldErrors.push(rule.message || `${field}长度不能超过${rule.maxLength}个字符`);
                }
                
                if (value && rule.pattern && !rule.pattern.test(value)) {
                    fieldErrors.push(rule.message || `${field}格式不正确`);
                }
            });
            
            if (fieldErrors.length > 0) {
                errors[field] = fieldErrors;
            }
        });
        
        return {
            isValid: Object.keys(errors).length === 0,
            errors
        };
    }
}

// 全局暴露
window.ValidationUtils = ValidationUtils;
```

## 8.6 性能优化

### 缓存策略

```javascript
// 实现智能缓存
class CacheManager {
    constructor() {
        this.cache = new Map();
        this.maxSize = 100;
        this.ttl = 5 * 60 * 1000; // 5分钟
    }
    
    set(key, value, customTTL = null) {
        const expiry = Date.now() + (customTTL || this.ttl);
        
        if (this.cache.size >= this.maxSize) {
            this.evictOldest();
        }
        
        this.cache.set(key, {
            value,
            expiry,
            accessCount: 0,
            lastAccess: Date.now()
        });
    }
    
    get(key) {
        const item = this.cache.get(key);
        
        if (!item) return null;
        
        if (Date.now() > item.expiry) {
            this.cache.delete(key);
            return null;
        }
        
        item.accessCount++;
        item.lastAccess = Date.now();
        
        return item.value;
    }
    
    evictOldest() {
        let oldestKey = null;
        let oldestTime = Infinity;
        
        for (const [key, item] of this.cache) {
            if (item.lastAccess < oldestTime) {
                oldestTime = item.lastAccess;
                oldestKey = key;
            }
        }
        
        if (oldestKey) {
            this.cache.delete(oldestKey);
        }
    }
}
```

### 防抖和节流

```javascript
// 高级防抖实现
function advancedDebounce(func, wait, options = {}) {
    let timeout;
    let lastArgs;
    let lastThis;
    let result;
    let lastCallTime;
    let lastInvokeTime = 0;
    
    const {
        leading = false,
        trailing = true,
        maxWait
    } = options;
    
    function invokeFunc(time) {
        const args = lastArgs;
        const thisArg = lastThis;
        
        lastArgs = lastThis = undefined;
        lastInvokeTime = time;
        result = func.apply(thisArg, args);
        return result;
    }
    
    function shouldInvoke(time) {
        const timeSinceLastCall = time - lastCallTime;
        const timeSinceLastInvoke = time - lastInvokeTime;
        
        return (lastCallTime === undefined ||
                timeSinceLastCall >= wait ||
                timeSinceLastCall < 0 ||
                (maxWait !== undefined && timeSinceLastInvoke >= maxWait));
    }
    
    function debounced(...args) {
        const time = Date.now();
        const isInvoking = shouldInvoke(time);
        
        lastArgs = args;
        lastThis = this;
        lastCallTime = time;
        
        if (isInvoking) {
            if (timeout === undefined) {
                return leadingEdge(lastCallTime);
            }
            if (maxWait) {
                timeout = setTimeout(timerExpired, wait);
                return invokeFunc(lastCallTime);
            }
        }
        
        if (timeout === undefined) {
            timeout = setTimeout(timerExpired, wait);
        }
        
        return result;
    }
    
    function leadingEdge(time) {
        lastInvokeTime = time;
        timeout = setTimeout(timerExpired, wait);
        return leading ? invokeFunc(time) : result;
    }
    
    function timerExpired() {
        const time = Date.now();
        if (shouldInvoke(time)) {
            return trailingEdge(time);
        }
        timeout = setTimeout(timerExpired, remainingWait(time));
    }
    
    function trailingEdge(time) {
        timeout = undefined;
        
        if (trailing && lastArgs) {
            return invokeFunc(time);
        }
        lastArgs = lastThis = undefined;
        return result;
    }
    
    function remainingWait(time) {
        const timeSinceLastCall = time - lastCallTime;
        const timeSinceLastInvoke = time - lastInvokeTime;
        const timeWaiting = wait - timeSinceLastCall;
        
        return maxWait !== undefined
            ? Math.min(timeWaiting, maxWait - timeSinceLastInvoke)
            : timeWaiting;
    }
    
    debounced.cancel = function() {
        if (timeout !== undefined) {
            clearTimeout(timeout);
        }
        lastInvokeTime = 0;
        lastArgs = lastCallTime = lastThis = timeout = undefined;
    };
    
    debounced.flush = function() {
        return timeout === undefined ? result : trailingEdge(Date.now());
    };
    
    return debounced;
}
```

## 8.7 错误处理和调试

### 全局错误处理

```javascript
// 统一错误处理
class ErrorHandler {
    static init() {
        // 捕获未处理的错误
        window.addEventListener('error', this.handleError.bind(this));
        
        // 捕获Promise拒绝
        window.addEventListener('unhandledrejection', this.handlePromiseRejection.bind(this));
        
        // 捕获资源加载错误
        window.addEventListener('error', this.handleResourceError.bind(this), true);
    }
    
    static handleError(event) {
        const error = {
            message: event.message,
            filename: event.filename,
            lineno: event.lineno,
            colno: event.colno,
            stack: event.error?.stack,
            timestamp: new Date().toISOString(),
            userAgent: navigator.userAgent,
            url: window.location.href
        };
        
        this.logError(error);
        this.reportError(error);
    }
    
    static handlePromiseRejection(event) {
        const error = {
            type: 'unhandledrejection',
            reason: event.reason,
            timestamp: new Date().toISOString(),
            userAgent: navigator.userAgent,
            url: window.location.href
        };
        
        this.logError(error);
        this.reportError(error);
        
        // 阻止默认的控制台错误输出
        event.preventDefault();
    }
    
    static handleResourceError(event) {
        if (event.target !== window) {
            const error = {
                type: 'resource',
                element: event.target.tagName,
                source: event.target.src || event.target.href,
                timestamp: new Date().toISOString(),
                url: window.location.href
            };
            
            this.logError(error);
        }
    }
    
    static logError(error) {
        console.error('应用错误:', error);
        
        // 存储到本地存储用于调试
        const errors = JSON.parse(localStorage.getItem('app_errors') || '[]');
        errors.push(error);
        
        // 只保留最近50个错误
        if (errors.length > 50) {
            errors.splice(0, errors.length - 50);
        }
        
        localStorage.setItem('app_errors', JSON.stringify(errors));
    }
    
    static reportError(error) {
        // 发送错误报告到服务器
        if (window.AnalyticsService) {
            AnalyticsService.trackError(error);
        }
    }
    
    static getStoredErrors() {
        return JSON.parse(localStorage.getItem('app_errors') || '[]');
    }
    
    static clearStoredErrors() {
        localStorage.removeItem('app_errors');
    }
}

// 初始化错误处理
ErrorHandler.init();
```

## 8.8 总结

本章详细介绍了JavaScript交互逻辑的开发，包括：

### 核心特性

1. **模块化架构** - 清晰的代码组织结构
2. **API客户端** - 统一的HTTP请求处理
3. **组件系统** - 可复用的UI组件
4. **事件系统** - 灵活的事件处理机制
5. **状态管理** - 应用状态的统一管理
6. **错误处理** - 完善的错误捕获和处理
7. **性能优化** - 缓存、防抖、节流等优化技术
8. **用户体验** - 加载状态、通知系统等

### 最佳实践

1. **代码规范** - 统一的编码风格和命名规范
2. **错误处理** - 完善的错误捕获和用户友好的错误提示
3. **性能优化** - 合理使用缓存和优化技术
4. **用户体验** - 响应式设计和交互反馈
5. **可维护性** - 模块化设计和清晰的代码结构
6. **可扩展性** - 灵活的架构设计

### 下一步

在下一章中，我们将学习如何进行系统测试和部署，确保应用的质量和稳定性。
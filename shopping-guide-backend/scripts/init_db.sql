-- 创建数据库
CREATE DATABASE IF NOT EXISTS shopping_guide DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE shopping_guide;

-- 用户画像表
CREATE TABLE `user_profiles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` varchar(64) NOT NULL COMMENT '用户唯一标识（业务ID，如UUID/手机号等）',
  `preferred_style` varchar(50) NOT NULL COMMENT '偏好风格（如：xiaohongshu/dongyuhui）',
  `age` int(11) NOT NULL COMMENT '年龄',
  `gender` varchar(20) NOT NULL COMMENT '性别（如：male/female/other）',
  `interests` json DEFAULT NULL COMMENT '兴趣爱好列表',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_user_id` (`user_id`) COMMENT '用户ID唯一索引，加速查询并防止重复'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户画像表';
-- 用户表
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(64) PRIMARY KEY,
    username VARCHAR(128),
    email VARCHAR(255),
    phone VARCHAR(32),
    profile JSON COMMENT '用户画像',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 商品表
CREATE TABLE IF NOT EXISTS products (
    product_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(64) NOT NULL,
    sub_category VARCHAR(64),
    price DECIMAL(10,2) NOT NULL,
    stock INT DEFAULT 0,
    description TEXT,
    images JSON COMMENT '商品图片数组',
    attributes JSON COMMENT '商品属性',
    status TINYINT DEFAULT 1 COMMENT '1:上架 0:下架',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category, sub_category),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- 对话日志表
CREATE TABLE IF NOT EXISTS chat_logs (
    log_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    query TEXT NOT NULL COMMENT '用户输入',
    response TEXT NOT NULL COMMENT 'AI回复',
    tool_used VARCHAR(64) COMMENT '使用的工具',
    planner_result JSON COMMENT 'Planner返回结果',
    executor_result JSON COMMENT 'Executor返回结果',
    recommended_products JSON COMMENT '推荐的商品',
    latency_ms INT COMMENT '响应时长(毫秒)',
    tokens_used INT COMMENT 'Token消耗',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_session (session_id),
    INDEX idx_user (user_id),
    INDEX idx_tool (tool_used),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对话日志表';

-- Dify调用日志表
CREATE TABLE IF NOT EXISTS dify_call_logs (
    log_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    workflow_name VARCHAR(64) NOT NULL COMMENT '工作流名称',
    app_id VARCHAR(64) NOT NULL,
    workflow_run_id VARCHAR(128) COMMENT 'Dify运行ID',
    inputs JSON COMMENT '输入参数',
    outputs JSON COMMENT '输出结果',
    status VARCHAR(32) NOT NULL COMMENT '状态: success/error/timeout',
    error_message TEXT COMMENT '错误信息',
    latency_ms INT COMMENT '调用时长',
    tokens_used INT COMMENT 'Token消耗',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_workflow (workflow_name),
    INDEX idx_app (app_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Dify调用日志表';

-- 商品推荐记录表
CREATE TABLE IF NOT EXISTS product_recommendations (
    rec_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    product_id VARCHAR(64) NOT NULL,
    reason TEXT COMMENT '推荐理由',
    user_action VARCHAR(32) COMMENT '用户行为: view/click/add_cart/purchase',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_session (session_id),
    INDEX idx_user (user_id),
    INDEX idx_product (product_id),
    INDEX idx_action (user_action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品推荐记录表';

-- 插入测试数据
INSERT INTO products (product_id, name, category, sub_category, price, stock, description, status) VALUES
('bike-001', '山地自行车X1', '骑行', '自行车', 1299.00, 50, '适合山地骑行的专业自行车', 1),
('bike-002', '通勤自行车C1', '骑行', '自行车', 899.00, 30, '轻便舒适的城市通勤自行车', 1),
('bottle-001', '运动水壶500ml', '骑行', '水壶', 59.00, 100, '便携运动水壶', 1);

-- 测试数据1：年轻女性，偏好小红书风格，喜欢美妆和追剧
INSERT INTO `user_profiles` (`user_id`, `preferred_style`, `age`, `gender`, `interests`)
VALUES ('user_001', 'xiaohongshu', 23, 'female', '["美妆", "追剧", "奶茶"]');

-- 测试数据2：中年男性，偏好董宇辉风格，喜欢历史和运动
INSERT INTO `user_profiles` (`user_id`, `preferred_style`, `age`, `gender`, `interests`)
VALUES ('user_002', 'dongyuhui', 45, 'male', '["历史", "跑步", "书法"]');

-- 测试数据3：青少年，无明确偏好，喜欢游戏和动漫
INSERT INTO `user_profiles` (`user_id`, `preferred_style`, `age`, `gender`, `interests`)
VALUES ('user_003', '', 17, 'male', '["电竞", "动漫", "篮球"]');

-- 测试数据4：中年女性，偏好小红书，喜欢烹饪和园艺
INSERT INTO `user_profiles` (`user_id`, `preferred_style`, `age`, `gender`, `interests`)
VALUES ('user_004', 'xiaohongshu', 38, 'female', '["烘焙", "种花", "瑜伽"]');

-- 测试数据5：老年用户，偏好董宇辉，喜欢诗词和散步
INSERT INTO `user_profiles` (`user_id`, `preferred_style`, `age`, `gender`, `interests`)
VALUES ('user_005', 'dongyuhui', 62, 'female', '["古典诗词", "散步", "太极"]');

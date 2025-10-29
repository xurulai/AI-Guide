USE shopping_guide;

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
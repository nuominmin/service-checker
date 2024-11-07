import axios from 'axios';

// 创建 axios 实例
const apiClient = axios.create({
    baseURL: 'http://localhost:9000', // 替换为你的后端服务地址
    headers: {
        'Content-Type': 'application/json',
    },
});

// 获取服务状态
export const getServices = () => {
    return apiClient.get('/services');
};

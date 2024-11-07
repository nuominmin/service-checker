<template>
  <div class="container">
    <h1>服务状态</h1>
    <table class="service-table">
      <thead>
        <tr>
          <th>服务名称</th>
          <th>状态</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="service in services"
          :key="service.name"
          :class="getStatusClass(service.status)"
        >
          <td>{{ service.name }}</td>
          <td>
            <span :class="getTextClass(service.status)">
              {{ getStatusText(service.status) }}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { getServices } from '../api';

export default {
  data() {
    return {
      services: [],
      intervalId: null,
    };
  },
  created() {
    this.fetchServices();
  },
  mounted() {
    this.intervalId = setInterval(this.fetchServices, 5000);
  },
  beforeUnmount() {
    clearInterval(this.intervalId);
  },
  methods: {
    async fetchServices() {
      try {
        const response = await getServices();
        this.services = response.data.services;
      } catch (error) {
        console.error('获取服务状态失败:', error);
      }
    },

    getStatusClass(status) {
      switch (status) {
        case 0: return 'status-unknown';
        case 1: return 'status-healthy';
        case 2: return 'status-degraded';
        case 3: return 'status-unstable';
        case 4: return 'status-critical';
        case 5: return 'status-down';
        default: return '';
      }
    },

    getTextClass(status) {
      switch (status) {
        case 0: return 'text-unknown';
        case 1: return 'text-healthy';
        case 2: return 'text-degraded';
        case 3: return 'text-unstable';
        case 4: return 'text-critical';
        case 5: return 'text-down';
        default: return '';
      }
    },

    getStatusText(status) {
      switch (status) {
        case 0: return '未知';
        case 1: return '健康';
        case 2: return '降级';
        case 3: return '不稳定';
        case 4: return '严重';
        case 5: return '不可用';
        default: return '未知';
      }
    },
  },
};
</script>

<style scoped>
h1 {
  color: #333;
  text-align: center; /* 居中标题 */
}

.container {
  width: 80%; /* 设置表格容器的宽度 */
  margin: 0 auto; /* 将容器居中 */
}

.service-table {
  width: 100%; /* 表格宽度为容器的100% */
  border-collapse: collapse;
  margin-top: 20px;
  table-layout: fixed; /* 强制列宽一致 */
}

.service-table th,
.service-table td {
  padding: 12px;
  text-align: left;
  border: 1px solid #ddd;
}

.service-table th {
  background-color: #f1f1f1;
  font-weight: bold;
}

.status-unknown {
  background-color: #f0f0f0; /* 灰色背景 */
}

.status-healthy {
  background-color: #e0f7fa; /* 绿色背景 */
}

.status-degraded {
  background-color: #fff9c4; /* 黄色背景 */
}

.status-unstable {
  background-color: #ffcc80; /* 橙色背景 */
}

.status-critical {
  background-color: #ffebee; /* 红色背景 */
}

.status-down {
  background-color: #d32f2f; /* 深红色背景 */
}

.text-unknown {
  color: #757575; /* 灰色文字 */
}

.text-healthy {
  color: #388e3c; /* 绿色文字 */
}

.text-degraded {
  color: #f57f17; /* 黄色文字 */
}

.text-unstable {
  color: #f57c00; /* 橙色文字 */
}

.text-critical {
  color: #d32f2f; /* 红色文字 */
}

.text-down {
  color: #ffffff; /* 白色文字 */
}
</style>

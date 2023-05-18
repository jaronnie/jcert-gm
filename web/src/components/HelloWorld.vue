<template>
<div>
    <input type="file" ref="fileInput" @change="uploadFile">
    <ul>
      <li v-for="(file, index) in fileList" :key="index">
        {{ file.name }} (<a href="#" @click.prevent="downloadFile(file.name)">下载</a>)
      </li>
    </ul>
  </div>
</template>

<script>
import axios from "axios";
axios.defaults.baseURL = 'http://172.22.71.244:9999';


export default {
  data() {
    return {
      fileList: [],
    };
  },
  methods: {
    async uploadFile() {
      // 获取文件数据
      const file = this.$refs.fileInput.files[0];
      if (!file) return;

      // 创建 FormData 对象，用于将文件发送到服务器
      const formData = new FormData();
      formData.append("file", file);

      try {
        // 发送请求
        const response = await axios.post("/api/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        console.log(response.data);

        // 将新文件添加到列表中
        this.fileList.push({ name: response.data.Filename });
      } catch (error) {
        console.error(error);
      }
    },

    async downloadFile(filename) {
      try {
        // 发送请求
        const response = await axios.get(`/api/download/${filename}`, {
          responseType: "blob",
        });

        console.log(response.data);

        // 下载文件
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
      } catch (error) {
        console.error(error);
      }
    },
  },
};
</script>

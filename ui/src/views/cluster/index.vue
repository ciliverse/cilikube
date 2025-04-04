<template>
    <div class="dashboard-container">
      <h1>集群概览</h1>
  
      <!-- Other dashboard components -->
      <el-row :gutter="20">
        <el-col :span="24"> <!-- Or adjust span as needed -->
          <ResourceSummaryCard />
        </el-col>
      </el-row>
  
      <!-- More dashboard components -->
  
    </div>
  </template>
  
  <script setup lang="ts">
  import ResourceSummaryCard from '@/components/ClusterSummary/ResourceSummaryCard.vue'; // Adjust import path
  
  // Other dashboard logic
  </script>
  
  <style scoped>
  .dashboard-container {
    padding: 20px;
  }
  /* Add spacing if needed */
  .el-row {
      margin-bottom: 20px;
  }
  </style>
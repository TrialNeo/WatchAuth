<script lang="ts">
import {defineComponent} from 'vue'
import {POPCONFIRM_CONFIG, TABLE_CONFIG} from '@/config/elementConfig.ts'

export default defineComponent({
  name: 'ToolBar',
  computed: {
    POPCONFIRM_CONFIG() {
      return POPCONFIRM_CONFIG
    },
    TABLE_CONFIG() {
      return TABLE_CONFIG
    },
  },
})
</script>

<template>
  <el-table-column :align="TABLE_CONFIG.align" fixed="right" label="操作" min-width="130">
    <template #default="">
      <el-button v-permission="['role:edit']" link type="primary"> 编辑 </el-button>
      <el-popconfirm
        :placement="POPCONFIRM_CONFIG.placement"
        :width="POPCONFIRM_CONFIG.width"
        title="确定要删除选中的应用吗？"
      >
        <template #reference>
          <el-button v-permission="['role:delete']" link type="danger"> 删除</el-button>
        </template>
      </el-popconfirm>
    </template>
  </el-table-column>
</template>

<style lang="scss" scoped>


</style>

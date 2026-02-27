<script lang="ts">
import {defineComponent} from 'vue'
import {TABLE_CONFIG} from '@/config/elementConfig.ts'

export default defineComponent({
  name: 'MachineManager',
  computed: {
    TABLE_CONFIG() {
      return TABLE_CONFIG
    },
  },
})
const list = async () => {

}

</script>

<template>
  <div>
    <el-table :border="TABLE_CONFIG.border"  show-overflow-tooltip>
      <el-table-column :align="TABLE_CONFIG.align" type="selection" width="55" />
      <el-table-column :align="TABLE_CONFIG.align" fixed label="序号" type="index" width="55" />
      <el-table-column />
    </el-table>
  </div>
</template>

<style lang="scss" scoped></style>

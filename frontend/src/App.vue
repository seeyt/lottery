
<template>
    <div class="app-container">
        <div class="header" ref="header">
            <el-form>
            <div style="display: flex;">
                <el-form-item label="抽取人数" label-width="80px" >
                    <el-input-number  placeholder="抽奖人数" 
                    @blur="numberChange"
                    :min="1" :controls="false" size="small" v-model="num" style="width: 100px;"></el-input-number>
                </el-form-item>
                <el-form-item label="总人数" label-width="80px" v-if="total">
                    {{total}}
                </el-form-item>
                <el-form-item  label-width="20px" v-if="filePath">
                    <el-button @click="clean" type="primary" >重置抽奖</el-button>
                </el-form-item>
                <el-form-item label-width="80px" label="抽奖名单" style="flex: 1;" v-if="!filePath">
                    <div class="file">
                        <div class="file-path"  @click="selectFile">
                            <div class="file-name">{{ "点击选择抽奖名单" }}</div>
                        </div>
                    </div>
                </el-form-item>
            </div>
        </el-form>
            <div class="buttons">
                <el-button @click="lottery" type="primary"  v-if="filePath">{{tableData.length === 0 ? "抽取" : "再次抽取"}}</el-button>
                <el-button @click="save" v-if="filePath">导出获奖名单</el-button>
            </div>
        </div>
        <div class="table" :style="`height: ${tableHeight}px`">
            <el-table empty-text=" " :data="list" :max-height="tableHeight" v-if="tableData.length" >
                <el-table-column  prop="Name" label="推特名称" />
                <el-table-column  prop="Url" label="地址" />
            </el-table>
        </div>
        <div class="pagination"  ref="pagination" >
            <el-pagination
                v-model:current-page="page"
                :page-sizes="[10, 20, 30, 40]"
                v-model:page-size="size"
                background
                popper-class=" "
                layout="total, sizes, prev, pager, next, jumper"
                :total="tableData.length" />
        </div>
    </div>
</template>
<script setup>
import {Analyse, CleanFile, SelectFile, Lottery, GetNum, SetNum, SaveAs, GetTotal} from "../wailsjs/go/main/App";
import _ from 'lodash'
import {computed, onMounted, ref} from "vue";
import {ElLoading, ElMessage} from 'element-plus'

let filePath = ref("");
let pagination = ref();
let header = ref();
let tableHeight = ref(0);
let page = ref(1);
let size = ref(10);
let num = ref(0);
let total = ref(0);
let tableData = ref([]);

onMounted(async () => {
    computeHeight()
    window.addEventListener('resize', _.throttle(computeHeight, 1000));
    num.value = await GetNum()
   
})
const list  = computed(() => {
    return  tableData.value.slice((page.value-1)*size.value,page.value * size.value)
})

const computeHeight = () => {
    const windowHeight = window.innerHeight
    const headerHeight = header.value.clientHeight
    const paginationHeight = pagination.value.clientHeight
    tableHeight.value = windowHeight- headerHeight -paginationHeight
}

const selectFile = async  () => {
    try {
        filePath.value = await SelectFile()
        await analyse()
        total.value = await GetTotal()
    } catch (e) {
        ElMessage.error(e)
    }

}

const save = async () => {
    try {
        await SaveAs()
    } catch (e) {
        ElMessage.error(e)
    }

}


const analyse = async  () => {
    try {
        await Analyse()
    } catch (e) {
        ElMessage.error(e)
    }
}

const numberChange = async () => {
    try {
        if (num.value === null) num.value = 1;
        await SetNum(num.value)
    } catch (e) {
        ElMessage.error(e)
    }
}

const lottery = async () => {
    const loading = ElLoading.service({
        lock: true,
        text: '抽取中',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    try {
        let res = await Lottery()
        tableData.value = res.Data;
        loading.close();
    } catch (e) {
        ElMessage.error(e)
        loading.close()
    }
}

const clean = async() => {
    await CleanFile()
    filePath.value = ""
    tableData.value = []
    page.value = 1
    size.value = 10
    total.value = 0
}
</script>

<style>
.app-container {
    min-height: 100vh;
    min-width: 100vw;
    background: white;
    color: black;
    overflow-x: hidden;
}
.header {
    display: flex;
    padding: 10px 10px 0;
    box-sizing: border-box;
    flex-direction: column;
    height: 120px;
}
.header .buttons {
    display: flex;
    align-items: center;
}
.file {
    display: flex;
    align-items: center;
    flex: 1;
}
.file .file-path {
    -webkit-appearance: none;
    background-color: #fff;
    background-image: none;
    border-radius: 4px;
    border: 1px solid #dcdfe6;
    box-sizing: border-box;
    color: #606266;
    font-size: inherit;
    height: 32px;
    line-height: 32px;
    outline: none;
    padding: 0 15px;
    text-align: left;
    transition: border-color .2s cubic-bezier(.645,.045,.355,1);
    width: 100%;
    display: flex;
    align-items: center;
}
.file-path .file-name {
    flex: 1;
}
.table {
    padding: 0 10px;
    box-sizing: border-box;
    width: 100vw;
}
::-webkit-scrollbar{
    display:none;
    width: 0;
}
.pagination {
    display: flex;
    justify-content: flex-end;
    padding: 10px;
    height: 52px;
}
</style>

<template>
  <el-card>
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="邮件配置" name="dataSmtp">
        <el-form :model="dataSmtp" ref="dataSmtp" :rules="rules" label-width="100px" class="tab-one">
          <el-form-item label="服务器地址" prop="host">
            <el-input v-model="dataSmtp.host"></el-input>
          </el-form-item>
          <el-form-item label="服务器端口" prop="port">
            <el-input v-model.number="dataSmtp.port"></el-input>
          </el-form-item>
          <el-form-item label="用户名" prop="username">
            <el-input v-model="dataSmtp.username"></el-input>
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="dataSmtp.password"></el-input>
          </el-form-item>
          <el-form-item label="启用SSL" prop="use_ssl">
            <el-switch v-model="dataSmtp.use_ssl"></el-switch>
          </el-form-item>
          <el-form-item label="邮件from" prop="from">
            <el-input v-model="dataSmtp.from"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitForm('dataSmtp')">保存</el-button>
            <el-button @click="resetForm('dataSmtp')">重置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="其他设置" name="dataOther">
        <el-form :model="dataOther" ref="dataOther" :rules="rules" label-width="100px" class="tab-one">

          <el-form-item label="vpn对外地址" prop="link_addr">
            <el-input
                placeholder="请输入内容"
                v-model="dataOther.link_addr">
            </el-input>
          </el-form-item>

          <el-form-item label="Banner信息" prop="banner">
            <el-input
                type="textarea"
                :rows="5"
                placeholder="请输入内容"
                v-model="dataOther.banner">
            </el-input>
          </el-form-item>

          <el-form-item label="账户开通邮件" prop="account_mail">
            <el-input
                type="textarea"
                :rows="10"
                placeholder="请输入内容"
                v-model="dataOther.account_mail">
            </el-input>
          </el-form-item>

          <el-form-item label="邮件展示">
            <iframe
                width="500px"
                height="300px"
                :srcdoc="dataOther.account_mail">
            </iframe>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="submitForm('dataOther')">保存</el-button>
            <el-button @click="resetForm('dataOther')">重置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

    </el-tabs>
  </el-card>
</template>

<script>
import axios from "axios";

export default {
  name: "Other",
  created() {
    this.$emit('update:route_path', this.$route.path)
    this.$emit('update:route_name', ['基础信息', '其他设置'])
  },
  mounted() {
    this.getSmtp()
  },
  data() {
    return {
      activeName: 'dataSmtp',
      dataSmtp: {},
      dataOther: {},
      rules: {
        host: {required: true, message: '请输入服务器地址', trigger: 'blur'},
        port: [
          {required: true, message: '请输入服务器端口', trigger: 'blur'},
          {type: 'number', message: '请输入正确的服务器端口', trigger: ['blur', 'change']}
        ],
        issuer: {required: true, message: '请输入系统名称', trigger: 'blur'},
      },
    };
  },
  methods: {
    handleClick(tab, event) {
      window.console.log(tab.name, event);
      switch (tab.name) {
        case "dataSmtp":
          this.getSmtp()
          break
        case "dataOther":
          this.getOther()
          break
      }
    },
    getSmtp() {
      axios.get('/set/other/smtp').then(resp => {
        let rdata = resp.data
        console.log(rdata)
        if (rdata.code !== 0) {
          this.$message.error(rdata.msg);
          return;
        }
        this.dataSmtp = rdata.data
      }).catch(error => {
        this.$message.error('哦，请求出错');
        console.log(error);
      });
    },
    getOther() {
      axios.get('/set/other').then(resp => {
        let rdata = resp.data
        console.log(rdata)
        if (rdata.code !== 0) {
          this.$message.error(rdata.msg);
          return;
        }
        this.dataOther = rdata.data
      }).catch(error => {
        this.$message.error('哦，请求出错');
        console.log(error);
      });
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (!valid) {
          alert('error submit!');
        }

        switch (formName) {
          case "dataSmtp":
            axios.post('/set/other/smtp/edit', this.dataSmtp).then(resp => {
              var rdata = resp.data
              console.log(rdata);
              if (rdata.code === 0) {
                this.$message.success(rdata.msg);
              } else {
                this.$message.error(rdata.msg);
              }

            })
            break;
          case "dataOther":
            axios.post('/set/other/edit', this.dataOther).then(resp => {
              var rdata = resp.data
              console.log(rdata);
              if (rdata.code === 0) {
                this.$message.success(rdata.msg);
              } else {
                this.$message.error(rdata.msg);
              }
            })
            break;
        }

      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    }
  },
}
</script>

<style scoped>
.tab-one {
  width: 600px;
}

</style>

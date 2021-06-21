<template>
  <v-app class="ofd-reader">
    <div class="top-bar">
      <public-top-bar title="OFD Reader"></public-top-bar>
    </div>
    <div class="main-content">
      <history-table class="histories-box" v-if="operate == 'history'" @historyClick="openHistoryFile($event)">
      </history-table>
      <file-content class="file-content-box" v-if="operate != 'history'" ref='OfdRender'></file-content>
    </div>
    <div class="bottom-bar">
      <v-bottom-navigation v-model="operate">
        <v-btn value="history">
          <span>History</span>
          <v-icon>mdi-file-clock</v-icon>
        </v-btn>
        <v-btn value="open" @click="openNewFile">
          <span>Open</span>
          <v-icon>mdi-file-eye</v-icon>
          <input type="file" hidden ref="upload" @change="showNewFile" accept=".ofd">
        </v-btn>
        <v-btn value="print" :disabled="operate == 'history'" @click="printOfd">
          <span>Print</span>
          <v-icon>mdi-printer</v-icon>
        </v-btn>
      </v-bottom-navigation>
    </div>
  </v-app>
</template>

<script>
  import publicTopBar from "@/components/publicTopBar"
  import historyTable from "@/components/ofd_reader/historyTable"
  import fileContent from "@/components/ofd_reader/fileContent"
  export default {
    name: 'OfdReader',
    components: {
      'public-top-bar': publicTopBar,
      'history-table': historyTable,
      'file-content': fileContent,
    },
    data: () => ({
      operate: "history",
      screenWidth: document.body.clientWidth - 20,
    }),
    methods: {
      openHistoryFile(item) {
        this.$data.operate = "open";
        console.log(item)
      },
      openNewFile() {
        let uploadBtn = this.$refs.upload;
        uploadBtn.click();
      },
      async showNewFile(e) {
        let file = e.target.files[0];
        let ofdReader = this;
        if (file) {
          this.$refs.OfdRender.RenderOfdFile(file, this.$data.screenWidth)
        }
      },
      async printOfd() {
        this.$refs.OfdRender.print();
      }
    }

  }
</script>

<style scoped>
  .top-bar {
    position: absolute;
    width: 100%;
    z-index: 99999;
    top: 0;
  }

  .main-content {
    margin-top: 3rem;
    height: calc(100vh - 6.25rem);
    overflow: hidden;
    padding: 0.625rem;
  }

  .bottom-bar {
    position: absolute;
    width: 100%;
    z-index: 99999;
    bottom: 0;
  }

  .histories-box,
  .file-content-box {
    position: relative;
    width: 100%;
    height: 100%;
    overflow-y: auto;
    margin: auto;
  }
</style>

<template>
  <div>
    <div v-if="loading" class="loading-box">
      <v-progress-circular :size="100" :width="10" color="purple" indeterminate></v-progress-circular>
    </div>
    <div id="content" v-show="loaded" ref="contentDiv">
    </div>
    <v-snackbar v-model="printErr" :timeount="1000" centered>There is no content to print
      <template v-slot:action="{ attrs }">
        <v-btn color="blue" text v-bind="attrs" @click="printErr = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script>
  import {
    parseOfdDocument,
    renderOfd,
    renderOfdByScale,
    digestCheck,
    getPageScale,
    setPageScale
  } from "@/utils/ofd/ofd";
  export default {
    name: "fileContent",
    data: () => ({
      loading: false,
      loaded: false,
      currOfd: null,
      printErr: false,
    }),
    methods: {
      displayOfdDiv(divs) {
        this.scale = getPageScale();
        let contentDiv = document.getElementById('content');
        contentDiv.innerHTML = '';
        for (const div of divs) {
          div.setAttribute("class", "ofd-page")
          contentDiv.appendChild(div)
        }
        // for (let ele of document.getElementsByName('seal_img_div')) {
        //   this.addEventOnSealDiv(ele, JSON.parse(ele.dataset.sesSignature), JSON.parse(ele.dataset.signedInfo));
        // }
      },
      RenderOfdFile(file, width) {
        if (!file) {
          return
        }
        this.$data.loading = true;
        this.$data.loaded = false;
        let that = this;
        let t = new Date().getTime();
        this.loading = true;
        parseOfdDocument({
          ofd: file,
          success(res) {
            that.currOfd = res[0];
            that.pageCount = res[0].pages.length;
            const divs = renderOfd(width, res[0]);
            that.displayOfdDiv(divs);
            that.loading = false;
            that.loaded = true;
            console.log("loaded")
          },
          fail(error) {
            that.loading = false;
            that.loaded = false;
            console.log(error)
          }
        });
      },
      print() {
        if (!this.currOfd) {
          this.$data.printErr = true;
          return
        }
        let dom = this.$refs["contentDiv"];
        let childs = dom.children;
        this.loading = true;
        let list = [];
        let i = 0;
        for (let ele of childs) {
          list.push(ele.cloneNode(true));
          this.percentage = i / childs.length;
        }
        if (list.length > 0) {
          var mywindow = window.open("打印窗口", "_blank");
          //给新打开的标签页添加画布内容（这里指的是id=div2img元素的内容）
          var documentBody = mywindow.document.body;
          console.log(list.length)
          for (let canvas of list) {
            documentBody.appendChild(canvas);
          }
          this.loading = false;
          //焦点移到新打开的标签页
          mywindow.focus();
          //执行打印的方法（注意打印方法是打印的当前窗口内的元素，所以前面才新建一个窗口：print()--打印当前窗口的内容。）
          mywindow.print();
          //操作完成之后关闭当前标签页（点击确定或者取消都会关闭）
          mywindow.close();
        }
      },
    }
  }
</script>

<style>
  .loading-box {
    position: relative;
    margin-top: calc(50vh - 8rem);
  }

  #content {
    background-color: #2C3E50;
    padding: 0.625rem;
    justify-content: center;
  }

  #content .ofd-page {
    margin: 0 auto;
  }
</style>

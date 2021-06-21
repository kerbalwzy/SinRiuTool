<template>
  <div>
    <v-toolbar flat>
      <v-btn color="warning" dark class="mb-2" @click="dialogDelete = 1">
        Delete
      </v-btn>
      <v-spacer></v-spacer>
      <v-text-field v-model="search" append-icon="mdi-magnify" label="Search" single-line hide-details
        @keyup.enter.native="searchEnter()">
      </v-text-field>
      <v-dialog v-model="dialogDelete" max-width="550px">
        <v-card>
          <v-card-title class="headline" v-if="selected.length == 0">Please select at least one!
          </v-card-title>
          <v-card-title class="headline" v-if="selected.length > 0">Are you sure you want to delete this item?
          </v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="dialogDelete = false">Cancel</v-btn>
            <v-btn color="blue darken-1" text @click="deleteItemConfirm" v-if="selected.length > 0">OK</v-btn>
            <v-spacer></v-spacer>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-toolbar>
    <v-data-table :headers="headers" :items="histories" item-key="name" class="elevation-1" :search="search"
      :loading="loading" loading-text="Loading... Please wait" v-model="selected" show-select
      @dblclick:row="historyClick">
    </v-data-table>
  </div>
</template>
<script>
  export default {
    name: "historyTable",
    data: () => ({
      search: '',
      dialogDelete: false,
      selected: [],
      histories: function() {
        let res = [];
        for (let i = 0; i < 20; i++) {
          res.push({
            name: 'KitKat' + i,
            upload_at: "2021-05-28 12:21:00",
            sign_str: "xxxxxxx"
          })
        }
        return res
      }()
    }),
    computed: {
      headers() {
        return [{
            text: 'FileName',
            align: 'start',
            sortable: false,
            value: 'name',
          },
          {
            text: 'UploadTime',
            sortable: false,
            value: 'upload_at',
            align: "center",
            width: "240px"
          }
        ]
      },
      loading() {
        return this.histories.length == 0
      }
    },
    methods: {
      deleteItemConfirm() {
        console.log("deleteItemConfirm")
      },
      historyClick(_, value) {
        this.$emit("historyClick", value.item)
      },
      searchEnter(event) {
        console.log(this.$data.search)
      }
    },
  }
</script>

<style>
</style>

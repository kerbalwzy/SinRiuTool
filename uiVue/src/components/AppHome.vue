<template>
  <v-app class="app-home">
    <div class="top-bar">
      <public-top-bar title="App Home"></public-top-bar>
    </div>
    <div class="app-list-box">
      <v-container fluid>
        <v-row dense justify="center">
          <v-col v-for="card in cards" :key="card.title" cols="auto">
            <v-card @click="$router.push({path:card.path})">
              <v-img :src="card.src" class="white--text align-end" gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
                height="220px" width="250px">
                <v-card-title v-text="card.title"></v-card-title>
              </v-img>
              <v-card-actions>
                <v-spacer></v-spacer>
                <div v-for="icon in card.icons">
                  <v-icon>{{icon}}</v-icon>
                </div>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </div>
  </v-app>
</template>

<script>
  import publicTopBar from "@/components/publicTopBar"
  export default {
    name: 'AppHome',
    components: {
      'public-top-bar': publicTopBar,
    },
    data() {
      return {
        cards: [{
            title: 'OFD Reader',
            path: "/ofd-reader",
            src: require("@/assets/ofd-reader-logo.png"),
            icons: ["mdi-file-clock", "mdi-file-eye", "mdi-printer"],

          },
          {
            title: 'PayPal Robot',
            path: "/paypal-auto",
            src: require("@/assets/paypal-auto-logo.png"),
            icons: ["mdi-robot-happy", "mdi-shield-alert", "mdi-monitor-screenshot", ]
          },
        ],
      }
    },

  }
</script>

<style scoped>
  .top-bar {
    position: absolute;
    width: 100%;
    z-index: 99999;
    top: 0;
  }

  .app-list-box {
    margin-top: 4rem;
  }
</style>

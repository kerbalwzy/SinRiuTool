<template>
  <v-toolbar dense>
    <v-toolbar-title>{{title}}</v-toolbar-title>
    <v-spacer></v-spacer>
    <v-btn icon>
      <v-icon  @click="$router.push({path:'/'})">mdi-apps</v-icon>
    </v-btn>
  </v-toolbar>
</template>

<script>
  export default {
    name: 'publicTopBar',
    props: ['title'],
  }
</script>

<style>
</style>

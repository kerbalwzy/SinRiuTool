import Vue from 'vue'
import Router from 'vue-router'
import AppHome from '@/components/AppHome'
import OfdReader from '@/components/OfdReader'
import PayPalAuto from '@/components/PayPalAuto'

Vue.use(Router)

export default new Router({
  routes: [{
      path: '/',
      name: 'AppHome',
      component: AppHome,
    },
    {
      path: '/ofd-reader',
      name: 'OfdReader',
      component: OfdReader,
    },
    {
      path: '/paypal-auto',
      name: 'PayPalAuto',
      component: PayPalAuto,
    }
  ]
})

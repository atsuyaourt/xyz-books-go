import 'virtual:uno.css'
import 'vuetify/styles'
import './style.css'

import { createApp } from 'vue'
import App from './App.vue'

import { VueQueryPlugin } from '@tanstack/vue-query'

import { setupLayouts } from 'virtual:generated-layouts'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router/auto'

import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const vuetify = createVuetify({
  components,
  directives,
})

const recursiveLayouts = (route: RouteRecordRaw): RouteRecordRaw => {
  if (route.children) {
    for (let i = 0; i < route.children.length; i++) {
      route.children[i] = recursiveLayouts(route.children[i])
    }

    return route
  }

  return setupLayouts([route])[0]
}

const router = createRouter({
  history: createWebHistory(),
  extendRoutes: (routes) => routes.map((route) => recursiveLayouts(route)),
})

const app = createApp(App)
app.use(router)
app.use(vuetify)
app.use(VueQueryPlugin)
app.mount('#app')

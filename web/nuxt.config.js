// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },

  // app/ srcDir layout (pages, app.vue, ... under app/).
  future: { compatibilityVersion: 4 },

  // Node.js server runtime is required — @grpc/grpc-js needs net/http2,
  // not available on edge/serverless presets.
  nitro: {
    preset: 'node-server',
  },

  modules: ['@pinia/nuxt'],

  css: ['~/assets/css/main.css'],

  components: [{ path: '~/components', pathPrefix: false }],

  runtimeConfig: {
    // Private — server routes only (server/api/**), never reaches the
    // browser bundle. Overridable with NUXT_GRPC_BASE_URL, etc.
    grpc: {
      baseUrl: '',
      timeoutMs: 15000,
    },
    backendHttpUrl: '',
    public: {
      appName: 'TomLights',
    },
  },
})

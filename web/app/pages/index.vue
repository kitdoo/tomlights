<script setup>
const config = useRuntimeConfig()

const { data: status, refresh, pending } = await useFetch('/api/status')
</script>

<template>
  <main>
    <h1>{{ config.public.appName }}</h1>

    <section class="panel">
      <h2>Connections</h2>
      <div class="row">
        <span>Web (Nuxt)</span>
        <span :class="status?.web.ok ? 'ok' : 'bad'">{{ status?.web.ok ? 'OK' : '—' }}</span>
      </div>
      <div class="row">
        <span>Backend gRPC</span>
        <span :class="status?.grpc.ok ? 'ok' : 'bad'">{{ status?.grpc.status }}</span>
      </div>
      <div class="row">
        <span>Backend HTTP + MongoDB</span>
        <span :class="status?.http.ok ? 'ok' : 'bad'">{{ status?.http.status }}</span>
      </div>
      <button :disabled="pending" @click="refresh()">Refresh</button>
    </section>
  </main>
</template>

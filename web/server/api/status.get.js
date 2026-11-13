// GET /api/status — connectivity check for the whole chain:
// browser -> Nitro -> backend gRPC (grpc.health.v1) and backend HTTP
// (/internal/healthz/details, which includes the MongoDB check).
export default defineEventHandler(async () => {
  const config = useRuntimeConfig()

  const grpcStatus = await grpcCall(
    getServiceClient('health.proto', 'grpc.health.v1.Health'),
    'Check',
    { service: '' },
  )
    .then((res) => ({ ok: res.status === 'SERVING', status: res.status }))
    .catch((err) => ({ ok: false, status: 'UNREACHABLE', error: err.details || err.message }))

  let httpStatus = { ok: false, status: 'NOT_CONFIGURED' }
  if (config.backendHttpUrl) {
    httpStatus = await $fetch(`${config.backendHttpUrl}/internal/healthz/details`, { timeout: 5000 })
      .then((details) => ({ ok: true, status: 'SERVING', details }))
      .catch((err) => ({
        ok: false,
        status: err.statusCode ? `HTTP ${err.statusCode}` : 'UNREACHABLE',
        details: err.data,
      }))
  }

  return {
    web: { ok: true },
    grpc: grpcStatus,
    http: httpStatus,
  }
})

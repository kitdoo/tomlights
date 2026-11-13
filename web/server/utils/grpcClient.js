import path from 'node:path'
import * as grpc from '@grpc/grpc-js'
import protoLoader from '@grpc/proto-loader'

// Node gRPC client, alive only inside server/api/** (Nitro). Browser code
// must never import this file.
//
// Resolved from process.cwd(), not import.meta.url: Nitro bundles this file
// into a single chunk, so a module-relative path breaks in production.
// process.cwd() is web/ in dev and /app in the Docker image, where the
// Dockerfile copies proto/ next to .output/.
const PROTO_DIR = path.resolve(process.cwd(), 'proto')

const packageDefCache = new Map()
const serviceClientCache = new Map()

function loadPackageDefinition(protoFile) {
  if (packageDefCache.has(protoFile)) return packageDefCache.get(protoFile)
  const definition = protoLoader.loadSync(protoFile, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
    includeDirs: [PROTO_DIR],
  })
  const packageDefinition = grpc.loadPackageDefinition(definition)
  packageDefCache.set(protoFile, packageDefinition)
  return packageDefinition
}

/**
 * @param {string} protoFile file name under web/proto/, e.g. 'health.proto'
 * @param {string} packagePath dotted proto package + service, e.g.
 *   'grpc.health.v1.Health'
 */
export function getServiceClient(protoFile, packagePath) {
  const cacheKey = `${protoFile}:${packagePath}`
  if (serviceClientCache.has(cacheKey)) return serviceClientCache.get(cacheKey)

  const config = useRuntimeConfig()
  const packageDefinition = loadPackageDefinition(path.join(PROTO_DIR, protoFile))
  const ServiceCtor = packagePath.split('.').reduce((obj, key) => obj?.[key], packageDefinition)
  if (!ServiceCtor) {
    throw new Error(`gRPC service not found: ${packagePath} in ${protoFile}`)
  }

  const baseUrl = config.grpc.baseUrl
  if (!baseUrl) {
    throw new Error('runtimeConfig.grpc.baseUrl is not set (env NUXT_GRPC_BASE_URL)')
  }

  const client = new ServiceCtor(baseUrl, grpc.credentials.createInsecure())
  serviceClientCache.set(cacheKey, client)
  return client
}

/**
 * Promisified unary call with optional bearer token metadata and a
 * deadline from runtimeConfig.grpc.timeoutMs.
 * @param {import('@grpc/grpc-js').Client} client
 * @param {string} methodName
 * @param {object} request
 * @param {{ token?: string }} [opts]
 */
export function grpcCall(client, methodName, request, opts = {}) {
  const config = useRuntimeConfig()
  const metadata = new grpc.Metadata()
  if (opts.token) {
    metadata.set('authorization', `Bearer ${opts.token}`)
  }
  const deadline = new Date(Date.now() + (Number(config.grpc.timeoutMs) || 15000))

  return new Promise((resolve, reject) => {
    client[methodName](request, metadata, { deadline }, (err, response) => {
      if (err) {
        reject(err)
        return
      }
      resolve(response)
    })
  })
}

import * as grpc from '@grpc/grpc-js'
import { H3Error } from 'h3'

// Maps a raw @grpc/grpc-js error to the single JSON error shape every
// server route returns on failure: { error: { code, message, field? } }.
const CODE_NAMES = {
  [grpc.status.NOT_FOUND]: 'NOT_FOUND',
  [grpc.status.ALREADY_EXISTS]: 'ALREADY_EXISTS',
  [grpc.status.ABORTED]: 'ABORTED',
  [grpc.status.FAILED_PRECONDITION]: 'FAILED_PRECONDITION',
  [grpc.status.INVALID_ARGUMENT]: 'INVALID_ARGUMENT',
  [grpc.status.UNAUTHENTICATED]: 'UNAUTHENTICATED',
  [grpc.status.PERMISSION_DENIED]: 'PERMISSION_DENIED',
  [grpc.status.UNIMPLEMENTED]: 'UNIMPLEMENTED',
  [grpc.status.UNAVAILABLE]: 'UNAVAILABLE',
  [grpc.status.DEADLINE_EXCEEDED]: 'UNAVAILABLE',
}

const HTTP_STATUS = {
  NOT_FOUND: 404,
  ALREADY_EXISTS: 409,
  ABORTED: 409,
  FAILED_PRECONDITION: 412,
  INVALID_ARGUMENT: 400,
  UNAUTHENTICATED: 401,
  PERMISSION_DENIED: 403,
  UNIMPLEMENTED: 501,
  UNAVAILABLE: 503,
  INTERNAL: 500,
}

/**
 * @param {import('@grpc/grpc-js').ServiceError} grpcError
 * @returns {H3Error}
 */
export function mapGrpcError(grpcError) {
  const code = CODE_NAMES[grpcError.code] || 'INTERNAL'
  const message = grpcError.details || grpcError.message || 'Unexpected error'
  const field = extractFieldFromMessage(code, message)

  const error = new H3Error(message)
  error.statusCode = HTTP_STATUS[code] || 500
  error.data = { error: { code, message, ...(field ? { field } : {}) } }
  return error
}

// InvalidArgument messages from buf.validate follow "<field>: <reason>".
function extractFieldFromMessage(code, message) {
  if (code !== 'INVALID_ARGUMENT') return undefined
  const match = /^([a-zA-Z0-9_.]+):/.exec(message)
  return match ? match[1] : undefined
}

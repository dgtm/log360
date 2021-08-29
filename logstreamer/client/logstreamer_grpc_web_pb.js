/**
 * @fileoverview gRPC-Web generated client stub for logstreamer
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.logstreamer = require('./logstreamer_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.logstreamer.LogStreamerClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.logstreamer.LogStreamerPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.logstreamer.LogRequest,
 *   !proto.logstreamer.LogResponse>}
 */
const methodDescriptor_LogStreamer_ProcessRequest = new grpc.web.MethodDescriptor(
  '/logstreamer.LogStreamer/ProcessRequest',
  grpc.web.MethodType.UNARY,
  proto.logstreamer.LogRequest,
  proto.logstreamer.LogResponse,
  /**
   * @param {!proto.logstreamer.LogRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.logstreamer.LogResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.logstreamer.LogRequest,
 *   !proto.logstreamer.LogResponse>}
 */
const methodInfo_LogStreamer_ProcessRequest = new grpc.web.AbstractClientBase.MethodInfo(
  proto.logstreamer.LogResponse,
  /**
   * @param {!proto.logstreamer.LogRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.logstreamer.LogResponse.deserializeBinary
);


/**
 * @param {!proto.logstreamer.LogRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.logstreamer.LogResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.logstreamer.LogResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.logstreamer.LogStreamerClient.prototype.processRequest =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/logstreamer.LogStreamer/ProcessRequest',
      request,
      metadata || {},
      methodDescriptor_LogStreamer_ProcessRequest,
      callback);
};


/**
 * @param {!proto.logstreamer.LogRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.logstreamer.LogResponse>}
 *     Promise that resolves to the response
 */
proto.logstreamer.LogStreamerPromiseClient.prototype.processRequest =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/logstreamer.LogStreamer/ProcessRequest',
      request,
      metadata || {},
      methodDescriptor_LogStreamer_ProcessRequest);
};


module.exports = proto.logstreamer;


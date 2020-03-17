/* eslint-disable */
// source: clobquote.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var modelcommon_pb = require('./modelcommon_pb.js');
goog.object.extend(proto, modelcommon_pb);
goog.exportSymbol('proto.model.ClobLine', null, global);
goog.exportSymbol('proto.model.ClobQuote', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.model.ClobLine = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.ClobLine, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.ClobLine.displayName = 'proto.model.ClobLine';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.model.ClobQuote = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.model.ClobQuote.repeatedFields_, null);
};
goog.inherits(proto.model.ClobQuote, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.ClobQuote.displayName = 'proto.model.ClobQuote';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.model.ClobLine.prototype.toObject = function(opt_includeInstance) {
  return proto.model.ClobLine.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.ClobLine} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ClobLine.toObject = function(includeInstance, msg) {
  var f, obj = {
    size: (f = msg.getSize()) && modelcommon_pb.Decimal64.toObject(includeInstance, f),
    price: (f = msg.getPrice()) && modelcommon_pb.Decimal64.toObject(includeInstance, f),
    entryid: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.model.ClobLine}
 */
proto.model.ClobLine.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.ClobLine;
  return proto.model.ClobLine.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.ClobLine} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.ClobLine}
 */
proto.model.ClobLine.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new modelcommon_pb.Decimal64;
      reader.readMessage(value,modelcommon_pb.Decimal64.deserializeBinaryFromReader);
      msg.setSize(value);
      break;
    case 2:
      var value = new modelcommon_pb.Decimal64;
      reader.readMessage(value,modelcommon_pb.Decimal64.deserializeBinaryFromReader);
      msg.setPrice(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEntryid(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.model.ClobLine.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.ClobLine.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.ClobLine} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ClobLine.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSize();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      modelcommon_pb.Decimal64.serializeBinaryToWriter
    );
  }
  f = message.getPrice();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      modelcommon_pb.Decimal64.serializeBinaryToWriter
    );
  }
  f = message.getEntryid();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional Decimal64 size = 1;
 * @return {?proto.model.Decimal64}
 */
proto.model.ClobLine.prototype.getSize = function() {
  return /** @type{?proto.model.Decimal64} */ (
    jspb.Message.getWrapperField(this, modelcommon_pb.Decimal64, 1));
};


/**
 * @param {?proto.model.Decimal64|undefined} value
 * @return {!proto.model.ClobLine} returns this
*/
proto.model.ClobLine.prototype.setSize = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.model.ClobLine} returns this
 */
proto.model.ClobLine.prototype.clearSize = function() {
  return this.setSize(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.model.ClobLine.prototype.hasSize = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Decimal64 price = 2;
 * @return {?proto.model.Decimal64}
 */
proto.model.ClobLine.prototype.getPrice = function() {
  return /** @type{?proto.model.Decimal64} */ (
    jspb.Message.getWrapperField(this, modelcommon_pb.Decimal64, 2));
};


/**
 * @param {?proto.model.Decimal64|undefined} value
 * @return {!proto.model.ClobLine} returns this
*/
proto.model.ClobLine.prototype.setPrice = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.model.ClobLine} returns this
 */
proto.model.ClobLine.prototype.clearPrice = function() {
  return this.setPrice(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.model.ClobLine.prototype.hasPrice = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string entryId = 3;
 * @return {string}
 */
proto.model.ClobLine.prototype.getEntryid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.ClobLine} returns this
 */
proto.model.ClobLine.prototype.setEntryid = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.model.ClobQuote.repeatedFields_ = [2,3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.model.ClobQuote.prototype.toObject = function(opt_includeInstance) {
  return proto.model.ClobQuote.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.ClobQuote} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ClobQuote.toObject = function(includeInstance, msg) {
  var f, obj = {
    listingid: jspb.Message.getFieldWithDefault(msg, 1, 0),
    bidsList: jspb.Message.toObjectList(msg.getBidsList(),
    proto.model.ClobLine.toObject, includeInstance),
    offersList: jspb.Message.toObjectList(msg.getOffersList(),
    proto.model.ClobLine.toObject, includeInstance),
    streaminterrupted: jspb.Message.getBooleanFieldWithDefault(msg, 4, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.model.ClobQuote}
 */
proto.model.ClobQuote.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.ClobQuote;
  return proto.model.ClobQuote.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.ClobQuote} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.ClobQuote}
 */
proto.model.ClobQuote.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setListingid(value);
      break;
    case 2:
      var value = new proto.model.ClobLine;
      reader.readMessage(value,proto.model.ClobLine.deserializeBinaryFromReader);
      msg.addBids(value);
      break;
    case 3:
      var value = new proto.model.ClobLine;
      reader.readMessage(value,proto.model.ClobLine.deserializeBinaryFromReader);
      msg.addOffers(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setStreaminterrupted(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.model.ClobQuote.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.ClobQuote.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.ClobQuote} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ClobQuote.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getListingid();
  if (f !== 0) {
    writer.writeInt32(
      1,
      f
    );
  }
  f = message.getBidsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto.model.ClobLine.serializeBinaryToWriter
    );
  }
  f = message.getOffersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.model.ClobLine.serializeBinaryToWriter
    );
  }
  f = message.getStreaminterrupted();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
};


/**
 * optional int32 listingId = 1;
 * @return {number}
 */
proto.model.ClobQuote.prototype.getListingid = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.ClobQuote} returns this
 */
proto.model.ClobQuote.prototype.setListingid = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * repeated ClobLine bids = 2;
 * @return {!Array<!proto.model.ClobLine>}
 */
proto.model.ClobQuote.prototype.getBidsList = function() {
  return /** @type{!Array<!proto.model.ClobLine>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.ClobLine, 2));
};


/**
 * @param {!Array<!proto.model.ClobLine>} value
 * @return {!proto.model.ClobQuote} returns this
*/
proto.model.ClobQuote.prototype.setBidsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.model.ClobLine=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ClobLine}
 */
proto.model.ClobQuote.prototype.addBids = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.model.ClobLine, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.ClobQuote} returns this
 */
proto.model.ClobQuote.prototype.clearBidsList = function() {
  return this.setBidsList([]);
};


/**
 * repeated ClobLine offers = 3;
 * @return {!Array<!proto.model.ClobLine>}
 */
proto.model.ClobQuote.prototype.getOffersList = function() {
  return /** @type{!Array<!proto.model.ClobLine>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.ClobLine, 3));
};


/**
 * @param {!Array<!proto.model.ClobLine>} value
 * @return {!proto.model.ClobQuote} returns this
*/
proto.model.ClobQuote.prototype.setOffersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.model.ClobLine=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ClobLine}
 */
proto.model.ClobQuote.prototype.addOffers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.model.ClobLine, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.ClobQuote} returns this
 */
proto.model.ClobQuote.prototype.clearOffersList = function() {
  return this.setOffersList([]);
};


/**
 * optional bool streamInterrupted = 4;
 * @return {boolean}
 */
proto.model.ClobQuote.prototype.getStreaminterrupted = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.model.ClobQuote} returns this
 */
proto.model.ClobQuote.prototype.setStreaminterrupted = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


goog.object.extend(exports, proto.model);

// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf type {@code Common.SecAltIDGrp}
 */
public  final class SecAltIDGrp extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:Common.SecAltIDGrp)
    SecAltIDGrpOrBuilder {
private static final long serialVersionUID = 0L;
  // Use SecAltIDGrp.newBuilder() to construct.
  private SecAltIDGrp(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private SecAltIDGrp() {
    securityAltId_ = "";
    securityAltIdSource_ = 0;
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new SecAltIDGrp();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private SecAltIDGrp(
      com.google.protobuf.CodedInputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    this();
    if (extensionRegistry == null) {
      throw new java.lang.NullPointerException();
    }
    com.google.protobuf.UnknownFieldSet.Builder unknownFields =
        com.google.protobuf.UnknownFieldSet.newBuilder();
    try {
      boolean done = false;
      while (!done) {
        int tag = input.readTag();
        switch (tag) {
          case 0:
            done = true;
            break;
          case 10: {
            java.lang.String s = input.readStringRequireUtf8();

            securityAltId_ = s;
            break;
          }
          case 16: {
            int rawValue = input.readEnum();

            securityAltIdSource_ = rawValue;
            break;
          }
          default: {
            if (!parseUnknownField(
                input, unknownFields, extensionRegistry, tag)) {
              done = true;
            }
            break;
          }
        }
      }
    } catch (com.google.protobuf.InvalidProtocolBufferException e) {
      throw e.setUnfinishedMessage(this);
    } catch (java.io.IOException e) {
      throw new com.google.protobuf.InvalidProtocolBufferException(
          e).setUnfinishedMessage(this);
    } finally {
      this.unknownFields = unknownFields.build();
      makeExtensionsImmutable();
    }
  }
  public static final com.google.protobuf.Descriptors.Descriptor
      getDescriptor() {
    return org.fixprotocol.components.Common.internal_static_Common_SecAltIDGrp_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return org.fixprotocol.components.Common.internal_static_Common_SecAltIDGrp_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            org.fixprotocol.components.SecAltIDGrp.class, org.fixprotocol.components.SecAltIDGrp.Builder.class);
  }

  public static final int SECURITY_ALT_ID_FIELD_NUMBER = 1;
  private volatile java.lang.Object securityAltId_;
  /**
   * <code>string security_alt_id = 1 [(.fix.tag) = 455, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The securityAltId.
   */
  public java.lang.String getSecurityAltId() {
    java.lang.Object ref = securityAltId_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      securityAltId_ = s;
      return s;
    }
  }
  /**
   * <code>string security_alt_id = 1 [(.fix.tag) = 455, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for securityAltId.
   */
  public com.google.protobuf.ByteString
      getSecurityAltIdBytes() {
    java.lang.Object ref = securityAltId_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      securityAltId_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int SECURITY_ALT_ID_SOURCE_FIELD_NUMBER = 2;
  private int securityAltIdSource_;
  /**
   * <code>.Common.SecurityAltIDSourceEnum security_alt_id_source = 2 [(.fix.tag) = 456, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The enum numeric value on the wire for securityAltIdSource.
   */
  public int getSecurityAltIdSourceValue() {
    return securityAltIdSource_;
  }
  /**
   * <code>.Common.SecurityAltIDSourceEnum security_alt_id_source = 2 [(.fix.tag) = 456, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The securityAltIdSource.
   */
  public org.fixprotocol.components.SecurityAltIDSourceEnum getSecurityAltIdSource() {
    @SuppressWarnings("deprecation")
    org.fixprotocol.components.SecurityAltIDSourceEnum result = org.fixprotocol.components.SecurityAltIDSourceEnum.valueOf(securityAltIdSource_);
    return result == null ? org.fixprotocol.components.SecurityAltIDSourceEnum.UNRECOGNIZED : result;
  }

  private byte memoizedIsInitialized = -1;
  @java.lang.Override
  public final boolean isInitialized() {
    byte isInitialized = memoizedIsInitialized;
    if (isInitialized == 1) return true;
    if (isInitialized == 0) return false;

    memoizedIsInitialized = 1;
    return true;
  }

  @java.lang.Override
  public void writeTo(com.google.protobuf.CodedOutputStream output)
                      throws java.io.IOException {
    if (!getSecurityAltIdBytes().isEmpty()) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 1, securityAltId_);
    }
    if (securityAltIdSource_ != org.fixprotocol.components.SecurityAltIDSourceEnum.SECURITY_ALT_ID_SOURCE_UNSPECIFIED.getNumber()) {
      output.writeEnum(2, securityAltIdSource_);
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!getSecurityAltIdBytes().isEmpty()) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, securityAltId_);
    }
    if (securityAltIdSource_ != org.fixprotocol.components.SecurityAltIDSourceEnum.SECURITY_ALT_ID_SOURCE_UNSPECIFIED.getNumber()) {
      size += com.google.protobuf.CodedOutputStream
        .computeEnumSize(2, securityAltIdSource_);
    }
    size += unknownFields.getSerializedSize();
    memoizedSize = size;
    return size;
  }

  @java.lang.Override
  public boolean equals(final java.lang.Object obj) {
    if (obj == this) {
     return true;
    }
    if (!(obj instanceof org.fixprotocol.components.SecAltIDGrp)) {
      return super.equals(obj);
    }
    org.fixprotocol.components.SecAltIDGrp other = (org.fixprotocol.components.SecAltIDGrp) obj;

    if (!getSecurityAltId()
        .equals(other.getSecurityAltId())) return false;
    if (securityAltIdSource_ != other.securityAltIdSource_) return false;
    if (!unknownFields.equals(other.unknownFields)) return false;
    return true;
  }

  @java.lang.Override
  public int hashCode() {
    if (memoizedHashCode != 0) {
      return memoizedHashCode;
    }
    int hash = 41;
    hash = (19 * hash) + getDescriptor().hashCode();
    hash = (37 * hash) + SECURITY_ALT_ID_FIELD_NUMBER;
    hash = (53 * hash) + getSecurityAltId().hashCode();
    hash = (37 * hash) + SECURITY_ALT_ID_SOURCE_FIELD_NUMBER;
    hash = (53 * hash) + securityAltIdSource_;
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static org.fixprotocol.components.SecAltIDGrp parseFrom(
      com.google.protobuf.CodedInputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }

  @java.lang.Override
  public Builder newBuilderForType() { return newBuilder(); }
  public static Builder newBuilder() {
    return DEFAULT_INSTANCE.toBuilder();
  }
  public static Builder newBuilder(org.fixprotocol.components.SecAltIDGrp prototype) {
    return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
  }
  @java.lang.Override
  public Builder toBuilder() {
    return this == DEFAULT_INSTANCE
        ? new Builder() : new Builder().mergeFrom(this);
  }

  @java.lang.Override
  protected Builder newBuilderForType(
      com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
    Builder builder = new Builder(parent);
    return builder;
  }
  /**
   * Protobuf type {@code Common.SecAltIDGrp}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:Common.SecAltIDGrp)
      org.fixprotocol.components.SecAltIDGrpOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return org.fixprotocol.components.Common.internal_static_Common_SecAltIDGrp_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return org.fixprotocol.components.Common.internal_static_Common_SecAltIDGrp_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              org.fixprotocol.components.SecAltIDGrp.class, org.fixprotocol.components.SecAltIDGrp.Builder.class);
    }

    // Construct using org.fixprotocol.components.SecAltIDGrp.newBuilder()
    private Builder() {
      maybeForceBuilderInitialization();
    }

    private Builder(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      super(parent);
      maybeForceBuilderInitialization();
    }
    private void maybeForceBuilderInitialization() {
      if (com.google.protobuf.GeneratedMessageV3
              .alwaysUseFieldBuilders) {
      }
    }
    @java.lang.Override
    public Builder clear() {
      super.clear();
      securityAltId_ = "";

      securityAltIdSource_ = 0;

      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return org.fixprotocol.components.Common.internal_static_Common_SecAltIDGrp_descriptor;
    }

    @java.lang.Override
    public org.fixprotocol.components.SecAltIDGrp getDefaultInstanceForType() {
      return org.fixprotocol.components.SecAltIDGrp.getDefaultInstance();
    }

    @java.lang.Override
    public org.fixprotocol.components.SecAltIDGrp build() {
      org.fixprotocol.components.SecAltIDGrp result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public org.fixprotocol.components.SecAltIDGrp buildPartial() {
      org.fixprotocol.components.SecAltIDGrp result = new org.fixprotocol.components.SecAltIDGrp(this);
      result.securityAltId_ = securityAltId_;
      result.securityAltIdSource_ = securityAltIdSource_;
      onBuilt();
      return result;
    }

    @java.lang.Override
    public Builder clone() {
      return super.clone();
    }
    @java.lang.Override
    public Builder setField(
        com.google.protobuf.Descriptors.FieldDescriptor field,
        java.lang.Object value) {
      return super.setField(field, value);
    }
    @java.lang.Override
    public Builder clearField(
        com.google.protobuf.Descriptors.FieldDescriptor field) {
      return super.clearField(field);
    }
    @java.lang.Override
    public Builder clearOneof(
        com.google.protobuf.Descriptors.OneofDescriptor oneof) {
      return super.clearOneof(oneof);
    }
    @java.lang.Override
    public Builder setRepeatedField(
        com.google.protobuf.Descriptors.FieldDescriptor field,
        int index, java.lang.Object value) {
      return super.setRepeatedField(field, index, value);
    }
    @java.lang.Override
    public Builder addRepeatedField(
        com.google.protobuf.Descriptors.FieldDescriptor field,
        java.lang.Object value) {
      return super.addRepeatedField(field, value);
    }
    @java.lang.Override
    public Builder mergeFrom(com.google.protobuf.Message other) {
      if (other instanceof org.fixprotocol.components.SecAltIDGrp) {
        return mergeFrom((org.fixprotocol.components.SecAltIDGrp)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(org.fixprotocol.components.SecAltIDGrp other) {
      if (other == org.fixprotocol.components.SecAltIDGrp.getDefaultInstance()) return this;
      if (!other.getSecurityAltId().isEmpty()) {
        securityAltId_ = other.securityAltId_;
        onChanged();
      }
      if (other.securityAltIdSource_ != 0) {
        setSecurityAltIdSourceValue(other.getSecurityAltIdSourceValue());
      }
      this.mergeUnknownFields(other.unknownFields);
      onChanged();
      return this;
    }

    @java.lang.Override
    public final boolean isInitialized() {
      return true;
    }

    @java.lang.Override
    public Builder mergeFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      org.fixprotocol.components.SecAltIDGrp parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (org.fixprotocol.components.SecAltIDGrp) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }

    private java.lang.Object securityAltId_ = "";
    /**
     * <code>string security_alt_id = 1 [(.fix.tag) = 455, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The securityAltId.
     */
    public java.lang.String getSecurityAltId() {
      java.lang.Object ref = securityAltId_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        securityAltId_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string security_alt_id = 1 [(.fix.tag) = 455, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The bytes for securityAltId.
     */
    public com.google.protobuf.ByteString
        getSecurityAltIdBytes() {
      java.lang.Object ref = securityAltId_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        securityAltId_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string security_alt_id = 1 [(.fix.tag) = 455, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The securityAltId to set.
     * @return This builder for chaining.
     */
    public Builder setSecurityAltId(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      securityAltId_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>string security_alt_id = 1 [(.fix.tag) = 455, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return This builder for chaining.
     */
    public Builder clearSecurityAltId() {
      
      securityAltId_ = getDefaultInstance().getSecurityAltId();
      onChanged();
      return this;
    }
    /**
     * <code>string security_alt_id = 1 [(.fix.tag) = 455, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The bytes for securityAltId to set.
     * @return This builder for chaining.
     */
    public Builder setSecurityAltIdBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      securityAltId_ = value;
      onChanged();
      return this;
    }

    private int securityAltIdSource_ = 0;
    /**
     * <code>.Common.SecurityAltIDSourceEnum security_alt_id_source = 2 [(.fix.tag) = 456, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The enum numeric value on the wire for securityAltIdSource.
     */
    public int getSecurityAltIdSourceValue() {
      return securityAltIdSource_;
    }
    /**
     * <code>.Common.SecurityAltIDSourceEnum security_alt_id_source = 2 [(.fix.tag) = 456, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The enum numeric value on the wire for securityAltIdSource to set.
     * @return This builder for chaining.
     */
    public Builder setSecurityAltIdSourceValue(int value) {
      securityAltIdSource_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>.Common.SecurityAltIDSourceEnum security_alt_id_source = 2 [(.fix.tag) = 456, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The securityAltIdSource.
     */
    public org.fixprotocol.components.SecurityAltIDSourceEnum getSecurityAltIdSource() {
      @SuppressWarnings("deprecation")
      org.fixprotocol.components.SecurityAltIDSourceEnum result = org.fixprotocol.components.SecurityAltIDSourceEnum.valueOf(securityAltIdSource_);
      return result == null ? org.fixprotocol.components.SecurityAltIDSourceEnum.UNRECOGNIZED : result;
    }
    /**
     * <code>.Common.SecurityAltIDSourceEnum security_alt_id_source = 2 [(.fix.tag) = 456, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The securityAltIdSource to set.
     * @return This builder for chaining.
     */
    public Builder setSecurityAltIdSource(org.fixprotocol.components.SecurityAltIDSourceEnum value) {
      if (value == null) {
        throw new NullPointerException();
      }
      
      securityAltIdSource_ = value.getNumber();
      onChanged();
      return this;
    }
    /**
     * <code>.Common.SecurityAltIDSourceEnum security_alt_id_source = 2 [(.fix.tag) = 456, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return This builder for chaining.
     */
    public Builder clearSecurityAltIdSource() {
      
      securityAltIdSource_ = 0;
      onChanged();
      return this;
    }
    @java.lang.Override
    public final Builder setUnknownFields(
        final com.google.protobuf.UnknownFieldSet unknownFields) {
      return super.setUnknownFields(unknownFields);
    }

    @java.lang.Override
    public final Builder mergeUnknownFields(
        final com.google.protobuf.UnknownFieldSet unknownFields) {
      return super.mergeUnknownFields(unknownFields);
    }


    // @@protoc_insertion_point(builder_scope:Common.SecAltIDGrp)
  }

  // @@protoc_insertion_point(class_scope:Common.SecAltIDGrp)
  private static final org.fixprotocol.components.SecAltIDGrp DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new org.fixprotocol.components.SecAltIDGrp();
  }

  public static org.fixprotocol.components.SecAltIDGrp getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<SecAltIDGrp>
      PARSER = new com.google.protobuf.AbstractParser<SecAltIDGrp>() {
    @java.lang.Override
    public SecAltIDGrp parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new SecAltIDGrp(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<SecAltIDGrp> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<SecAltIDGrp> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public org.fixprotocol.components.SecAltIDGrp getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}

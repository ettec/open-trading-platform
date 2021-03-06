// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf type {@code Common.PtysSubGrp}
 */
public  final class PtysSubGrp extends
    com.google.protobuf.GeneratedMessageV3 implements
    // @@protoc_insertion_point(message_implements:Common.PtysSubGrp)
    PtysSubGrpOrBuilder {
private static final long serialVersionUID = 0L;
  // Use PtysSubGrp.newBuilder() to construct.
  private PtysSubGrp(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
    super(builder);
  }
  private PtysSubGrp() {
    partySubId_ = "";
  }

  @java.lang.Override
  @SuppressWarnings({"unused"})
  protected java.lang.Object newInstance(
      UnusedPrivateParameter unused) {
    return new PtysSubGrp();
  }

  @java.lang.Override
  public final com.google.protobuf.UnknownFieldSet
  getUnknownFields() {
    return this.unknownFields;
  }
  private PtysSubGrp(
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

            partySubId_ = s;
            break;
          }
          case 16: {
            int rawValue = input.readEnum();
            partySubIdTypeUnionCase_ = 2;
            partySubIdTypeUnion_ = rawValue;
            break;
          }
          case 29: {
            partySubIdTypeUnionCase_ = 3;
            partySubIdTypeUnion_ = input.readFixed32();
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
    return org.fixprotocol.components.Common.internal_static_Common_PtysSubGrp_descriptor;
  }

  @java.lang.Override
  protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internalGetFieldAccessorTable() {
    return org.fixprotocol.components.Common.internal_static_Common_PtysSubGrp_fieldAccessorTable
        .ensureFieldAccessorsInitialized(
            org.fixprotocol.components.PtysSubGrp.class, org.fixprotocol.components.PtysSubGrp.Builder.class);
  }

  private int partySubIdTypeUnionCase_ = 0;
  private java.lang.Object partySubIdTypeUnion_;
  public enum PartySubIdTypeUnionCase
      implements com.google.protobuf.Internal.EnumLite,
          com.google.protobuf.AbstractMessage.InternalOneOfEnum {
    PARTY_SUB_ID_TYPE(2),
    PARTY_SUB_ID_TYPE_RESERVED4000PLUS(3),
    PARTYSUBIDTYPEUNION_NOT_SET(0);
    private final int value;
    private PartySubIdTypeUnionCase(int value) {
      this.value = value;
    }
    /**
     * @param value The number of the enum to look for.
     * @return The enum associated with the given number.
     * @deprecated Use {@link #forNumber(int)} instead.
     */
    @java.lang.Deprecated
    public static PartySubIdTypeUnionCase valueOf(int value) {
      return forNumber(value);
    }

    public static PartySubIdTypeUnionCase forNumber(int value) {
      switch (value) {
        case 2: return PARTY_SUB_ID_TYPE;
        case 3: return PARTY_SUB_ID_TYPE_RESERVED4000PLUS;
        case 0: return PARTYSUBIDTYPEUNION_NOT_SET;
        default: return null;
      }
    }
    public int getNumber() {
      return this.value;
    }
  };

  public PartySubIdTypeUnionCase
  getPartySubIdTypeUnionCase() {
    return PartySubIdTypeUnionCase.forNumber(
        partySubIdTypeUnionCase_);
  }

  public static final int PARTY_SUB_ID_FIELD_NUMBER = 1;
  private volatile java.lang.Object partySubId_;
  /**
   * <code>string party_sub_id = 1 [(.fix.tag) = 523, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The partySubId.
   */
  public java.lang.String getPartySubId() {
    java.lang.Object ref = partySubId_;
    if (ref instanceof java.lang.String) {
      return (java.lang.String) ref;
    } else {
      com.google.protobuf.ByteString bs = 
          (com.google.protobuf.ByteString) ref;
      java.lang.String s = bs.toStringUtf8();
      partySubId_ = s;
      return s;
    }
  }
  /**
   * <code>string party_sub_id = 1 [(.fix.tag) = 523, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for partySubId.
   */
  public com.google.protobuf.ByteString
      getPartySubIdBytes() {
    java.lang.Object ref = partySubId_;
    if (ref instanceof java.lang.String) {
      com.google.protobuf.ByteString b = 
          com.google.protobuf.ByteString.copyFromUtf8(
              (java.lang.String) ref);
      partySubId_ = b;
      return b;
    } else {
      return (com.google.protobuf.ByteString) ref;
    }
  }

  public static final int PARTY_SUB_ID_TYPE_FIELD_NUMBER = 2;
  /**
   * <code>.Common.PartySubIDTypeEnum party_sub_id_type = 2 [(.fix.tag) = 803, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The enum numeric value on the wire for partySubIdType.
   */
  public int getPartySubIdTypeValue() {
    if (partySubIdTypeUnionCase_ == 2) {
      return (java.lang.Integer) partySubIdTypeUnion_;
    }
    return 0;
  }
  /**
   * <code>.Common.PartySubIDTypeEnum party_sub_id_type = 2 [(.fix.tag) = 803, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The partySubIdType.
   */
  public org.fixprotocol.components.PartySubIDTypeEnum getPartySubIdType() {
    if (partySubIdTypeUnionCase_ == 2) {
      @SuppressWarnings("deprecation")
      org.fixprotocol.components.PartySubIDTypeEnum result = org.fixprotocol.components.PartySubIDTypeEnum.valueOf(
          (java.lang.Integer) partySubIdTypeUnion_);
      return result == null ? org.fixprotocol.components.PartySubIDTypeEnum.UNRECOGNIZED : result;
    }
    return org.fixprotocol.components.PartySubIDTypeEnum.PARTY_SUB_ID_TYPE_UNSPECIFIED;
  }

  public static final int PARTY_SUB_ID_TYPE_RESERVED4000PLUS_FIELD_NUMBER = 3;
  /**
   * <code>fixed32 party_sub_id_type_reserved4000plus = 3 [(.fix.tag) = 803, (.fix.type) = DATATYPE_RESERVED4000PLUS, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The partySubIdTypeReserved4000plus.
   */
  public int getPartySubIdTypeReserved4000Plus() {
    if (partySubIdTypeUnionCase_ == 3) {
      return (java.lang.Integer) partySubIdTypeUnion_;
    }
    return 0;
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
    if (!getPartySubIdBytes().isEmpty()) {
      com.google.protobuf.GeneratedMessageV3.writeString(output, 1, partySubId_);
    }
    if (partySubIdTypeUnionCase_ == 2) {
      output.writeEnum(2, ((java.lang.Integer) partySubIdTypeUnion_));
    }
    if (partySubIdTypeUnionCase_ == 3) {
      output.writeFixed32(
          3, (int)((java.lang.Integer) partySubIdTypeUnion_));
    }
    unknownFields.writeTo(output);
  }

  @java.lang.Override
  public int getSerializedSize() {
    int size = memoizedSize;
    if (size != -1) return size;

    size = 0;
    if (!getPartySubIdBytes().isEmpty()) {
      size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, partySubId_);
    }
    if (partySubIdTypeUnionCase_ == 2) {
      size += com.google.protobuf.CodedOutputStream
        .computeEnumSize(2, ((java.lang.Integer) partySubIdTypeUnion_));
    }
    if (partySubIdTypeUnionCase_ == 3) {
      size += com.google.protobuf.CodedOutputStream
        .computeFixed32Size(
            3, (int)((java.lang.Integer) partySubIdTypeUnion_));
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
    if (!(obj instanceof org.fixprotocol.components.PtysSubGrp)) {
      return super.equals(obj);
    }
    org.fixprotocol.components.PtysSubGrp other = (org.fixprotocol.components.PtysSubGrp) obj;

    if (!getPartySubId()
        .equals(other.getPartySubId())) return false;
    if (!getPartySubIdTypeUnionCase().equals(other.getPartySubIdTypeUnionCase())) return false;
    switch (partySubIdTypeUnionCase_) {
      case 2:
        if (getPartySubIdTypeValue()
            != other.getPartySubIdTypeValue()) return false;
        break;
      case 3:
        if (getPartySubIdTypeReserved4000Plus()
            != other.getPartySubIdTypeReserved4000Plus()) return false;
        break;
      case 0:
      default:
    }
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
    hash = (37 * hash) + PARTY_SUB_ID_FIELD_NUMBER;
    hash = (53 * hash) + getPartySubId().hashCode();
    switch (partySubIdTypeUnionCase_) {
      case 2:
        hash = (37 * hash) + PARTY_SUB_ID_TYPE_FIELD_NUMBER;
        hash = (53 * hash) + getPartySubIdTypeValue();
        break;
      case 3:
        hash = (37 * hash) + PARTY_SUB_ID_TYPE_RESERVED4000PLUS_FIELD_NUMBER;
        hash = (53 * hash) + getPartySubIdTypeReserved4000Plus();
        break;
      case 0:
      default:
    }
    hash = (29 * hash) + unknownFields.hashCode();
    memoizedHashCode = hash;
    return hash;
  }

  public static org.fixprotocol.components.PtysSubGrp parseFrom(
      java.nio.ByteBuffer data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(
      java.nio.ByteBuffer data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(
      com.google.protobuf.ByteString data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(
      com.google.protobuf.ByteString data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(byte[] data)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(
      byte[] data,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws com.google.protobuf.InvalidProtocolBufferException {
    return PARSER.parseFrom(data, extensionRegistry);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input, extensionRegistry);
  }
  public static org.fixprotocol.components.PtysSubGrp parseDelimitedFrom(java.io.InputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input);
  }
  public static org.fixprotocol.components.PtysSubGrp parseDelimitedFrom(
      java.io.InputStream input,
      com.google.protobuf.ExtensionRegistryLite extensionRegistry)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(
      com.google.protobuf.CodedInputStream input)
      throws java.io.IOException {
    return com.google.protobuf.GeneratedMessageV3
        .parseWithIOException(PARSER, input);
  }
  public static org.fixprotocol.components.PtysSubGrp parseFrom(
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
  public static Builder newBuilder(org.fixprotocol.components.PtysSubGrp prototype) {
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
   * Protobuf type {@code Common.PtysSubGrp}
   */
  public static final class Builder extends
      com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
      // @@protoc_insertion_point(builder_implements:Common.PtysSubGrp)
      org.fixprotocol.components.PtysSubGrpOrBuilder {
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return org.fixprotocol.components.Common.internal_static_Common_PtysSubGrp_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return org.fixprotocol.components.Common.internal_static_Common_PtysSubGrp_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              org.fixprotocol.components.PtysSubGrp.class, org.fixprotocol.components.PtysSubGrp.Builder.class);
    }

    // Construct using org.fixprotocol.components.PtysSubGrp.newBuilder()
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
      partySubId_ = "";

      partySubIdTypeUnionCase_ = 0;
      partySubIdTypeUnion_ = null;
      return this;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.Descriptor
        getDescriptorForType() {
      return org.fixprotocol.components.Common.internal_static_Common_PtysSubGrp_descriptor;
    }

    @java.lang.Override
    public org.fixprotocol.components.PtysSubGrp getDefaultInstanceForType() {
      return org.fixprotocol.components.PtysSubGrp.getDefaultInstance();
    }

    @java.lang.Override
    public org.fixprotocol.components.PtysSubGrp build() {
      org.fixprotocol.components.PtysSubGrp result = buildPartial();
      if (!result.isInitialized()) {
        throw newUninitializedMessageException(result);
      }
      return result;
    }

    @java.lang.Override
    public org.fixprotocol.components.PtysSubGrp buildPartial() {
      org.fixprotocol.components.PtysSubGrp result = new org.fixprotocol.components.PtysSubGrp(this);
      result.partySubId_ = partySubId_;
      if (partySubIdTypeUnionCase_ == 2) {
        result.partySubIdTypeUnion_ = partySubIdTypeUnion_;
      }
      if (partySubIdTypeUnionCase_ == 3) {
        result.partySubIdTypeUnion_ = partySubIdTypeUnion_;
      }
      result.partySubIdTypeUnionCase_ = partySubIdTypeUnionCase_;
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
      if (other instanceof org.fixprotocol.components.PtysSubGrp) {
        return mergeFrom((org.fixprotocol.components.PtysSubGrp)other);
      } else {
        super.mergeFrom(other);
        return this;
      }
    }

    public Builder mergeFrom(org.fixprotocol.components.PtysSubGrp other) {
      if (other == org.fixprotocol.components.PtysSubGrp.getDefaultInstance()) return this;
      if (!other.getPartySubId().isEmpty()) {
        partySubId_ = other.partySubId_;
        onChanged();
      }
      switch (other.getPartySubIdTypeUnionCase()) {
        case PARTY_SUB_ID_TYPE: {
          setPartySubIdTypeValue(other.getPartySubIdTypeValue());
          break;
        }
        case PARTY_SUB_ID_TYPE_RESERVED4000PLUS: {
          setPartySubIdTypeReserved4000Plus(other.getPartySubIdTypeReserved4000Plus());
          break;
        }
        case PARTYSUBIDTYPEUNION_NOT_SET: {
          break;
        }
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
      org.fixprotocol.components.PtysSubGrp parsedMessage = null;
      try {
        parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        parsedMessage = (org.fixprotocol.components.PtysSubGrp) e.getUnfinishedMessage();
        throw e.unwrapIOException();
      } finally {
        if (parsedMessage != null) {
          mergeFrom(parsedMessage);
        }
      }
      return this;
    }
    private int partySubIdTypeUnionCase_ = 0;
    private java.lang.Object partySubIdTypeUnion_;
    public PartySubIdTypeUnionCase
        getPartySubIdTypeUnionCase() {
      return PartySubIdTypeUnionCase.forNumber(
          partySubIdTypeUnionCase_);
    }

    public Builder clearPartySubIdTypeUnion() {
      partySubIdTypeUnionCase_ = 0;
      partySubIdTypeUnion_ = null;
      onChanged();
      return this;
    }


    private java.lang.Object partySubId_ = "";
    /**
     * <code>string party_sub_id = 1 [(.fix.tag) = 523, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The partySubId.
     */
    public java.lang.String getPartySubId() {
      java.lang.Object ref = partySubId_;
      if (!(ref instanceof java.lang.String)) {
        com.google.protobuf.ByteString bs =
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        partySubId_ = s;
        return s;
      } else {
        return (java.lang.String) ref;
      }
    }
    /**
     * <code>string party_sub_id = 1 [(.fix.tag) = 523, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The bytes for partySubId.
     */
    public com.google.protobuf.ByteString
        getPartySubIdBytes() {
      java.lang.Object ref = partySubId_;
      if (ref instanceof String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        partySubId_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }
    /**
     * <code>string party_sub_id = 1 [(.fix.tag) = 523, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The partySubId to set.
     * @return This builder for chaining.
     */
    public Builder setPartySubId(
        java.lang.String value) {
      if (value == null) {
    throw new NullPointerException();
  }
  
      partySubId_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>string party_sub_id = 1 [(.fix.tag) = 523, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return This builder for chaining.
     */
    public Builder clearPartySubId() {
      
      partySubId_ = getDefaultInstance().getPartySubId();
      onChanged();
      return this;
    }
    /**
     * <code>string party_sub_id = 1 [(.fix.tag) = 523, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The bytes for partySubId to set.
     * @return This builder for chaining.
     */
    public Builder setPartySubIdBytes(
        com.google.protobuf.ByteString value) {
      if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
      
      partySubId_ = value;
      onChanged();
      return this;
    }

    /**
     * <code>.Common.PartySubIDTypeEnum party_sub_id_type = 2 [(.fix.tag) = 803, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The enum numeric value on the wire for partySubIdType.
     */
    public int getPartySubIdTypeValue() {
      if (partySubIdTypeUnionCase_ == 2) {
        return ((java.lang.Integer) partySubIdTypeUnion_).intValue();
      }
      return 0;
    }
    /**
     * <code>.Common.PartySubIDTypeEnum party_sub_id_type = 2 [(.fix.tag) = 803, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The enum numeric value on the wire for partySubIdType to set.
     * @return This builder for chaining.
     */
    public Builder setPartySubIdTypeValue(int value) {
      partySubIdTypeUnionCase_ = 2;
      partySubIdTypeUnion_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>.Common.PartySubIDTypeEnum party_sub_id_type = 2 [(.fix.tag) = 803, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The partySubIdType.
     */
    public org.fixprotocol.components.PartySubIDTypeEnum getPartySubIdType() {
      if (partySubIdTypeUnionCase_ == 2) {
        @SuppressWarnings("deprecation")
        org.fixprotocol.components.PartySubIDTypeEnum result = org.fixprotocol.components.PartySubIDTypeEnum.valueOf(
            (java.lang.Integer) partySubIdTypeUnion_);
        return result == null ? org.fixprotocol.components.PartySubIDTypeEnum.UNRECOGNIZED : result;
      }
      return org.fixprotocol.components.PartySubIDTypeEnum.PARTY_SUB_ID_TYPE_UNSPECIFIED;
    }
    /**
     * <code>.Common.PartySubIDTypeEnum party_sub_id_type = 2 [(.fix.tag) = 803, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The partySubIdType to set.
     * @return This builder for chaining.
     */
    public Builder setPartySubIdType(org.fixprotocol.components.PartySubIDTypeEnum value) {
      if (value == null) {
        throw new NullPointerException();
      }
      partySubIdTypeUnionCase_ = 2;
      partySubIdTypeUnion_ = value.getNumber();
      onChanged();
      return this;
    }
    /**
     * <code>.Common.PartySubIDTypeEnum party_sub_id_type = 2 [(.fix.tag) = 803, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return This builder for chaining.
     */
    public Builder clearPartySubIdType() {
      if (partySubIdTypeUnionCase_ == 2) {
        partySubIdTypeUnionCase_ = 0;
        partySubIdTypeUnion_ = null;
        onChanged();
      }
      return this;
    }

    /**
     * <code>fixed32 party_sub_id_type_reserved4000plus = 3 [(.fix.tag) = 803, (.fix.type) = DATATYPE_RESERVED4000PLUS, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return The partySubIdTypeReserved4000plus.
     */
    public int getPartySubIdTypeReserved4000Plus() {
      if (partySubIdTypeUnionCase_ == 3) {
        return (java.lang.Integer) partySubIdTypeUnion_;
      }
      return 0;
    }
    /**
     * <code>fixed32 party_sub_id_type_reserved4000plus = 3 [(.fix.tag) = 803, (.fix.type) = DATATYPE_RESERVED4000PLUS, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @param value The partySubIdTypeReserved4000plus to set.
     * @return This builder for chaining.
     */
    public Builder setPartySubIdTypeReserved4000Plus(int value) {
      partySubIdTypeUnionCase_ = 3;
      partySubIdTypeUnion_ = value;
      onChanged();
      return this;
    }
    /**
     * <code>fixed32 party_sub_id_type_reserved4000plus = 3 [(.fix.tag) = 803, (.fix.type) = DATATYPE_RESERVED4000PLUS, (.fix.field_added) = VERSION_FIX_4_4];</code>
     * @return This builder for chaining.
     */
    public Builder clearPartySubIdTypeReserved4000Plus() {
      if (partySubIdTypeUnionCase_ == 3) {
        partySubIdTypeUnionCase_ = 0;
        partySubIdTypeUnion_ = null;
        onChanged();
      }
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


    // @@protoc_insertion_point(builder_scope:Common.PtysSubGrp)
  }

  // @@protoc_insertion_point(class_scope:Common.PtysSubGrp)
  private static final org.fixprotocol.components.PtysSubGrp DEFAULT_INSTANCE;
  static {
    DEFAULT_INSTANCE = new org.fixprotocol.components.PtysSubGrp();
  }

  public static org.fixprotocol.components.PtysSubGrp getDefaultInstance() {
    return DEFAULT_INSTANCE;
  }

  private static final com.google.protobuf.Parser<PtysSubGrp>
      PARSER = new com.google.protobuf.AbstractParser<PtysSubGrp>() {
    @java.lang.Override
    public PtysSubGrp parsePartialFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return new PtysSubGrp(input, extensionRegistry);
    }
  };

  public static com.google.protobuf.Parser<PtysSubGrp> parser() {
    return PARSER;
  }

  @java.lang.Override
  public com.google.protobuf.Parser<PtysSubGrp> getParserForType() {
    return PARSER;
  }

  @java.lang.Override
  public org.fixprotocol.components.PtysSubGrp getDefaultInstanceForType() {
    return DEFAULT_INSTANCE;
  }

}


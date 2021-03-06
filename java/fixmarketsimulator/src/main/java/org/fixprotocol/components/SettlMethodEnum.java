// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf enum {@code Common.SettlMethodEnum}
 */
public enum SettlMethodEnum
    implements com.google.protobuf.ProtocolMessageEnum {
  /**
   * <code>SETTL_METHOD_UNSPECIFIED = 0;</code>
   */
  SETTL_METHOD_UNSPECIFIED(0),
  /**
   * <code>SETTL_METHOD_CASH_SETTLEMENT_REQUIRED = 1 [(.fix.enum_value) = "C", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  SETTL_METHOD_CASH_SETTLEMENT_REQUIRED(1),
  /**
   * <code>SETTL_METHOD_PHYSICAL_SETTLEMENT_REQUIRED = 2 [(.fix.enum_value) = "P", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  SETTL_METHOD_PHYSICAL_SETTLEMENT_REQUIRED(2),
  UNRECOGNIZED(-1),
  ;

  /**
   * <code>SETTL_METHOD_UNSPECIFIED = 0;</code>
   */
  public static final int SETTL_METHOD_UNSPECIFIED_VALUE = 0;
  /**
   * <code>SETTL_METHOD_CASH_SETTLEMENT_REQUIRED = 1 [(.fix.enum_value) = "C", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  public static final int SETTL_METHOD_CASH_SETTLEMENT_REQUIRED_VALUE = 1;
  /**
   * <code>SETTL_METHOD_PHYSICAL_SETTLEMENT_REQUIRED = 2 [(.fix.enum_value) = "P", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  public static final int SETTL_METHOD_PHYSICAL_SETTLEMENT_REQUIRED_VALUE = 2;


  public final int getNumber() {
    if (this == UNRECOGNIZED) {
      throw new java.lang.IllegalArgumentException(
          "Can't get the number of an unknown enum value.");
    }
    return value;
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   * @deprecated Use {@link #forNumber(int)} instead.
   */
  @java.lang.Deprecated
  public static SettlMethodEnum valueOf(int value) {
    return forNumber(value);
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   */
  public static SettlMethodEnum forNumber(int value) {
    switch (value) {
      case 0: return SETTL_METHOD_UNSPECIFIED;
      case 1: return SETTL_METHOD_CASH_SETTLEMENT_REQUIRED;
      case 2: return SETTL_METHOD_PHYSICAL_SETTLEMENT_REQUIRED;
      default: return null;
    }
  }

  public static com.google.protobuf.Internal.EnumLiteMap<SettlMethodEnum>
      internalGetValueMap() {
    return internalValueMap;
  }
  private static final com.google.protobuf.Internal.EnumLiteMap<
      SettlMethodEnum> internalValueMap =
        new com.google.protobuf.Internal.EnumLiteMap<SettlMethodEnum>() {
          public SettlMethodEnum findValueByNumber(int number) {
            return SettlMethodEnum.forNumber(number);
          }
        };

  public final com.google.protobuf.Descriptors.EnumValueDescriptor
      getValueDescriptor() {
    return getDescriptor().getValues().get(ordinal());
  }
  public final com.google.protobuf.Descriptors.EnumDescriptor
      getDescriptorForType() {
    return getDescriptor();
  }
  public static final com.google.protobuf.Descriptors.EnumDescriptor
      getDescriptor() {
    return org.fixprotocol.components.Common.getDescriptor().getEnumTypes().get(24);
  }

  private static final SettlMethodEnum[] VALUES = values();

  public static SettlMethodEnum valueOf(
      com.google.protobuf.Descriptors.EnumValueDescriptor desc) {
    if (desc.getType() != getDescriptor()) {
      throw new java.lang.IllegalArgumentException(
        "EnumValueDescriptor is not for this type.");
    }
    if (desc.getIndex() == -1) {
      return UNRECOGNIZED;
    }
    return VALUES[desc.getIndex()];
  }

  private final int value;

  private SettlMethodEnum(int value) {
    this.value = value;
  }

  // @@protoc_insertion_point(enum_scope:Common.SettlMethodEnum)
}


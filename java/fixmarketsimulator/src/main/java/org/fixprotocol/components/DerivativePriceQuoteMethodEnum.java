// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf enum {@code Common.DerivativePriceQuoteMethodEnum}
 */
public enum DerivativePriceQuoteMethodEnum
    implements com.google.protobuf.ProtocolMessageEnum {
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_UNSPECIFIED = 0;</code>
   */
  DERIVATIVE_PRICE_QUOTE_METHOD_UNSPECIFIED(0),
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_INDEX = 1 [(.fix.enum_value) = "INX", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  DERIVATIVE_PRICE_QUOTE_METHOD_INDEX(1),
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_INTEREST_RATE_INDEX = 2 [(.fix.enum_value) = "INT", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  DERIVATIVE_PRICE_QUOTE_METHOD_INTEREST_RATE_INDEX(2),
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_STANDARD = 3 [(.fix.enum_value) = "STD", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  DERIVATIVE_PRICE_QUOTE_METHOD_STANDARD(3),
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_PERCENT_OF_PAR = 4 [(.fix.enum_value) = "PCTPAR", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 83];</code>
   */
  DERIVATIVE_PRICE_QUOTE_METHOD_PERCENT_OF_PAR(4),
  UNRECOGNIZED(-1),
  ;

  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_UNSPECIFIED = 0;</code>
   */
  public static final int DERIVATIVE_PRICE_QUOTE_METHOD_UNSPECIFIED_VALUE = 0;
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_INDEX = 1 [(.fix.enum_value) = "INX", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  public static final int DERIVATIVE_PRICE_QUOTE_METHOD_INDEX_VALUE = 1;
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_INTEREST_RATE_INDEX = 2 [(.fix.enum_value) = "INT", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  public static final int DERIVATIVE_PRICE_QUOTE_METHOD_INTEREST_RATE_INDEX_VALUE = 2;
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_STANDARD = 3 [(.fix.enum_value) = "STD", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 52];</code>
   */
  public static final int DERIVATIVE_PRICE_QUOTE_METHOD_STANDARD_VALUE = 3;
  /**
   * <code>DERIVATIVE_PRICE_QUOTE_METHOD_PERCENT_OF_PAR = 4 [(.fix.enum_value) = "PCTPAR", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 83];</code>
   */
  public static final int DERIVATIVE_PRICE_QUOTE_METHOD_PERCENT_OF_PAR_VALUE = 4;


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
  public static DerivativePriceQuoteMethodEnum valueOf(int value) {
    return forNumber(value);
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   */
  public static DerivativePriceQuoteMethodEnum forNumber(int value) {
    switch (value) {
      case 0: return DERIVATIVE_PRICE_QUOTE_METHOD_UNSPECIFIED;
      case 1: return DERIVATIVE_PRICE_QUOTE_METHOD_INDEX;
      case 2: return DERIVATIVE_PRICE_QUOTE_METHOD_INTEREST_RATE_INDEX;
      case 3: return DERIVATIVE_PRICE_QUOTE_METHOD_STANDARD;
      case 4: return DERIVATIVE_PRICE_QUOTE_METHOD_PERCENT_OF_PAR;
      default: return null;
    }
  }

  public static com.google.protobuf.Internal.EnumLiteMap<DerivativePriceQuoteMethodEnum>
      internalGetValueMap() {
    return internalValueMap;
  }
  private static final com.google.protobuf.Internal.EnumLiteMap<
      DerivativePriceQuoteMethodEnum> internalValueMap =
        new com.google.protobuf.Internal.EnumLiteMap<DerivativePriceQuoteMethodEnum>() {
          public DerivativePriceQuoteMethodEnum findValueByNumber(int number) {
            return DerivativePriceQuoteMethodEnum.forNumber(number);
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
    return org.fixprotocol.components.Common.getDescriptor().getEnumTypes().get(186);
  }

  private static final DerivativePriceQuoteMethodEnum[] VALUES = values();

  public static DerivativePriceQuoteMethodEnum valueOf(
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

  private DerivativePriceQuoteMethodEnum(int value) {
    this.value = value;
  }

  // @@protoc_insertion_point(enum_scope:Common.DerivativePriceQuoteMethodEnum)
}


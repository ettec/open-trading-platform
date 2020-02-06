// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf enum {@code Common.ProductEnum}
 */
public enum ProductEnum
    implements com.google.protobuf.ProtocolMessageEnum {
  /**
   * <code>PRODUCT_UNSPECIFIED = 0;</code>
   */
  PRODUCT_UNSPECIFIED(0),
  /**
   * <code>PRODUCT_AGENCY = 1 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_AGENCY(1),
  /**
   * <code>PRODUCT_COMMODITY = 2 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_COMMODITY(2),
  /**
   * <code>PRODUCT_CORPORATE = 3 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_CORPORATE(3),
  /**
   * <code>PRODUCT_CURRENCY = 4 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_CURRENCY(4),
  /**
   * <code>PRODUCT_EQUITY = 5 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_EQUITY(5),
  /**
   * <code>PRODUCT_GOVERNMENT = 6 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_GOVERNMENT(6),
  /**
   * <code>PRODUCT_INDEX = 7 [(.fix.enum_value) = "7", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_INDEX(7),
  /**
   * <code>PRODUCT_LOAN = 8 [(.fix.enum_value) = "8", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_LOAN(8),
  /**
   * <code>PRODUCT_MONEYMARKET = 9 [(.fix.enum_value) = "9", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_MONEYMARKET(9),
  /**
   * <code>PRODUCT_MORTGAGE = 10 [(.fix.enum_value) = "10", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_MORTGAGE(10),
  /**
   * <code>PRODUCT_MUNICIPAL = 11 [(.fix.enum_value) = "11", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_MUNICIPAL(11),
  /**
   * <code>PRODUCT_OTHER = 12 [(.fix.enum_value) = "12", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  PRODUCT_OTHER(12),
  /**
   * <code>PRODUCT_FINANCING = 13 [(.fix.enum_value) = "13", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  PRODUCT_FINANCING(13),
  UNRECOGNIZED(-1),
  ;

  /**
   * <code>PRODUCT_UNSPECIFIED = 0;</code>
   */
  public static final int PRODUCT_UNSPECIFIED_VALUE = 0;
  /**
   * <code>PRODUCT_AGENCY = 1 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_AGENCY_VALUE = 1;
  /**
   * <code>PRODUCT_COMMODITY = 2 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_COMMODITY_VALUE = 2;
  /**
   * <code>PRODUCT_CORPORATE = 3 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_CORPORATE_VALUE = 3;
  /**
   * <code>PRODUCT_CURRENCY = 4 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_CURRENCY_VALUE = 4;
  /**
   * <code>PRODUCT_EQUITY = 5 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_EQUITY_VALUE = 5;
  /**
   * <code>PRODUCT_GOVERNMENT = 6 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_GOVERNMENT_VALUE = 6;
  /**
   * <code>PRODUCT_INDEX = 7 [(.fix.enum_value) = "7", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_INDEX_VALUE = 7;
  /**
   * <code>PRODUCT_LOAN = 8 [(.fix.enum_value) = "8", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_LOAN_VALUE = 8;
  /**
   * <code>PRODUCT_MONEYMARKET = 9 [(.fix.enum_value) = "9", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_MONEYMARKET_VALUE = 9;
  /**
   * <code>PRODUCT_MORTGAGE = 10 [(.fix.enum_value) = "10", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_MORTGAGE_VALUE = 10;
  /**
   * <code>PRODUCT_MUNICIPAL = 11 [(.fix.enum_value) = "11", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_MUNICIPAL_VALUE = 11;
  /**
   * <code>PRODUCT_OTHER = 12 [(.fix.enum_value) = "12", (.fix.enum_added) = VERSION_FIX_4_3];</code>
   */
  public static final int PRODUCT_OTHER_VALUE = 12;
  /**
   * <code>PRODUCT_FINANCING = 13 [(.fix.enum_value) = "13", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int PRODUCT_FINANCING_VALUE = 13;


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
  public static ProductEnum valueOf(int value) {
    return forNumber(value);
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   */
  public static ProductEnum forNumber(int value) {
    switch (value) {
      case 0: return PRODUCT_UNSPECIFIED;
      case 1: return PRODUCT_AGENCY;
      case 2: return PRODUCT_COMMODITY;
      case 3: return PRODUCT_CORPORATE;
      case 4: return PRODUCT_CURRENCY;
      case 5: return PRODUCT_EQUITY;
      case 6: return PRODUCT_GOVERNMENT;
      case 7: return PRODUCT_INDEX;
      case 8: return PRODUCT_LOAN;
      case 9: return PRODUCT_MONEYMARKET;
      case 10: return PRODUCT_MORTGAGE;
      case 11: return PRODUCT_MUNICIPAL;
      case 12: return PRODUCT_OTHER;
      case 13: return PRODUCT_FINANCING;
      default: return null;
    }
  }

  public static com.google.protobuf.Internal.EnumLiteMap<ProductEnum>
      internalGetValueMap() {
    return internalValueMap;
  }
  private static final com.google.protobuf.Internal.EnumLiteMap<
      ProductEnum> internalValueMap =
        new com.google.protobuf.Internal.EnumLiteMap<ProductEnum>() {
          public ProductEnum findValueByNumber(int number) {
            return ProductEnum.forNumber(number);
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
    return org.fixprotocol.components.Common.getDescriptor().getEnumTypes().get(12);
  }

  private static final ProductEnum[] VALUES = values();

  public static ProductEnum valueOf(
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

  private ProductEnum(int value) {
    this.value = value;
  }

  // @@protoc_insertion_point(enum_scope:Common.ProductEnum)
}

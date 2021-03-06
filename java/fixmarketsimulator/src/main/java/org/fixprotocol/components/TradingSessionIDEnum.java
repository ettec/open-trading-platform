// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf enum {@code Common.TradingSessionIDEnum}
 */
public enum TradingSessionIDEnum
    implements com.google.protobuf.ProtocolMessageEnum {
  /**
   * <code>TRADING_SESSION_ID_UNSPECIFIED = 0;</code>
   */
  TRADING_SESSION_ID_UNSPECIFIED(0),
  /**
   * <code>TRADING_SESSION_ID_AFTERNOON = 1 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  TRADING_SESSION_ID_AFTERNOON(1),
  /**
   * <code>TRADING_SESSION_ID_AFTER_HOURS = 2 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  TRADING_SESSION_ID_AFTER_HOURS(2),
  /**
   * <code>TRADING_SESSION_ID_DAY = 3 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  TRADING_SESSION_ID_DAY(3),
  /**
   * <code>TRADING_SESSION_ID_EVENING = 4 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  TRADING_SESSION_ID_EVENING(4),
  /**
   * <code>TRADING_SESSION_ID_HALF_DAY = 5 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  TRADING_SESSION_ID_HALF_DAY(5),
  /**
   * <code>TRADING_SESSION_ID_MORNING = 6 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  TRADING_SESSION_ID_MORNING(6),
  UNRECOGNIZED(-1),
  ;

  /**
   * <code>TRADING_SESSION_ID_UNSPECIFIED = 0;</code>
   */
  public static final int TRADING_SESSION_ID_UNSPECIFIED_VALUE = 0;
  /**
   * <code>TRADING_SESSION_ID_AFTERNOON = 1 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  public static final int TRADING_SESSION_ID_AFTERNOON_VALUE = 1;
  /**
   * <code>TRADING_SESSION_ID_AFTER_HOURS = 2 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  public static final int TRADING_SESSION_ID_AFTER_HOURS_VALUE = 2;
  /**
   * <code>TRADING_SESSION_ID_DAY = 3 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  public static final int TRADING_SESSION_ID_DAY_VALUE = 3;
  /**
   * <code>TRADING_SESSION_ID_EVENING = 4 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  public static final int TRADING_SESSION_ID_EVENING_VALUE = 4;
  /**
   * <code>TRADING_SESSION_ID_HALF_DAY = 5 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  public static final int TRADING_SESSION_ID_HALF_DAY_VALUE = 5;
  /**
   * <code>TRADING_SESSION_ID_MORNING = 6 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 58];</code>
   */
  public static final int TRADING_SESSION_ID_MORNING_VALUE = 6;


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
  public static TradingSessionIDEnum valueOf(int value) {
    return forNumber(value);
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   */
  public static TradingSessionIDEnum forNumber(int value) {
    switch (value) {
      case 0: return TRADING_SESSION_ID_UNSPECIFIED;
      case 1: return TRADING_SESSION_ID_AFTERNOON;
      case 2: return TRADING_SESSION_ID_AFTER_HOURS;
      case 3: return TRADING_SESSION_ID_DAY;
      case 4: return TRADING_SESSION_ID_EVENING;
      case 5: return TRADING_SESSION_ID_HALF_DAY;
      case 6: return TRADING_SESSION_ID_MORNING;
      default: return null;
    }
  }

  public static com.google.protobuf.Internal.EnumLiteMap<TradingSessionIDEnum>
      internalGetValueMap() {
    return internalValueMap;
  }
  private static final com.google.protobuf.Internal.EnumLiteMap<
      TradingSessionIDEnum> internalValueMap =
        new com.google.protobuf.Internal.EnumLiteMap<TradingSessionIDEnum>() {
          public TradingSessionIDEnum findValueByNumber(int number) {
            return TradingSessionIDEnum.forNumber(number);
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
    return org.fixprotocol.components.Common.getDescriptor().getEnumTypes().get(111);
  }

  private static final TradingSessionIDEnum[] VALUES = values();

  public static TradingSessionIDEnum valueOf(
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

  private TradingSessionIDEnum(int value) {
    this.value = value;
  }

  // @@protoc_insertion_point(enum_scope:Common.TradingSessionIDEnum)
}


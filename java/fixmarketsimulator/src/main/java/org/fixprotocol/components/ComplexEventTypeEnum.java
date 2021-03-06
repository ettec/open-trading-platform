// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf enum {@code Common.ComplexEventTypeEnum}
 */
public enum ComplexEventTypeEnum
    implements com.google.protobuf.ProtocolMessageEnum {
  /**
   * <code>COMPLEX_EVENT_TYPE_UNSPECIFIED = 0;</code>
   */
  COMPLEX_EVENT_TYPE_UNSPECIFIED(0),
  /**
   * <code>COMPLEX_EVENT_TYPE_CAPPED = 1 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_CAPPED(1),
  /**
   * <code>COMPLEX_EVENT_TYPE_KNOCK_IN_UP = 2 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_KNOCK_IN_UP(2),
  /**
   * <code>COMPLEX_EVENT_TYPE_KNOCK_OUT_DOWN = 3 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_KNOCK_OUT_DOWN(3),
  /**
   * <code>COMPLEX_EVENT_TYPE_KNOCK_OUT_UP = 4 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_KNOCK_OUT_UP(4),
  /**
   * <code>COMPLEX_EVENT_TYPE_KOCK_IN_DOWN = 5 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_KOCK_IN_DOWN(5),
  /**
   * <code>COMPLEX_EVENT_TYPE_RESET_BARRIER = 6 [(.fix.enum_value) = "8", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_RESET_BARRIER(6),
  /**
   * <code>COMPLEX_EVENT_TYPE_ROLLING_BARRIER = 7 [(.fix.enum_value) = "9", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_ROLLING_BARRIER(7),
  /**
   * <code>COMPLEX_EVENT_TYPE_TRIGGER = 8 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_TRIGGER(8),
  /**
   * <code>COMPLEX_EVENT_TYPE_UNDERLYING = 9 [(.fix.enum_value) = "7", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  COMPLEX_EVENT_TYPE_UNDERLYING(9),
  UNRECOGNIZED(-1),
  ;

  /**
   * <code>COMPLEX_EVENT_TYPE_UNSPECIFIED = 0;</code>
   */
  public static final int COMPLEX_EVENT_TYPE_UNSPECIFIED_VALUE = 0;
  /**
   * <code>COMPLEX_EVENT_TYPE_CAPPED = 1 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_CAPPED_VALUE = 1;
  /**
   * <code>COMPLEX_EVENT_TYPE_KNOCK_IN_UP = 2 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_KNOCK_IN_UP_VALUE = 2;
  /**
   * <code>COMPLEX_EVENT_TYPE_KNOCK_OUT_DOWN = 3 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_KNOCK_OUT_DOWN_VALUE = 3;
  /**
   * <code>COMPLEX_EVENT_TYPE_KNOCK_OUT_UP = 4 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_KNOCK_OUT_UP_VALUE = 4;
  /**
   * <code>COMPLEX_EVENT_TYPE_KOCK_IN_DOWN = 5 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_KOCK_IN_DOWN_VALUE = 5;
  /**
   * <code>COMPLEX_EVENT_TYPE_RESET_BARRIER = 6 [(.fix.enum_value) = "8", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_RESET_BARRIER_VALUE = 6;
  /**
   * <code>COMPLEX_EVENT_TYPE_ROLLING_BARRIER = 7 [(.fix.enum_value) = "9", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_ROLLING_BARRIER_VALUE = 7;
  /**
   * <code>COMPLEX_EVENT_TYPE_TRIGGER = 8 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_TRIGGER_VALUE = 8;
  /**
   * <code>COMPLEX_EVENT_TYPE_UNDERLYING = 9 [(.fix.enum_value) = "7", (.fix.enum_added) = VERSION_FIX_5_0SP1, (.fix.enum_added_ep) = 92];</code>
   */
  public static final int COMPLEX_EVENT_TYPE_UNDERLYING_VALUE = 9;


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
  public static ComplexEventTypeEnum valueOf(int value) {
    return forNumber(value);
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   */
  public static ComplexEventTypeEnum forNumber(int value) {
    switch (value) {
      case 0: return COMPLEX_EVENT_TYPE_UNSPECIFIED;
      case 1: return COMPLEX_EVENT_TYPE_CAPPED;
      case 2: return COMPLEX_EVENT_TYPE_KNOCK_IN_UP;
      case 3: return COMPLEX_EVENT_TYPE_KNOCK_OUT_DOWN;
      case 4: return COMPLEX_EVENT_TYPE_KNOCK_OUT_UP;
      case 5: return COMPLEX_EVENT_TYPE_KOCK_IN_DOWN;
      case 6: return COMPLEX_EVENT_TYPE_RESET_BARRIER;
      case 7: return COMPLEX_EVENT_TYPE_ROLLING_BARRIER;
      case 8: return COMPLEX_EVENT_TYPE_TRIGGER;
      case 9: return COMPLEX_EVENT_TYPE_UNDERLYING;
      default: return null;
    }
  }

  public static com.google.protobuf.Internal.EnumLiteMap<ComplexEventTypeEnum>
      internalGetValueMap() {
    return internalValueMap;
  }
  private static final com.google.protobuf.Internal.EnumLiteMap<
      ComplexEventTypeEnum> internalValueMap =
        new com.google.protobuf.Internal.EnumLiteMap<ComplexEventTypeEnum>() {
          public ComplexEventTypeEnum findValueByNumber(int number) {
            return ComplexEventTypeEnum.forNumber(number);
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
    return org.fixprotocol.components.Common.getDescriptor().getEnumTypes().get(203);
  }

  private static final ComplexEventTypeEnum[] VALUES = values();

  public static ComplexEventTypeEnum valueOf(
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

  private ComplexEventTypeEnum(int value) {
    this.value = value;
  }

  // @@protoc_insertion_point(enum_scope:Common.ComplexEventTypeEnum)
}


// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

/**
 * Protobuf enum {@code Common.SettlPartySubIDTypeEnum}
 */
public enum SettlPartySubIDTypeEnum
    implements com.google.protobuf.ProtocolMessageEnum {
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_UNSPECIFIED = 0;</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_UNSPECIFIED(0),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_APPLICATION = 1 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_APPLICATION(1),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_BIC = 2 [(.fix.enum_value) = "16", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_BIC(2),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NAME = 3 [(.fix.enum_value) = "23", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NAME(3),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NUMBER = 4 [(.fix.enum_value) = "15", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NUMBER(4),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CONTACT_NAME = 5 [(.fix.enum_value) = "9", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_CONTACT_NAME(5),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CSDPARTICIPANT_MEMBER_CODE = 6 [(.fix.enum_value) = "17", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_CSDPARTICIPANT_MEMBER_CODE(6),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_DEPARTMENT = 7 [(.fix.enum_value) = "24", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_DEPARTMENT(7),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_EMAIL_ADDRESS = 8 [(.fix.enum_value) = "8", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_EMAIL_ADDRESS(8),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FAX_NUMBER = 9 [(.fix.enum_value) = "21", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_FAX_NUMBER(9),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FIRM = 10 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_FIRM(10),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FULL_LEGAL_NAME_OF_FIRM = 11 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_FULL_LEGAL_NAME_OF_FIRM(11),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FUND_ACCOUNT_NAME = 12 [(.fix.enum_value) = "19", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_FUND_ACCOUNT_NAME(12),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_LOCATION_DESK = 13 [(.fix.enum_value) = "25", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_LOCATION_DESK(13),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_PERSON = 14 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_PERSON(14),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_PHONE_NUMBER = 15 [(.fix.enum_value) = "7", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_PHONE_NUMBER(15),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_POSITION_ACCOUNT_TYPE = 16 [(.fix.enum_value) = "26", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_POSITION_ACCOUNT_TYPE(16),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_POSTAL_ADDRESS = 17 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_POSTAL_ADDRESS(17),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS = 18 [(.fix.enum_value) = "18", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS(18),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS_FOR_CONFIRMATION = 19 [(.fix.enum_value) = "12", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS_FOR_CONFIRMATION(19),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NAME = 20 [(.fix.enum_value) = "14", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NAME(20),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NUMBER = 21 [(.fix.enum_value) = "11", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NUMBER(21),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGULATORY_STATUS = 22 [(.fix.enum_value) = "13", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_REGULATORY_STATUS(22),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NAME = 23 [(.fix.enum_value) = "22", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NAME(23),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NUMBER = 24 [(.fix.enum_value) = "10", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NUMBER(24),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SYSTEM = 25 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_SYSTEM(25),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_TELEX_NUMBER = 26 [(.fix.enum_value) = "20", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_TELEX_NUMBER(26),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SECURITY_LOCATE_ID = 27 [(.fix.enum_value) = "27", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 1];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_SECURITY_LOCATE_ID(27),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_ELIGIBLE_COUNTERPARTY = 28 [(.fix.enum_value) = "29", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_ELIGIBLE_COUNTERPARTY(28),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_EXECUTION_VENUE = 29 [(.fix.enum_value) = "32", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_EXECUTION_VENUE(29),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_LOCATION = 30 [(.fix.enum_value) = "31", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_LOCATION(30),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_MARKET_MAKER = 31 [(.fix.enum_value) = "28", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_MARKET_MAKER(31),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_PROFESSIONAL_CLIENT = 32 [(.fix.enum_value) = "30", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_PROFESSIONAL_CLIENT(32),
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CURRENCY_DELIVERY_IDENTIFIER = 33 [(.fix.enum_value) = "33", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 44];</code>
   */
  SETTL_PARTY_SUB_ID_TYPE_CURRENCY_DELIVERY_IDENTIFIER(33),
  UNRECOGNIZED(-1),
  ;

  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_UNSPECIFIED = 0;</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_UNSPECIFIED_VALUE = 0;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_APPLICATION = 1 [(.fix.enum_value) = "4", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_APPLICATION_VALUE = 1;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_BIC = 2 [(.fix.enum_value) = "16", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_BIC_VALUE = 2;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NAME = 3 [(.fix.enum_value) = "23", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NAME_VALUE = 3;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NUMBER = 4 [(.fix.enum_value) = "15", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NUMBER_VALUE = 4;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CONTACT_NAME = 5 [(.fix.enum_value) = "9", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_CONTACT_NAME_VALUE = 5;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CSDPARTICIPANT_MEMBER_CODE = 6 [(.fix.enum_value) = "17", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_CSDPARTICIPANT_MEMBER_CODE_VALUE = 6;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_DEPARTMENT = 7 [(.fix.enum_value) = "24", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_DEPARTMENT_VALUE = 7;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_EMAIL_ADDRESS = 8 [(.fix.enum_value) = "8", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_EMAIL_ADDRESS_VALUE = 8;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FAX_NUMBER = 9 [(.fix.enum_value) = "21", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_FAX_NUMBER_VALUE = 9;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FIRM = 10 [(.fix.enum_value) = "1", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_FIRM_VALUE = 10;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FULL_LEGAL_NAME_OF_FIRM = 11 [(.fix.enum_value) = "5", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_FULL_LEGAL_NAME_OF_FIRM_VALUE = 11;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_FUND_ACCOUNT_NAME = 12 [(.fix.enum_value) = "19", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_FUND_ACCOUNT_NAME_VALUE = 12;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_LOCATION_DESK = 13 [(.fix.enum_value) = "25", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_LOCATION_DESK_VALUE = 13;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_PERSON = 14 [(.fix.enum_value) = "2", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_PERSON_VALUE = 14;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_PHONE_NUMBER = 15 [(.fix.enum_value) = "7", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_PHONE_NUMBER_VALUE = 15;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_POSITION_ACCOUNT_TYPE = 16 [(.fix.enum_value) = "26", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_POSITION_ACCOUNT_TYPE_VALUE = 16;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_POSTAL_ADDRESS = 17 [(.fix.enum_value) = "6", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_POSTAL_ADDRESS_VALUE = 17;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS = 18 [(.fix.enum_value) = "18", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS_VALUE = 18;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS_FOR_CONFIRMATION = 19 [(.fix.enum_value) = "12", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS_FOR_CONFIRMATION_VALUE = 19;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NAME = 20 [(.fix.enum_value) = "14", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NAME_VALUE = 20;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NUMBER = 21 [(.fix.enum_value) = "11", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NUMBER_VALUE = 21;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_REGULATORY_STATUS = 22 [(.fix.enum_value) = "13", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_REGULATORY_STATUS_VALUE = 22;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NAME = 23 [(.fix.enum_value) = "22", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NAME_VALUE = 23;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NUMBER = 24 [(.fix.enum_value) = "10", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NUMBER_VALUE = 24;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SYSTEM = 25 [(.fix.enum_value) = "3", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_SYSTEM_VALUE = 25;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_TELEX_NUMBER = 26 [(.fix.enum_value) = "20", (.fix.enum_added) = VERSION_FIX_4_4];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_TELEX_NUMBER_VALUE = 26;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_SECURITY_LOCATE_ID = 27 [(.fix.enum_value) = "27", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 1];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_SECURITY_LOCATE_ID_VALUE = 27;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_ELIGIBLE_COUNTERPARTY = 28 [(.fix.enum_value) = "29", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_ELIGIBLE_COUNTERPARTY_VALUE = 28;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_EXECUTION_VENUE = 29 [(.fix.enum_value) = "32", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_EXECUTION_VENUE_VALUE = 29;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_LOCATION = 30 [(.fix.enum_value) = "31", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_LOCATION_VALUE = 30;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_MARKET_MAKER = 31 [(.fix.enum_value) = "28", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_MARKET_MAKER_VALUE = 31;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_PROFESSIONAL_CLIENT = 32 [(.fix.enum_value) = "30", (.fix.enum_added) = VERSION_FIX_4_4, (.fix.enum_added_ep) = 26];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_PROFESSIONAL_CLIENT_VALUE = 32;
  /**
   * <code>SETTL_PARTY_SUB_ID_TYPE_CURRENCY_DELIVERY_IDENTIFIER = 33 [(.fix.enum_value) = "33", (.fix.enum_added) = VERSION_FIX_5_0, (.fix.enum_added_ep) = 44];</code>
   */
  public static final int SETTL_PARTY_SUB_ID_TYPE_CURRENCY_DELIVERY_IDENTIFIER_VALUE = 33;


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
  public static SettlPartySubIDTypeEnum valueOf(int value) {
    return forNumber(value);
  }

  /**
   * @param value The numeric wire value of the corresponding enum entry.
   * @return The enum associated with the given numeric wire value.
   */
  public static SettlPartySubIDTypeEnum forNumber(int value) {
    switch (value) {
      case 0: return SETTL_PARTY_SUB_ID_TYPE_UNSPECIFIED;
      case 1: return SETTL_PARTY_SUB_ID_TYPE_APPLICATION;
      case 2: return SETTL_PARTY_SUB_ID_TYPE_BIC;
      case 3: return SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NAME;
      case 4: return SETTL_PARTY_SUB_ID_TYPE_CASH_ACCOUNT_NUMBER;
      case 5: return SETTL_PARTY_SUB_ID_TYPE_CONTACT_NAME;
      case 6: return SETTL_PARTY_SUB_ID_TYPE_CSDPARTICIPANT_MEMBER_CODE;
      case 7: return SETTL_PARTY_SUB_ID_TYPE_DEPARTMENT;
      case 8: return SETTL_PARTY_SUB_ID_TYPE_EMAIL_ADDRESS;
      case 9: return SETTL_PARTY_SUB_ID_TYPE_FAX_NUMBER;
      case 10: return SETTL_PARTY_SUB_ID_TYPE_FIRM;
      case 11: return SETTL_PARTY_SUB_ID_TYPE_FULL_LEGAL_NAME_OF_FIRM;
      case 12: return SETTL_PARTY_SUB_ID_TYPE_FUND_ACCOUNT_NAME;
      case 13: return SETTL_PARTY_SUB_ID_TYPE_LOCATION_DESK;
      case 14: return SETTL_PARTY_SUB_ID_TYPE_PERSON;
      case 15: return SETTL_PARTY_SUB_ID_TYPE_PHONE_NUMBER;
      case 16: return SETTL_PARTY_SUB_ID_TYPE_POSITION_ACCOUNT_TYPE;
      case 17: return SETTL_PARTY_SUB_ID_TYPE_POSTAL_ADDRESS;
      case 18: return SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS;
      case 19: return SETTL_PARTY_SUB_ID_TYPE_REGISTERED_ADDRESS_FOR_CONFIRMATION;
      case 20: return SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NAME;
      case 21: return SETTL_PARTY_SUB_ID_TYPE_REGISTRATION_NUMBER;
      case 22: return SETTL_PARTY_SUB_ID_TYPE_REGULATORY_STATUS;
      case 23: return SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NAME;
      case 24: return SETTL_PARTY_SUB_ID_TYPE_SECURITIES_ACCOUNT_NUMBER;
      case 25: return SETTL_PARTY_SUB_ID_TYPE_SYSTEM;
      case 26: return SETTL_PARTY_SUB_ID_TYPE_TELEX_NUMBER;
      case 27: return SETTL_PARTY_SUB_ID_TYPE_SECURITY_LOCATE_ID;
      case 28: return SETTL_PARTY_SUB_ID_TYPE_ELIGIBLE_COUNTERPARTY;
      case 29: return SETTL_PARTY_SUB_ID_TYPE_EXECUTION_VENUE;
      case 30: return SETTL_PARTY_SUB_ID_TYPE_LOCATION;
      case 31: return SETTL_PARTY_SUB_ID_TYPE_MARKET_MAKER;
      case 32: return SETTL_PARTY_SUB_ID_TYPE_PROFESSIONAL_CLIENT;
      case 33: return SETTL_PARTY_SUB_ID_TYPE_CURRENCY_DELIVERY_IDENTIFIER;
      default: return null;
    }
  }

  public static com.google.protobuf.Internal.EnumLiteMap<SettlPartySubIDTypeEnum>
      internalGetValueMap() {
    return internalValueMap;
  }
  private static final com.google.protobuf.Internal.EnumLiteMap<
      SettlPartySubIDTypeEnum> internalValueMap =
        new com.google.protobuf.Internal.EnumLiteMap<SettlPartySubIDTypeEnum>() {
          public SettlPartySubIDTypeEnum findValueByNumber(int number) {
            return SettlPartySubIDTypeEnum.forNumber(number);
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
    return org.fixprotocol.components.Common.getDescriptor().getEnumTypes().get(120);
  }

  private static final SettlPartySubIDTypeEnum[] VALUES = values();

  public static SettlPartySubIDTypeEnum valueOf(
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

  private SettlPartySubIDTypeEnum(int value) {
    this.value = value;
  }

  // @@protoc_insertion_point(enum_scope:Common.SettlPartySubIDTypeEnum)
}


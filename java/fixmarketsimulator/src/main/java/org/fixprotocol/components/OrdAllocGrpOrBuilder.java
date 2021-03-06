// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

public interface OrdAllocGrpOrBuilder extends
    // @@protoc_insertion_point(interface_extends:Common.OrdAllocGrp)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>string cl_ord_id = 1 [(.fix.tag) = 11, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The clOrdId.
   */
  java.lang.String getClOrdId();
  /**
   * <code>string cl_ord_id = 1 [(.fix.tag) = 11, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for clOrdId.
   */
  com.google.protobuf.ByteString
      getClOrdIdBytes();

  /**
   * <code>string list_id = 2 [(.fix.tag) = 66, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The listId.
   */
  java.lang.String getListId();
  /**
   * <code>string list_id = 2 [(.fix.tag) = 66, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for listId.
   */
  com.google.protobuf.ByteString
      getListIdBytes();

  /**
   * <code>repeated .Common.NestedParties2 nested_parties2 = 3 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  java.util.List<org.fixprotocol.components.NestedParties2> 
      getNestedParties2List();
  /**
   * <code>repeated .Common.NestedParties2 nested_parties2 = 3 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  org.fixprotocol.components.NestedParties2 getNestedParties2(int index);
  /**
   * <code>repeated .Common.NestedParties2 nested_parties2 = 3 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  int getNestedParties2Count();
  /**
   * <code>repeated .Common.NestedParties2 nested_parties2 = 3 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  java.util.List<? extends org.fixprotocol.components.NestedParties2OrBuilder> 
      getNestedParties2OrBuilderList();
  /**
   * <code>repeated .Common.NestedParties2 nested_parties2 = 3 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  org.fixprotocol.components.NestedParties2OrBuilder getNestedParties2OrBuilder(
      int index);

  /**
   * <code>.fix.Decimal64 order_avg_px = 4 [(.fix.tag) = 799, (.fix.type) = DATATYPE_PRICE, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return Whether the orderAvgPx field is set.
   */
  boolean hasOrderAvgPx();
  /**
   * <code>.fix.Decimal64 order_avg_px = 4 [(.fix.tag) = 799, (.fix.type) = DATATYPE_PRICE, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The orderAvgPx.
   */
  org.fixprotocol.components.Fix.Decimal64 getOrderAvgPx();
  /**
   * <code>.fix.Decimal64 order_avg_px = 4 [(.fix.tag) = 799, (.fix.type) = DATATYPE_PRICE, (.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  org.fixprotocol.components.Fix.Decimal64OrBuilder getOrderAvgPxOrBuilder();

  /**
   * <code>.fix.Decimal64 order_booking_qty = 5 [(.fix.tag) = 800, (.fix.type) = DATATYPE_QTY, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return Whether the orderBookingQty field is set.
   */
  boolean hasOrderBookingQty();
  /**
   * <code>.fix.Decimal64 order_booking_qty = 5 [(.fix.tag) = 800, (.fix.type) = DATATYPE_QTY, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The orderBookingQty.
   */
  org.fixprotocol.components.Fix.Decimal64 getOrderBookingQty();
  /**
   * <code>.fix.Decimal64 order_booking_qty = 5 [(.fix.tag) = 800, (.fix.type) = DATATYPE_QTY, (.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  org.fixprotocol.components.Fix.Decimal64OrBuilder getOrderBookingQtyOrBuilder();

  /**
   * <code>string order_id = 6 [(.fix.tag) = 37, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The orderId.
   */
  java.lang.String getOrderId();
  /**
   * <code>string order_id = 6 [(.fix.tag) = 37, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for orderId.
   */
  com.google.protobuf.ByteString
      getOrderIdBytes();

  /**
   * <code>.fix.Decimal64 order_qty = 7 [(.fix.tag) = 38, (.fix.type) = DATATYPE_QTY, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return Whether the orderQty field is set.
   */
  boolean hasOrderQty();
  /**
   * <code>.fix.Decimal64 order_qty = 7 [(.fix.tag) = 38, (.fix.type) = DATATYPE_QTY, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The orderQty.
   */
  org.fixprotocol.components.Fix.Decimal64 getOrderQty();
  /**
   * <code>.fix.Decimal64 order_qty = 7 [(.fix.tag) = 38, (.fix.type) = DATATYPE_QTY, (.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  org.fixprotocol.components.Fix.Decimal64OrBuilder getOrderQtyOrBuilder();

  /**
   * <code>string secondary_cl_ord_id = 8 [(.fix.tag) = 526, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The secondaryClOrdId.
   */
  java.lang.String getSecondaryClOrdId();
  /**
   * <code>string secondary_cl_ord_id = 8 [(.fix.tag) = 526, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for secondaryClOrdId.
   */
  com.google.protobuf.ByteString
      getSecondaryClOrdIdBytes();

  /**
   * <code>string secondary_order_id = 9 [(.fix.tag) = 198, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The secondaryOrderId.
   */
  java.lang.String getSecondaryOrderId();
  /**
   * <code>string secondary_order_id = 9 [(.fix.tag) = 198, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for secondaryOrderId.
   */
  com.google.protobuf.ByteString
      getSecondaryOrderIdBytes();
}

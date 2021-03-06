// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

public interface NestedParties2OrBuilder extends
    // @@protoc_insertion_point(interface_extends:Common.NestedParties2)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>string nested2party_id = 1 [(.fix.tag) = 757, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The nested2partyId.
   */
  java.lang.String getNested2PartyId();
  /**
   * <code>string nested2party_id = 1 [(.fix.tag) = 757, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The bytes for nested2partyId.
   */
  com.google.protobuf.ByteString
      getNested2PartyIdBytes();

  /**
   * <code>.Common.Nested2PartyIDSourceEnum nested2party_id_source = 2 [(.fix.tag) = 758, (.fix.type) = DATATYPE_CHAR, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The enum numeric value on the wire for nested2partyIdSource.
   */
  int getNested2PartyIdSourceValue();
  /**
   * <code>.Common.Nested2PartyIDSourceEnum nested2party_id_source = 2 [(.fix.tag) = 758, (.fix.type) = DATATYPE_CHAR, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The nested2partyIdSource.
   */
  org.fixprotocol.components.Nested2PartyIDSourceEnum getNested2PartyIdSource();

  /**
   * <code>.Common.Nested2PartyRoleEnum nested2party_role = 3 [(.fix.tag) = 759, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The enum numeric value on the wire for nested2partyRole.
   */
  int getNested2PartyRoleValue();
  /**
   * <code>.Common.Nested2PartyRoleEnum nested2party_role = 3 [(.fix.tag) = 759, (.fix.type) = DATATYPE_INT, (.fix.field_added) = VERSION_FIX_4_4];</code>
   * @return The nested2partyRole.
   */
  org.fixprotocol.components.Nested2PartyRoleEnum getNested2PartyRole();

  /**
   * <code>repeated .Common.NstdPtys2SubGrp nstd_ptys2sub_grp = 4 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  java.util.List<org.fixprotocol.components.NstdPtys2SubGrp> 
      getNstdPtys2SubGrpList();
  /**
   * <code>repeated .Common.NstdPtys2SubGrp nstd_ptys2sub_grp = 4 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  org.fixprotocol.components.NstdPtys2SubGrp getNstdPtys2SubGrp(int index);
  /**
   * <code>repeated .Common.NstdPtys2SubGrp nstd_ptys2sub_grp = 4 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  int getNstdPtys2SubGrpCount();
  /**
   * <code>repeated .Common.NstdPtys2SubGrp nstd_ptys2sub_grp = 4 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  java.util.List<? extends org.fixprotocol.components.NstdPtys2SubGrpOrBuilder> 
      getNstdPtys2SubGrpOrBuilderList();
  /**
   * <code>repeated .Common.NstdPtys2SubGrp nstd_ptys2sub_grp = 4 [(.fix.field_added) = VERSION_FIX_4_4];</code>
   */
  org.fixprotocol.components.NstdPtys2SubGrpOrBuilder getNstdPtys2SubGrpOrBuilder(
      int index);
}

// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

package org.fixprotocol.components;

public interface DerivativeSecurityXMLOrBuilder extends
    // @@protoc_insertion_point(interface_extends:Common.DerivativeSecurityXML)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>string derivative_security_xml = 1 [(.fix.tag) = 1283, (.fix.type) = DATATYPE_DATA, (.fix.field_added) = VERSION_FIX_5_0];</code>
   * @return The derivativeSecurityXml.
   */
  java.lang.String getDerivativeSecurityXml();
  /**
   * <code>string derivative_security_xml = 1 [(.fix.tag) = 1283, (.fix.type) = DATATYPE_DATA, (.fix.field_added) = VERSION_FIX_5_0];</code>
   * @return The bytes for derivativeSecurityXml.
   */
  com.google.protobuf.ByteString
      getDerivativeSecurityXmlBytes();

  /**
   * <code>sfixed64 derivative_security_xml_len = 2 [(.fix.tag) = 1282, (.fix.type) = DATATYPE_LENGTH, (.fix.field_added) = VERSION_FIX_5_0];</code>
   * @return The derivativeSecurityXmlLen.
   */
  long getDerivativeSecurityXmlLen();

  /**
   * <code>string derivative_security_xml_schema = 3 [(.fix.tag) = 1284, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_5_0];</code>
   * @return The derivativeSecurityXmlSchema.
   */
  java.lang.String getDerivativeSecurityXmlSchema();
  /**
   * <code>string derivative_security_xml_schema = 3 [(.fix.tag) = 1284, (.fix.type) = DATATYPE_STRING, (.fix.field_added) = VERSION_FIX_5_0];</code>
   * @return The bytes for derivativeSecurityXmlSchema.
   */
  com.google.protobuf.ByteString
      getDerivativeSecurityXmlSchemaBytes();
}
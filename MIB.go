// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// MIB
const (
	SYSTEM_OID               string = "1.3.6.1.2.1.1"
	SYS_DESCR_OID            string = "1.3.6.1.2.1.1.1"
	SYS_UPTIME_OID           string = "1.3.6.1.2.1.1.3"
	SYS_CONTACT_OID          string = "1.3.6.1.2.1.1.4"
	SYS_NAME_OID             string = "1.3.6.1.2.1.1.5"
	SYS_LOCATION_OID         string = "1.3.6.1.2.1.1.6"
	SYS_SERVICES_OID         string = "1.3.6.1.2.1.1.7"
	INTERFACES_OID           string = "1.3.6.1.2.1.2"
	IF_NUMBER_OID            string = "1.3.6.1.2.1.2.1" // number of interfaces
	IF_TABLE_OID             string = "1.3.6.1.2.1.2.2"
	IP_ADDR_TABLE_OID        string = "1.3.6.1.2.1.4.20"
	IP_AD_ENT_ADDR_OID       string = "1.3.6.1.2.1.4.20.1.1"
	IP_AD_ENT_IF_INDEX       string = "1.3.6.1.2.1.4.20.1.2"
	IP_AD_ENT_NET_MASK       string = "1.3.6.1.2.1.4.20.1.3"
	IP_AD_ENT_BCAST_ADDR     string = "1.3.6.1.2.1.4.20.1.4"
	IP_AD_ENT_REASM_MAX_SIZE string = "1.3.6.1.2.1.4.20.1.5"
	IP_ROUTE_TABLE_OID       string = "1.3.6.1.2.1.4.21"
	IP_ROUTE_DEST_OID        string = "1.3.6.1.2.1.4.21.1.1"
	IP_ROUTE_IF_INDEX_OID    string = "1.3.6.1.2.1.4.21.1.2"
	IP_ROUTE_NEXT_HOP_OID    string = "1.3.6.1.2.1.4.21.1.7"
	IF_SPEED                 string = "1.3.6.1.2.1.2.2.1.5"
	IF_OUT_OCTETS            string = "1.3.6.1.2.1.2.2.1.16"
	IFPHYSADDRESS_OID        string = ".1.3.6.1.2.1.2.2.1.6"
)

// Asn1BER is the type of the SNMP PDU
type Asn1BER byte

// Asn1BER's - http://www.ietf.org/rfc/rfc1442.txt
const (
	END_OF_CONTENTS    Asn1BER = 0x00
	UNKOWN_TYPE        Asn1BER = 0x00
	BOOLEAN            Asn1BER = 0x01
	INTEGER            Asn1BER = 0x02
	BIT_STRING         Asn1BER = 0x03
	OCTET_STRING       Asn1BER = 0x04
	NULL               Asn1BER = 0x05
	OBJECT_IDENTIFIER  Asn1BER = 0x06
	OBJECT_DESCRIPTION Asn1BER = 0x07
	IP_ADDRESS         Asn1BER = 0x40
	COUNTER32          Asn1BER = 0x41
	GAUGE32            Asn1BER = 0x42
	TIME_TICKS         Asn1BER = 0x43
	OPAQUE             Asn1BER = 0x44
	NSAP_ADDRESS       Asn1BER = 0x45
	COUNTER64          Asn1BER = 0x46
	UINTEGER32         Asn1BER = 0x47
	OPAQUE_FLOAT       Asn1BER = 0x78
	OPAQUE_DOUBLE      Asn1BER = 0x79
	NO_SUCH_OBJECT     Asn1BER = 0x80
	NO_SUCH_INSTANCE   Asn1BER = 0x81
	END_OF_MIB_VIEW    Asn1BER = 0x82
)

package main

// MIB
const systemOID string = "1.3.6.1.2.1.1"
const sysDescrOID string = "1.3.6.1.2.1.1.1"
const sysUpTimeOID string = "1.3.6.1.2.1.1.3"
const sysContactOID string = "1.3.6.1.2.1.1.4"
const sysNameOID string = "1.3.6.1.2.1.1.5"
const sysLocationOID string = "1.3.6.1.2.1.1.6"
const sysServicesOID string = "1.3.6.1.2.1.1.7"
const interfacesOID string = "1.3.6.1.2.1.2"
const ifNumberOID string = "1.3.6.1.2.1.2.1" // number of interfaces
const ifTableOID string = "1.3.6.1.2.1.2.2"
const ipAddrTableOID string = "1.3.6.1.2.1.4.20"
const ipAdEntAddrOID string = "1.3.6.1.2.1.4.20.1.1"
const ipAdEntIfIndex string = "1.3.6.1.2.1.4.20.1.2"
const ipAdEntNetMask string = "1.3.6.1.2.1.4.20.1.3"
const ipAdEntBcastAddr string = "1.3.6.1.2.1.4.20.1.4"
const ipAdEntReasmMaxSize string = "1.3.6.1.2.1.4.20.1.5"
const ipRouteTableOID string = "1.3.6.1.2.1.4.21"
const ipRouteDestOID string = "1.3.6.1.2.1.4.21.1.1"
const ipRouteIfIndexOID string = "1.3.6.1.2.1.4.21.1.2"
const ipRouteNextHopOID string = "1.3.6.1.2.1.4.21.1.7"
const ifSpeed string = "1.3.6.1.2.1.2.2.1.5"
const ifOutOctets string = "1.3.6.1.2.1.2.2.1.16"

// Asn1BER is the type of the SNMP PDU
type Asn1BER byte

// Asn1BER's - http://www.ietf.org/rfc/rfc1442.txt
const (
	EndOfContents     Asn1BER = 0x00
	UnknownType       Asn1BER = 0x00
	Boolean           Asn1BER = 0x01
	Integer           Asn1BER = 0x02
	BitString         Asn1BER = 0x03
	OctetString       Asn1BER = 0x04
	Null              Asn1BER = 0x05
	ObjectIdentifier  Asn1BER = 0x06
	ObjectDescription Asn1BER = 0x07
	IPAddress         Asn1BER = 0x40
	Counter32         Asn1BER = 0x41
	Gauge32           Asn1BER = 0x42
	TimeTicks         Asn1BER = 0x43
	Opaque            Asn1BER = 0x44
	NsapAddress       Asn1BER = 0x45
	Counter64         Asn1BER = 0x46
	Uinteger32        Asn1BER = 0x47
	OpaqueFloat       Asn1BER = 0x78
	OpaqueDouble      Asn1BER = 0x79
	NoSuchObject      Asn1BER = 0x80
	NoSuchInstance    Asn1BER = 0x81
	EndOfMibView      Asn1BER = 0x82
)

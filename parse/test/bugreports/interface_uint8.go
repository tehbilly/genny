// Code generated by genny. DO NOT EDIT.
// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package bugreports

type InterfaceUint8 interface {
	DoSomthingUint8()
}

// Call calls a method on an instance of generic interface.
// Targets github issue 49
func CallWithUint8(i InterfaceUint8) {
	i.DoSomthingUint8()
}

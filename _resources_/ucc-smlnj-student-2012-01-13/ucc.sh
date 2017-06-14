#!/bin/sh
ARCH=unknown
OPSYS=unknown
case `uname -s` in
SunOS)	case `uname -r` in
	5.*)	OPSYS=solaris;;
	esac
	case `uname -m` in
	sun4*)	ARCH=sparc;;
	i86pc)	ARCH=x86;;
	esac;;
Linux)	OPSYS=linux
	case `uname -m` in
	i[3456]86)	ARCH=x86;;
	esac;;
esac
exec sml @SMLload=ucc.$ARCH-$OPSYS $*

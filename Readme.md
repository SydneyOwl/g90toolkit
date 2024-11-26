# G90Tools

## Introduction

This software allows users to:

+ Modify the embedded boot image/text in the G90 firmware.
+ Encrypt/Decrypt firmware use user provided key
+ Take a brief look at the firmware info
+ Flash firmware into G90 series Rigs (g90updatefw integrated)
+ ......

> [!important]  
> The official Xiegu firmware is encrypted using the AES-256 algorithm, making it nearly impossible to decrypt without
> the key (unless you own a quantum computer). Fortunately, there are open-source methods available to extract the
> encryption key using an ST-LINK debugger and OpenOCD tools. Due to copyright restrictions, I cannot provide the key
> here, but you can extract the encryption key and decrypt the firmware using the methods outlined
> in [G90Tools](https://github.com/OpenHamradioFirmware/G90Tools) (which also inspired the approach for modifying the
> startup images in this software) or by finding the key shared by others online.

## Usage

TODO

## Many thanks to ...

- [G90Tools](https://github.com/OpenHamradioFirmware/G90Tools) (kbeckmann, GitHub) *Tools and guides for analyzing Xiegu
  G90 firmware*
- [Bootloader extraction procedure from Xiegu G90 processors](https://radiochief.ru/radio/protsedura-izvlecheniya-bootloader-iz-xiegu-g90/) (
  Denis Dubov, Radiochief.ru magazine 06/2022) *Dumping firmware and bootloaders*
- [BBFW](https://github.com/fventuri/BBFW) (Franco Venturi, GitHub) *BBFW utilities and tools*
- [g90updatefw](https://github.com/DaleFarnsworth/g90updatefw) (Dale Farnsworth, GitHub)  *Xiegu G90 and Xiego G106
  Firmware Updater*

## Disclaimer

- No warranty is provided. Any damage caused by using this tool is your own responsibility.
- The purpose of this tool is to help users modify the startup screen more conveniently, rather than to harm Xiegu's
  interests.